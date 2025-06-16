package usecase

import (
	"errors"
	"net/http"
	"testing"

	dclickhouse "github.com/0LuigiCode0/hezzl/internal/domain/clickhouse"
	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
	dusecase "github.com/0LuigiCode0/hezzl/internal/domain/usecase"
	"github.com/0LuigiCode0/tshort/tshort"
	"github.com/stretchr/testify/assert"
)

var errFoo = errors.New("errFoo")

func TestCreateGood(t *testing.T) {
	env := initTest[dusecase.GoodResp](t)

	var projectId int
	var data *dusecase.CreateGoodReq

	ts := tshort.Init(func(t *testing.T) {
		env.wantRes = nil
		env.wantErr = nil
	}, ".", "projectId.invalid", "data.invalid", "InsertGood")

	ts.AddStage("projectId.invalid", func() {
		projectId = -3
		env.wantErr = dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, nil, dusecase.ErrorDetails{"projectid": projectId})
	})
	ts.AddStage("data.invalid", func() {
		projectId = 3
		data = &dusecase.CreateGoodReq{Name: ""}
		env.wantErr = dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, data.Validate(), nil)
	})

	ts.AddStage("InsertGood", func() {
		projectId = 3
		data = &dusecase.CreateGoodReq{Name: "hello"}
	}, "@InsertGood.error", "@InsertGood.success")
	{
		ts.AddStage("@InsertGood.error", func() {
			env.wantErr = dusecase.NewError(http.StatusInternalServerError, errCreateGood, errFoo, nil)
			env.pg.EXPECT().InsertGood(env.ctx, &dpostgres.InsertGood{
				ProjectId: projectId,
				Name:      data.Name,
			}).Return(nil, errFoo).Once()
		})
		ts.AddStage("@InsertGood.success", func() {
			env.wantRes = &dusecase.GoodResp{Name: data.Name, ProjectId: projectId}
			env.pg.EXPECT().InsertGood(env.ctx, &dpostgres.InsertGood{
				ProjectId: projectId,
				Name:      data.Name,
			}).Return(&dpostgres.Good{Name: data.Name, ProjectId: projectId}, nil).Once()
		}, "background")
	}
	ts.AddStage("background", func() {}, "@background.error", "@background.success")
	{
		ts.AddStage("@background.error", func() {
			env.nats.EXPECT().PushGoodsLog(env.ctx, &dclickhouse.LogEventGood{Name: data.Name, ProjectId: projectId}).Return(errFoo).Once()
			env.nats.EXPECT().PushGoodsLog(env.ctx, &dclickhouse.LogEventGood{Name: data.Name, ProjectId: projectId}).Return(nil).Once()
		})
		ts.AddStage("@background.success", func() {
			env.nats.EXPECT().PushGoodsLog(env.ctx, &dclickhouse.LogEventGood{Name: data.Name, ProjectId: projectId}).Return(nil).Once()
		})
	}

	ts.Run(t, func(t *testing.T) {
		res, err := env.u.CreateGood(env.ctx, projectId, data)
		env.u.wg.Wait()
		assert.Equal(t, []any{res, err}, []any{env.wantRes, env.wantErr})
	})
}

func TestUpdateGood(t *testing.T) {
	env := initTest[dusecase.GoodResp](t)

	var id int
	var data *dusecase.UpdateGoodReq

	ts := tshort.Init(func(t *testing.T) {
		env.wantRes = nil
		env.wantErr = nil
	}, ".", "id.invalid", "data.invalid", "UpdateGood")

	ts.AddStage("id.invalid", func() {
		id = -3
		env.wantErr = dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, nil, dusecase.ErrorDetails{"id": id})
	})
	ts.AddStage("data.invalid", func() {
		id = 3
		data = &dusecase.UpdateGoodReq{Name: "", Description: ""}
		env.wantErr = dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, data.Validate(), nil)
	})

	ts.AddStage("UpdateGood", func() {
		id = 3
		data = &dusecase.UpdateGoodReq{Name: "hello", Description: "world"}
	}, "@UpdateGood.error", "@UpdateGood.notfound", "@UpdateGood.success")
	{
		ts.AddStage("@UpdateGood.error", func() {
			env.wantErr = dusecase.NewError(http.StatusInternalServerError, errUpdateGood, errFoo, nil)
			env.pg.EXPECT().UpdateGood(env.ctx, &dpostgres.UpdateGood{
				Id:          id,
				Name:        data.Name,
				Description: data.Description,
			}).Return(nil, errFoo).Once()
		})
		ts.AddStage("@UpdateGood.notfound", func() {
			env.wantErr = dusecase.NewError(3, errNotFound, nil, nil)
			env.pg.EXPECT().UpdateGood(env.ctx, &dpostgres.UpdateGood{
				Id:          id,
				Name:        data.Name,
				Description: data.Description,
			}).Return(nil, nil).Once()
		})
		ts.AddStage("@UpdateGood.success", func() {
			env.wantRes = &dusecase.GoodResp{Name: data.Name, Id: id, Description: data.Description}
			env.pg.EXPECT().UpdateGood(env.ctx, &dpostgres.UpdateGood{
				Id:          id,
				Name:        data.Name,
				Description: data.Description,
			}).Return(&dpostgres.Good{Name: data.Name, Id: id, Description: data.Description}, nil).Once()
		}, "background")
	}
	ts.AddStage("background", func() {}, "@background.error", "@background.success")
	{
		ts.AddStage("@background.error", func() {
			env.nats.EXPECT().PushGoodsLog(env.ctx, &dclickhouse.LogEventGood{Name: data.Name, Id: id, Description: data.Description}).Return(errFoo).Once()
			env.nats.EXPECT().PushGoodsLog(env.ctx, &dclickhouse.LogEventGood{Name: data.Name, Id: id, Description: data.Description}).Return(nil).Once()

			env.redis.EXPECT().DeleteAllWithGood(env.ctx, id).Return(errFoo).Once()
			env.redis.EXPECT().DeleteAllWithGood(env.ctx, id).Return(nil).Once()
		})
		ts.AddStage("@background.success", func() {
			env.nats.EXPECT().PushGoodsLog(env.ctx, &dclickhouse.LogEventGood{Name: data.Name, Id: id, Description: data.Description}).Return(nil).Once()

			env.redis.EXPECT().DeleteAllWithGood(env.ctx, id).Return(nil).Once()
		})
	}

	ts.Run(t, func(t *testing.T) {
		res, err := env.u.UpdateGood(env.ctx, id, data)
		env.u.wg.Wait()
		assert.Equal(t, []any{res, err}, []any{env.wantRes, env.wantErr})
	})
}

func TestRemoveGood(t *testing.T) {
	env := initTest[dusecase.RemoveGoodResp](t)

	var id int

	ts := tshort.Init(func(t *testing.T) {
		env.wantRes = nil
		env.wantErr = nil
	}, ".", "id.invalid", "RemoveGood")

	ts.AddStage("id.invalid", func() {
		id = -3
		env.wantErr = dusecase.NewError(http.StatusBadRequest, consts.ErrFieldValid, nil, dusecase.ErrorDetails{"id": id})
	})

	ts.AddStage("RemoveGood", func() {
		id = 3
	}, "@RemoveGood.error", "@RemoveGood.notfound", "@RemoveGood.success")
	{
		ts.AddStage("@RemoveGood.error", func() {
			env.wantErr = dusecase.NewError(http.StatusInternalServerError, errRemoveGood, errFoo, nil)
			env.pg.EXPECT().RemoveGood(env.ctx, id).Return(nil, errFoo).Once()
		})
		ts.AddStage("@RemoveGood.notfound", func() {
			env.wantErr = dusecase.NewError(3, errNotFound, nil, nil)
			env.pg.EXPECT().RemoveGood(env.ctx, id).Return(nil, nil).Once()
		})
		ts.AddStage("@RemoveGood.success", func() {
			env.wantRes = &dusecase.RemoveGoodResp{Id: id, Removed: true}
			env.pg.EXPECT().RemoveGood(env.ctx, id).Return(&dpostgres.Good{Id: id, Removed: true}, nil).Once()
		}, "background")
	}
	ts.AddStage("background", func() {}, "@background.error", "@background.success")
	{
		ts.AddStage("@background.error", func() {
			env.nats.EXPECT().PushGoodsLog(env.ctx, &dclickhouse.LogEventGood{Id: id, Removed: true}).Return(errFoo).Once()
			env.nats.EXPECT().PushGoodsLog(env.ctx, &dclickhouse.LogEventGood{Id: id, Removed: true}).Return(nil).Once()

			env.redis.EXPECT().DeleteAllWithGood(env.ctx, id).Return(errFoo).Once()
			env.redis.EXPECT().DeleteAllWithGood(env.ctx, id).Return(nil).Once()
		})
		ts.AddStage("@background.success", func() {
			env.nats.EXPECT().PushGoodsLog(env.ctx, &dclickhouse.LogEventGood{Id: id, Removed: true}).Return(nil).Once()

			env.redis.EXPECT().DeleteAllWithGood(env.ctx, id).Return(nil).Once()
		})
	}

	ts.Run(t, func(t *testing.T) {
		res, err := env.u.RemoveGood(env.ctx, id)
		env.u.wg.Wait()
		assert.Equal(t, []any{res, err}, []any{env.wantRes, env.wantErr})
	})
}

func TestGetGoodList(t *testing.T) {
	env := initTest[dusecase.GetGoodsResp](t)

	limit := 10
	offset := 1
	isDB := false

	ts := tshort.Init(func(t *testing.T) {
		env.wantRes = nil
		env.wantErr = nil
		isDB = false
	}, ".", "RedisGetGoods")

	ts.AddStage("RedisGetGoods", func() {}, "@RedisGetGoods.error", "@RedisGetGoods.success")
	{
		ts.AddStage("@RedisGetGoods.error", func() {
			env.redis.EXPECT().GetGoods(env.ctx, limit, offset).Return(nil, errFoo).Once()
		}, "PGGetGoods")
		ts.AddStage("@RedisGetGoods.success", func() {
			env.redis.EXPECT().GetGoods(env.ctx, limit, offset).Return([]*dpostgres.Good{{Id: 4}}, nil).Once()
		}, "GetMeta")
	}

	ts.AddStage("PGGetGoods", func() {}, "@PGGetGoods.error", "@PGGetGoods.success")
	{
		ts.AddStage("@PGGetGoods.error", func() {
			env.wantErr = dusecase.NewError(http.StatusInternalServerError, errGetGoods, errFoo, nil)
			env.pg.EXPECT().GetGoods(env.ctx, limit, offset).Return(nil, errFoo).Once()
		})
		ts.AddStage("@PGGetGoods.success", func() {
			isDB = true
			env.pg.EXPECT().GetGoods(env.ctx, limit, offset).Return([]*dpostgres.Good{{Id: 4}}, nil).Once()
		}, "GetMeta")
	}

	ts.AddStage("GetMeta", func() {}, "@GetMeta.error", "@GetMeta.success")
	{
		ts.AddStage("@GetMeta.error", func() {
			env.wantErr = dusecase.NewError(http.StatusInternalServerError, errGetGoods, errFoo, nil)
			env.pg.EXPECT().GetGoodsMeta(env.ctx).Return(nil, errFoo).Once()
		})
		ts.AddStage("@GetMeta.success", func() {
			env.wantRes = &dusecase.GetGoodsResp{Meta: &dusecase.Meta{Total: 1, Limit: limit, Offset: offset}, Goods: []*dusecase.GoodResp{{Id: 4}}}
			env.pg.EXPECT().GetGoodsMeta(env.ctx).Return(&dpostgres.Meta{Total: 1}, nil).Once()
		}, "background")
	}

	ts.AddStage("background", func() {}, "@background.error", "@background.success")
	{
		ts.AddStage("@background.error", func() {
			if isDB {
				env.redis.EXPECT().PushGoods(env.ctx, limit, offset, []*dpostgres.Good{{Id: 4}}).Return(errFoo).Once()
				env.redis.EXPECT().PushGoods(env.ctx, limit, offset, []*dpostgres.Good{{Id: 4}}).Return(nil).Once()
			}
		})
		ts.AddStage("@background.success", func() {
			if isDB {
				env.redis.EXPECT().PushGoods(env.ctx, limit, offset, []*dpostgres.Good{{Id: 4}}).Return(nil).Once()
			}
		})
	}

	ts.Run(t, func(t *testing.T) {
		res, err := env.u.GetGoodsList(env.ctx, 0, 0)
		env.u.wg.Wait()
		assert.Equal(t, []any{res, err}, []any{env.wantRes, env.wantErr})
	})
}
