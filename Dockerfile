FROM golang:1.20-alpine

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o my-app ./cmd/orderserver/main.go

CMD ["./my-app"]