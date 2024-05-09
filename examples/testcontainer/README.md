## Testcontainers

Testcontainers adalah library Go yang memungkinkan Anda untuk dengan mudah membuat dan mengelola kontainer Docker dalam
unit test. Ini memungkinkan testing yang konsisten dan dapat diulang, serta isolasi environment testing.

Berikut adalah beberapa contoh unit test tech stack yang di uji menggunakan Testcontainers :
- Database
  - [Clickhouse](https://github.com/ciazhar/go-zhar/blob/master/examples/clickhouse/crud-testcontainers/internal/repository/clickhouse_repository_test.go)
  - [Redis](https://github.com/ciazhar/go-zhar/blob/master/examples/redis/crud-testcontainers/internal/repository/redis_repository_test.go)
- Message Broker
  - [Rabbitmq](https://github.com/ciazhar/go-zhar/blob/master/examples/rabbitmq/publish-consume-testcontainers/internal/service/basic_service_test.go)