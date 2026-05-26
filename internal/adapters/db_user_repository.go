package adapters

import (
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/db/shared"
)

// DBUserRepository contains common fields for MySQL and SQLite user repositories.
// Both MySQL and SQLite use the same database/sql-based implementation.
type DBUserRepository struct {
	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewDBUserRepository creates a new DBUserRepository with the given database and converter type.
func NewDBUserRepository(db shared.DBTX, dbType string) *DBUserRepository {
	return &DBUserRepository{
		db:         db,
		converters: converters.NewConverterSet(dbType),
	}
}
