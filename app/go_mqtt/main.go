package main

import (
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zhaoguanproxy main.go

// nohup ./proxy -t tcp://172.16.101.74:8022 > runoob.log 2>&1 &

type MqttMsg struct {
	topic   string `json:"topic"`
	payload []byte `json:"payload"`
}

var (
	proxySourceAddr string
	sourceUserName  string
	sourcePassword  string

	proxyTargetAddr string
	targetUserName  string
	targetPassword  string
	msgChan         = make(chan MqttMsg, 16)
)

func init() {
	flag.StringVar(&proxySourceAddr, "s", fmt.Sprintf("wss://%s:%d/mqtt", "wss8084.megahealth.cn", 443), "")
	flag.StringVar(&sourceUserName, "su", "user_zhbl02", "")
	flag.StringVar(&sourcePassword, "sp", "123456", "")
	flag.StringVar(&proxyTargetAddr, "t", fmt.Sprintf("tcp://%s:%d", "172.16.101.74", 8022), "")
	//flag.StringVar(&proxyTargetAddr, "t", fmt.Sprintf("tcp://%s:%d", "127.0.0.1", 1883), "")
	flag.StringVar(&targetUserName, "tu", "admin", "")
	flag.StringVar(&targetPassword, "tp", "admin", "")
	flag.Parse()
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	mqttMessage := MqttMsg{
		topic:   msg.Topic(),
		payload: msg.Payload(),
	}
	msgChan <- mqttMessage
	//收到消息后转发消息

}

var messageTargetPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

}

var connectSourceHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Source Connected")
}

var connectTargetHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Target Connected")
}

var connectSourceLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Source Connect lost: %v", err)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("reconnect error: %v", err)
	}
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Source Connect lost: %v at %s", err, time.Now().Format("2006-01-02 15:04:05 MST Mon"))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Source Reconnect error: %v", err)
	}
	fmt.Printf(" Reconnectd ")
}

func main() {
	log.Info(fmt.Sprintf("proxy form %s to %s at %s", proxySourceAddr, proxyTargetAddr, time.Now().Format("2006-01-02 15:04:05 MST Mon")))
	opts := mqtt.NewClientOptions()
	opts.AddBroker(proxySourceAddr)
	//opts.SetClientID("267A69A74F31E4FB56B61DCACB85C7A472B4D27141C19918211E85A416D7F836")
	opts.SetClientID("D4D03F21C67373686D4F8246A720870B6125FBC974559173A870926EDE5202EF")
	opts.SetUsername(sourceUserName)
	opts.SetPassword(sourcePassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectSourceHandler
	opts.OnConnectionLost = connectLostHandler
	//opts.SetCleanSession(true)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(connectSourceLostHandler)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topics := []string{
		"test",
		"web/180/point",
		"web/180/fall",
		"web/J01MT05B2112000126/upline",
		"web/J01MT05B2112000126/downline",
	}
	for _, v := range topics {
		go client.Subscribe(v, 1, nil)
	}

	go proxyMsg()
	select {}
}

func proxyMsg() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(proxyTargetAddr)
	opts.SetClientID("J01MT05B2112000126")
	opts.SetUsername(targetUserName)
	opts.SetPassword(targetPassword)
	opts.SetDefaultPublishHandler(messageTargetPubHandler)
	opts.OnConnect = connectTargetHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetCleanSession(true)
	opts.SetKeepAlive(30 * time.Second)
	opts.AutoReconnect = true
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		select {
		case data := <-msgChan:
			publish(client, data)
		default:
		}
	}

}

/**
发送消息
*/
func publish(client mqtt.Client, data MqttMsg) {
	client.Publish(fmt.Sprintf("%s%s", "/proxy/", data.topic), 0, false, data.payload)
}
