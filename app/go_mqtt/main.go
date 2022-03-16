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

	//flag.StringVar(&proxySourceAddr, "s", fmt.Sprintf("tcp://%s:%d", "127.0.0.1", 1883), "")
	//flag.StringVar(&sourceUserName, "su", "admin", "")
	//flag.StringVar(&sourcePassword, "sp", "admin", "")

	//flag.StringVar(&proxyTargetAddr, "t", fmt.Sprintf("tcp://%s:%d", "172.16.101.74", 8022), "")
	flag.StringVar(&proxyTargetAddr, "t", fmt.Sprintf("tcp://%s:%d", "127.0.0.1", 1883), "")
	flag.StringVar(&targetUserName, "tu", "admin", "")
	flag.StringVar(&targetPassword, "tp", "admin", "")
	flag.Parse()
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Info("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	mqttMessage := MqttMsg{
		topic:   msg.Topic(),
		payload: msg.Payload(),
	}
	msgChan <- mqttMessage
	//收到消息后转发消息

}

var messageTargetPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Info("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

}

var connectSourceHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info("Source Connected")
}

var connectTargetHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info("Target Connected")
}

var connectSourceLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Info("Source Connect lost: %v at %s", err, time.Now().Format("2006-01-02 15:04:05 MST Mon"))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error("reconnect error: %v", err)
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

	err, client := subscrible()
	if err != nil || client == nil {
		log.Info("失败重连...！")
		for {
			err, _ := subscrible()
			if err == nil {
				break
			}
			time.Sleep(10 * time.Second)
		}
		log.Info("重连成功！")
	}

	go proxyMsg()

	go monitor(client)

	select {}
}

//监控短信处理
func monitor(client mqtt.Client) {

	tick := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-tick.C:
			sourceMonitor(client)
		}
	}

}

func sourceMonitor(client mqtt.Client) {
	if !client.IsConnected() {
		log.Info("监控到掉线重连")
		client.Connect()
	}

}

func subscrible() (error, mqtt.Client) {

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
	opts.SetCleanSession(true)
	opts.SetKeepAlive(30 * time.Second)
	//opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(connectSourceLostHandler)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error("connect error:" + token.Error().Error())
		return token.Error(), nil
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
	return nil, client
}

func proxyMsg() error {
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
		log.Error("目的borker 连接失败")
		return token.Error()
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
