up:
	docker-compose up
.PHONY:up

stop:
	docker-compose stop
.PHONY:stop

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

test:
	go test ./... -v
.PHONY:test

run-serv:
	go run ./cmd/orderserver/main.go
.PHONY:run-serv

stan-push:
	go run ./cmd/orderpub/main.go
.PHONY:stan-push

stan-push-file:
	go run ./cmd/orderpub/main.go ./mocks/model.json
.PHONY:stan-push-file