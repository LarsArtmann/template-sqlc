// Package mysql provides MySQL-specific database adapter implementations.
package mysql

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/db/shared"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// UserRepository implements UserRepository for MySQL
// This adapts MySQL-specific types to domain interfaces.
type UserRepository struct {
	*adapters.SharedUserRepository

	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewUserRepository creates a new MySQL user repository.
func NewUserRepository(db shared.DBTX) repositories.UserRepository {
	return &UserRepository{
		SharedUserRepository: adapters.NewSharedUserRepository("MySQL"),
		db:                   db,
		converters:           converters.NewConverterSet(converters.DbTypeMySQL),
	}
}

// Delete soft deletes a user from MySQL.
func (r *UserRepository) Delete(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Activate activates a user in MySQL.
func (r *UserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user in MySQL.
func (r *UserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user in MySQL.
func (r *UserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role in MySQL.
func (r *UserRepository) ChangeRole(
	ctx context.Context,
	id entities.UserID,
	role entities.UserRole,
) error {
	return adapters.ChangeRoleWithValidation(ctx, r, id, role, "ChangeRole")
}
