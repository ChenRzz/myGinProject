package domain

import (
	"context"
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:100;unique"`
	Password  string `gorm:"size:255"`
	Email     string `gorm:"size:100;unique"`
	CreatedAt time.Time
}

type UserRepository interface {
	FindByUsername(username string) (*User, error)
	FindByUserID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(user *User) error
	Update(user *User) error
}
type SessionRepository interface {
	CreateSession(context.Context, *User) (string, error)
	FindBySessionID(context.Context, string) (uint, error)
}
