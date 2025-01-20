package main

import (
	"log"
	"my_gin_project/application"
	"my_gin_project/consumerMQ"
	"my_gin_project/infrastructure"
	"my_gin_project/webinterface"
)

func main() {
	infrastructure.InitDB()
	infrastructure.InitRedis()
	userRepository := infrastructure.NewUserRepository()
	sessionRepository := infrastructure.NewSessionManger()
	nameserver := []string{"172.18.0.4:9876"}
	log.Printf("Using RocketMQ NameServer: %v", nameserver)
	eventPublisher, err := infrastructure.NewRocketMQPublisher(nameserver)
	if err != nil {
		panic(err)
	}
	userService := application.NewUserService(eventPublisher, userRepository, sessionRepository)
	webHandler := webinterface.NewWebHandler(userService)
	router := webinterface.SetupRouter(webHandler)
	go func() {
		log.Println("启动 RocketMQ 消费者...")
		consumerMQ.StartUserConsumer(nameserver, userService)
	}()
	err = router.Run(":8081")
	if err != nil {
		return
	}
}
