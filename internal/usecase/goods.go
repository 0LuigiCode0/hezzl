package usecase

import (
	"context"
	"log"
	"net/http"

	"github.com/0LuigiCode0/hezzl/config"
	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	"github.com/0LuigiCode0/hezzl/internal/domain/conv"
	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
	dusecase "github.com/0LuigiCode0/hezzl/internal/domain/usecase"
)

func (u *_usecase) CreateGood(ctx context.Context, projectId int, data *dusecase.CreateGoodReq) (*dusecase.GoodResp, *dusecase.ErrorResp) {
	if projectId <= 0 {
		return nil, dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, nil, dusecase.ErrorDetails{"projectid": projectId})
	}
	err := data.Validate()
	if err != nil {
		return nil, dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, err, nil)
	}

	good, err := u.pg.InsertGood(ctx, &dpostgres.InsertGood{
		ProjectId: projectId,
		Name:      data.Name,
	})
	if err != nil {
		return nil, dusecase.NewError(http.StatusInternalServerError, errCreateGood, err, nil)
	}

	u.withRetry(config.Cfg.Nats.RetryCount, func(ctx context.Context) error {
		err = u.nats.PushGoodsLog(ctx, conv.GoodPgToCh(good))
		if err != nil {
			log.Print("creteGood", err)
		}
		return err
	})

	return conv.GoodPgToResp(good), nil
}

func (u *_usecase) UpdateGood(ctx context.Context, id int, data *dusecase.UpdateGoodReq) (*dusecase.GoodResp, *dusecase.ErrorResp) {
	if id <= 0 {
		return nil, dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, nil, dusecase.ErrorDetails{"id": id})
	}
	err := data.Validate()
	if err != nil {
		return nil, dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, err, nil)
	}

	good, err := u.pg.UpdateGood(ctx, &dpostgres.UpdateGood{
		Id:          id,
		Description: data.Description,
		Name:        data.Name,
	})
	if err != nil {
		return nil, dusecase.NewError(http.StatusInternalServerError, errUpdateGood, err, nil)
	}
	if good == nil {
		return nil, dusecase.NewError(3, errNotFound, nil, nil)
	}

	u.withRetry(config.Cfg.Nats.RetryCount, func(ctx context.Context) error {
		err = u.nats.PushGoodsLog(ctx, conv.GoodPgToCh(good))
		if err != nil {
			log.Print("updateGood", err)
		}
		return err
	})
	u.withRetry(config.Cfg.Redis.RetryCount, func(ctx context.Context) error {
		err = u.redis.DeleteAllWithGood(ctx, good.Id)
		if err != nil {
			log.Print(err)
		}
		return err
	})

	return conv.GoodPgToResp(good), nil
}

func (u *_usecase) RemoveGood(ctx context.Context, id int) (*dusecase.RemoveGoodResp, *dusecase.ErrorResp) {
	if id <= 0 {
		return nil, dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, nil, dusecase.ErrorDetails{"id": id})
	}

	good, err := u.pg.RemoveGood(ctx, id)
	if err != nil {
		return nil, dusecase.NewError(http.StatusInternalServerError, errRemoveGood, err, nil)
	}
	if good == nil {
		return nil, dusecase.NewError(3, errNotFound, nil, nil)
	}

	u.withRetry(config.Cfg.Nats.RetryCount, func(ctx context.Context) error {
		err = u.nats.PushGoodsLog(ctx, conv.GoodPgToCh(good))
		if err != nil {
			log.Print("removeGood", err)
		}
		return err
	})
	u.withRetry(config.Cfg.Redis.RetryCount, func(ctx context.Context) error {
		err = u.redis.DeleteAllWithGood(ctx, good.Id)
		if err != nil {
			log.Print(err)
		}
		return err
	})

	return conv.GoodPgToRemoveResp(good), nil
}

func (u *_usecase) GetGoodsList(ctx context.Context, limit, offset int) (*dusecase.GetGoodsResp, *dusecase.ErrorResp) {
	if limit <= 0 {
		limit = 10
	}
	if offset <= 0 {
		offset = 1
	}

	var isDB bool
	goods, err := u.redis.GetGoods(ctx, limit, offset)
	if err != nil {
		log.Print(err)

		goods, err = u.pg.GetGoods(ctx, limit, offset)
		if err != nil {
			return nil, dusecase.NewError(http.StatusInternalServerError, errGetGoods, err, nil)
		}
		isDB = true
	}

	meta, err := u.pg.GetGoodsMeta(ctx)
	if err != nil {
		return nil, dusecase.NewError(http.StatusInternalServerError, errGetGoods, err, nil)
	}

	if isDB {
		u.withRetry(config.Cfg.Redis.RetryCount, func(ctx context.Context) error {
			err = u.redis.PushGoods(ctx, limit, offset, goods)
			if err != nil {
				log.Print(err)
			}
			return err
		})
	}

	return conv.GoodsPgToRespMeta(meta, goods, limit, offset), nil
}
