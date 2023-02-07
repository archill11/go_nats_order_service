package repository

import (
	"errors"
	"myapp/internal/entity"
)

// OrderRepo - хранилище заказов.
type OrderRepo interface {
	AddOrder(orderUID, jsonOrder string) error       // добавить заказ
	Get(orderUID string) (entity.OrderEntity, error) // получить заказ по id
	GetAll() ([]entity.OrderEntity, error)           // получить все заказы
}

var (
	ErrAlreadyExists = errors.New("order already exists")
	ErrNotFound      = errors.New("order not found")
)
