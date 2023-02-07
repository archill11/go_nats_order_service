package http

import (
	"net/http"
	"time"

	"myapp/internal/service"

	"github.com/gorilla/mux"
)

type APIServer struct {
	http.Server
	service *service.OrderService
	router  *mux.Router
}

func New(port string, service *service.OrderService) (*APIServer, error) {

	router := mux.NewRouter()
	server := &APIServer{
		Server: http.Server{
			Addr:         ":" + port,
			Handler:      router,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 90 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		service: service,
		router:  router,
	}

	server.router.HandleFunc("/", server.getAllOrders).Methods(http.MethodGet)      // GET /
	server.router.HandleFunc("/{uid}", server.getOrderById).Methods(http.MethodGet) // GET /uid

	return server, nil
}
