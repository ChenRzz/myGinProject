package consumerMQ

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"my_gin_project/application"
)

type KafkaUserConsumer struct {
	userApplication application.IUserApplication
}

func NewKafkaUserConsumer() Consumer {
	return &KafkaUserConsumer{userApplication: application.NewUserApplication()}
}

func (K *KafkaUserConsumer) Start() {
	brokers := []string{"localhost:9092"}
	topic := "user-register"
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true // 让 Consumer 返回错误信息

	// 3. 创建 Kafka 消费者
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close() // 关闭消费者

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Failed to get partitions: %v", err)
	}
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Failed to start partition consumer: %v", err)
		}
		defer pc.Close()

		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				var RegisReq RegisterReqMQ
				err := json.Unmarshal(msg.Value, &RegisReq)
				if err != nil {
					log.Printf("Failed to decode JSON: %v", err)
					continue
				}
				err = K.userApplication.Register(RegisReq.Username, RegisReq.Password, RegisReq.Email)
				if err != nil {
					log.Printf("Failed to register user", err)
					continue
				}
				log.Printf("用户%v已注册", RegisReq.Username)
			}
		}(pc)
	}
	select {}
}
