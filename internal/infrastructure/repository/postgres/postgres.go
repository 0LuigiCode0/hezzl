package rpostgres

import (
	"context"
	"fmt"
	"log"

	"github.com/0LuigiCode0/hezzl/config"
	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
	"github.com/0LuigiCode0/hezzl/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type _IConn interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

type _db struct {
	db _IConn
}

type _IDB interface {
	GetGoods(ctx context.Context, limit, offset int) ([]*dpostgres.Good, error)
	GetGoodsMeta(ctx context.Context) (*dpostgres.Meta, error)
	InsertGood(ctx context.Context, in *dpostgres.InsertGood) (*dpostgres.Good, error)
	UpdateGood(ctx context.Context, in *dpostgres.UpdateGood) (*dpostgres.Good, error)
	RemoveGood(ctx context.Context, id int) (*dpostgres.Good, error)
}

type _postgres struct {
	_db
	conn *pgx.Conn
}

type IRepoPostgres interface {
	_IDB

	Begin(ctx context.Context) (ITx, error)
}

func InitRepoPostgres(ctx context.Context) (IRepoPostgres, error) {
	pgConn, err := pgx.Connect(ctx, config.Cfg.Postgres.URL)
	if err != nil {
		return nil, fmt.Errorf(consts.ErrOpenConnect, err)
	}
	utils.AddShutdown(func() {
		if err := pgConn.Close(context.Background()); err != nil {
			log.Printf("ошибка закрытия соединения postgres: %s", err)
		} else {
			log.Print("postgres закрыт")
		}
	})

	return &_postgres{conn: pgConn, _db: _db{pgConn}}, nil
}
