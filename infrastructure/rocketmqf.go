package infrastructure

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"my_gin_project/domain"
)

type RocketmqPublisher struct {
	producer rocketmq.Producer
}

func NewRocketMQPublisher(nameServers []string) (*RocketmqPublisher, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(nameServers),
		producer.WithRetry(2),
	)
	if err != nil {
		return nil, err
	}

	err = p.Start()
	if err != nil {
		return nil, err
	}

	return &RocketmqPublisher{producer: p}, nil
}

// Publish 发送事件到 RocketMQ
func (r *RocketmqPublisher) Publish(event domain.Event) error {
	msg := &primitive.Message{
		Topic: event.Name,
		Body:  event.Body,
	}
	_, err := r.producer.SendSync(context.Background(), msg)
	return err
}
