package sqlite

import (
	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/db/shared"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// UserRepository implements UserRepository for SQLite
// This adapts SQLite-specific types to domain interfaces.
type UserRepository struct {
	*adapters.BaseUserRepository

	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewUserRepository creates a new SQLite user repository.
func NewUserRepository(db shared.DBTX) repositories.UserRepository {
	return &UserRepository{
		BaseUserRepository: adapters.NewBaseUserRepository("SQLite"),
		db:                 db,
		converters:         converters.NewConverterSet(converters.DbTypeSQLite),
	}
}
