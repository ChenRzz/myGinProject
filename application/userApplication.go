package application

import (
	"context"
	"my_gin_project/domain/entity"
	"my_gin_project/domain/service"
	"my_gin_project/infrastructure"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IUserApplication interface {
	//RegisterPublish(username, password, email string) error
	Register(username, password, email string) error
	Login(c *gin.Context, username, password string) (string, error)
	ChangeUserPassword(userID uint, oldPassword, newPassword string) error
	GetUserInfo(userID uint) (*entity.User, error)
	IsLogged(ctx context.Context, sessionID string) (uint, error)
}

type userApplication struct {
	userService service.IUserService
}

/*func (u userApplication) RegisterPublish(username, password, email string) error {
	return u.userService.RegisterPublish(username, password, email)
}*/

func (u userApplication) Register(username, password, email string) error {
	db := infrastructure.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		err := u.userService.CheckUserIsExist(tx, username, email)
		if err != nil {
			return err
		}
		return u.userService.Register(tx, username, password, email)
	})
}

func (u userApplication) Login(c *gin.Context, username, password string) (string, error) {
	//TODO implement me
	db := infrastructure.GetDB()
	var sessionID string
	err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		sessionID, err = u.userService.Login(tx, c, username, password)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (u userApplication) ChangeUserPassword(userID uint, oldPassword, newPassword string) error {
	//TODO implement me
	db := infrastructure.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		err := u.userService.ChangeUserPassword(tx, userID, oldPassword, newPassword)
		if err != nil {
			return err
		}
		return nil
	})
}

func (u userApplication) GetUserInfo(userID uint) (*entity.User, error) {
	//TODO implement me
	db := infrastructure.GetDB()
	var newUser *entity.User
	err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		newUser, err = u.userService.GetUserInfo(tx, userID)
		if err != nil {
			return err
		}
		return nil
	})
	return newUser, err
}

func (u userApplication) IsLogged(ctx context.Context, sessionID string) (uint, error) {
	//TODO implement me
	userID, err := u.userService.IsLogged(ctx, sessionID)
	if err != nil {
		return 0, err
	}
	return userID, nil

}

func NewUserApplication() IUserApplication {
	return &userApplication{
		userService: service.NewUserService(),
	}
}
