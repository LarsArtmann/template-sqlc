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
	*adapters.DBUserRepository
}

// NewUserRepository creates a new SQLite user repository.
func NewUserRepository(db shared.DBTX) repositories.UserRepository {
	return &UserRepository{
		BaseUserRepository: adapters.NewBaseUserRepository("SQLite"),
		DBUserRepository:   adapters.NewDBUserRepository(db, converters.DbTypeSQLite),
	}
}
