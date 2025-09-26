# 🧪 Go Cache Benchmark

This project benchmarks popular in-memory caching libraries in Go and compares them with **Redis** as an external cache.
It helps evaluate performance trade-offs between local in-process caches and a network-based distributed cache.

---

## 📦 Libraries Benchmarked

* [BigCache](https://github.com/allegro/bigcache)
* [FreeCache](https://github.com/coocood/freecache)
* [Ristretto](https://github.com/dgraph-io/ristretto)
* [Golang Built-in `map` + `sync.RWMutex`]
* [Redis](https://redis.io/) (via Docker)

---

## 🛠️ Setup

### 1. Clone Repo

```bash
git clone https://github.com/ciazhar/go-zhar.git
cd benchmarks/cache
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Start Redis (via Docker Compose)

```bash
make run-docker
```

This will start a Redis container accessible at `localhost:6379`.

---

## ▶️ Run Benchmarks

```bash
make bench
```

Example output (numbers will vary by machine):

```
go test -bench=. -benchmem cmd/bench_test.go
goos: darwin
goarch: arm64
cpu: Apple M1
BenchmarkRedisSet-8       	   11706	     88148 ns/op	     273 B/op	       8 allocs/op
BenchmarkRedisGet-8       	   13410	     90978 ns/op	     192 B/op	       6 allocs/op
BenchmarkMapSet-8         	 5208207	       270.8 ns/op	     152 B/op	       2 allocs/op
BenchmarkMapGet-8         	198241964	         6.193 ns/op	       0 B/op	       0 allocs/op
BenchmarkRistrettoSet-8   	 4045575	       385.0 ns/op	     269 B/op	       4 allocs/op
BenchmarkRistrettoGet-8   	30681480	        39.02 ns/op	       6 B/op	       0 allocs/op
BenchmarkBigCacheSet-8    	 4244593	       349.5 ns/op	      71 B/op	       2 allocs/op
BenchmarkBigCacheGet-8    	33033036	        31.92 ns/op	       8 B/op	       1 allocs/op
BenchmarkFreeCacheSet-8   	 4750627	       346.9 ns/op	      66 B/op	       2 allocs/op
BenchmarkFreeCacheGet-8   	12978537	        91.93 ns/op	       5 B/op	       1 allocs/op
PASS
ok  	command-line-arguments	17.898s
```

---

## 📊 What’s Measured?

* **Set performance** → How fast can the cache store values.
* **Get performance** → How fast can the cache retrieve values.
* **Memory allocations** → How efficient is memory usage.

---

## 🏗️ Docker Compose (Redis)

```yaml
version: "3.9"

services:
  redis:
    image: redis:7-alpine
    container_name: redis-benchmark
    ports:
      - "6379:6379"
    restart: unless-stopped
```

---

## 🔑 Key Takeaways

* **Local caches (BigCache, FreeCache, Ristretto, map+mutex)** are much faster than Redis for single-node apps.
* **Redis** is slower due to network I/O, but it enables **shared state** across multiple services.
* Choice depends on use case:

    * Single service → use local in-memory caches.
    * Distributed system → use Redis.

---

## 🚀 Next Improvements

* Add **TTL eviction benchmark** for local caches.
* Compare performance under **concurrent load (high goroutines)**.
* Add **real workload simulations** (e.g., 80% reads, 20% writes).