package redis_series

//Redis consumer, 订阅topic, 并从bus 中读取property和事件上报数据， 然后持久化到redis 时间序列数据库中

import (
	"fmt"
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	"go-lab/app/go-iot/message_bus"
	"strings"
)

type Consumer interface {
	Consume(topic string)
}

type RedisConsumer struct {
	cli *redistimeseries.Client
}

// Consume  topic 为： iot/+/+/up   全部
func (rc *RedisConsumer) Consume(topic string) error {

	ch, err := message_bus.Subscribe(topic)
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-ch:
			saveMsg(rc.cli, msg)
		default:

		}
	}

}

func saveMsg(cli *redistimeseries.Client, msg message_bus.Message) {
	parts := strings.Split(msg.Topic, "/")
	productId := parts[1]
	deviceId := parts[2]

	key := fmt.Sprintf("sensor:%s:%s", productId, deviceId)

}
