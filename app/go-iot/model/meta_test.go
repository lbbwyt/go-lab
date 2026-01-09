package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMeta_Property(t *testing.T) {
	meta := &Meta{
		Schema:     "https://iot.aliyun.com/schema.json",
		ProductKey: "ac001",
		Version:    "1.0",
		Properties: []Property{
			{
				Identifier: "temperature",
				Name:       "当前温度",
				AccessMode: "r",
				Required:   true,
				DataType: DataType{
					Type: "float",
					Specs: Specs{
						Min:  float64Ptr(-50),
						Max:  float64Ptr(100),
						Unit: "℃",
					},
				},
			},
			{
				Identifier: "switch",
				Name:       "电源开关",
				AccessMode: "rw",
				DataType:   DataType{Type: "bool"},
			},
		},
		Services: []Service{
			{
				Identifier: "reboot",
				Name:       "重启设备",
				CallType:   "async",
				Input:      nil,
				Output:     nil,
			},
		},
		Events: []Event{
			{
				Identifier: "overTemp",
				Name:       "高温告警",
				Type:       "alert",
				Output: []Param{
					{Identifier: "temp", Name: "温度", DataType: DataType{Type: "float"}},
				},
			},
		},
	}

	bs, _ := json.Marshal(meta)
	fmt.Println(string(bs))
}

func float64Ptr(v float64) *float64 {
	return &v
}
