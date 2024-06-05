// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
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
	Update(ctx context.Context, username string, order models.Order) error
	GetAll(ctx context.Context, username string) (models.Orders, error)
}

type BalanceManager interface {
	Accrue(ctx context.Context, accrue models.Accrue) error
}

type Settings struct {
	OrderRepo   OrderRepository
	Balance     BalanceManager
	AccrualAddr string
}

type OrderManager struct {
	repo        OrderRepository
	balance     BalanceManager
	accrualAddr string
	client      *resty.Client
	limiter     *rate.Limiter
	eg          *errgroup.Group
}

func New(settings Settings) *OrderManager {
	return &OrderManager{
		repo:        settings.OrderRepo,
		balance:     settings.Balance,
		accrualAddr: settings.AccrualAddr,
		client:      resty.New(),
		limiter:     rate.NewLimiter(rate.Inf, 0),
	}
}

func (om *OrderManager) Get(ctx context.Context, username string) (models.Orders, error) {
	return om.repo.GetAll(ctx, username)
}

func (om *OrderManager) Register(ctx context.Context, username string, orderID uint64) error {
	// Добавление заказа в репо
	return om.updateOrder(ctx, username, orderID, true)
}

func (om *OrderManager) updateOrder(
	ctx context.Context,
	username string,
	orderID uint64,
	async bool,
) error {
	// Рассчет баллов через внешнюю систему
	order, err := om.calculateOrder(ctx, orderID)
	if err != nil {
		return err
	}

	// Проверка статуса рассчета баллов
	if order.Status == models.StatusOrderNew || order.Status == models.StatusOrderProcessing {
		// Ждем какое-то время
		time.Sleep(1 * time.Second) // TODO: Костыль, переделать

		// Повторно обновляем заказ
		// в другой горутине, если необходимо
		if async {
			om.eg.Go(func() error {
				return om.updateOrder(ctx, username, orderID, false)
			})
		} else {
			return om.updateOrder(ctx, username, orderID, false)
		}
	} else if order.Status == models.StatusOrderProcessed && order.Accrual != nil {
		// Начисление баллов
		err = om.balance.Accrue(ctx, models.Accrue{
			Username: username,
			Amount:   *order.Accrual,
			OrderID:  orderID,
		})
		if err != nil {
			return err
		}
	}

	// Обновляем статус в репо
	return om.repo.Update(ctx, username, order)
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
