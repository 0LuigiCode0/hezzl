package usecase

import (
	"context"
	"testing"

	"github.com/0LuigiCode0/hezzl/config"
	dusecase "github.com/0LuigiCode0/hezzl/internal/domain/usecase"
	mnats "github.com/0LuigiCode0/hezzl/internal/infrastructure/broker/nats/mocks"
	mpostgres "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/postgres/mocks"
	mredis "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/redis/mocks"
	"github.com/0LuigiCode0/hezzl/internal/utils"
)

type envTest[T any] struct {
	ctx context.Context

	u *_usecase

	pg    *mpostgres.IRepoPostgres
	redis *mredis.IRedis
	nats  *mnats.IStream

	wantRes *T
	wantErr *dusecase.ErrorResp
}

func initTest[T any](t *testing.T) *envTest[T] {
	utils.IsTimeMock = true
	config.ParseConfig("")
	ctx := context.Background()
	pg := mpostgres.NewIRepoPostgres(t)
	redis := mredis.NewIRedis(t)
	nats := mnats.NewIStream(t)

	return &envTest[T]{
		ctx: ctx,
		u: &_usecase{
			ctx:   ctx,
			pg:    pg,
			redis: redis,
			nats:  nats,
		},
		pg:    pg,
		redis: redis,
		nats:  nats,
	}
}
