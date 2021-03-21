package event

import (
	"fmt"
	"testing"
	"time"
)

func TestEventBus_Public(t *testing.T) {
	eventBus := NewEventBus()

	var (
		topic = "data.test"
		data  = "sdf"
	)

	go publish(topic, data, eventBus)

	ch := make(chan DataEvent)
	eventBus.Subscribe(topic, ch)
	for {
		select {
		case d := <-ch:
			fmt.Println(fmt.Sprintf("receive : %v", d))
		}
	}
}

func publish(topic string, data string, eb *EventBus) {

	tick := time.NewTicker(time.Second * 1)
	for {
		select {
		case t := <-tick.C:
			fmt.Println("start publish")
			go eventBus.Public(topic, fmt.Sprintf("%s:%s", data, t.String()))
		}
	}
}
