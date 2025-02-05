package webinterface

import (
	"github.com/gin-gonic/gin"
	"my_gin_project/application"
	"net/http"
	"regexp"
)

type WebHandler struct {
	uerService *application.UserService
}

func NewWebHandler(service *application.UserService) *WebHandler {
	return &WebHandler{uerService: service}
}
func (h *WebHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usernameRegex := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(usernameRegex, req.Username)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名只能包含字母、数字和下划线"})
		return
	}
	passwordRegex := `^[a-zA-Z0-9]+$`
	matched, err = regexp.MatchString(passwordRegex, req.Password)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码只能包含字母和数字"})
		return
	}
	err = h.uerService.RegisterPublish(req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已加入注册队列"})
}

func (h *WebHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	SessionID, err := h.uerService.Login(c, req.Username, req.Password)
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
	passwordRegex := `^[a-zA-Z0-9]+$`
	matched, err := regexp.MatchString(passwordRegex, req.NewPassword)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码只能包含字母和数字"})
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
	err = h.uerService.ChangeUserPassword(uid, req.OldPassword, req.NewPassword)
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
	user, err := h.uerService.GetUserInfo(uid)
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
