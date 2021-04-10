package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"go-lab/pkg/pubsub"
	"go-lab/pkg/serialization"
	"time"
)

func init() {
	// 设置全局序列化反序列化器的类型: json/msgpackage,默认为json
	serialization.SetSerializer(string(pubsub.EncodingTypeJSON))
}

type Consumer struct {
	ready   chan bool
	message chan *sarama.ConsumerMessage
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the Consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a Consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for msg := range claim.Messages() {
		consumer.message <- msg
		session.MarkMessage(msg, "")
	}
	return nil
}

type kafkaPubSub struct {
	client        sarama.Client
	asyncProducer sarama.AsyncProducer
	consumer      sarama.Consumer
	consumerGroup sarama.ConsumerGroup
}

// endpoints：kafka集群地址
// groupId：分组id，如果不传，则初始话consumer客户端；如果传了groupId,则初始化consumerGroup
func NewKafkaPubSubClient(endpoints []string, groupId string, conf *sarama.Config) (pubsub.PubSubClient, error) {
	// 创建时建议从 client 开始，这样在错误处理时可以通过 client.RefreshMetadata() 来马上刷新 metadata，从而快速重试下
	client, err := sarama.NewClient(endpoints, conf)
	if err != nil {
		return nil, err
	}

	asyncProducer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	pubSubClient := &kafkaPubSub{
		client:        client,
		asyncProducer: asyncProducer,
	}

	if groupId == "" {
		consumer, err := sarama.NewConsumerFromClient(client)
		if err != nil {
			return nil, err
		}
		pubSubClient.consumer = consumer
	} else {
		consumerGroup, err := sarama.NewConsumerGroupFromClient(groupId, client)
		if err != nil {
			return nil, err
		}
		pubSubClient.consumerGroup = consumerGroup
	}

	return pubSubClient, nil
}

func (p *kafkaPubSub) Publish(ctx context.Context, publication *pubsub.Publication, opts ...pubsub.PubSubOptPublish) error {
	data, err := serialization.Marshal(publication)
	if err != nil {
		return fmt.Errorf("unable to encode publication. message dropped: %s", err.Error())
	}

	msg := &sarama.ProducerMessage{
		Topic: publication.Topic,
		Key:   nil,
		Value: sarama.ByteEncoder(data),
	}

	select {
	case p.asyncProducer.Input() <- msg:

	case <-time.After(500 * time.Millisecond):
		err = errors.New("producer message timeout error")
		// 这里加上超时设置，就是为了避免遇到发送 kafka 出现网络问题导致 buffer 满了，从而堵塞住，影响外部调用接口
		// 到这里，刷新下 metadata，然后重试一次
		if e := p.client.RefreshMetadata(publication.Topic); e != nil {
			return e
		} else {
			select {
			case p.asyncProducer.Input() <- msg:
				err = nil // 重试成功需要重置 err
			default: // 避免堵塞
			}
		}
	}

	// 发布结果处理
	select {
	case msg := <-p.asyncProducer.Errors():
		return msg.Err // 出错，直接返回error

	case <-p.asyncProducer.Successes():
		// 发布成功，什么也不做

	default: // 避免堵塞
	}

	return err
}

func (p *kafkaPubSub) Subscribe(ctx context.Context, pubs chan *pubsub.Publication, errors chan error, topics []string, opts ...pubsub.PubSubOptSubscribe) {
	if p.consumerGroup != nil {
		go p.consumerGroupListen(ctx, pubs, errors, topics)
	}

	if p.consumer != nil {
		for _, topic := range topics {
			go p.consumeListen(ctx, pubs, errors, topic)
		}
	}
}

func (p *kafkaPubSub) Connect(ctx context.Context) error {
	return nil
}

func (p *kafkaPubSub) Disconnect() (retErr error) {
	if p.consumerGroup != nil {
		if err := p.consumerGroup.Close(); err != nil {
			retErr = err
		}
	}

	if p.consumer != nil {
		if err := p.consumer.Close(); err != nil {
			retErr = err
		}
	}

	if err := p.asyncProducer.Close(); err != nil {
		retErr = err
	}

	if err := p.client.Close(); err != nil {
		retErr = err
	}

	return
}

func (p *kafkaPubSub) consumerGroupListen(ctx context.Context, pubs chan *pubsub.Publication, errors chan error, topics []string, opts ...pubsub.PubSubOptSubscribe) {
	consumer := Consumer{
		ready:   make(chan bool),
		message: make(chan *sarama.ConsumerMessage),
	}

	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the Consumer session will need to be
			// recreated to get the new claims
			if err := p.consumerGroup.Consume(ctx, topics, &consumer); err != nil {
				errors <- err
			}

			// check if context was cancelled, signaling that the Consumer should stop
			if ctx.Err() != nil {
				return
			}

			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the Consumer has been set up
	fmt.Println("Sarama Consumer up and running!...")

	for {
		select {
		case err := <-p.consumerGroup.Errors():
			errors <- err

		case msg := <-consumer.message:
			publication := &pubsub.Publication{}
			if err := serialization.Unmarshal(msg.Value, publication); err != nil {
				errors <- err
				continue
			}
			publication.Topic = msg.Topic
			pubs <- publication

		case <-ctx.Done():
			// close kafka asyncProducer、consumerGroup、client
			if err := p.Disconnect(); err != nil {
				errors <- err
			}
			return
		}
	}
}

func (p *kafkaPubSub) consumeListen(ctx context.Context, pubs chan *pubsub.Publication, errors chan error, topic string, opts ...pubsub.PubSubOptSubscribe) {
	// 拿到对应主题下所有分区
	partitionList, err := p.consumer.Partitions(topic)
	if err != nil {
		errors <- err
		return
	}

	// 遍历所有分区，异步消费
	for _, partition := range partitionList {
		// 从最新的位置开始订阅
		partitionConsumer, err := p.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			errors <- err
			return
		}
		go func(sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					publication := &pubsub.Publication{}
					if err := serialization.Unmarshal(msg.Value, publication); err != nil {
						errors <- err
						continue
					}
					publication.Topic = msg.Topic
					pubs <- publication

				case <-ctx.Done():
					// close kafka asyncProducer、consumer、consumerGroup、client
					if err := p.Disconnect(); err != nil {
						errors <- err
					}
					if err := partitionConsumer.Close(); err != nil {
						errors <- err
					}
					return
				}
			}
		}(partitionConsumer)
	}

}
