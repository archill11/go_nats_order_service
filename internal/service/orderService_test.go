package service

import (
	"myapp/internal/entity"
	"myapp/internal/model"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_parseOrderToJson(t *testing.T) {
	service := &OrderService{}
	ord := entity.NewOrderEntity(
		"b563feb7b2b84b6test",
		`{"entry": "WBIL", "items": [{"rid": "ab4219087a764ae0btest", "name": "Mascaras", "sale": 30, "size": "0", "brand": "Vivienne Sabo", "nm_id": 2389212, "price": 453, "status": 202, "chrt_id": 9934930, "total_price": 317, "track_number": "WBILMTESTTRACK"}], "sm_id": 99, "locale": "en", "payment": {"bank": "alpha", "amount": 1817, "currency": "USD", "provider": "wbpay", "custom_fee": 0, "payment_dt": 1637907727, "request_id": "", "goods_total": 317, "transaction": "b563feb7b2b84b6test", "delivery_cost": 1500}, "delivery": {"zip": "2639809", "city": "Kiryat Mozkin", "name": "Test Testov", "email": "test@gmail.com", "phone": "+9720000000", "region": "Kraiot", "address": "Ploshad Mira 15"}, "shardkey": "9", "oof_shard": "1", "order_uid": "b563feb7b2b84b6test", "customer_id": "test", "date_created": "2021-11-26T06:22:19Z", "track_number": "WBILMTESTTRACK", "delivery_service": "meest", "internal_signature": ""}`,
	)
	modelFromEntity := service.parseOrderToJson(ord)
	timeCreated, _ := time.Parse("2006-01-02T15:04:05Z", "2021-11-26T06:22:19Z")
	modelOrder := model.Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			ZIP:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []model.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       timeCreated,
		OOFShard:          "1",
	}

	if diff := cmp.Diff(modelFromEntity[0], modelOrder); diff != "" {
		t.Error(diff)
	}

}
