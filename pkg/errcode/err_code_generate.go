//go:generate stringer -type ErrCode -linecomment -output err_code_string.go

package errcode

//每次定义错误码的时候，同时需要添加描述信息。而且描述信息经常会忘。
//go generate + stringer 的方式定义错误码

//stringer并不是 Go 自带的工具，需要手动安装。可以执行下面的命令安装：
//$ go get golang.org/x/tools/cmd/stringer
//安装完成后使用go generate 在对应errcode包目录下生成错误码文件：err_code_string.go
type ErrCode int

const (
	ERR_CODE_OK             ErrCode = 0 // OK
	ERR_CODE_INVALID_PARAMS ErrCode = 1 // 无效参数
	ERR_CODE_TIMEOUT        ErrCode = 2 // 超时
)
