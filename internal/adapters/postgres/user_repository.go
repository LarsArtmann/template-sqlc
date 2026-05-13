// Package postgres provides PostgreSQL-specific database adapter implementations.
package postgres

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
	"github.com/jackc/pgx/v5"
)

// UserRepository implements UserRepository for PostgreSQL
// This adapts PostgreSQL-specific types to domain interfaces.
type UserRepository struct {
	*adapters.NotImplementedUserRepository

	pool       pgx.Tx
	converters *converters.ConverterSet
}

// NewUserRepository creates a new PostgreSQL user repository.
func NewUserRepository(pool pgx.Tx) repositories.UserRepository {
	return &UserRepository{
		NotImplementedUserRepository: adapters.NewNotImplementedUserRepository("PostgreSQL"),
		pool:                         pool,
		converters:                   converters.NewConverterSet(converters.DbTypePostgres),
	}
}

// Delete soft deletes a user from PostgreSQL.
func (r *UserRepository) Delete(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// List retrieves users with pagination from PostgreSQL.
func (r *UserRepository) List(
	ctx context.Context,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.ListWithPagination(ctx, r, status, limit, offset, "List")
}

// Search searches users by query in PostgreSQL using FTS.
func (r *UserRepository) Search(
	ctx context.Context,
	query string,
	status entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	return adapters.SearchWithValidation(ctx, r, query, status, limit, "Search")
}

// SearchByTags searches users by tags in PostgreSQL using GIN index.
func (r *UserRepository) SearchByTags(
	ctx context.Context,
	tags []string,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.SearchByTagsWithValidation(ctx, r, tags, status, limit, offset, "SearchByTags")
}

// ChangeStatus changes user status in PostgreSQL.
func (r *UserRepository) ChangeStatus(
	ctx context.Context,
	id entities.UserID,
	status entities.UserStatus,
) error {
	return adapters.ChangeStatusWithValidation(ctx, r, id, status, "ChangeStatus")
}

// Activate activates a user in PostgreSQL.
func (r *UserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user in PostgreSQL.
func (r *UserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user in PostgreSQL.
func (r *UserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role in PostgreSQL.
func (r *UserRepository) ChangeRole(
	ctx context.Context,
	id entities.UserID,
	role entities.UserRole,
) error {
	return adapters.ChangeRoleWithValidation(ctx, r, id, role, "ChangeRole")
}
