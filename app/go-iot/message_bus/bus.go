package message_bus

// Bus 内存总线
type Bus interface {
	Publish(topic string, payload interface{}) error
	Subscribe(topic string) (interface{}, error)
	Close() error
}

var (
	memBus = NewMemBus()
)

func Publish(topic string, data map[string]interface{}) error {

	return memBus.Publish(topic, data)
}

func Subscribe(topic string) (<-chan Message, error) {
	r, err := memBus.Subscribe(topic)
	if err != nil {
		return nil, err
	}
	return r.(<-chan Message), nil
}
