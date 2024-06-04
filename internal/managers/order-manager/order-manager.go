package ordermanager

import (
	"context"
	"fmt"

	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
	"github.com/go-resty/resty/v2"
)

type OrderRepository interface {
	Add(ctx context.Context, username string, order models.Order) error
	GetAll(ctx context.Context, username string) (models.Orders, error)
}

type Settings struct {
	OrderRepo   OrderRepository
	AccrualAddr string
}

type OrderManager struct {
	repo        OrderRepository
	accrualAddr string
	client      *resty.Client
}

func New(settings Settings) *OrderManager {
	return &OrderManager{
		repo:        settings.OrderRepo,
		accrualAddr: settings.AccrualAddr,
		client:      resty.New(),
	}
}

func (om *OrderManager) Register(ctx context.Context, username string, orderID uint64) error {
	// Рассчет баллов через внешнюю систему
	// TODO: Так как рассчет балов может происходить непредвиденное время,
	// необходимо сделать механизм повтора запросов перед добавлением заказа в репозиторий.
	// Либо добавлять заказ с промежуточным статусом, а потом изменять его по окончании рассчетов.
	order, err := om.calculateOrder(orderID)
	if err != nil {
		return err
	}

	// Добавление заказа в репо
	return om.repo.Add(ctx, username, order)
}

func (om *OrderManager) Get(ctx context.Context, username string) (models.Orders, error) {
	return om.repo.GetAll(ctx, username)
}

func (om *OrderManager) calculateOrder(orderID uint64) (models.Order, error) {
	// TODO: Сервис рассчета баллов имеет ограничения на кол-во запросов в минуту.
	// Необходимо налету понимать ограничение и изменять лимитер.
	// Запрос к сервису рассчета балов
	resp, err := om.client.R().Get(fmt.Sprintf("%s/api/orders/%d", om.accrualAddr, orderID))
	if err != nil {
		return models.Order{}, err
	}

	// Unmarshal ответа
	var order models.Order
	err = order.UnmarshalJSON(resp.Body())
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}
