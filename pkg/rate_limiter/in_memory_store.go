package rate_limiter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type InMemoryStore struct {
	cache *cache.Cache
}

func NewInMemoryStore(defaultExpiration, cleanupInterval time.Duration) *InMemoryStore {
	return &InMemoryStore{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Get Generic Get that deserializes into target struct
func (s *InMemoryStore) Get(key string, out interface{}) (bool, error) {
	val, found := s.cache.Get(key)
	if !found {
		return false, nil
	}

	// Convert interface{} to JSON → then to target struct
	b, err := json.Marshal(val)
	if err != nil {
		return false, fmt.Errorf("failed to marshal cached value: %w", err)
	}

	if err := json.Unmarshal(b, out); err != nil {
		return false, fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return true, nil
}

// Set Set any value (struct, slice, etc)
func (s *InMemoryStore) Set(key string, value interface{}, ttl time.Duration) error {
	s.cache.Set(key, value, ttl)
	return nil
}

// Delete Optional helper if you want to remove entry manually
func (s *InMemoryStore) Delete(key string) {
	s.cache.Delete(key)
}

func (s *InMemoryStore) Type() StorageType {
	return InMemory
}
