package rredis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/0LuigiCode0/hezzl/config"
	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
)

func (r *_redis) PushGoods(ctx context.Context, limit, offset int, goods []*dpostgres.Good) error {
	key := hashKey(strconv.Itoa(limit), strconv.Itoa(offset))
	tx := r.conn.TxPipeline()

	err := jsonSet(tx, ctx, key, ".", goods, config.Cfg.Redis.Expire)
	if err != nil {
		tx.Discard()
		return err
	}

	for _, v := range goods {
		err = sAdd(tx, ctx, idxId+strconv.Itoa(v.Id), config.Cfg.Redis.Expire, key)
		if err != nil {
			tx.Discard()
			return err
		}
	}

	_, err = tx.Exec(ctx)
	if err != nil {
		return fmt.Errorf(errTx, err)
	}

	return nil
}

func (r *_redis) GetGoods(ctx context.Context, limit, offset int) ([]*dpostgres.Good, error) {
	key := hashKey(strconv.Itoa(limit), strconv.Itoa(offset))

	res, err := jsonGet[[]*dpostgres.Good](r.conn, ctx, key, ".", config.Cfg.Redis.Expire)
	if err != nil || res == nil {
		return nil, err
	}
	for _, v := range *res {
		r.conn.Expire(ctx, idxId+strconv.Itoa(v.Id), config.Cfg.Redis.Expire)
	}

	return *res, nil
}

func (r *_redis) DeleteAllWithGood(ctx context.Context, id int) error {
	idx := idxId + strconv.Itoa(id)
	keys, err := sGet(r.conn, ctx, idx)
	if err != nil || keys == nil {
		return err
	}

	tx := r.conn.TxPipeline()
	err = deleteKeys(tx, ctx, keys...)
	if err != nil {
		tx.Discard()
		return err
	}
	err = sRem(tx, ctx, idx, keys...)
	if err != nil {
		tx.Discard()
		return err
	}

	_, err = tx.Exec(ctx)
	if err != nil {
		return fmt.Errorf(errTx, err)
	}

	return nil
}
