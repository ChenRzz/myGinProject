package router

import (
	"my_gin_project/webinterface/handler"
	"my_gin_project/webinterface/mw"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.WebHandler) *gin.Engine {
	r := gin.Default()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
	}
	protectedGroup := r.Group("/protected")
	protectedGroup.Use(mw.NewAuthMiddleware().Handle())
	{
		protectedGroup.POST("/changePassword", userHandler.ChangePassword)
		protectedGroup.GET("/userinfo", userHandler.GetUserInfo)
	}
	return r
}
