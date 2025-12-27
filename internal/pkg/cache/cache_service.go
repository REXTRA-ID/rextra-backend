package cache

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

var ErrCacheMiss = errors.New("key not found")

type Service interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context, pattern string) error
}

type memoryCache struct {
	data *sync.Map
}

// NewMemory creates an in-memory cache implementation.
func NewMemory() Service {
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

// containsPattern performs a simple substring match; avoid heavy globbing to keep it lightweight.
func containsPattern(key, pattern string) bool {
	return pattern == "" || strings.Contains(key, pattern)
}
