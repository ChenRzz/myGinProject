package main

import (
	"my_gin_project/application"
	"my_gin_project/infrastructure"
	"my_gin_project/webinterface"
)

func main() {
	infrastructure.InitDB()
	infrastructure.InitRedis()
	userRepository := infrastructure.NewUserRepository()
	sessionRepository := infrastructure.NewSessionManger()
	userService := application.NewUserService(userRepository, sessionRepository)
	webHandler := webinterface.NewWebHandler(userService)
	router := webinterface.SetupRouter(webHandler)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
