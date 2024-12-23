package repository

import (
	"context"
	"database/sql"
	"github.com/ciazhar/go-start-small/examples/audit_trail/internal/model"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

var clickhouseDDL = `CREATE TABLE IF NOT EXISTS audit_log (
    event_type String,
    user_id String,
    timestamp DateTime64(3, 'UTC'),
    payload String DEFAULT '' -- Optional, if payload is not always set
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (event_type, user_id, timestamp);
`

type ClickHouseRepository struct {
	db *sql.DB
}

func NewClickHouseRepository(db *sql.DB) *ClickHouseRepository {
	return &ClickHouseRepository{db: db}
}

func (r *ClickHouseRepository) InsertAuditLog(ctx context.Context, log model.AuditLog) error {
	query := `
		INSERT INTO audit_log (event_type, user_id, timestamp, payload)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, log.EventType, log.UserID, log.Timestamp, log.Payload)
	if err != nil {
		logger.LogError(ctx, err, "Error inserting into ClickHouse", nil)
		return err
	}
	return nil
}
