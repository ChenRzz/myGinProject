package main

import (
	"log"
	"my_gin_project/config"
	"my_gin_project/consumerMQ"
	"my_gin_project/utils"
	"my_gin_project/webinterface/handler"
	"my_gin_project/webinterface/router"
)

func main() {
	webHandler := handler.NewWebHandler()
	router := router.SetupRouter(webHandler)

	utils.SafeGo(func() {
		log.Println("启动 RocketMQ 消费者...")
		consumerMQ.NewUserConsumer(config.Nameserver).Start()
	})

	err := router.Run(":8081")
	if err != nil {
		return
	}
}
