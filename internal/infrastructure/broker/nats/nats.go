package bnats

import (
	"context"
	"fmt"
	"log"

	"github.com/0LuigiCode0/hezzl/config"
	dclickhouse "github.com/0LuigiCode0/hezzl/internal/domain/clickhouse"
	"github.com/0LuigiCode0/hezzl/internal/utils"
	"github.com/nats-io/nats.go"
)

type _nats struct {
	conn *nats.Conn
}

type INats interface {
	Push(v *dclickhouse.LogEventGood) error
}

func InitNats(ctx context.Context) (INats, error) {
	natsConn, err := nats.Connect(
		config.Cfg.Nats.Addr,
		nats.UserInfo(config.Cfg.Nats.User, config.Cfg.Nats.Pwd),
		nats.Name(config.ServiceName))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %w", err)
	}

	utils.AddShutdown(func() {
		natsConn.Close()
		log.Print("nats закрыт")
	})

	return &_nats{conn: natsConn}, nil
}
