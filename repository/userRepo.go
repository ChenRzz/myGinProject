package repository

import (
	"errors"
	"my_gin_project/domain/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	baseRepo1[entity.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) FindByUsername(db *gorm.DB, username string) (*entity.User, error) {
	var user entity.User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
func (repo *UserRepository) FindByEmail(db *gorm.DB, email string) (*entity.User, error) {
	var user entity.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
