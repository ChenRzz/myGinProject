package repository

import (
	"context"
	"my_gin_project/domain/entity"
)

type SessionRepository interface {
	CreateSession(context.Context, *entity.User) (string, error)
	FindBySessionID(context.Context, string) (uint, error)
}
