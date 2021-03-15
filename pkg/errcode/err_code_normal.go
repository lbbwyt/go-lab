package errcode

import "fmt"

// 使用传统方式定义错误码
const (
	ERR_CODE_OK_Normal             = 0 // OK
	ERR_CODE_INVALID_PARAMS_Normal = 1 // 无效参数
	ERR_CODE_TIMEOUT_Normal        = 2 // 超时
)

// 定义错误码与描述信息的映射
var mapErrDesc = map[int]string{
	ERR_CODE_OK_Normal:             "OK",
	ERR_CODE_INVALID_PARAMS_Normal: "无效参数",
	ERR_CODE_TIMEOUT_Normal:        "超时",
}

// 根据错误码返回描述信息
func GetDescription(errCode int) string {
	if desc, exist := mapErrDesc[errCode]; exist {
		return desc
	}
	return fmt.Sprintf("error code: %d", errCode)
}
