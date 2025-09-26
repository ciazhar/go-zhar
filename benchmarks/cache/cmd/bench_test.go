package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/allegro/bigcache/v3"
	"github.com/coocood/freecache"
	"github.com/dgraph-io/ristretto"
	"github.com/redis/go-redis/v9"
)

// --- Redis client ---
var (
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx = context.Background()
)

// --- Simple Go map ---
type MapCache struct {
	data map[string]string
}

func NewMapCache() *MapCache {
	return &MapCache{data: make(map[string]string)}
}
func (c *MapCache) Set(key, value string) {
	c.data[key] = value
}
func (c *MapCache) Get(key string) (string, bool) {
	val, ok := c.data[key]
	return val, ok
}

// --- Ristretto ---
var ristrettoCache *ristretto.Cache

// --- BigCache ---
var bigCache *bigcache.BigCache

// --- FreeCache ---
var freeCache *freecache.Cache

// --- Init caches ---
func init() {
	// Ristretto
	rc, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // 1GB
		BufferItems: 64,      // recommended value
	})
	ristrettoCache = rc

	// BigCache
	bc, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * 60)) // 10 min
	bigCache = bc

	// FreeCache (512MB)
	freeCache = freecache.NewCache(512 * 1024 * 1024)
}

// --- Benchmarks ---

// Redis
func BenchmarkRedisSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		redisClient.Set(ctx, fmt.Sprintf("key:%d", i), "value", 0)
	}
}

func BenchmarkRedisGet(b *testing.B) {
	redisClient.Set(ctx, "benchkey", "value", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		redisClient.Get(ctx, "benchkey")
	}
}

// Go map
func BenchmarkMapSet(b *testing.B) {
	cache := NewMapCache()
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("key:%d", i), "value")
	}
}

func BenchmarkMapGet(b *testing.B) {
	cache := NewMapCache()
	cache.Set("benchkey", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("benchkey")
	}
}

// Ristretto
func BenchmarkRistrettoSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ristrettoCache.Set(fmt.Sprintf("key:%d", i), "value", 1)
	}
}

func BenchmarkRistrettoGet(b *testing.B) {
	ristrettoCache.Set("benchkey", "value", 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ristrettoCache.Get("benchkey")
	}
}

// BigCache
func BenchmarkBigCacheSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bigCache.Set(fmt.Sprintf("key:%d", i), []byte("value"))
	}
}

func BenchmarkBigCacheGet(b *testing.B) {
	bigCache.Set("benchkey", []byte("value"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bigCache.Get("benchkey")
	}
}

// FreeCache
func BenchmarkFreeCacheSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		freeCache.Set([]byte(fmt.Sprintf("key:%d", i)), []byte("value"), 0)
	}
}

func BenchmarkFreeCacheGet(b *testing.B) {
	freeCache.Set([]byte("benchkey"), []byte("value"), 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		freeCache.Get([]byte("benchkey"))
	}
}
