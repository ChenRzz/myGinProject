package consumerMQ

import (
	"context"
	"encoding/json"
	"log"
	"my_gin_project/application"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

// StartUserConsumer 启动用户注册消费者
func StartUserConsumer(nameServers []string, userService *application.UserService) {
	// 获取 RocketMQ NameServer
	log.Printf("RocketMQ Consumer 连接 NameServer: %s", nameServers[0])

	// 创建 RocketMQ Consumer
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("user-consumer-group"),
		consumer.WithNameServer(nameServers),
	)
	if err != nil {
		log.Fatalf("创建 RocketMQ Consumer 失败: %s", err.Error())
	}

	// 订阅 "user-registered" 事件
	err = c.Subscribe("user_register", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			var userData struct {
				Username string
				Password string
				Email    string
			}
			err = json.Unmarshal(msg.Body, &userData)
			if err != nil {
				return consumer.ConsumeRetryLater, err
			}

			// 调用 UserService 处理注册逻辑
			err := userService.Register(userData.Username, userData.Password, userData.Email)
			if err != nil {
				log.Printf("注册用户失败: %s", err.Error())
				return consumer.ConsumeRetryLater, err
			}

			log.Printf("成功存储用户: %+v", userData)
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		log.Fatalf("订阅用户注册事件失败: %s", err.Error())
	}

	err = c.Start()
	if err != nil {
		log.Fatalf("启动 Consumer 失败: %s", err.Error())
	}
	defer c.Shutdown()
	select {} // 让 Goroutine 保持运行
}
