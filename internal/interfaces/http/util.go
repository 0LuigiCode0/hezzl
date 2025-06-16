package ihttp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	dusecase "github.com/0LuigiCode0/hezzl/internal/domain/usecase"
)

func jsonParse[T any](data io.Reader) (*T, error) {
	buf, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf(errReadBody, err)
	}

	out := new(T)
	return out, json.Unmarshal(buf, out)
}

func writeJson(w http.ResponseWriter, statusCode int, data any) {
	payload, err := json.Marshal(data)
	if err != nil {
		errResp := dusecase.NewError(500, consts.ErrJsonMarshal, err, nil)
		log.Print(errResp)
		payload, err = json.Marshal(errResp)
		if err != nil {
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(payload)
}

func writeErrorLog(w http.ResponseWriter, err *dusecase.ErrorResp) {
	log.Print(prefix, err)

	writeJson(w, mapCode(err.Code), err)
}

func mapCode(code int) int {
	switch code {
	case 3:
		return http.StatusNotFound
	default:
		return code
	}
}
