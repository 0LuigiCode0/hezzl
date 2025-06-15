package ihttp

import (
	"context"
	"net/http"
	"sync"

	bnats "github.com/0LuigiCode0/hezzl/internal/infrastructure/broker/nats"
	rpostgres "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/postgres"
	rredis "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/redis"
)

type _handler struct {
	wg  sync.WaitGroup
	ctx context.Context

	pg    rpostgres.IRepoPostgres
	redis rredis.IRedis
	nats  bnats.IStream
}

const (
	urlGoodCreate = "POST /good/create"
	urlGoodUpdate = "PATCH /good/update"
	urlGoodRemove = "DELETE /good/remove"
	urlGoodList   = "GET /goods/list"
)

func InitHandler(ctx context.Context, pg rpostgres.IRepoPostgres, redis rredis.IRedis, nats bnats.IStream) http.Handler {
	handler := new(_handler)
	handler.ctx = ctx
	handler.pg = pg
	handler.redis = redis
	handler.nats = nats

	mux := http.NewServeMux()

	mux.HandleFunc(urlGoodCreate, handler.createGood)
	mux.HandleFunc(urlGoodUpdate, handler.updateGood)
	mux.HandleFunc(urlGoodRemove, handler.removeGood)
	mux.HandleFunc(urlGoodList, handler.getGoodsList)

	return mux
}
