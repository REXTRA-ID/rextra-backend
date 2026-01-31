package cache

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var ErrCacheMiss = errors.New("key not found")

// CacheService abstracts caching operations used across the app.
type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context, pattern string) error
	Close() error
}

// Service is kept as alias for backward compatibility with earlier naming.
type Service = CacheService

type memoryCache struct {
	data *sync.Map
}

// New initializes a cache service based on environment configuration.
// It prioritizes Redis if configured, otherwise falls back to in-memory cache.
func New() CacheService {
	return initCache()
}

// NewMemory creates an in-memory cache implementation.
func NewMemory() CacheService {
	return &memoryCache{
		data: &sync.Map{},
	}
}

func (c *memoryCache) Get(_ context.Context, key string) (interface{}, error) {
	value, ok := c.data.Load(key)
	if !ok {
		return nil, ErrCacheMiss
	}
	return value, nil
}

func (c *memoryCache) Set(_ context.Context, key string, value interface{}, _ time.Duration) error {
	c.data.Store(key, value)
	return nil
}

func (c *memoryCache) Delete(_ context.Context, key string) error {
	c.data.Delete(key)
	return nil
}

func (c *memoryCache) Clear(_ context.Context, pattern string) error {
	c.data.Range(func(k, _ interface{}) bool {
		if keyStr, ok := k.(string); ok && (pattern == "" || containsPattern(keyStr, pattern)) {
			c.data.Delete(k)
		}
		return true
	})
	return nil
}

func (c *memoryCache) Close() error {
	return nil
}

// containsPattern performs a simple substring match; avoid heavy globbing to keep it lightweight.
func containsPattern(key, pattern string) bool {
	return pattern == "" || strings.Contains(key, pattern)
}

func initCache() CacheService {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0 // Default to DB 0

	// Try to initialize Redis cache
	if redisHost != "" {
		redisCache, err := NewRedis(redisHost, redisPort, redisPassword, redisDB)
		if err != nil {
			log.Printf("⚠️  Failed to connect to Redis at %s:%s, falling back to memory cache: %v", redisHost, redisPort, err)
			return NewMemory()
		}
		log.Printf("✅ Redis cache initialized successfully at %s:%s", redisHost, redisPort)
		return redisCache
	}

	// If REDIS_HOST not set, use memory cache
	log.Println("ℹ️  REDIS_HOST not configured, using memory cache")
	return NewMemory()
}
