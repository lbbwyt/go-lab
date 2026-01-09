package gateway

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"go-lab/app/go-iot/message_bus"
	"go-lab/app/go-iot/protocol"
	"strings"
)

type MqttTransport struct {
	cli     mqtt.Client
	handler func(ctx context.Context, msg *RawMessage)
}

func NewMqttTransport(opts *mqtt.ClientOptions) *MqttTransport {
	t := &MqttTransport{cli: mqtt.NewClient(opts)}
	t.SetHandler(handler)
	return t
}

// Listen 上行topic：  iot/{productKey}/{deviceName}/up, 通过可以productKey， 找到对应的codec， 来解码

func (m *MqttTransport) Listen(addr string) error {
	if token := m.cli.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	m.cli.Subscribe("iot/+/up", 1, func(c mqtt.Client, msg mqtt.Message) {
		parts := strings.Split(msg.Topic(), "/")
		m.handler(context.TODO(), &RawMessage{
			Network: "mqtt",
			Product: parts[1],
			Device:  parts[2],
			Payload: msg.Payload(),
		})
	})
	return nil
}

func (m *MqttTransport) SetHandler(h func(ctx context.Context, msg *RawMessage)) {
	m.handler = h
}

func (m *MqttTransport) Close() error {
	m.cli.Disconnect(250)
	return nil
}

func handler(ctx context.Context, rawMsg *RawMessage) {
	productId := rawMsg.Product
	deviceId := rawMsg.Device

	topic := fmt.Sprintf(message_bus.UpTopic, productId, deviceId)

	m := make(map[string]interface{})
	protocolId := getProtocolId(productId)
	codec, ok := protocol.GetProtocol(protocolId)
	if !ok {
		log.Warnf("codec does not exist， product:%s, protocol:%s", productId, protocolId)
		return
	}

	m, err := codec.Decode(rawMsg.Payload)
	if err != nil {
		log.Errorf("decode err : %s ", err.Error())
		return
	}

	message_bus.Publish(topic, m)
}

func getProtocolId(productID string) string {
	//todo 通过产品id , 获取对应的协议ID
	return ""
}
