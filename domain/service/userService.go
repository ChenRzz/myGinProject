package service

import (
	"context"
	"errors"
	"my_gin_project/domain/entity"
	"my_gin_project/repository"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserService interface {
	// @crz RegisterPublish 这个其实应该在 adapter
	//RegisterPublish(username, password, email string) error
	Register(db *gorm.DB, username, password, email string) error
	CheckUserIsExist(db *gorm.DB, username, email string) error
	Login(db *gorm.DB, ctx context.Context, username, password string) (string, error)
	ChangeUserPassword(db *gorm.DB, userid uint, oldPassword, newPassword string) error
	GetUserInfo(db *gorm.DB, userid uint) (*entity.User, error)
	IsLogged(ctx context.Context, sessionID string) (uint, error)
}

// test
type userService struct {
	//messagePublisher  infrastructure.EventProducer
	userRepository    repository.IUserRepository
	sessionRepository repository.ISessionRepository
}

func NewUserService() IUserService {
	return &userService{
		//messagePublisher:  infrastructure.NewRocketMQPublisher(config.Nameserver),
		userRepository:    repository.NewUserRepository(),
		sessionRepository: repository.NewSessionManger(),
	}
}

/*func (us *userService) RegisterPublish(username, password, email string) error {
	userdata := map[string]string{
		"username": username,
		"password": password,
		"email":    email,
	}
	// @crz Marshal、Unmarshal 都应该用 entity
	eventBody, err := json.Marshal(userdata)
	if err != nil {
		return err
	}
	newEvent := infrastructure.Event{
		Name: "user_register",
		Body: eventBody,
	}
	return us.messagePublisher.Publish(newEvent)
}*/

func (us *userService) Register(db *gorm.DB, username, password, email string) error {
	userRegex := `^[a-zA-Z0-9]+$`
	matched, err := regexp.MatchString(userRegex, username)
	if err != nil || !matched {
		return errors.New("用户名只能包含数字和字母")
	}
	passwordRegex := `^[a-zA-Z0-9]+$`
	matched, err = regexp.MatchString(passwordRegex, password)
	if err != nil || !matched {
		return errors.New("密码只能包含数字和字母")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var newUser = entity.User{Username: username, Password: string(hashPassword), Email: email}
	if err := us.userRepository.Create(db, &newUser); err != nil {
		return err
	}
	return nil
}

func (us *userService) CheckUserIsExist(db *gorm.DB, username, email string) error {
	if _, err := us.userRepository.FindByUsername(db, username); err == nil {
		return errors.New("Username has been registered")
	}
	if _, err := us.userRepository.FindByEmail(db, email); err == nil {
		return errors.New("Email has been registered")
	}
	return nil
}

func (us *userService) Login(db *gorm.DB, ctx context.Context, username, password string) (string, error) {
	user, err := us.userRepository.FindByUsername(db, username)
	if err != nil {
		return "", errors.New("User not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("Password is incorrect")
	}
	sessionID, err := us.sessionRepository.CreateSession(ctx, user)
	if err != nil {
		return "", errors.New("Session creation failed,logged in failed")
	}
	return sessionID, nil
}

func (us *userService) ChangeUserPassword(db *gorm.DB, userid uint, oldPassword, newPassword string) error {

	user, err := us.userRepository.FindByID(db, userid)
	if err != nil {
		return errors.New("User not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("Password is incorrect")
	}
	passwordRegex := `^[a-zA-Z0-9]+$`
	matched, err := regexp.MatchString(passwordRegex, newPassword)
	if err != nil || !matched {
		return errors.New("新密码只能包含数字和字母")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)
	err = us.userRepository.Update(db, user)
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) GetUserInfo(db *gorm.DB, userid uint) (*entity.User, error) {
	user, err := us.userRepository.FindByID(db, userid)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil

}

func (us *userService) IsLogged(ctx context.Context, sessionID string) (uint, error) {
	userid, err := us.sessionRepository.FindBySessionID(ctx, sessionID)
	if err != nil {
		return 0, errors.New("User not logged in")
	}
	return userid, nil
}
