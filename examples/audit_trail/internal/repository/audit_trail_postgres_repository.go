package repository

import (
	"context"
	"database/sql"
	"github.com/ciazhar/go-start-small/examples/audit_trail/internal/model"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

var postgresDDL = `CREATE TABLE audit_log (
    event_type TEXT NOT NULL,
    user_id TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    payload JSONB DEFAULT '{}' -- Optional, stored as JSONB for efficient querying
);
`

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) InsertAuditLog(ctx context.Context, log model.AuditLog) error {
	query := `
		INSERT INTO audit_log (event_type, user_id, timestamp, payload)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, log.EventType, log.UserID, log.Timestamp, log.Payload)
	if err != nil {
		logger.LogError(ctx, err, "Error inserting into PostgreSQLn", nil)
		return err
	}
	return nil
}
