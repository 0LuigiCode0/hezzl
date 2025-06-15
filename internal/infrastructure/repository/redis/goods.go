package rredis

import (
	"context"
	"crypto/sha1"
	"strconv"

	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
)

func (r *_redis) PushGoods(ctx context.Context, limit, offset int, value []*dpostgres.Good) error {
	hash := sha1.Sum([]byte(strconv.Itoa(limit) + "_" + strconv.Itoa(offset)))
	return jsonSet(r.conn, ctx, string(hash[:]), ".", value)
}

func (r *_redis) GetGoods(ctx context.Context, limit, offset int) ([]*dpostgres.Good, error) {
	hash := sha1.Sum([]byte(strconv.Itoa(limit) + "_" + strconv.Itoa(offset)))
	res, err := jsonGet[[]*dpostgres.Good](r.conn, ctx, string(hash[:]), ".")
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

func (r *_redis) DeleteGood(ctx context.Context, key string) error {
	return jsonDelete(r.conn, ctx, key, ".")
}
