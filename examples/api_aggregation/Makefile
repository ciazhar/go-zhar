run-dashboard:
	go run cmd/dashboard_service/main.go

run-order:
	go run cmd/order_service/main.go

run-user:
	go run cmd/user_service/main.go

run-product:
	go run cmd/product_service/main.go

run-pond:
	go run cmd/experimental_pond_service/main.go

run:
	curl -w "\nResponse time: %{time_total}s\n" -s -X GET http://localhost:3000/dashboard?userID=1

wrk-pond:
	wrk -t12 -c400 -d30s http://127.0.0.1:3004/dashboard?userID=1

wrk-dashboard:
	wrk -t12 -c400 -d30s http://127.0.0.1:3000/dashboard?userID=1


# Running 30s test @ http://127.0.0.1:3004/dashboard?userID=1
#   12 threads and 400 connections
#   Thread Stats   Avg      Stdev     Max   +/- Stdev
#     Latency   382.77ms  402.46ms   1.69s    78.38%
#     Req/Sec   166.80    289.73     1.14k    83.90%
#   43461 requests in 30.10s, 21.54MB read
#   Socket errors: connect 0, read 383, write 0, timeout 0
#   Non-2xx or 3xx responses: 362
# Requests/sec:   1443.93
# Transfer/sec:    732.68KB

# Running 30s test @ http://127.0.0.1:3000/dashboard?userID=1
#   12 threads and 400 connections
#   Thread Stats   Avg      Stdev     Max   +/- Stdev
#     Latency   368.57ms  409.01ms   1.96s    79.52%
#     Req/Sec   193.50    321.04     1.31k    81.47%
#   51504 requests in 30.10s, 25.57MB read
#   Socket errors: connect 0, read 513, write 0, timeout 0
#   Non-2xx or 3xx responses: 147
# Requests/sec:   1710.87
# Transfer/sec:    869.71KB