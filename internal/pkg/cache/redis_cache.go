package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(host string, port string, password string, db int) (CacheService, error) {
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		MaxRetries:   3,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Redis at %s: %w", addr, err)
	}

	return &redisCache{
		client: client,
		ttl:    24 * time.Hour,
	}, nil
}

func (c *redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	if c.client == nil {
		return nil, errors.New("redis client not initialized")
	}

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	// Unmarshal JSON to ensure data consistency
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, fmt.Errorf("redis json unmarshal error: %w", err)
	}
	return result, nil
}

// Set stores a value in Redis cache with TTL support.
func (c *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if c.client == nil {
		return errors.New("redis client not initialized")
	}

	if expiration == 0 {
		expiration = c.ttl
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Delete removes a key from Redis cache.
func (c *redisCache) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return errors.New("redis client not initialized")
	}

	return c.client.Del(ctx, key).Err()
}

// Clear removes all keys matching a pattern from Redis cache.
// Uses SCAN with pattern matching for safe deletion.
func (c *redisCache) Clear(ctx context.Context, pattern string) error {
	if c.client == nil {
		return errors.New("redis client not initialized")
	}

	var cursor uint64
	var keysDeleted int64

	for {
		// Use SCAN to find keys matching pattern
		keys, newCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("redis scan error: %w", err)
		}

		// Delete found keys in batches
		if len(keys) > 0 {
			deleted, err := c.client.Del(ctx, keys...).Result()
			if err != nil {
				return fmt.Errorf("redis delete error: %w", err)
			}
			keysDeleted += deleted
		}

		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

func (c *redisCache) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

// HealthCheck verifies Redis connection is healthy
func (c *redisCache) HealthCheck(ctx context.Context) error {
	if c.client == nil {
		return errors.New("redis client not initialized")
	}

	// Respect caller's context timeout if available
	return c.client.Ping(ctx).Err()
}

// SetWithTTL sets a value with custom TTL (convenience method)
func (c *redisCache) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.Set(ctx, key, value, ttl)
}

// GetMany retrieves multiple values from Redis (batch operation)
func (c *redisCache) GetMany(ctx context.Context, keys []string) (map[string]interface{}, error) {
	if c.client == nil {
		return nil, errors.New("redis client not initialized")
	}

	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	vals, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis mget error: %w", err)
	}

	result := make(map[string]interface{})
	for i, val := range vals {
		if val != nil {
			if strVal, ok := val.(string); ok {
				var parsed interface{}
				if err := json.Unmarshal([]byte(strVal), &parsed); err == nil {
					result[keys[i]] = parsed
				} else {
					result[keys[i]] = strVal
				}
			} else {
				result[keys[i]] = val
			}
		}
	}

	return result, nil
}

// DeleteMany removes multiple keys from Redis (batch operation)
func (c *redisCache) DeleteMany(ctx context.Context, keys []string) error {
	if c.client == nil {
		return errors.New("redis client not initialized")
	}

	if len(keys) == 0 {
		return nil
	}

	return c.client.Del(ctx, keys...).Err()
}

// GetWithTTL retrieves a value and its remaining TTL
func (c *redisCache) GetWithTTL(ctx context.Context, key string) (interface{}, time.Duration, error) {
	if c.client == nil {
		return nil, 0, errors.New("redis client not initialized")
	}

	val, err := c.Get(ctx, key)
	if err != nil {
		return nil, 0, err
	}

	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("redis ttl error: %w", err)
	}

	return val, ttl, nil
}

// Exists checks if a key exists in Redis
func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	if c.client == nil {
		return false, errors.New("redis client not initialized")
	}

	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists error: %w", err)
	}

	return result > 0, nil
}

// Keys returns all keys matching a pattern using SCAN (safe for large databases)
func (c *redisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	if c.client == nil {
		return nil, errors.New("redis client not initialized")
	}

	var cursor uint64
	var allKeys []string

	for {
		keys, newCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("redis scan error: %w", err)
		}
		allKeys = append(allKeys, keys...)
		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	return allKeys, nil
}

// Increment increments a numeric value in Redis with default TTL
func (c *redisCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	if c.client == nil {
		return 0, errors.New("redis client not initialized")
	}

	val, err := c.client.IncrBy(ctx, key, delta).Result()
	if err != nil {
		return 0, fmt.Errorf("redis increment error: %w", err)
	}

	if err := c.client.Expire(ctx, key, c.ttl).Err(); err != nil {
		fmt.Printf("warning: failed to set TTL on increment key %s: %v\n", key, err)
	}

	return val, nil
}

// Stats returns cache statistics
func (c *redisCache) Stats(ctx context.Context) (map[string]string, error) {
	if c.client == nil {
		return nil, errors.New("redis client not initialized")
	}

	info, err := c.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("redis info error: %w", err)
	}

	stats := make(map[string]string)
	for _, line := range strings.Split(info, "\n") {
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				stats[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	return stats, nil
}
