package ihttp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
)

func jsonParse[T any](data io.Reader) (*T, error) {
	buf, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf(consts.ErrReadBody, err)
	}

	out := new(T)
	return out, json.Unmarshal(buf, out)
}

func writeJson(w http.ResponseWriter, statusCode int, data any) {
	payload, err := json.Marshal(data)
	if err != nil {
		writeErrorFLog(w, 500, consts.ErrJsonMarshal, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(payload)
}

func writeErrorFLog(w http.ResponseWriter, statusCode int, format string, arg ...any) {
	err := fmt.Errorf(format, arg...)
	log.Print(err)

	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}
