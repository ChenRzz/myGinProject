package infrastructure

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type RocketmqPublisher struct {
	producer rocketmq.Producer
}

func NewRocketMQPublisher(nameServers []string) *RocketmqPublisher {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(nameServers),
		producer.WithRetry(2),
	)
	if err != nil {
		panic(err)
	}

	err = p.Start()
	if err != nil {
		panic(err)
	}

	return &RocketmqPublisher{producer: p}
}

// Publish 发送事件到 RocketMQ
func (r *RocketmqPublisher) Publish(event Event) error {
	msg := &primitive.Message{
		Topic: event.Name,
		Body:  event.Body,
	}
	_, err := r.producer.SendSync(context.Background(), msg)
	return err
}
