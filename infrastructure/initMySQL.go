package infrastructure

import (
	"fmt"
	"log"
	"my_gin_project/domain/entity"
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
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	fmt.Println("用户表迁移成功")
}

func GetDB() *gorm.DB {
	onceInitDB.Do(InitDB)
	return db
}
