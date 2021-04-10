package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"go-lab/pkg/pubsub"
	"testing"
	"time"
)

//下载zookeeper
//docker pull wurstmeister/zookeeper

//启动zookeeper
//docker run -d --name zookeeper --publish 2181:2181 --volume /etc/localtime:/etc/localtime wurstmeister/zookeeper

//启动kafka
//docker pull wurstmeister/kafka:2.13-2.6.0

//docker run -d --name kafka --publish 9092:9092 \
//--link zookeeper \
//--env KAFKA_ZOOKEEPER_CONNECT=192.168.0.102:2181 \
//--env KAFKA_ADVERTISED_HOST_NAME=192.168.0.102 \
//--env KAFKA_ADVERTISED_PORT=9092  \
//--volume /etc/localtime:/etc/localtime \
//wurstmeister/kafka:2.13-2.6.0

var (
	endpoint = []string{"192.168.0.102:9092"}
	topic    = "test1"
	groupId  = "testGroup"
)

type student struct {
	Name string
	Age  int
}

type StudentPayload struct {
	pubsub.EventBaseInfo
	Student student
}

func Test_kafkaPubSub_Publish(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = 5
	conf.Producer.RequiredAcks = sarama.WaitForAll

	client, err := NewKafkaPubSubClient(endpoint, groupId, conf)
	if err != nil {
		panic(err)
	}

	stu := StudentPayload{
		EventBaseInfo: pubsub.EventBaseInfo{
			DataSource: pubsub.DataSourceAppTest,
			EventType:  pubsub.EventCreate,
			TimeStamp:  time.Now().Unix(),
		},
		Student: student{
			Name: "张三",
			Age:  20,
		},
	}

	publication, err := pubsub.NewPublication(topic, pubsub.EncodingTypeJSON, &stu)
	if err != nil {
		panic(err)
	}

	if err := client.Publish(ctx, publication); err != nil {
		panic(err)
	}
	fmt.Println("publish success: ", stu)

	// 等待client发布完成
	time.Sleep(1 * time.Second)
}

func Test_kafkaConsumerGroup_Subscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion("2.6.0")
	if err != nil {
		panic(err)
	}
	// 指定 kafka server 版本号，避免因为版本导致协议不同，从而出现数据格式不一致的情况
	//conf.Version = sarama.V2_6_0_0
	conf.Version = version

	// 接受失败的 error，做一些错误处理
	conf.Consumer.Return.Errors = true
	// 指定此 Consumer group 第一次启动时从哪里开始消费
	// 默认是 OffsetNewest，代表从此 topic 的最新 offset 开始消费
	// 如果你的 Consumer group 是想回溯之前存储在 kafka 的所有消息，则修改为 OffsetOldest
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	// 配置后台协程刷新 metadata 的频率
	// 注意其实这里比 producer 配置的更频繁，具体原因后面单独讲
	conf.Metadata.RefreshFrequency = 3 * time.Minute

	client, err := NewKafkaPubSubClient(endpoint, groupId, conf)
	if err != nil {
		panic(err)
	}

	pubs := make(chan *pubsub.Publication)
	errors := make(chan error)

	client.Subscribe(ctx, pubs, errors, []string{topic})

	for {
		select {
		case pub := <-pubs:
			switch pub.Topic {
			case topic:
				fmt.Println("publication: ", pub)

				stu := StudentPayload{}
				if err := pub.Decode(&stu); err != nil {
					panic(err)
				}
				fmt.Println("stuPayload: ", stu)
			}

		case err := <-errors:
			if err != nil {
				fmt.Println("error: ", err)
				return
			}
		}
	}
}

func Test_kafkaConsumer_Subscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion("2.6.0")
	if err != nil {
		panic(err)
	}
	// 指定 kafka server 版本号，避免因为版本导致协议不同，从而出现数据格式不一致的情况
	//conf.Version = sarama.V2_6_0_0
	conf.Version = version

	// 接受失败的 error，做一些错误处理
	conf.Consumer.Return.Errors = true
	// 指定此 Consumer group 第一次启动时从哪里开始消费
	// 默认是 OffsetNewest，代表从此 topic 的最新 offset 开始消费
	// 如果你的 Consumer group 是想回溯之前存储在 kafka 的所有消息，则修改为 OffsetOldest
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	// 配置后台协程刷新 metadata 的频率
	// 注意其实这里比 producer 配置的更频繁，具体原因后面单独讲
	conf.Metadata.RefreshFrequency = 3 * time.Minute

	// 不传groupId则初始化consumer client；而非consumerGroupId
	client, err := NewKafkaPubSubClient(endpoint, "", conf)
	if err != nil {
		panic(err)
	}

	pubs := make(chan *pubsub.Publication)
	errors := make(chan error)

	client.Subscribe(ctx, pubs, errors, []string{topic})

	log.Println("consumer start!...")
	for {
		select {
		case pub := <-pubs:
			switch pub.Topic {
			case topic:
				fmt.Println("publication: ", pub)

				stu := StudentPayload{}
				if err := pub.Decode(&stu); err != nil {
					panic(err)
				}
				fmt.Println("stuPayload: ", stu)
			}

		case err := <-errors:
			if err != nil {
				fmt.Println("error: ", err)
				return
			}
		}
	}
}
