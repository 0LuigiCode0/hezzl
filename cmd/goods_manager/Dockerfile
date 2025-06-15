FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o service ./cmd/goods_manager

CMD ["./service"]