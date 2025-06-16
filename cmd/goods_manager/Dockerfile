FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o service ./cmd/goods_manager

FROM alpine:latest

WORKDIR /
COPY --from=builder /build/service .

CMD ["./service"]