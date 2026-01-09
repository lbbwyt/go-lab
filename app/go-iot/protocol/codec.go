package protocol

import "fmt"

// Codec 同时服务「上行解码」与「下行编码」
type Codec interface {
	// Decode 设备 → 云：把原始字节解析成 map[string]interface{}
	Decode([]byte) (map[string]interface{}, error)
	// Encode 云 → 设备：把 map 编码成原始字节
	Encode(map[string]interface{}) ([]byte, error)
	// Name 名称，用于产品元数据里绑定
	Name() string
}

var (
	codecMap = map[string]Codec{}
)

func RegisterProtocol(protocolId string, codec Codec) {
	codecMap[protocolId] = codec
}

func GetProtocol(protocolId string) (Codec, bool) {
	if v, ok := codecMap[protocolId]; ok {
		return v, true
	}
	return nil, false
}

func init() {
	fmt.Println("************************init json Codec************************")
	jsonCodec := &JSONCodec{}
	// 默认注册json 编解码器
	RegisterProtocol(jsonCodec.Name(), jsonCodec)
}
