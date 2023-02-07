package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"myapp/internal/config"
	api "myapp/internal/controller/http"
	stan "myapp/internal/controller/stan"
	cache "myapp/internal/repository/in_memory_cache"
	"myapp/internal/repository/postgres"
	"myapp/internal/service"
)

func main() {
	config := config.Get()

	pg, err := postgres.New(config) // БД
	if err != nil {
		log.Fatal(err)
	}
	defer logFnError(pg.CloseDb)

	cacheStore, err := cache.New(pg) // in-memory кэш
	if err != nil {
		log.Fatal(err)
	}

	service, err := service.New(cacheStore) // сервис общается с репозиторием
	if err != nil {
		log.Fatal()
	}

	natsListener, err := stan.New(config, service) // STAN
	if err != nil {
		log.Fatal()
	}
	defer logFnError(natsListener.Close)

	server, err := api.New(config.HttpServer.Port, service) // api server
	if err != nil {
		log.Fatal()
	}
	go logFnError(server.ListenAndServe)
	log.Printf("HTTP server is listening at %s", config.HttpServer.Port)

	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint
	log.Println("Service stopped...")
}

// Функция принимает аргументом функцию которая может вернуть error
// и логирует error
func logFnError(fn func() error) {
	if err := fn(); err != nil {
		log.Println(err)
	}
}
