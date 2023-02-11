package service

import (
	"encoding/json"
	"log"
	"myapp/internal/entity"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepo
}

func New(repository repository.OrderRepo) (*OrderService, error) {
	r := &OrderService{
		repo: repository,
	}
	return r, nil
}

// AddOrder Метод передает передает заказ в репозиторий
func (os *OrderService) GetAllOrders() ([]model.Order, error) {
	data, err := os.repo.GetAll()
	if err != nil {
		log.Printf("OrderService: could not get records from the storage: %s", err)
		return nil, err
	}
	orders := os.parseOrderToJson(data...)
	return orders, nil
}

// AddOrder Метод передает передает заказ в репозиторий
func (os *OrderService) GetOrderById(uid string) ([]model.Order, error) {
	data, err := os.repo.Get(uid)
	if err != nil {
		log.Printf("OrderService: could not get order %s: %s", uid, err)
		return nil, err
	}
	order := os.parseOrderToJson(data)
	return order, nil
}

// AddOrder Метод передает передает заказ в репозиторий
func (os *OrderService) AddOrder(orderUID, jsonOrder string) error {
	err := os.repo.AddOrder(orderUID, jsonOrder)
	if err != nil {
		return err
	}
	return nil
}

// Метод парсит заказы в Json
func (srv *OrderService) parseOrderToJson(data ...entity.OrderEntity) []model.Order {
	orders := make([]model.Order, 0, len(data))
	for _, dataOrder := range data {
		var order model.Order
		dataOrderUID := dataOrder.GetJsonData()

		if err := json.Unmarshal([]byte(dataOrderUID), &order); err != nil {
			panic("unreachable: " + err.Error())
		}
		orders = append(orders, order)
	}
	return orders
}
