package repository

import (
	"my_gin_project/domain/entity"

	"gorm.io/gorm"
)

type IUserRepository interface {
	IBaseRepo[entity.User]
	FindByUsername(db *gorm.DB, username string) (*entity.User, error)
	FindByEmail(db *gorm.DB, email string) (*entity.User, error)
}
