package rclickhouse

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/0LuigiCode0/hezzl/config"
	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	"github.com/0LuigiCode0/hezzl/internal/utils"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type _clickhouse struct {
	conn driver.Conn
}

type IClickHouse interface{}

func InitClickHouse(ctx context.Context) (IClickHouse, error) {
	clConn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.Cfg.ClickHouse.Addr},
		Auth: clickhouse.Auth{
			Database: config.Cfg.ClickHouse.DB,
			Username: config.Cfg.ClickHouse.User,
			Password: config.Cfg.ClickHouse.Pwd,
		},
		TLS: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		return nil, fmt.Errorf(consts.ErrOpenConnect, err)
	}
	err = clConn.Ping(ctx)
	if err != nil {
		clConn.Close()
		return nil, fmt.Errorf(consts.ErrPing, err)
	}

	utils.AddShutdown(func() {
		if err := clConn.Close(); err != nil {
			log.Printf("ошибка закрытия соединения clickhouse: %s", err)
		} else {
			log.Print("clickhouse закрыт")
		}
	})

	return &_clickhouse{conn: clConn}, nil
}
