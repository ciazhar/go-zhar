package rate_limiter

import (
	"sync"
	"time"
)

type tokenBucket struct {
	capacity   int
	tokens     float64
	lastRefill time.Time
	refillRate float64 // tokens per second
	mu         sync.Mutex
}

func (b *tokenBucket) Allow(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * b.refillRate
	if b.tokens > float64(b.capacity) {
		b.tokens = float64(b.capacity)
	}
	b.lastRefill = now

	if b.tokens >= 1 {
		b.tokens -= 1
		return true
	}
	return false
}
