package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"regexp"
	"time"
)

var db *gorm.DB
var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // Redis 密码（如果无密码，留空）
		DB:       0,                // 使用默认数据库
	})

	// 测试连接
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}
	fmt.Println("Redis 连接成功")
}
func initDB() {
	dsn := "gin_user:gin123@tcp(127.0.0.1:3306)/gin_project?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}
	fmt.Println("数据库连接成功！")
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	fmt.Println("用户表迁移成功")
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`    // 用户名不能为空
	Password string `json:"password" binding:"required"`    // 密码不能为空
	Email    string `json:"email" binding:"required,email"` // 邮箱不能为空，且必须符合邮箱格式
}
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:100;unique"`
	Password  string `gorm:"size:255"`
	Email     string `gorm:"size:100;unique"`
	CreatedAt time.Time
}

func LoginUser(c *gin.Context) {
	// 1. 获取请求数据并验证格式
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 检查用户是否存在
	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 3. 验证密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 4. 生成会话并保存到 Redis
	sessionID := uuid.NewString()                                    // 创建唯一的 session ID
	err = redisClient.Set(c, sessionID, user.ID, time.Hour*24).Err() // 保存 session，过期时间为 24 小时
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}

	// 5. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message":   "登录成功",
		"sessionID": sessionID,
	})
}

func RegisterUser(c *gin.Context) {
	// 1. 获取请求数据并进行参数绑定和验证
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usernameRegex := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(usernameRegex, req.Username)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名只能包含字母、数字和下划线"})
		return
	}

	// 校验密码是否合法
	passwordRegex := `^[a-zA-Z0-9]+$`
	matched, err = regexp.MatchString(passwordRegex, req.Password)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码只能包含字母和数字"})
		return
	}

	// 2. 检查用户名和邮箱是否已存在
	var existingUser User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	// 3. 加密密码（这里用 bcrypt，需安装包 "golang.org/x/crypto/bcrypt"）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 4. 创建新用户
	newUser := User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	// 5. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头中获取 sessionID
		sessionID := c.GetHeader("Authorization")
		if sessionID == "" {
			// 如果客户端没有提供 sessionID，拒绝访问
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供 sessionID"})
			c.Abort() // 阻止请求继续处理
			return
		}

		// 2. 验证 sessionID 是否存在于 Redis
		userID, err := redisClient.Get(context.Background(), sessionID).Result()
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

		// 3. 鉴权成功，将 userID 存入上下文供后续处理使用
		c.Set("userID", userID)
		c.Next() // 请求继续处理
	}
}

func GetUserInfo(c *gin.Context) {
	// 从上下文获取 userID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	// 从数据库中查询用户信息
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到"})
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func main() {
	initDB()
	initRedis()
	r := gin.Default()
	r.POST("/register", RegisterUser)
	r.POST("/login", LoginUser)
	protected := r.Group("/api")
	protected.Use(AuthMiddleware()) // 应用鉴权中间件

	// 示例受保护接口
	protected.GET("/userinfo", GetUserInfo)
	r.Run(":8080")
}
