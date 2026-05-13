// Package sqlite provides SQLite-specific database adapter implementations.
package sqlite

import (
	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/db/shared"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// SessionRepository implements SessionRepository for SQLite.
type SessionRepository struct {
	*adapters.NotImplementedSessionRepository

	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewSessionRepository creates a new SQLite session repository.
func NewSessionRepository(db shared.DBTX) repositories.SessionRepository {
	return &SessionRepository{
		NotImplementedSessionRepository: adapters.NewNotImplementedSessionRepository("SQLite"),
		db:                              db,
		converters:                      converters.NewConverterSet(converters.DbTypeSQLite),
	}
}
