package model

import (
	"errors"
	"net/mail"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
)

type (
	// Order DTO
	Order struct {
		OrderUID          string    `json:"order_uid"`
		TrackNumber       string    `json:"track_number"`
		Entry             string    `json:"entry"`
		Delivery          Delivery  `json:"delivery"`
		Payment           Payment   `json:"payment"`
		Items             []Item    `json:"items"`
		Locale            string    `json:"locale"`
		InternalSignature string    `json:"internal_signature"`
		Shardkey          string    `json:"shardkey"`
		SmID              int       `json:"sm_id"`
		DateCreated       time.Time `json:"date_created"`
		OOFShard          string    `json:"oof_shard"`
	}

	Delivery struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		ZIP     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	}

	Payment struct {
		Transaction  string `json:"transaction"`
		RequestID    string `json:"request_id"`
		Currency     string `json:"currency"`
		Provider     string `json:"provider"`
		Amount       int    `json:"amount"`
		PaymentDt    uint64 `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	}

	Item struct {
		ChrtID      uint64 `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		RID         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int    `json:"total_price"`
		NmID        uint64 `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int    `json:"status"`
	}
)

var (
	ErrEmptyOrderUID                = errors.New("empty order UID")
	ErrEmptyPaymentTransactionField = errors.New("empty payment transaction field")
	ErrInvalidEmail                 = errors.New("incorrect delivery email")
	ErrInvalidPhoneNumbe            = errors.New("incorrect delivery phone numbe")
)

func (o *Order) Validate() error {
	phoneNumRegex := regexp.MustCompile(`^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$`)

	var errs error
	if o.OrderUID == "" {
		errs = multierror.Append(errs, ErrEmptyOrderUID)
	}
	if o.Payment.Transaction == "" {
		errs = multierror.Append(errs, ErrEmptyPaymentTransactionField)
	}
	if _, err := mail.ParseAddress(o.Delivery.Email); err != nil {
		errs = multierror.Append(errs, ErrInvalidEmail)
	}
	if !phoneNumRegex.MatchString(o.Delivery.Phone) {
		errs = multierror.Append(errs, ErrInvalidPhoneNumbe)
	}
	return errs
}
