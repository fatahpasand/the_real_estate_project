// internal/adapters/repositories/redis/cache_repository.go
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *redisRepository {
	return &redisRepository{
		client: client,
	}
}

func (r *redisRepository) SetUserSession(ctx context.Context, userID uint, token string, expiration time.Duration) error {
	key := fmt.Sprintf("user_session:%d", userID)
	return r.client.Set(ctx, key, token, expiration).Err()
}

func (r *redisRepository) SetOTP(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.client.Set(ctx, "otp:"+key, value, expiration).Err()
}

func (r *redisRepository) GetOTP(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, "otp:"+key).Result()
}

func (r *redisRepository) DeleteOTP(ctx context.Context, key string) error {
	return r.client.Del(ctx, "otp:"+key).Err()
}
