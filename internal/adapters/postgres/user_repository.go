package postgres

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/adapters"
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
	ctx context.Context,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.ListWithPagination(r, ctx, status, limit, offset, "List")
}

// Search searches users by query in PostgreSQL using FTS.
func (r *PostgresUserRepository) Search(
	ctx context.Context,
	query string,
	status entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	return adapters.SearchWithValidation(r, ctx, query, status, limit, "Search")
}

// SearchByTags searches users by tags in PostgreSQL using GIN index.
func (r *PostgresUserRepository) SearchByTags(
	ctx context.Context,
	tags []string,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.SearchByTagsWithValidation(r, ctx, tags, status, limit, offset, "SearchByTags")
}

// ChangeStatus changes user status in PostgreSQL.
func (r *PostgresUserRepository) ChangeStatus(
	ctx context.Context,
	id entities.UserID,
	status entities.UserStatus,
) error {
	return adapters.ChangeStatusWithValidation(r, ctx, id, status, "ChangeStatus")
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
	ctx context.Context,
	id entities.UserID,
	role entities.UserRole,
) error {
	return adapters.ChangeRoleWithValidation(r, ctx, id, role, "ChangeRole")
}
