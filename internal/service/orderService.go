package service

import (
	"encoding/json"
	"log"
	"myapp/internal/entity"
	"myapp/internal/model"
	"myapp/internal/repository"
	"sync"
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

// AddOrder Метод передает заказ в репозиторий
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
		dataOrder := dataOrder.JsonData
		if err := json.Unmarshal([]byte(dataOrder), &order); err != nil {
			panic("unreachable: " + err.Error())
		}
		orders = append(orders, order)
	}
	return orders
}

var dataPool = sync.Pool{
	New: func() interface{} {
		return new(model.Order)
	},
}
func (srv *OrderService) parseOrderToJsonWithPool(data ...entity.OrderEntity) []model.Order {
	orders := make([]model.Order, 0, len(data))
	for _, dataOrder := range data {
		order := dataPool.Get().(*model.Order)
		order.OrderUID = ""
		if err := json.Unmarshal([]byte(dataOrder.JsonData), order); err != nil {
			panic("unreachable: " + err.Error())
		}
		orders = append(orders, *order)
		dataPool.Put(order)
	}
	return orders
}

func (os *OrderService) GetAllOrdersWithPool() ([]model.Order, error) {
	data, err := os.repo.GetAll()
	if err != nil {
		log.Printf("OrderService: could not get records from the storage: %s", err)
		return nil, err
	}
	orders := os.parseOrderToJsonWithPool(data...)
	return orders, nil
}
