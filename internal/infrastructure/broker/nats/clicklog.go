package bnats

import (
	"fmt"

	"github.com/0LuigiCode0/hezzl/config"
	dclickhouse "github.com/0LuigiCode0/hezzl/internal/domain/clickhouse"
)

func (n *_nats) Push(v *dclickhouse.LogEventGood) error {
	err := push(n, config.Cfg.Nats.Sub, v)
	if err != nil {
		return fmt.Errorf("ошибка записи nats: %w", err)
	}

	return nil
}
