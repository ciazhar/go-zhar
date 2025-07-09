package rate_limiter

import (
	"sync"
	"time"
)

type fixedWindow struct {
	requests map[string]int
	resetAt  map[string]time.Time
	limit    int
	window   time.Duration
	mu       sync.Mutex
}

func (fw *fixedWindow) Allow(key string) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	if expire, ok := fw.resetAt[key]; !ok || now.After(expire) {
		fw.requests[key] = 0
		fw.resetAt[key] = now.Add(fw.window)
	}

	if fw.requests[key] < fw.limit {
		fw.requests[key]++
		return true
	}
	return false
}
