package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/examples/audit_trail/internal/model"
	"github.com/ciazhar/go-start-small/examples/audit_trail/internal/repository"
	"github.com/ciazhar/go-start-small/pkg/kafka/consumer"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	postgresDB, err := sql.Open("postgres", "host=localhost port=5432 user=your_user password=your_password dbname=your_db sslmode=disable")
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to connect to PostgreSQL", nil)
	}

	clickhouseDB, err := sql.Open("clickhouse", "tcp://localhost:9000?username=default&password=&database=default")
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to connect to ClickHouse", nil)
	}

	postgresRepo := repository.NewPostgresRepository(postgresDB)
	clickhouseRepo := repository.NewClickHouseRepository(clickhouseDB)

	brokers := "localhost:8097,localhost:8098,localhost:8099"
	consumer.StartConsumerGroup(context.Background(), brokers, map[string]consumer.ConsumerGroup{
		model.AuditLogsTopic: {
			Topic:   model.AuditLogsTopic,
			GroupID: "audit-trail",
			Process: func(msg *sarama.ConsumerMessage) error {
				var auditLog model.AuditLog
				err = json.Unmarshal(msg.Value, &auditLog)
				if err != nil {
					logger.LogError(context.Background(), err, "Error unmarshalling message", nil)
					return err
				}

				err = postgresRepo.InsertAuditLog(context.Background(), auditLog)
				if err != nil {
					logger.LogError(context.Background(), err, "Error inserting into PostgreSQL", nil)
					return err
				}

				err = clickhouseRepo.InsertAuditLog(context.Background(), auditLog)
				if err != nil {
					logger.LogError(context.Background(), err, "Error inserting into ClickHouse", nil)
					return err
				}

				return nil
			},
		},
	}, &wg, "", true)

	wg.Wait()
}
