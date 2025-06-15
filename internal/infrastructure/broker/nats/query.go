package bnats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	"github.com/nats-io/nats.go"
)

func pushStream(s nats.JetStreamContext, ctx context.Context, subj string, in any) error {
	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf(consts.ErrJsonMarshal, err)
	}

	_, err = s.Publish(subj, data, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf(errPush, subj, err)
	}

	return nil
}
