package bk_res

type BkResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewErrBkResult(msg string) *BkResult {
	return &BkResult{
		Result:  false,
		Code:    -1,
		Message: msg,
		Data:    nil,
	}
}
