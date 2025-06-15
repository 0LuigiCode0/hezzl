package bnats

import (
	"context"
	"fmt"

	"github.com/0LuigiCode0/hezzl/config"
	dclickhouse "github.com/0LuigiCode0/hezzl/internal/domain/clickhouse"
)

func (s *_stream) PushGoodsLog(ctx context.Context, in *dclickhouse.LogEventGood) error {
	err := pushStream(s, ctx, config.Cfg.Nats.Subj, in)
	if err != nil {
		return fmt.Errorf(errInsertLog, err)
	}

	return nil
}
