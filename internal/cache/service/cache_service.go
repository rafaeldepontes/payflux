package service

import (
	"time"

	"github.com/rafaeldepontes/goplo/internal/cache"
)

type redis struct {
}

// TODO: impl Redis logic and connection.
func NewService() cache.Cache[string, string] {
	return redis{}
}

// Add implements [cache.Cache].
func (r redis) Add(key string, value string) {
	panic("unimplemented")
}

// AddWithTLS implements [cache.Cache].
func (r redis) AddWithTLS(key string, value string, time *time.Duration) {
	panic("unimplemented")
}

// Clear implements [cache.Cache].
func (r redis) Clear() {
	panic("unimplemented")
}

// FullClear implements [cache.Cache].
func (r redis) FullClear() {
	panic("unimplemented")
}

// Get implements [cache.Cache].
func (r redis) Get(key string) (string, bool) {
	panic("unimplemented")
}

// Remove implements [cache.Cache].
func (r redis) Remove(key string) {
	panic("unimplemented")
}

// Set implements [cache.Cache].
func (r redis) Set(key string, value string) {
	panic("unimplemented")
}
