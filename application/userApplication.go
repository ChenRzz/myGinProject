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
	RegisterPublish(username, password, email string) error
	Register(username, password, email string) error
	Login(c *gin.Context, username, password string) (string, error)
	ChangeUserPassword(userID uint, oldPassword, newPassword string) error
	GetUserInfo(userID uint) (*entity.User, error)
	IsLogged(ctx context.Context, sessionID string) (uint, error)
}

type userApplication struct {
	userService service.IUserService
}

func (u userApplication) RegisterPublish(username, password, email string) error {
	return u.userService.RegisterPublish(username, password, email)
}

func (u userApplication) Register(username, password, email string) error {
	db := infrastructure.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		err := u.userService.CheckUserIsExist(tx, username, email)
		if err != nil {
			return err
		}
		return u.userService.Register(db, username, password, email)
	})
}

func (u userApplication) Login(c *gin.Context, username, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (u userApplication) ChangeUserPassword(userID uint, oldPassword, newPassword string) error {
	//TODO implement me
	panic("implement me")
}

func (u userApplication) GetUserInfo(userID uint) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userApplication) IsLogged(ctx context.Context, sessionID string) (uint, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserApplication() IUserApplication {
	return &userApplication{
		userService: service.NewUserService(),
	}
}
