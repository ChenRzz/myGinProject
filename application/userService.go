package application

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"golang.org/x/crypto/bcrypt"
	"my_gin_project/domain"
)

type UserService struct {
	messagePublisher  domain.EventProducer
	userRepository    domain.UserRepository
	sessionRepository domain.SessionRepository
}

func NewUserService(messagePublisher domain.EventProducer, useRepository domain.UserRepository, sessionRepository domain.SessionRepository) *UserService {
	return &UserService{
		messagePublisher:  messagePublisher,
		userRepository:    useRepository,
		sessionRepository: sessionRepository,
	}
}

func (us *UserService) RegisterPublish(username, password, email string) error {
	userdata := map[string]string{
		"username": username,
		"password": password,
		"email":    email,
	}
	eventBody, err := json.Marshal(userdata)
	if err != nil {
		return err
	}
	newEvent := domain.Event{
		Name: "user_register",
		Body: eventBody,
	}
	return us.messagePublisher.Publish(newEvent)
}

func (us *UserService) Register(username, password, email string) error {
	if _, err := us.userRepository.FindByUsername(username); err == nil {
		return errors.New("Username has been registered")
	}
	if _, err := us.userRepository.FindByEmail(email); err == nil {
		return errors.New("Email has been registered")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var newUser domain.User = domain.User{Username: username, Password: string(hashPassword), Email: email}
	if err := us.userRepository.Create(&newUser); err != nil {
		return err
	}
	return nil
}

func (us *UserService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := us.userRepository.FindByUsername(username)
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

func (us *UserService) ChangeUserPassword(userid uint, oldPassword, newPassword string) error {
	user, err := us.userRepository.FindByUserID(userid)
	if err != nil {
		return errors.New("User not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("Password is incorrect")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)
	err = us.userRepository.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUserInfo(userid uint) (*domain.User, error) {
	user, err := us.userRepository.FindByUserID(userid)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil

}

func (us *UserService) IsLogged(ctx context.Context, sessionID string) (uint, error) {
	userid, err := us.sessionRepository.FindBySessionID(ctx, sessionID)
	if err != nil {
		return 0, errors.New("User not logged in")
	}
	return userid, nil
}
