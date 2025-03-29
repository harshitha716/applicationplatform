package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CacheClient interface {
	Get(ctx context.Context, key string, valuePtr interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)

	FormatKey(prefix string, id interface{}) (string, error)
	Client() *redis.Client
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) CacheClient {
	return &RedisCache{
		client: client,
	}
}

func (rc *RedisCache) Client() *redis.Client {
	return rc.client
}

func (rc *RedisCache) FormatKey(prefix string, id interface{}) (string, error) {
	if id == nil {
		return prefix, nil
	}

	var idStr string
	switch v := id.(type) {
	case string:
		idStr = v
	case int, int64, uint, uint64, float64:
		idStr = fmt.Sprintf("%v", v)
	default:
		idBytes, err := json.Marshal(id)
		if err != nil {
			return "", fmt.Errorf("failed to marshal key ID: %w", err)
		}
		idStr = string(idBytes)
	}

	return fmt.Sprintf("%s:%s", prefix, idStr), nil
}

func (rc *RedisCache) Get(ctx context.Context, key string, valuePtr interface{}) error {
	logger := getLoggerFromContext(ctx)
	// Prefix the key with organization ID
	_, _, userOrganizations := apicontext.GetAuthFromContext(ctx)
	if len(userOrganizations) > 0 {
		key = fmt.Sprintf("%s:%s", userOrganizations[0], key)
	}

	value, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		logger.Error("failed to get key from cache",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("cache get error: %w", err)
	}

	if err := json.Unmarshal([]byte(value), valuePtr); err != nil {
		logger.Error("failed to unmarshal cache value",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("unmarshal error: %w", err)
	}

	return nil
}

func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	logger := getLoggerFromContext(ctx)
	_, _, userOrganizations := apicontext.GetAuthFromContext(ctx)
	if len(userOrganizations) > 0 {
		key = fmt.Sprintf("%s:%s", userOrganizations[0], key)
	}

	valueJSON, err := json.Marshal(value)
	if err != nil {
		logger.Error("failed to marshal value for cache",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := rc.client.Set(ctx, key, valueJSON, expiry).Err(); err != nil {
		logger.Error("failed to set key in cache",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("cache set error: %w", err)
	}

	return nil
}

func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	logger := getLoggerFromContext(ctx)
	_, _, userOrganizations := apicontext.GetAuthFromContext(ctx)
	if len(userOrganizations) > 0 {
		key = fmt.Sprintf("%s:%s", userOrganizations[0], key)
	}

	if err := rc.client.Del(ctx, key).Err(); err != nil {
		logger.Error("failed to delete key from cache",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("cache delete error: %w", err)
	}

	return nil
}

func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	logger := getLoggerFromContext(ctx)
	_, _, userOrganizations := apicontext.GetAuthFromContext(ctx)
	if len(userOrganizations) > 0 {
		key = fmt.Sprintf("%s:%s", userOrganizations[0], key)
	}

	exists, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		logger.Error("failed to check key existence in cache",
			zap.String("key", key),
			zap.Error(err))
		return false, fmt.Errorf("cache exists error: %w", err)
	}

	return exists > 0, nil
}

func (rc *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	logger := getLoggerFromContext(ctx)
	_, _, userOrganizations := apicontext.GetAuthFromContext(ctx)
	if len(userOrganizations) > 0 {
		key = fmt.Sprintf("%s:%s", userOrganizations[0], key)
	}

	result, err := rc.client.Incr(ctx, key).Result()
	if err != nil {
		logger.Error("failed to increment key in cache",
			zap.String("key", key),
			zap.Error(err))
		return 0, fmt.Errorf("cache increment error: %w", err)
	}

	return result, nil
}

func (rc *RedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	logger := getLoggerFromContext(ctx)
	_, _, userOrganizations := apicontext.GetAuthFromContext(ctx)
	if len(userOrganizations) > 0 {
		key = fmt.Sprintf("%s:%s", userOrganizations[0], key)
	}

	result, err := rc.client.Decr(ctx, key).Result()
	if err != nil {
		logger.Error("failed to decrement key in cache",
			zap.String("key", key),
			zap.Error(err))
		return 0, fmt.Errorf("cache decrement error: %w", err)
	}

	return result, nil
}

func getLoggerFromContext(ctx context.Context) *zap.Logger {
	logger := apicontext.GetLoggerFromCtx(ctx)
	return logger
}
