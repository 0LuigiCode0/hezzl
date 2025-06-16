package ihttp

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/0LuigiCode0/hezzl/config"
	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	"github.com/0LuigiCode0/hezzl/internal/domain/conv"
	dhttp "github.com/0LuigiCode0/hezzl/internal/domain/http"
	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
)

func (h *_handler) createGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	projectID, err := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil {
		writeErrorFLog(w, http.StatusBadRequest, errParseParam, err)
		return
	}
	if projectID <= 0 {
		writeErrorFLog(w, http.StatusBadRequest, consts.ErrFieldValid, "projectId")
		return
	}

	data, err := jsonParse[dhttp.CreateGoodReq](r.Body)
	if err != nil {
		writeErrorFLog(w, http.StatusBadRequest, errParseBody, err)
		return
	}
	err = data.Validate()
	if err != nil {
		writeErrorFLog(w, http.StatusBadRequest, errParseBody, err)
		return
	}

	good, err := h.pg.InsertGood(r.Context(), &dpostgres.InsertGood{
		ProjectId: projectID,
		Name:      data.Name,
	})
	if err != nil {
		writeErrorLog(w, http.StatusInternalServerError, err)
		return
	}

	h.withRetry(config.Cfg.Nats.RetryCount, func(ctx context.Context) error {
		err = h.nats.PushGoodsLog(ctx, conv.GoodPgToCh(good))
		if err != nil {
			log.Print("creteGood", err)
		}
		return err
	})

	writeJson(w, http.StatusOK, conv.GoodPgToResp(good))
}

func (h *_handler) updateGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeErrorFLog(w, http.StatusBadRequest, errParseParam, err)
		return
	}
	if id <= 0 {
		writeErrorFLog(w, http.StatusBadRequest, consts.ErrFieldValid, "id")
		return
	}

	data, err := jsonParse[dhttp.UpdateGoodReq](r.Body)
	if err != nil {
		writeErrorFLog(w, http.StatusBadRequest, errParseBody, err)
		return
	}
	err = data.Validate()
	if err != nil {
		writeErrorFLog(w, http.StatusBadRequest, errParseBody, err)
		return
	}

	good, err := h.pg.UpdateGood(r.Context(), &dpostgres.UpdateGood{
		Id:          id,
		Description: data.Description,
		Name:        data.Name,
	})
	if err != nil {
		writeErrorLog(w, http.StatusInternalServerError, err)
		return
	}
	if good == nil {
		writeJson(w, http.StatusNotFound, dhttp.NewError(3, errNotFound, nil))
		return
	}

	h.withRetry(config.Cfg.Nats.RetryCount, func(ctx context.Context) error {
		err = h.nats.PushGoodsLog(ctx, conv.GoodPgToCh(good))
		if err != nil {
			log.Print("updateGood", err)
		}
		return err
	})
	h.withRetry(config.Cfg.Redis.RetryCount, func(ctx context.Context) error {
		err = h.redis.DeleteAllWithGood(ctx, good.Id)
		if err != nil {
			log.Print(err)
		}
		return err
	})

	writeJson(w, http.StatusOK, conv.GoodPgToResp(good))
}

func (h *_handler) removeGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeErrorFLog(w, http.StatusBadRequest, errParseParam, err)
		return
	}
	if id <= 0 {
		writeErrorFLog(w, http.StatusBadRequest, consts.ErrFieldValid, "id")
		return
	}

	good, err := h.pg.RemoveGood(r.Context(), id)
	if err != nil {
		writeErrorLog(w, http.StatusInternalServerError, err)
		return
	}
	if good == nil {
		writeJson(w, http.StatusNotFound, dhttp.NewError(3, errNotFound, nil))
		return
	}

	h.withRetry(config.Cfg.Nats.RetryCount, func(ctx context.Context) error {
		err = h.nats.PushGoodsLog(ctx, conv.GoodPgToCh(good))
		if err != nil {
			log.Print("removeGood", err)
		}
		return err
	})
	h.withRetry(config.Cfg.Redis.RetryCount, func(ctx context.Context) error {
		err = h.redis.DeleteAllWithGood(ctx, good.Id)
		if err != nil {
			log.Print(err)
		}
		return err
	})

	writeJson(w, http.StatusOK, conv.GoodPgToRemoveResp(good))
}

func (h *_handler) getGoodsList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		writeErrorFLog(w, http.StatusBadRequest, errParseParam, err)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		writeErrorFLog(w, http.StatusBadRequest, errParseParam, err)
		return
	}
	if limit <= 0 {
		limit = 10
	}
	if offset <= 0 {
		offset = 1
	}

	var isDB bool
	goods, err := h.redis.GetGoods(r.Context(), limit, offset)
	if err != nil {
		log.Print(err)

		goods, err = h.pg.GetGoods(r.Context(), limit, offset)
		if err != nil {
			writeErrorLog(w, http.StatusInternalServerError, err)
			return
		}
		isDB = true
	}

	meta, err := h.pg.GetGoodsMeta(r.Context())
	if err != nil {
		writeErrorLog(w, http.StatusInternalServerError, err)
		return
	}

	if isDB {
		h.withRetry(config.Cfg.Redis.RetryCount, func(ctx context.Context) error {
			err = h.redis.PushGoods(ctx, limit, offset, goods)
			if err != nil {
				log.Print(err)
			}
			return err
		})
	}

	writeJson(w, http.StatusOK, conv.GoodsPgToRespMeta(meta, goods, limit, offset))
}
