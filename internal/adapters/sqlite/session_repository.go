package sqlite

import (
	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/db/shared"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// SQLiteSessionRepository implements SessionRepository for SQLite.
type SQLiteSessionRepository struct {
	*adapters.NotImplementedSessionRepository

	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewSQLiteSessionRepository creates a new SQLite session repository.
func NewSQLiteSessionRepository(db shared.DBTX) repositories.SessionRepository {
	return &SQLiteSessionRepository{
		NotImplementedSessionRepository: adapters.NewNotImplementedSessionRepository("SQLite"),
		db:                              db,
		converters:                      converters.NewConverterSet(converters.DbTypeSQLite),
	}
}
