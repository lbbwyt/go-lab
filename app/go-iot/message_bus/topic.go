package message_bus

import "fmt"

const (
	UpTopic = "iot/%s/%s/up"
)

var (
	testUpTopic = fmt.Sprintf(UpTopic, "product001", "device001")
)

// iot/{productKey}/{deviceName}/up
// iot/{productKey}/{deviceName}/down
//
//iot/+/+/up   全部
//iot/ac001/+/down 下发给ac001产品下的所有设备
