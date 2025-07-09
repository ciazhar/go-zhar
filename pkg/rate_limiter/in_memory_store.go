package rate_limiter

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type InMemoryStore struct {
	cache *cache.Cache
}

func NewInMemoryStore(defaultExpiration, cleanupInterval time.Duration) *InMemoryStore {
	return &InMemoryStore{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (s *InMemoryStore) Get(key string) (interface{}, bool) {
	return s.cache.Get(key)
}

func (s *InMemoryStore) Set(key string, value interface{}, ttl time.Duration) {
	s.cache.Set(key, value, ttl)
}

func (s *InMemoryStore) Delete(key string) {
	s.cache.Delete(key)
}
