package main

// orderpub публикует заказы в nats cluster

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

const (
	clusterName = "cluster-L0"
	clientID    = "orderPub"
	subject     = "orders"
)

func main() {
	sc, err := stan.Connect(clusterName, clientID)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()
	log.Println("Successfully connected to nats-streaming")

	var order string
	if len(os.Args) < 2 {
		pushAlotData(sc, 1, 100)
	} else {
		order = readFromFile()
		stanPublish(sc, order)
	}
}

func stanPublish(sc stan.Conn, order string) {
	if err := sc.Publish(subject, []byte(order)); err != nil {
		log.Fatal(err)
	}
	log.Println("Order sent")
}

func readFromFile() string {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	res, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order successfully read from the file %s\n", os.Args[1])
	return string(res)
}

func pushAlotData(sc stan.Conn, countStart int, countFinish int) {
	for i := countStart; i <= countFinish; i++ {
		mockOrder := fmt.Sprintf(
			`{"order_uid":"%d__563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"b563feb7b2b84b6test","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":"meest","shardkey":"9","sm_id":99,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`,
			i,
		)
		stanPublish(sc, mockOrder)
	}
}
