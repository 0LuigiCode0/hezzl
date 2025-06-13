package ihttp

import (
	"log"
	"net/http"
	"strconv"

	"github.com/0LuigiCode0/hezzl/internal/domain/conv"
	dhttp "github.com/0LuigiCode0/hezzl/internal/domain/http"
	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
)

func (h *handler) createGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	projectID, err := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil {
		writeErrorFLog(w, 400, errParseParam, err)
		return
	}

	data, err := jsonParse[dhttp.CreateGoodReq](r.Body)
	if err != nil {
		writeErrorFLog(w, 400, errParseBody, err)
		return
	}
	err = data.Validate()
	if err != nil {
		writeErrorFLog(w, 400, errParseBody, err)
		return
	}

	good, err := h.pg.InsertGood(r.Context(), &dpostgres.InsertGood{
		ProjectId: projectID,
		Name:      data.Name,
	})
	if err != nil {
		writeErrorFLog(w, 500, errInsertGood, err)
		return
	}

	err = h.nats.Push(conv.GoodPgToCh(good))
	if err != nil {
		log.Printf(errInsertLog, "creteGood", err)
	}

	writeJson(w, 200, conv.GoodPgToResp(good))
}

func (h *handler) updateGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeErrorFLog(w, 400, errParseParam, err)
		return
	}

	data, err := jsonParse[dhttp.UpdateGoodReq](r.Body)
	if err != nil {
		writeErrorFLog(w, 400, errParseBody, err)
		return
	}
	err = data.Validate()
	if err != nil {
		writeErrorFLog(w, 400, errParseBody, err)
		return
	}

	good, err := h.pg.UpdateGood(r.Context(), &dpostgres.UpdateGood{
		Id:          id,
		Description: data.Description,
		Name:        data.Name,
	})
	if err != nil {
		writeErrorFLog(w, 500, errUpdateGood, err)
		return
	}
	if good == nil {
		writeJson(w, 404, dhttp.NewError(3, errNotFound, nil))
		return
	}

	err = h.nats.Push(conv.GoodPgToCh(good))
	if err != nil {
		log.Printf(errInsertLog, "updateGood", err)
	}

	writeJson(w, 200, conv.GoodPgToResp(good))
}

func (h *handler) removeGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeErrorFLog(w, 400, errParseParam, err)
		return
	}

	good, err := h.pg.RemoveGood(r.Context(), id)
	if err != nil {
		writeErrorFLog(w, 500, errRemoveGood, err)
		return
	}
	if good == nil {
		writeJson(w, 404, dhttp.NewError(3, errNotFound, nil))
		return
	}

	err = h.nats.Push(conv.GoodPgToCh(good))
	if err != nil {
		log.Printf(errInsertLog, "removeGood", err)
	}

	writeJson(w, 200, conv.GoodPgToRemoveResp(good))
}

func (h *handler) getGoodsList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		writeErrorFLog(w, 400, errParseParam, err)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		writeErrorFLog(w, 400, errParseParam, err)
		return
	}

	goods, err := h.pg.GetGoods(r.Context(), limit, offset)
	if err != nil {
		writeErrorFLog(w, 500, errGetGoods, err)
		return
	}
	meta, err := h.pg.GetGoodsMeta(r.Context())
	if err != nil {
		writeErrorFLog(w, 500, errGetGoods, err)
		return
	}

	writeJson(w, 200, conv.GoodsPgToRespMeta(meta, goods, limit, offset))
}
