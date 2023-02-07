package model

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func Test_Validate(t *testing.T) {

	timeCreated, _ := time.Parse("2006-01-02T15:04:05Z", "2021-11-26T06:22:19Z")
	validModelOrder := Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			ZIP:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: Payment{
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
		Items: []Item{
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

	err := validModelOrder.Validate()
	if err != nil {
		t.Error(err)
	}

}

func Test_ValidateError(t *testing.T) {

	invalidPhoneNumber := Order{
		OrderUID: "b563feb7b2b84b6test",
		Delivery: Delivery{
			Phone: "020000000", // <- invalid phone number
			Email: "test@gmail.com",
		},
	}
	invalidEmail := Order{
		OrderUID: "b563feb7b2b84b6test",
		Delivery: Delivery{
			Phone: "+9720000000",
			Email: "testg456maom", // <- invalid Email
		},
	}
	emptyOrderUID := Order{
		OrderUID: "", // <- invalid order UID
		Delivery: Delivery{
			Phone: "+9720000000",
			Email: "test@gmail.com",
		},
	}
	emptyPaymentTransactionField := Order{
		OrderUID: "b563feb7b2b84b6test",
		Delivery: Delivery{
			Phone: "+9720000000",
			Email: "test@gmail.com",
		},
		Payment: Payment{
			Transaction: "", // <- invalid payment transaction field
		},
	}

	cases := []struct {
		name   string
		in     *Order
		expErr error
	}{
		{
			name:   "bad_phone_number",
			in:     &invalidPhoneNumber,
			expErr: ErrInvalidPhoneNumbe,
		},
		{
			name:   "bad_email",
			in:     &invalidEmail,
			expErr: ErrInvalidEmail,
		},
		{
			name:   "bad_payment_transaction_field",
			in:     &emptyPaymentTransactionField,
			expErr: ErrEmptyPaymentTransactionField,
		},
		{
			name:   "bad_order_UID",
			in:     &emptyOrderUID,
			expErr: ErrEmptyOrderUID,
		},
	}

	for i, tCase := range cases {
		err := tCase.in.Validate()
		if !errors.Is(err, tCase.expErr) {
			fmt.Println(i)
			t.Errorf("invalid order model passed validation! %v", err)
		}
	}

}
