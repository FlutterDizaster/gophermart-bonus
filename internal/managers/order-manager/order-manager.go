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
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
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
	eg          *errgroup.Group
}

func New(settings Settings) *OrderManager {
	om := &OrderManager{
		repo:        settings.OrderRepo,
		accrualAddr: settings.AccrualAddr,
		client:      resty.New(),
		limiter:     rate.NewLimiter(rate.Inf, 0),
	}

	return om
}

func (om *OrderManager) Start(ctx context.Context) error {
	updateList, err := om.repo.GetNotUpdatedOrders(ctx)
	if err != nil {
		return err
	}

	for i := range updateList {
		orderID := updateList[i]
		// Создание задачи на обновление метрики
		om.eg.Go(func() error {
			return om.sheduleUpdate(ctx, orderID)
		})
	}

	<-ctx.Done()
	return om.eg.Wait()
}

func (om *OrderManager) Get(ctx context.Context, userID uint64) (models.Orders, error) {
	return om.repo.GetAllOrders(ctx, userID)
}

func (om *OrderManager) Register(ctx context.Context, userID uint64, orderID uint64) error {
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
	om.eg.Go(func() error {
		return om.sheduleUpdate(ctx, orderID)
	})

	return nil
}

func (om *OrderManager) sheduleUpdate(ctx context.Context, orderID uint64) error {
	// Рассчет баллов через внешнюю систему
	order, err := om.calculateOrder(ctx, orderID)
	if err != nil {
		return err
	}

	// Проверка статуса рассчета баллов
	if order.Status == models.StatusOrderNew || order.Status == models.StatusOrderProcessing {
		// Повторно обновляем заказ
		return om.sheduleUpdate(ctx, orderID)
	}

	// Обновляем статус в репо
	return om.repo.UpdateOrder(ctx, order)
}

func (om *OrderManager) calculateOrder(ctx context.Context, orderID uint64) (models.Order, error) {
	// Ожидание свободного токена лимитера
	err := om.limiter.Wait(ctx)
	if err != nil {
		return models.Order{}, err
	}

	// Подготовка запроса
	req := om.client.R().SetContext(ctx)

	// Запрос к сервису рассчета балов
	resp, err := req.Get(fmt.Sprintf("%s/api/orders/%d", om.accrualAddr, orderID))
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
	om.limiter.SetLimitAt(
		time.Now().Add(time.Duration(delay)*time.Second),
		rate.Every(time.Minute/time.Duration(requests)),
	)

	return nil
}
