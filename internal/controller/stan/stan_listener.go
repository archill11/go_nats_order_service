package stan

import (
	"encoding/json"
	"log"
	"myapp/internal/config"
	"myapp/internal/model"
	"myapp/internal/service"

	"github.com/hashicorp/go-multierror"
	"github.com/nats-io/stan.go"
)

type StanListener struct {
	stanConn stan.Conn
	sub      stan.Subscription
	service  *service.OrderService
}

func New(conf config.Config, service *service.OrderService) (StanListener, error) {

	stanConf := conf.StanListener

	stanConn, err := stan.Connect(stanConf.ClusterName, stanConf.ClientId)
	if err != nil {
		return StanListener{}, err
	}

	sl := StanListener{
		stanConn: stanConn,
		service:  service,
	}

	sub, err := stanConn.Subscribe(stanConf.Subject, sl.msgHandler, stan.DurableName(stanConf.DurableName))
	if err != nil {
		return StanListener{}, err
	}
	sl.sub = sub

	log.Println("STAN started")
	return sl, nil
}

// Close Метод закрывает потоковую подписку и соединение nats.
func (nl StanListener) Close() (retErr error) {
	if err := nl.sub.Close(); err != nil {
		retErr = multierror.Append(retErr, err)
	}
	if err := nl.stanConn.Close(); err != nil {
		retErr = multierror.Append(retErr, err)
	}
	return
}

// Метод отправляет все входящие заказы в хранилище.
func (nl StanListener) msgHandler(msg *stan.Msg) {
	var order model.Order

	if err := json.Unmarshal(msg.Data, &order); err != nil {
		log.Printf("StanListener: ERR: order rejected: incorrect order type: %s", err)
		return
	}
	if err := order.Validate(); err != nil {
		log.Printf("StanListener: ERR: order rejected: invalid order: %s", err)
		return
	}
	if err := nl.service.AddOrder(order.OrderUID, string(msg.Data)); err != nil {
		log.Printf("StanListener: ERR: order rejected: %s", err)
		return
	}
	log.Printf("StanListener: order %q received and transfer to repository", order.OrderUID)
}
