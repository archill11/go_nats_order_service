package entity

// OrderEntity - сущность заказа в БД.
type OrderEntity struct {
	OrderUID string
	JsonData string
}

func NewOrderEntity(orderUID, jsonData string) OrderEntity {
	o := OrderEntity{
		OrderUID: orderUID,
		JsonData: jsonData,
	}
	return o
}

