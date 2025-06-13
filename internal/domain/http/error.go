package dhttp

type errorResp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Details any    `json:"details"`
}

func NewError(code int, msg string, details any) *errorResp {
	return &errorResp{
		Code:    code,
		Msg:     msg,
		Details: details,
	}
}
