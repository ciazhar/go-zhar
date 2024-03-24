# Benchmark In Memory Database

## Run Docker Compose
```bash
cd ../../deployments/dragonfly && docker-compose up
```
```bash
cd ../../deployments/redis && docker-compose up
```
## Run Benchmark
```bash
    go test -bench=. -benchmem bench_test.go
```