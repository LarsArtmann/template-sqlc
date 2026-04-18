// Package shared provides common types for database/sql-based packages.
// This is used by both mysql and sqlite packages to avoid code duplication.
package shared

import (
	"context"
	"database/sql"
)

// DBTX is the common interface for database/sql operations.
// Both MySQL and SQLite use this interface.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// BaseQueries contains the common fields for Queries struct.
type BaseQueries struct {
	DB DBTX
	tx *sql.Tx
}

// NewBase creates a new BaseQueries instance.
func NewBase(db DBTX) *BaseQueries {
	return &BaseQueries{DB: db}
}

// WithTx returns a new BaseQueries with the given transaction.
func (q *BaseQueries) WithTx(tx *sql.Tx) *BaseQueries {
	return &BaseQueries{
		DB: tx,
		tx: tx,
	}
}
