package rpostgres

import (
	"context"

	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
	"github.com/jackc/pgx/v5"
)

const (
	qGetGoodsList = "select * from goods where removed=false limit $1 offset $2"
	qGetGoodsMeta = "select count(id) as 'total', (select count(id) from goods where removed=true) as 'removed' from goods"
	qInsertGoods  = "insert into goods (project_id,name) values (@project_id,@name) RETURNING *"
	qUpdateGoods  = "update goods set name=@name, description=@description where id=@id and removed=false RETURNING *"
	qDeleteGoods  = "update goods set removed=true where id=$1 and removed=false RETURNING *"
)

func (q *_db) GetGoods(ctx context.Context, limit, offset int) ([]*dpostgres.Good, error) {
	return queryRows[dpostgres.Good](q.db, ctx, qGetGoodsList, limit, offset)
}

func (q *_db) GetGoodsMeta(ctx context.Context) (*dpostgres.Meta, error) {
	return query[dpostgres.Meta](q.db, ctx, qGetGoodsMeta)
}

func (q *_db) InsertGood(ctx context.Context, in *dpostgres.InsertGood) (*dpostgres.Good, error) {
	return query[dpostgres.Good](q.db, ctx, qInsertGoods, pgx.NamedArgs{
		"project_id": in.ProjectId,
		"name":       in.Name,
	})
}

func (q *_db) UpdateGood(ctx context.Context, in *dpostgres.UpdateGood) (*dpostgres.Good, error) {
	return query[dpostgres.Good](q.db, ctx, qUpdateGoods, pgx.NamedArgs{
		"id":          in.Id,
		"name":        in.Name,
		"description": in.Description,
	})
}

func (q *_db) RemoveGood(ctx context.Context, id int) (*dpostgres.Good, error) {
	return query[dpostgres.Good](q.db, ctx, qDeleteGoods, id)
}
