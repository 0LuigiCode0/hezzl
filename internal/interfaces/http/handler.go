package ihttp

import (
	"context"
	"net/http"

	bnats "github.com/0LuigiCode0/hezzl/internal/infrastructure/broker/nats"
	rpostgres "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/postgres"
	rredis "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/redis"
	"github.com/0LuigiCode0/hezzl/internal/usecase"
)

type _handler struct {
	usecase usecase.IUsecase
}

const (
	urlGoodCreate = "POST /good/create"
	urlGoodUpdate = "PATCH /good/update"
	urlGoodRemove = "DELETE /good/remove"
	urlGoodList   = "GET /goods/list"
)

func InitHandler(ctx context.Context, pg rpostgres.IRepoPostgres, redis rredis.IRedis, nats bnats.IStream) http.Handler {
	handler := new(_handler)
	handler.usecase = usecase.InitUsecase(ctx, pg, redis, nats)

	mux := http.NewServeMux()

	mux.HandleFunc(urlGoodCreate, handler.createGood)
	mux.HandleFunc(urlGoodUpdate, handler.updateGood)
	mux.HandleFunc(urlGoodRemove, handler.removeGood)
	mux.HandleFunc(urlGoodList, handler.getGoodsList)

	return mux
}
