package rpostgres

import (
	"context"
	"fmt"

	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
	"github.com/jackc/pgx/v5"
)

const (
	qGetGoodsList = "select * from goods order by created_at asc limit $1 offset $2"
	qGetGoodsMeta = "select count(id) as total, (select count(id) from goods where removed=true) as removed from goods"
	qInsertGoods  = "insert into goods (project_id,name) values (@project_id,@name) RETURNING *"
	qUpdateGoods  = "update goods set name=@name, description=@description where id=@id and removed=false RETURNING *"
	qDeleteGoods  = "update goods set removed=true where id=$1 and removed=false RETURNING *"
)

func (r *_repo) GetGoods(ctx context.Context, limit, offset int) ([]*dpostgres.Good, error) {
	res, err := queryRows[dpostgres.Good](r.db, ctx, qGetGoodsList, limit, offset)
	if err != nil {
		return nil, fmt.Errorf(errGetGoods, err)
	}

	return res, nil
}

func (r *_repo) GetGoodsMeta(ctx context.Context) (*dpostgres.Meta, error) {
	res, err := query[dpostgres.Meta](r.db, ctx, qGetGoodsMeta)
	if err != nil {
		return nil, fmt.Errorf(errGetGoods, err)
	}

	return res, nil
}

func (r *_repo) InsertGood(ctx context.Context, in *dpostgres.InsertGood) (*dpostgres.Good, error) {
	res, err := query[dpostgres.Good](r.db, ctx, qInsertGoods, pgx.NamedArgs{
		"project_id": in.ProjectId,
		"name":       in.Name,
	})
	if err != nil {
		return nil, fmt.Errorf(errGetGoods, err)
	}

	return res, nil
}

func (r *_repo) UpdateGood(ctx context.Context, in *dpostgres.UpdateGood) (*dpostgres.Good, error) {
	res, err := query[dpostgres.Good](r.db, ctx, qUpdateGoods, pgx.NamedArgs{
		"id":          in.Id,
		"name":        in.Name,
		"description": in.Description,
	})
	if err != nil {
		return nil, fmt.Errorf(errGetGoods, err)
	}

	return res, nil
}

func (r *_repo) RemoveGood(ctx context.Context, id int) (*dpostgres.Good, error) {
	res, err := query[dpostgres.Good](r.db, ctx, qDeleteGoods, id)
	if err != nil {
		return nil, fmt.Errorf(errGetGoods, err)
	}

	return res, nil
}
