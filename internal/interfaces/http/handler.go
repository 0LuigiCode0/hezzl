package ihttp

import (
	"net/http"

	bnats "github.com/0LuigiCode0/hezzl/internal/infrastructure/broker/nats"
	rpostgres "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/postgres"
	rredis "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/redis"
)

type handler struct {
	pg    rpostgres.IRepoPostgres
	redis rredis.IRedis
	nats  bnats.INats
}

const (
	urlGoodCreate = "POST /good/create"
	urlGoodUpdate = "PATCH /good/update"
	urlGoodRemove = "DELETE /good/remove"
	urlGoodList   = "GET /goods/list"
)

func InitHandler(pg rpostgres.IRepoPostgres, redis rredis.IRedis, nats bnats.INats) http.Handler {
	handler := new(handler)
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
