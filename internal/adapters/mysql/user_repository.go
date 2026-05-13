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
	*adapters.NotImplementedUserRepository

	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewUserRepository creates a new MySQL user repository.
func NewUserRepository(db shared.DBTX) repositories.UserRepository {
	return &UserRepository{
		NotImplementedUserRepository: adapters.NewNotImplementedUserRepository("MySQL"),
		db:                           db,
		converters:                   converters.NewConverterSet(converters.DbTypeMySQL),
	}
}

// Delete soft deletes a user from MySQL.
func (r *UserRepository) Delete(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// List retrieves users with pagination from MySQL.
func (r *UserRepository) List(
	ctx context.Context,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.ListWithPagination(ctx, r, status, limit, offset, "List")
}

// Search searches users by query in MySQL using FULLTEXT.
func (r *UserRepository) Search(
	ctx context.Context,
	query string,
	status entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	return adapters.SearchWithValidation(ctx, r, query, status, limit, "Search")
}

// SearchByTags searches users by tags in MySQL using JSON operations.
func (r *UserRepository) SearchByTags(
	ctx context.Context,
	tags []string,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.SearchByTagsWithValidation(ctx, r, tags, status, limit, offset, "SearchByTags")
}

// ChangeStatus changes user status in MySQL.
func (r *UserRepository) ChangeStatus(
	ctx context.Context,
	id entities.UserID,
	status entities.UserStatus,
) error {
	return adapters.ChangeStatusWithValidation(ctx, r, id, status, "ChangeStatus")
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
