package infrastructure

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"my_gin_project/domain"
	"strconv"
	"time"
)

type RedisSession struct {
	redisClient *redis.Client
}

func NewSessionManger() *RedisSession {
	return &RedisSession{redisClient: GetRedisClient1()}
}

func (sm *RedisSession) CreateSession(ctx context.Context, user *domain.User) (string, error) {
	sessionID := uuid.NewString() // 创建唯一的 session ID
	err := sm.redisClient.Set(ctx, sessionID, user.ID, time.Hour*24).Err()
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (sm *RedisSession) FindBySessionID(ctx context.Context, sessionID string) (uint, error) {
	userID, err := sm.redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		return 0, err
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return 0, err
	}
	uintUserID := uint(userIDInt)
	return uintUserID, nil
}
