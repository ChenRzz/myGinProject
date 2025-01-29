package infrastructure

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	onceInitDB sync.Once
)

func InitDB() {
	dsn := "gin_user:gin123@tcp(127.0.0.1:3306)/gin_project?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}
	fmt.Println("数据库连接成功！")
}

func GetDB() *gorm.DB {
	onceInitDB.Do(InitDB)
	return db.Session(&gorm.Session{Context: &gin.Context{}})
}
