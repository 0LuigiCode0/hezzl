package ihttp

import (
	"net/http"
	"strconv"

	dusecase "github.com/0LuigiCode0/hezzl/internal/domain/usecase"
)

func (h *_handler) createGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	projectID, err := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusBadRequest, errParseParam, err, dusecase.ErrorDetails{"param": "projectId"}))
		return
	}

	data, err := jsonParse[dusecase.CreateGoodReq](r.Body)
	if err != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusBadRequest, errParseBody, err, nil))
		return
	}

	good, errResp := h.usecase.CreateGood(r.Context(), projectID, data)
	if errResp != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusInternalServerError, errResp.Msg, err, nil))
		return
	}

	writeJson(w, http.StatusOK, good)
}

func (h *_handler) updateGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusBadRequest, errParseParam, err, dusecase.ErrorDetails{"param": "id"}))
		return
	}

	data, err := jsonParse[dusecase.UpdateGoodReq](r.Body)
	if err != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusBadRequest, errParseParam, err, nil))
		return
	}

	good, errResp := h.usecase.UpdateGood(r.Context(), id, data)
	if errResp != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusInternalServerError, errResp.Msg, err, nil))
		return
	}

	writeJson(w, http.StatusOK, good)
}

func (h *_handler) removeGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusBadRequest, errParseParam, err, dusecase.ErrorDetails{"param": "id"}))
		return
	}

	good, errResp := h.usecase.RemoveGood(r.Context(), id)
	if errResp != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusInternalServerError, errResp.Msg, err, nil))
		return
	}

	writeJson(w, http.StatusOK, good)
}

func (h *_handler) getGoodsList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusBadRequest, errParseParam, err, dusecase.ErrorDetails{"param": "limit"}))
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusBadRequest, errParseParam, err, dusecase.ErrorDetails{"param": "offset"}))
		return
	}

	goods, errResp := h.usecase.GetGoodsList(r.Context(), limit, offset)
	if errResp != nil {
		writeErrorLog(w, dusecase.NewError(http.StatusInternalServerError, errResp.Msg, err, nil))
		return
	}

	writeJson(w, http.StatusOK, goods)
}
