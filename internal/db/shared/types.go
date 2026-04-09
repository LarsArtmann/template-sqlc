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
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// BaseQueries contains the common fields for Queries struct.
type BaseQueries struct {
	db DBTX
	tx *sql.Tx
}

// NewBase creates a new BaseQueries instance.
func NewBase(db DBTX) *BaseQueries {
	return &BaseQueries{db: db}
}

// WithTx returns a new BaseQueries with the given transaction.
func (q *BaseQueries) WithTx(tx *sql.Tx) *BaseQueries {
	return &BaseQueries{
		db: tx,
		tx: tx,
	}
}
