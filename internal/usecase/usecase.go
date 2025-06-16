package usecase

import (
	"context"
	"sync"

	dusecase "github.com/0LuigiCode0/hezzl/internal/domain/usecase"
	bnats "github.com/0LuigiCode0/hezzl/internal/infrastructure/broker/nats"
	rpostgres "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/postgres"
	rredis "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/redis"
	"github.com/0LuigiCode0/hezzl/internal/utils"
)

type _usecase struct {
	wg  sync.WaitGroup
	ctx context.Context

	pg    rpostgres.IRepoPostgres
	redis rredis.IRedis
	nats  bnats.IStream
}

type IUsecase interface {
	CreateGood(ctx context.Context, projectId int, data *dusecase.CreateGoodReq) (*dusecase.GoodResp, *dusecase.ErrorResp)
	UpdateGood(ctx context.Context, id int, data *dusecase.UpdateGoodReq) (*dusecase.GoodResp, *dusecase.ErrorResp)
	RemoveGood(ctx context.Context, id int) (*dusecase.RemoveGoodResp, *dusecase.ErrorResp)
	GetGoodsList(ctx context.Context, limit, offset int) (*dusecase.GetGoodsResp, *dusecase.ErrorResp)
}

func InitUsecase(ctx context.Context, pg rpostgres.IRepoPostgres, redis rredis.IRedis, nats bnats.IStream) IUsecase {
	u := &_usecase{
		ctx:   ctx,
		pg:    pg,
		redis: redis,
		nats:  nats,
	}

	utils.AddShutdown(func() {
		u.wg.Wait()
	})

	return u
}
