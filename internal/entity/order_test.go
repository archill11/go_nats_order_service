package entity

import "testing"

func Test_OrderEntityGetters(t *testing.T) {
	mockUID := "b563feb7b2b84b6test"
	mockJsonData := `{"entry": "WBIL", "items": [{"rid": "ab4219087a764ae0btest", "name": "Mascaras", "sale": 30, "size": "0", "brand": "Vivienne Sabo", "nm_id": 2389212, "price": 453, "status": 202, "chrt_id": 9934930, "total_price": 317, "track_number": "WBILMTESTTRACK"}], "sm_id": 99, "locale": "en", "payment": {"bank": "alpha", "amount": 1817, "currency": "USD", "provider": "wbpay", "custom_fee": 0, "payment_dt": 1637907727, "request_id": "", "goods_total": 317, "transaction": "b563feb7b2b84b6test", "delivery_cost": 1500}, "delivery": {"zip": "2639809", "city": "Kiryat Mozkin", "name": "Test Testov", "email": "test@gmail.com", "phone": "+9720000000", "region": "Kraiot", "address": "Ploshad Mira 15"}, "shardkey": "9", "oof_shard": "1", "order_uid": "b563feb7b2b84b6test", "customer_id": "test", "date_created": "2021-11-26T06:22:19Z", "track_number": "WBILMTESTTRACK", "delivery_service": "meest", "internal_signature": ""}`

	ord := NewOrderEntity(
		mockUID,
		mockJsonData,
	)
	if ord.OrderUID != mockUID {
		t.Errorf("order UID %v not assert mock UID %v", ord.OrderUID, mockUID)
	}
	if ord.JsonData != mockJsonData {
		t.Errorf("order JsonData %v not assert mock JsonData  %v", ord.JsonData, mockJsonData)
	}

	changedMockUID := "89__b563feb7b2b84b6test"
	changedMockJsonData := `{"entry": "WBIL", "items": [{"rid": "88__ab4219087a764ae0btest", "name": "Mascaras", "sale": 55, "size": "0", "brand": "Vivienne Sabo", "nm_id": 2389212, "price": 453, "status": 202, "chrt_id": 9934930, "total_price": 317, "track_number": "WBILMTESTTRACK"}], "sm_id": 99, "locale": "en", "payment": {"bank": "alpha", "amount": 1817, "currency": "USD", "provider": "wbpay", "custom_fee": 0, "payment_dt": 1637907727, "request_id": "", "goods_total": 317, "transaction": "b563feb7b2b84b6test", "delivery_cost": 1500}, "delivery": {"zip": "2639809", "city": "Kiryat Mozkin", "name": "Test Testov", "email": "test2@gmail.com", "phone": "+9720000000", "region": "Kraiot", "address": "Ploshad Mira 15"}, "shardkey": "9", "oof_shard": "1", "order_uid": "b563feb7b2b84b6test", "customer_id": "test", "date_created": "2021-11-26T06:22:19Z", "track_number": "WBILMTESTTRACK", "delivery_service": "meest", "internal_signature": ""}`

	ord.OrderUID = changedMockUID
	if ord.OrderUID != changedMockUID {
		t.Errorf("%v not assert %v after set SetUID mrthod", ord.OrderUID, changedMockUID)
	}
	ord.JsonData = changedMockJsonData
	if ord.JsonData != changedMockJsonData {
		t.Errorf("order JsonData %v not assert mock JsonData  %v after set JsonData mrthod", ord.JsonData, changedMockJsonData)
	}
}
