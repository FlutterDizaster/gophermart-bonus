// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License.
//
// This file uses a third-party package chi licensed under MIT License.
// This file uses a third-party package sync licensed under BSD-3-Clause License.
// This file uses a third-party package time licensed under BSD-3-Clause License.
//
// See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package ordermanager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

type OrderRepository interface {
	AddOrder(ctx context.Context, userID uint64, order models.Order) error
	UpdateOrder(ctx context.Context, order models.Order) error
	GetAllOrders(ctx context.Context, userID uint64) (models.Orders, error)
	GetNotUpdatedOrders(ctx context.Context) ([]uint64, error)
}

type Settings struct {
	OrderRepo   OrderRepository
	AccrualAddr string
}

type OrderManager struct {
	repo        OrderRepository
	accrualAddr string
	client      *resty.Client
	limiter     *rate.Limiter
	wg          *sync.WaitGroup
	errCh       chan error
	updateCtx   context.Context
}

func New(settings Settings) *OrderManager {
	om := &OrderManager{
		repo:        settings.OrderRepo,
		accrualAddr: settings.AccrualAddr,
		client:      resty.New(),
		limiter:     rate.NewLimiter(rate.Inf, 0),
		wg:          &sync.WaitGroup{},
		errCh:       make(chan error, 1),
	}

	return om
}

func (om *OrderManager) Start(ctx context.Context) error {
	slog.Debug("Starting order manager service")

	om.updateCtx = ctx

	updateList, err := om.repo.GetNotUpdatedOrders(ctx)
	if err != nil {
		slog.Error("order manager error", slog.Any("error", err))
		return err
	}

	for i := range updateList {
		orderID := updateList[i]
		// Создание задачи на обновление метрики
		om.wg.Add(1)
		go func() {
			err = om.sheduleUpdate(om.updateCtx, orderID)
			if err != nil {
				om.errCh <- err
			}
			om.wg.Done()
		}()
	}

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case err = <-om.errCh:
			slog.Error("order manager error", slog.Any("error", err))
		}
	}

	om.wg.Wait()
	return nil
}

func (om *OrderManager) Get(ctx context.Context, userID uint64) (models.Orders, error) {
	slog.Debug("getting orders for user", slog.Uint64("user id", userID))
	orders, err := om.repo.GetAllOrders(ctx, userID)

	for i, order := range orders {
		order.StringID = strconv.FormatUint(order.ID, 10)
		orders[i] = order
	}

	return orders, err
}

func (om *OrderManager) Register(ctx context.Context, userID uint64, orderID uint64) error {
	slog.Debug(
		"Adding order to the repo",
		slog.Uint64("order id", orderID),
		slog.Uint64("user id", userID),
	)
	// Добавление заказа в репо
	order := models.Order{
		ID:     orderID,
		Status: models.StatusOrderNew,
	}

	err := om.repo.AddOrder(ctx, userID, order)
	if err != nil {
		return err
	}

	// Создание задачи на обновление метрики
	om.wg.Add(1)
	go func() {
		err = om.sheduleUpdate(om.updateCtx, orderID)
		if err != nil {
			om.errCh <- err
		}
		om.wg.Done()
	}()

	return nil
}

func (om *OrderManager) sheduleUpdate(ctx context.Context, orderID uint64) error {
	slog.Debug("sheduling order update", slog.Uint64("order id", orderID))
	// Рассчет баллов через внешнюю систему
	order, err := om.calculateOrder(ctx, orderID)
	if err != nil {
		slog.Error("error calculating order", slog.Any("error", err))
		return err
	}
	slog.Debug("Order calculated", slog.String("status", string(order.Status)))

	// Проверка статуса рассчета баллов
	if order.Status == models.StatusOrderNew || order.Status == models.StatusOrderProcessing {
		// Повторно обновляем заказ
		return om.sheduleUpdate(ctx, orderID)
	}

	// Обновляем статус в репо
	err = om.repo.UpdateOrder(ctx, order)
	if err != nil {
		slog.Error("error updating order", slog.Any("error", err))
	}
	return err
}

func (om *OrderManager) calculateOrder(ctx context.Context, orderID uint64) (models.Order, error) {
	slog.Debug("Calculating order")
	// Ожидание свободного токена лимитера
	err := om.limiter.Wait(ctx)
	if err != nil {
		return models.Order{}, err
	}

	// Подготовка запроса
	req := om.client.R().SetContext(ctx)

	// Запрос к сервису рассчета балов
	resp, err := req.Get(fmt.Sprintf("http://%s/api/orders/%d", om.accrualAddr, orderID))
	if err != nil {
		return models.Order{}, err
	}

	// Проверка статуса ответа
	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusTooManyRequests {
			// Обновление лимитера
			err = om.updateLimiter(resp)
			if err != nil {
				return models.Order{}, err
			}

			// повтор запроса с увеличенным лимитом
			return om.calculateOrder(ctx, orderID)
		}
		return models.Order{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// Unmarshal ответа
	var order models.Order
	err = order.UnmarshalJSON(resp.Body())
	if err != nil {
		return models.Order{}, err
	}

	// фикс для несоответствия имен для поля order.id в приложении и в сервисе accrue
	order.ID = orderID

	return order, nil
}

func (om *OrderManager) updateLimiter(resp *resty.Response) error {
	// Установка нулевого лимита, чтобы приостановить запросы
	om.limiter.SetLimit(0)

	// Получение времени ожидания перед повторным запросом
	dealyStr := resp.Header().Get("Retry-After")
	delay, err := strconv.ParseInt(dealyStr, 10, 64)
	if err != nil {
		return err
	}

	// Парсинг тела ответа для установки новых параметров лимитера
	if len(resp.Body()) == 0 {
		return errors.New("empty body")
	}
	body := string(resp.Body())

	// Получение кол-ва запросов в период period из тела ответа
	body, ok := strings.CutPrefix(body, "No more than ")
	if !ok {
		return errors.New("failed to cut prefix")
	}
	body, ok = strings.CutSuffix(body, " requests per minute allowed")
	if !ok {
		return errors.New("failed to cut suffix")
	}
	requests, err := strconv.ParseInt(body, 10, 64)
	if err != nil {
		return err
	}

	// Установка нового лимита
	om.limiter.SetBurst(1)
	om.limiter.SetLimitAt(
		time.Now().Add(time.Duration(delay)*time.Second),
		rate.Every(time.Minute/time.Duration(requests)),
	)

	return nil
}
