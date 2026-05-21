package sqlite

import (
	"context"
	"fmt"

	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/db/shared"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// UserRepository implements UserRepository for SQLite
// This adapts SQLite-specific types to domain interfaces.
type UserRepository struct {
	*adapters.SharedUserRepository

	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewUserRepository creates a new SQLite user repository.
func NewUserRepository(db shared.DBTX) repositories.UserRepository {
	return &UserRepository{
		SharedUserRepository: adapters.NewSharedUserRepository("SQLite"),
		db:                   db,
		converters:           converters.NewConverterSet(converters.DbTypeSQLite),
	}
}

// Delete soft deletes a user from SQLite.
func (r *UserRepository) Delete(ctx context.Context, id entities.UserID) error {
	err := r.ChangeStatus(ctx, id, entities.UserStatusInactive)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	return nil
}

// Activate activates a user in SQLite.
func (r *UserRepository) Activate(ctx context.Context, id entities.UserID) error {
	err := r.ChangeStatus(ctx, id, entities.UserStatusActive)
	if err != nil {
		return fmt.Errorf("Activate: %w", err)
	}

	return nil
}

// Deactivate deactivates a user in SQLite.
func (r *UserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	err := r.ChangeStatus(ctx, id, entities.UserStatusInactive)
	if err != nil {
		return fmt.Errorf("Deactivate: %w", err)
	}

	return nil
}

// Suspend suspends a user in SQLite.
func (r *UserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	err := r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
	if err != nil {
		return fmt.Errorf("Suspend: %w", err)
	}

	return nil
}

// ChangeRole changes user role in SQLite.
func (r *UserRepository) ChangeRole(
	ctx context.Context,
	id entities.UserID,
	role entities.UserRole,
) error {
	return adapters.ChangeRoleWithValidation(ctx, r, id, role, "ChangeRole")
}
