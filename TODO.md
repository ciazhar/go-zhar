# TODO

- Semua code gk boleh pake interface karena bikin lemot
- config di implement embed biar kalo build bisa ambil (gk bagus kalo multi stage)
- implement context
- semua struct return pointer
- return langsung, gk perlu return err
- better to return pointer for smaller cost, only dont use it if u dont want the var change after pass 

- Barcode
- Clickhouse Explore data type
- Clickhouse 300 Juta Data
- Fiber File Server
- Grpc Clean Arch + Swagger
- Htmx
- Kafka : SASL Scram Client
- Kafka : Transactional Producer
- Kafka : Open Telemetry
- Kafka : Testcontainers
- Line Bot
- Telegram Bot
- WA Bot
- MongoDB : Explore data type
- MongoDB : Paging & Sorting
- MongoDB : Testcontainers
- MongoDB : Aggregate
- MongoDB : Change StreamQ
- MongoDB : Index
- MongoDB : Import From Xlsx / Csv
- Postgres : Explore data type
- Postgres : Testcontainers
- Postgres : Trx
- Prometheus
- RabbitMQ : Delayed Message
- Redis : Explore data type
- Redis : TTL Listener
- Scrap
- Sentry
- Export Xls / Csv
- OpenTelemetry & Jaeger
- CLI
- Gmap
- Firebase Push Notification
- SocketIO



### Table

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go

```

</td><td>

```go

```

</td></tr>
<tr><td>

</td></tr>
</tbody></table>



## List gk mudeng

- Verify Interface Compliance
- Receivers and Interfaces (sebenere ngerti kalo ambil data pake pointer ada cost e cuman setelah di bench gk ngaruhu)