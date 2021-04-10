package pubsub

import "context"

type PubSubOptPublish func(interface{})

type PubSubOptSubscribe func(interface{})

type PubSubClient interface {
	Publish(ctx context.Context, publication *Publication, opts ...PubSubOptPublish) error
	Subscribe(ctx context.Context, pubs chan *Publication, errors chan error, topics []string, opts ...PubSubOptSubscribe)
	Connect(ctx context.Context) error
	Disconnect() error
}
