package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
}

func NewKafkaProducer() *KafkaProducer {
	brokers := []string{"localhost:9092"}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // 让 Producer 确保消息已发送

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}
	return &KafkaProducer{Producer: producer}
}
func (k *KafkaProducer) Publish(event *RegisterEvent) error {
	eventData, err := json.Marshal(event.Body)
	if err != nil {
		log.Fatalf("Failed to encode JSON: %v", err)
	}
	message := &sarama.ProducerMessage{
		Topic: event.Topic,
		Value: sarama.StringEncoder(eventData),
	}
	partition, offset, err := k.Producer.SendMessage(message)
	if err != nil {
		log.Fatalf("发布消息失败: %v", err)
		return err
	}
	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	return nil
}
