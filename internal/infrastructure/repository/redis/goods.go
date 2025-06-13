package rredis

import (
	"context"

	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
)

func (r *_redis) PushGood(ctx context.Context, key string, value *dpostgres.Good) error {
	return jsonSet(r.conn, ctx, key, ".", value)
}

func (r *_redis) GetGoods(ctx context.Context, keys ...string) ([]*dpostgres.Good, error) {
	return jsonGetList[dpostgres.Good](r.conn, ctx, ".", keys...)
}

func (r *_redis) DeleteGoods(ctx context.Context, key string) error {
	return jsonDelete(r.conn, ctx, key, ".")
}
