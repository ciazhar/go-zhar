// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	ID        int32            `json:"id"`
	Name      pgtype.Text      `json:"name"`
	Price     pgtype.Numeric   `json:"price"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}
