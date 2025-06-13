package rredis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	"github.com/redis/go-redis/v9"
)

func jsonSet(r *redis.Client, ctx context.Context, key string, path string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf(consts.ErrJsonMarshal, err)
	}

	res := r.JSONSet(ctx, key, path, payload)
	if res.Err() != nil {
		return fmt.Errorf(errCreateKey, key, res.Err())
	}

	return nil
}

func jsonGet[T any](r *redis.Client, ctx context.Context, key string, path string) (*T, error) {
	res := r.JSONGet(ctx, key, path)
	if res.Err() != nil {
		return nil, fmt.Errorf(errRead, res.Err())
	}

	out := new(T)
	err := json.Unmarshal([]byte(res.Val()), out)
	if err != nil {
		return nil, fmt.Errorf(consts.ErrJsonUnmarshal, err)
	}

	return out, nil
}

func jsonGetList[T any](r *redis.Client, ctx context.Context, path string, keys ...string) ([]*T, error) {
	res := r.JSONMGet(ctx, path, keys...)
	if res.Err() != nil {
		return nil, fmt.Errorf(errRead, res.Err())
	}

	out := make([]*T, len(res.Val()))
	for i, v := range res.Val() {
		err := json.Unmarshal([]byte(v.(string)), out[i])
		if err != nil {
			return nil, fmt.Errorf(consts.ErrJsonUnmarshal, err)
		}
	}

	return out, nil
}

func jsonDelete(r *redis.Client, ctx context.Context, key string, path string) error {
	res := r.JSONDel(ctx, key, path)
	if res.Err() != nil {
		return fmt.Errorf(errDeleteKey, key, res.Err())
	}

	return nil
}
