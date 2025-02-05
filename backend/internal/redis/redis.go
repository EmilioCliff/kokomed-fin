package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/redis/go-redis/v9"
)

var _ services.CacheService = (*redisCache)(nil)

type redisCache struct {
	client *redis.Client
}

func NewCacheClient(addr, password string, db int) services.CacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &redisCache{client: rdb}
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *redisCache) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil 
	} else if err != nil {
		return false, err
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return false, err
	}
	return true, nil 
}

func (r *redisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

func (r *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}