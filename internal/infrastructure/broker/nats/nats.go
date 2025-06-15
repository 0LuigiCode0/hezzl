package bnats

import (
	"context"
	"fmt"
	"log"

	"github.com/0LuigiCode0/hezzl/config"
	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	"github.com/0LuigiCode0/hezzl/internal/utils"
	"github.com/nats-io/nats.go"
)

type _nats struct {
	conn *nats.Conn
}

type INats interface {
	CreateStream(ctx context.Context, name string, subs ...string) (IStream, error)
	ConnectStream(ctx context.Context, name string) (IStream, error)
}

func InitNats(ctx context.Context) (INats, error) {
	natsConn, err := nats.Connect(
		config.Cfg.Nats.Addr,
		nats.UserInfo(config.Cfg.Nats.User, config.Cfg.Nats.Pwd),
		nats.Name(config.ServiceName))
	if err != nil {
		return nil, fmt.Errorf(consts.ErrOpenConnect, err)
	}

	utils.AddShutdown(func() {
		natsConn.Close()
		log.Print("nats закрыт")
	})

	return &_nats{conn: natsConn}, nil
}
