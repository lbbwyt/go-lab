package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-lab/pkg/utils"
	"testing"
	"time"
)

type TestSendData struct {
	DeviceId    string `json:"device_id"`
	DataType    string `json:"data_type"`
	ImageBase64 string `json:"image_base64"`
	Tmp         int64  `json:"tmp"`
}

type TestRevData struct {
	DeviceId string    `json:"device_id"`
	Content  string    `json:"content"`
	Labels   VlmLabels `json:"labels"`
}

type VlmLabels struct {
	DataType string       `json:"data_type"`
	Data     []LabelValue `json:"data"`
}

type LabelValue struct {
	Type  string `json:"type"`
	Value string `json:"labels"`
}

func TestNewWebSocketClient(t *testing.T) {

	var singalChan = make(chan int, 1)

	client := NewWebSocketClient("ws://172.19.202.20:18765/ws", &Config{
		AutoReconnect: true,
		ReconnectWait: 3 * time.Second,
		BufferSize:    50,
		Origin:        "http://localhost:8080",
	})

	if err := client.Connect(); err != nil {
		log.Fatal("Initial connect failed:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println()

	// 启动自动重连管理器
	go client.AutoReconnect(ctx)

	fmt.Println(fmt.Printf("开始发送图片At：%s ", time.Now().Format("2006-01-02 15:04:05")))

	// 消息发送
	go func() {

		base64Str, err := utils.CompassFileToBase64("./cam01_1752031153.jpg")

		if err != nil {
			log.Errorf("ToBase64 error， Err： %v", err)
			return
		}

		d := &TestSendData{
			DeviceId:    "cam01",
			DataType:    "image",
			ImageBase64: base64Str,
			Tmp:         time.Now().Unix(),
		}

		err = client.SendJSON(d)
		if err != nil {
			client.isConnected = false
			log.Errorf("发送异常， Err： %v", err)
		}
		time.Sleep(time.Second)

		for {
			select {
			case <-singalChan:
				if client.isConnected {
					err = client.SendJSON(d)
					if err != nil {
						client.isConnected = false
						log.Errorf("发送异常， Err： %v", err)
					}
				}
			}
		}

	}()

	// 接收消息
	for {
		select {
		case msg := <-client.Receive():
			rv := &TestRevData{}
			err := json.Unmarshal(msg, rv)
			if err != nil {
				log.Errorf("Unmarshal error， Err： %v", err)
				continue
			}

			singalChan <- 1

			log.Printf("Received: %s At %s", rv.DeviceId, time.Now().Format("2006-01-02 15:04:05"))
		case err := <-client.Errors():
			log.Printf("Received Error: %v", err)
		}
	}
}
