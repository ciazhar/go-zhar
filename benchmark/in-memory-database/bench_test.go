package main

import (
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/redis"
	"testing"
)

var d *redis.Redis
var r *redis.Redis

func init() {

	l := logger.Init(
		logger.Config{ConsoleLoggingEnabled: true})
	d = redis.Init("127.0.0.1", 6379, "", l)
	r = redis.Init("127.0.0.1", 6377, "", l)
}

func BenchmarkDragonflySet(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Set(fmt.Sprintf("key-%d", i), "value1", 0)
	}
}
func BenchmarkRedisSet(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Set(fmt.Sprintf("key-%d", i), "value1", 0)
	}
}

func BenchmarkDragonflyGet(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Get(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkRedisGet(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Get(fmt.Sprintf("key-%d", i))
	}
}
