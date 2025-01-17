package webinterface

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *WebHandler) *gin.Engine {
	r := gin.Default()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
	}
	protectedGroup := r.Group("/protected")
	protectedGroup.Use(AuthMiddleware(userHandler))
	{
		protectedGroup.POST("/changePassword", userHandler.ChangePassword)
		protectedGroup.GET("/userinfo", userHandler.GetUserInfo)
	}
	return r
}
