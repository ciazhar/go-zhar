package postgres

import (
	"context"
	"database/sql"
	"errors"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Close()
	Config() *pgxpool.Config
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Reset()
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Stat() *pgxpool.Stat
}

// --- MockRow for QueryRow ---
type MockRow struct {
	Values []any
	Err    error
}

func (m MockRow) Scan(dest ...any) error {
	if m.Err != nil {
		return m.Err
	}
	if len(dest) != len(m.Values) {
		return errors.New("mismatched number of destinations")
	}

	for i := range dest {
		if dest[i] == nil {
			continue
		}

		dv := reflect.ValueOf(dest[i])
		if dv.Kind() != reflect.Ptr {
			return errors.New("destination not a pointer")
		}

		if m.Values[i] == nil {
			// handle sql.Null*
			switch d := dest[i].(type) {
			case *sql.NullString:
				*d = sql.NullString{Valid: false}
			case *sql.NullInt64:
				*d = sql.NullInt64{Valid: false}
			case *sql.NullBool:
				*d = sql.NullBool{Valid: false}
			case *sql.NullTime:
				*d = sql.NullTime{Valid: false}
			default:
				// leave as nil
			}
			continue
		}

		v := reflect.ValueOf(m.Values[i])
		if v.Type().AssignableTo(dv.Elem().Type()) {
			dv.Elem().Set(v)
		} else {
			return errors.New("type mismatch in Scan")
		}
	}
	return nil
}

// --- MockRows for Query (multi-row) ---
type MockRows struct {
	Rows   []MockRow
	Idx    int
	ErrVal error
	Closed bool
}

func (m *MockRows) Next() bool {
	if m.ErrVal != nil || m.Closed {
		return false
	}
	if m.Idx < len(m.Rows) {
		m.Idx++
		return true
	}
	return false
}

func (m *MockRows) Scan(dest ...any) error {
	if m.ErrVal != nil {
		return m.ErrVal
	}
	if m.Idx == 0 || m.Idx > len(m.Rows) {
		return errors.New("no current row")
	}
	return m.Rows[m.Idx-1].Scan(dest...)
}

func (m *MockRows) Close() { m.Closed = true }

func (m *MockRows) Err() error                                   { return m.ErrVal }
func (m *MockRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (m *MockRows) Values() ([]any, error) {
	if m.Idx == 0 || m.Idx > len(m.Rows) {
		return nil, errors.New("no current row")
	}
	return m.Rows[m.Idx-1].Values, nil
}
func (m *MockRows) RawValues() [][]byte { return nil }
func (m *MockRows) Conn() *pgx.Conn     { return nil }

// Ensure MockRows implements pgx.Rows
var _ pgx.Rows = (*MockRows)(nil)
