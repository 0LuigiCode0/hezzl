package bnats

import (
	"context"
	"fmt"
	"log"

	"github.com/0LuigiCode0/hezzl/internal/utils"
	"github.com/nats-io/nats.go"
)

type _consumer struct {
	sub  *nats.Subscription
	subj string
}

type IConsumer interface {
	Fetch(ctx context.Context, batch int) ([]*nats.Msg, error)
}

func (s *_stream) Subscribe(ctx context.Context, subj string, nameClient string) (IConsumer, error) {
	sub, err := s.js.PullSubscribe(subj, nameClient, nats.AckExplicit(), nats.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf(errSubscribe, err)
	}
	utils.AddShutdown(func() {
		if err := sub.Drain(); err != nil {
			log.Printf(errUnsubscribe, err)
		} else {
			log.Print(prefix + "подписчик nats закрыт")
		}
	})

	return &_consumer{sub: sub, subj: subj}, nil
}

func (c *_consumer) Fetch(ctx context.Context, batch int) ([]*nats.Msg, error) {
	msgs, err := c.sub.Fetch(batch, nats.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf(errFetch, c.subj, err)
	}

	return msgs, nil
}
