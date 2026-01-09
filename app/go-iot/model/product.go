package model

type Product struct {
	Id                uint   `json:"id"                description:""`
	Key               string `json:"key"               description:"产品标识"`
	Name              string `json:"name"              description:"产品名称"`
	MessageProtocol   string `json:"messageProtocol"   description:"消息协议"`
	TransportProtocol string `json:"transportProtocol" description:"传输协议: MQTT,COAP,UDP"`
	ProtocolId        uint   `json:"protocolId"        description:"协议id"`
	DeviceType        string `json:"deviceType"        description:"设备类型: 网关，设备，子设备"`
	Metadata          Meta   `json:"metadata"          description:"物模型"`
}
