package mw

import (
	"my_gin_project/application"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AuthMiddleware struct {
	userApplication application.IUserApplication
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		userApplication: application.NewUserApplication(),
	}
}

func (a *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// 1. 从请求头中获取 sessionID
		sessionID := c.GetHeader("Authorization")
		if sessionID == "" {
			// 如果客户端没有提供 sessionID，拒绝访问
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供 sessionID"})
			c.Abort() // 阻止请求继续处理
			return
		}

		// 2. 验证 sessionID 是否存在于 Redis
		userID, err := a.userApplication.IsLogged(ctx, sessionID)
		if err == redis.Nil {
			// Redis 中找不到对应的 sessionID，表示登录失效
			c.JSON(http.StatusUnauthorized, gin.H{"error": "sessionID 无效或已过期"})
			c.Abort()
			return
		} else if err != nil {
			// Redis 出现其他错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis 服务错误"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next() // 请求继续处理
	}
}
