package dusecase

import (
	"fmt"
)

type ErrorResp struct {
	Code    int          `json:"code"`
	Msg     string       `json:"msg"`
	Details ErrorDetails `json:"details"`
}

type ErrorDetails map[string]any

func NewError(code int, msg string, err error, details ErrorDetails) *ErrorResp {
	if details == nil {
		details = make(ErrorDetails)
	}
	if err != nil {
		details["err"] = err
	}
	return &ErrorResp{
		Code:    code,
		Msg:     msg,
		Details: details,
	}
}

func (e *ErrorResp) String() string {
	return fmt.Sprintf(e.Msg+":%v", e.Details)
}
