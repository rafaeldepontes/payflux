package service

import (
	"context"
	"time"

	"github.com/rafaeldepontes/ledger/internal/cache"
	cr "github.com/rafaeldepontes/ledger/pkg/cache"
	"github.com/redis/go-redis/v9"
)

const (
	_                = iota
	DefaultCacheTime = iota + 47
)

type cacheSvc struct {
	db *redis.Client
}

// NewService returns a new instance of the cache service.
func NewService() cache.Cache[string, string] {
	return cacheSvc{
		db: cr.GetCache(),
	}
}

// Add adds something to cache for 48 hours.
func (r cacheSvc) Add(key string, value string) {
	r.AddWithTTL(time.Duration(DefaultCacheTime*time.Hour), key, value)
}

// AddWithTTL adds something to cache with a specified TTL.
// If multiple values are provided, only the first one is used for simple KV storage.
func (r cacheSvc) AddWithTTL(t time.Duration, key string, value ...string) {
	if len(value) == 0 {
		return
	}
	ctx := context.Background()
	r.db.Set(ctx, key, value[0], t)
}

// Clear clears the current database.
func (r cacheSvc) Clear() {
	ctx := context.Background()
	r.db.FlushDB(ctx)
}

// FullClear clears all databases.
func (r cacheSvc) FullClear() {
	ctx := context.Background()
	r.db.FlushAll(ctx)
}

// Get gets the value if any and also returns a boolean to check if it exists.
func (r cacheSvc) Get(key string) (string, bool) {
	ctx := context.Background()
	val, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

// Remove removes the value from cache.
func (r cacheSvc) Remove(key string) {
	ctx := context.Background()
	r.db.Del(ctx, key)
}

// Set updates the cache value and also refreshes the TTL.
func (r cacheSvc) Set(key string, value string) {
	r.Add(key, value)
}
