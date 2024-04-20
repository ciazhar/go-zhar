# Apache Kafka

> Apache Kafka adalah open-source distributed event streaming platform yang terkenal karena arsitekturnya yang
> memiliki throughput tinggi dan toleransi kesalahan. Operasinya didasarkan pada sistem commit log terdistribusi, di
> mana data diatur ke dalam topik dan dibagi menjadi partisi, memungkinkan producer untuk memublikasikan pesan dan
> consumer untuk subscribe secara real-time. Kafka memiliki scalability, durability, dan pengiriman pesan dengan
> low-latency, menjadikannya ideal untuk menangani volume data besar di berbagai aplikasi. Ekosistemnya mencakup
> connector untuk integrasi yang seamless dengan berbagai source dan destination.

## Use Case
- [Simple Consumer Producer](https://github.com/ciazhar/go-zhar/tree/master/examples/kafka/simple-consumer-producer)
- [HTTP to Kafka Producer (Sync & Async)](https://github.com/ciazhar/go-zhar/tree/master/examples/kafka/sync-async-producer)
- [Consumer Group](https://github.com/ciazhar/go-zhar/tree/master/examples/kafka/consumer-group)
- [Custom Consumer Group (Group By Key, Windowing, Kafka Producer)](https://github.com/ciazhar/go-zhar/tree/master/examples/kafka/custom-consumer-group)
- SASL Scram Client
- Transactional Producer
- Open Telemetry
- Kafka Stream