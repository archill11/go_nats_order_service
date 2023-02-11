package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"myapp/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
)

// Метод возвращает все заказы
func (srv *APIServer) getAllOrders(rw http.ResponseWriter, r *http.Request) {

	orders, err := srv.service.GetAllOrders()
	if err != nil {
		http.Error(rw, "Something went wrong...", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(orders)
}

// Метод возвращает заказы по uid
func (srv *APIServer) getOrderById(rw http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"] // http param
	if !ok {
		log.Println("APIServer: unreachable: no key provided")
		http.Error(rw, "Something went wrong...", http.StatusInternalServerError)
		return
	}
	order, err := srv.service.GetOrderById(uid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(rw, fmt.Sprintf("Order %s not found.", uid), http.StatusNotFound)
			return
		}
		http.Error(rw, "Something went wrong...", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(order)
}
