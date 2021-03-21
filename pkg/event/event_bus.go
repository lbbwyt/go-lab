package event

import "sync"

var (
	eventBusOnce sync.Once
	eventBus     *EventBus
)

type EventBus struct {
	subscribers map[string]DataChannelSlice // key 为topic， 一个topic会有多个消费者订阅
	rm          sync.Mutex
}

func NewEventBus() *EventBus {
	if eventBus != nil {
		return eventBus
	}
	eventBusOnce.Do(func() {
		eventBus = &EventBus{
			subscribers: make(map[string]DataChannelSlice, 0),
			rm:          sync.Mutex{},
		}
	})
	return eventBus
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()
}

func (eb *EventBus) Public(topic string, data interface{}) {
	eb.rm.Lock()
	if chans, found := eb.subscribers[topic]; found {
		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, slice DataChannelSlice) {
			for _, ch := range slice {
				ch <- data
			}
		}(DataEvent{
			Data:  data,
			Topic: topic,
		}, channels)
	}
	eb.rm.Unlock()
}
