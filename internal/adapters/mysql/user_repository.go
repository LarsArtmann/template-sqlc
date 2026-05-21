// Package mysql provides MySQL-specific database adapter implementations.
package mysql

import (
	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/db/shared"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// UserRepository implements UserRepository for MySQL
// This adapts MySQL-specific types to domain interfaces.
type UserRepository struct {
	*adapters.BaseUserRepository

	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewUserRepository creates a new MySQL user repository.
func NewUserRepository(db shared.DBTX) repositories.UserRepository {
	return &UserRepository{
		BaseUserRepository: adapters.NewBaseUserRepository("MySQL"),
		db:                 db,
		converters:         converters.NewConverterSet(converters.DbTypeMySQL),
	}
}
