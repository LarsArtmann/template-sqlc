// Package postgres provides PostgreSQL-specific database adapter implementations.
package postgres

import (
	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
	"github.com/jackc/pgx/v5"
)

// UserRepository implements UserRepository for PostgreSQL
// This adapts PostgreSQL-specific types to domain interfaces.
type UserRepository struct {
	*adapters.BaseUserRepository

	pool       pgx.Tx
	converters *converters.ConverterSet
}

// NewUserRepository creates a new PostgreSQL user repository.
func NewUserRepository(pool pgx.Tx) repositories.UserRepository {
	return &UserRepository{
		BaseUserRepository: adapters.NewBaseUserRepository("PostgreSQL"),
		pool:               pool,
		converters:         converters.NewConverterSet(converters.DbTypePostgres),
	}
}
