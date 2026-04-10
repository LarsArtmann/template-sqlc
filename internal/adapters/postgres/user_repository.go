package postgres

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/validation"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// PostgresUserRepository implements UserRepository for PostgreSQL
// This adapts PostgreSQL-specific types to domain interfaces.
type PostgresUserRepository struct {
	*adapters.NotImplementedUserRepository
	pool       any
	converters any
}

// NewPostgresUserRepository creates a new PostgreSQL user repository.
func NewPostgresUserRepository(pool any) repositories.UserRepository {
	return &PostgresUserRepository{
		NotImplementedUserRepository: adapters.NewNotImplementedUserRepository("PostgreSQL"),
		pool:                         pool,
		converters:                   nil,
	}
}

// Delete soft deletes a user from PostgreSQL.
func (r *PostgresUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// List retrieves users with pagination from PostgreSQL.
func (r *PostgresUserRepository) List(
	_ context.Context,
	_ entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	err := validation.ValidatePagination(limit, offset)
	if err != nil {
		return nil, err
	}

	return nil, r.NotImplemented("List")
}

// Search searches users by query in PostgreSQL using FTS.
func (r *PostgresUserRepository) Search(
	_ context.Context,
	query string,
	_ entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	err := validation.ValidateSearchQuery(query, limit)
	if err != nil {
		return nil, err
	}

	return nil, r.NotImplemented("Search")
}

// SearchByTags searches users by tags in PostgreSQL using GIN index.
func (r *PostgresUserRepository) SearchByTags(
	_ context.Context,
	tags []string,
	_ entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	err := validation.ValidateTags(tags)
	if err != nil {
		return nil, err
	}

	if offset < 0 {
		offset = 0
	}
	_ = offset

	return nil, r.NotImplemented("SearchByTags")
}

// ChangeStatus changes user status in PostgreSQL.
func (r *PostgresUserRepository) ChangeStatus(
	_ context.Context,
	_ entities.UserID,
	status entities.UserStatus,
) error {
	return validation.ValidateAndExecute(status, "status", "invalid user status", func() error {
		return r.NotImplemented("ChangeStatus")
	})
}

// Activate activates a user in PostgreSQL.
func (r *PostgresUserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user in PostgreSQL.
func (r *PostgresUserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user in PostgreSQL.
func (r *PostgresUserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role in PostgreSQL.
func (r *PostgresUserRepository) ChangeRole(
	_ context.Context,
	_ entities.UserID,
	role entities.UserRole,
) error {
	return validation.ValidateAndExecute(role, "role", "invalid user role", func() error {
		return r.NotImplemented("ChangeRole")
	})
}
