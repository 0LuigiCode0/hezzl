package rredis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	"github.com/redis/go-redis/v9"
)

func jsonSet(r redis.Cmdable, ctx context.Context, key string, path string, value any, expire time.Duration) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf(consts.ErrJsonMarshal, err)
	}

	tx := r.TxPipeline()
	res := tx.JSONSet(ctx, key, path, payload)
	if res.Err() != nil {
		tx.Discard()
		return fmt.Errorf(errCreateKey, key, res.Err())
	}
	resExp := tx.Expire(ctx, key, expire)
	if resExp.Err() != nil {
		tx.Discard()
		return fmt.Errorf(errExpire, key, resExp.Err())
	}
	_, err = tx.Exec(ctx)
	if err != nil {
		return fmt.Errorf(errTx, err)
	}

	return nil
}

func jsonGet[T any](r redis.Cmdable, ctx context.Context, key string, path string, expire time.Duration) (*T, error) {
	res := r.JSONGet(ctx, key, path)
	if res.Err() != nil {
		return nil, fmt.Errorf(errReadKey, key, res.Err())
	}

	out := new(T)
	err := json.Unmarshal([]byte(res.Val()), out)
	if err != nil {
		return nil, fmt.Errorf(consts.ErrJsonUnmarshal, err)
	}

	r.Expire(ctx, key, expire)

	return out, nil
}

func sAdd(r redis.Cmdable, ctx context.Context, key string, expire time.Duration, values ...string) error {
	tx := r.TxPipeline()
	res := tx.SAdd(ctx, key, values)
	if res.Err() != nil {
		return fmt.Errorf(errAddSet, key, res.Err())
	}
	resExp := tx.Expire(ctx, key, expire)
	if resExp.Err() != nil {
		tx.Discard()
		return fmt.Errorf(errExpire, key, resExp.Err())
	}
	_, err := tx.Exec(ctx)
	if err != nil {
		return fmt.Errorf(errTx, err)
	}

	return nil
}

func sGet(r redis.Cmdable, ctx context.Context, key string) ([]string, error) {
	res := r.SMembers(ctx, key)
	if res.Err() != nil {
		return nil, fmt.Errorf(errCreateKey, key, res.Err())
	}

	return res.Result()
}

func sRem(r redis.Cmdable, ctx context.Context, key string, values ...string) error {
	res := r.SRem(ctx, key, values)
	if res.Err() != nil {
		return fmt.Errorf(errRemSet, key, res.Err())
	}

	return nil
}

func deleteKeys(r redis.Cmdable, ctx context.Context, keys ...string) error {
	res := r.Del(ctx, keys...)
	if res.Err() != nil {
		return fmt.Errorf(errDeleteKey, keys, res.Err())
	}

	return nil
}
