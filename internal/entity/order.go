package entity

// OrderEntity - сущность заказа в БД.
type OrderEntity struct {
	orderUID string
	jsonData string
}

func NewOrderEntity(orderUID, jsonData string) OrderEntity {
	o := OrderEntity{
		orderUID: orderUID,
		jsonData: jsonData,
	}
	return o
}

func (oe *OrderEntity) GetUIDPtr() *string {
	return &oe.orderUID
}

func (oe *OrderEntity) GetUID() string {
	return oe.orderUID
}

func (oe *OrderEntity) SetUID(uid string) {
	oe.orderUID = uid
}

func (oe *OrderEntity) GetJsonDataPtr() *string {
	return &oe.jsonData
}

func (oe *OrderEntity) GetJsonData() string {
	return oe.jsonData
}

func (oe *OrderEntity) SetJsonData(jd string) {
	oe.jsonData = jd
}
