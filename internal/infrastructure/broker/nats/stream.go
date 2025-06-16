package bnats

import (
	"context"
	"fmt"
	"time"

	dclickhouse "github.com/0LuigiCode0/hezzl/internal/domain/clickhouse"
	"github.com/nats-io/nats.go"
)

type _stream struct {
	js nats.JetStreamContext
}

//go:generate mockery --name IStream --outpkg mnats --recursive --with-expecter
type IStream interface {
	Subscribe(ctx context.Context, subj string, nameClient string) (IConsumer, error)

	PushGoodsLog(ctx context.Context, in *dclickhouse.LogEventGood) error
}

func (n *_nats) CreateStream(ctx context.Context, name string, subs ...string) (IStream, error) {
	js, err := n.conn.JetStream(nats.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf(errStreamCtx, err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      name,
		Subjects:  subs,
		Retention: nats.WorkQueuePolicy,
		MaxAge:    time.Hour * 24 * 7,
		Storage:   nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return nil, fmt.Errorf(errAddStream, err)
	}

	return &_stream{js: js}, nil
}

func (n *_nats) ConnectStream(ctx context.Context, name string) (IStream, error) {
	js, err := n.conn.JetStream(nats.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf(errStreamCtx, err)
	}

	_, err = js.StreamInfo(name)
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return nil, fmt.Errorf(errInfoStream, err)
	}

	return &_stream{js: js}, nil
}
