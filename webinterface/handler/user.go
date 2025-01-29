package handler

import (
	"github.com/gin-gonic/gin"
	"my_gin_project/application"
	"my_gin_project/infrastructure"
	"net/http"
)

type WebHandler struct {
	userApplication application.IUserApplication
}

func NewWebHandler() *WebHandler {
	return &WebHandler{userApplication: application.NewUserApplication()}
}

func (h *WebHandler) Register(c *gin.Context) {
	var reqEvent infrastructure.RegisterEvent
	if err := c.ShouldBind(&reqEvent.Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reqEvent.Topic = "user-register"

	registerProducer := infrastructure.NewKafkaProducer()
	defer registerProducer.Producer.Close()
	err := registerProducer.Publish(&reqEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已加入注册队列"})
}

func (h *WebHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	SessionID, err := h.userApplication.Login(c, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "登录成功",
		"sessionID": SessionID,
	})
}

func (h *WebHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到userID"})
	}
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID类型错误"})
	}
	err := h.userApplication.ChangeUserPassword(uid, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功！"})
}

func (h *WebHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到userID"})
	}
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID类型错误"})
	}
	user, err := h.userApplication.GetUserInfo(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "未找到用户"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
