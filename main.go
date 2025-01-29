package main

import (
	"log"
	"my_gin_project/consumerMQ"
	"my_gin_project/utils"
	"my_gin_project/webinterface/handler"
	"my_gin_project/webinterface/router"
)

func main() {
	webHandler := handler.NewWebHandler()
	routers := router.SetupRouter(webHandler)

	utils.SafeGo(func() {
		log.Println("启动 Kafka 消费者...")
		consumerMQ.NewKafkaUserConsumer().Start()
	})

	err := routers.Run(":8081")
	if err != nil {
		return
	}
}
