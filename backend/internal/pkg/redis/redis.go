package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/aryasatyawa/bayarin/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisClient struct {
	*redis.Client
}

// NewRedisClient creates new Redis client
func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info().Msg("✅ Redis connected successfully")

	return &RedisClient{Client: client}, nil
}

// Close closes Redis connection
func (r *RedisClient) Close() error {
	log.Info().Msg("Closing Redis connection...")
	return r.Client.Close()
}

// Health checks Redis health
func (r *RedisClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := r.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

// SetWithExpiry sets key with expiration
func (r *RedisClient) SetWithExpiry(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Set(ctx, key, value, expiration).Err()
}

// GetString gets string value by key
func (r *RedisClient) GetString(ctx context.Context, key string) (string, error) {
	return r.Get(ctx, key).Result()
}

// Delete deletes key
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.Del(ctx, keys...).Err()
}

// Exists checks if key exists
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
