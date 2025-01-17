package infrastructure

import (
	"errors"
	"gorm.io/gorm"
	"my_gin_project/domain"
)

type userRepository struct {
	db *gorm.DB
}

// 依赖注入
func NewUserRepository() domain.UserRepository {
	return &userRepository{db: GetDB()}
}

// 实现接口：
// FindByUsername(username string) (*User, error)
//
//	FindByEmail(email string) (*User, error)
//	Create(user *User) error
//	Update(user *User) error
//FindByUserID(id uint) (*User, error)

func (repo *userRepository) FindByUserID(id uint) (*domain.User, error) {
	var user domain.User
	err := repo.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (repo *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := repo.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
func (repo *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (repo *userRepository) Create(user *domain.User) error {
	err := repo.db.Create(user).Error
	if err != nil {
		return errors.New("Failed to create user")
	}
	return nil
}
func (repo *userRepository) Update(user *domain.User) error {
	err := repo.db.Save(user).Error
	if err != nil {
		return errors.New("Failed to update user")
	}
	return nil
}
