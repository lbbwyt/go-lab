package errcode

import (
	"fmt"
	"testing"
)

type monster struct {
	Name string
}

func TestErrCode_String(t *testing.T) {
	fmt.Println(ERR_CODE_INVALID_PARAMS)
	fmt.Println(fmt.Sprintf("错误码%s：%d", ERR_CODE_INVALID_PARAMS.String(), ERR_CODE_INVALID_PARAMS))
	fmt.Println(monster{Name: ERR_CODE_INVALID_PARAMS.String()})
}
