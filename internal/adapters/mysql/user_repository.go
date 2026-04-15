package mysql

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/db/shared"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// MySQLUserRepository implements UserRepository for MySQL
// This adapts MySQL-specific types to domain interfaces.
type MySQLUserRepository struct {
	*adapters.NotImplementedUserRepository

	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewMySQLUserRepository creates a new MySQL user repository.
func NewMySQLUserRepository(db shared.DBTX) repositories.UserRepository {
	return &MySQLUserRepository{
		NotImplementedUserRepository: adapters.NewNotImplementedUserRepository("MySQL"),
		db:                           db,
		converters:                   converters.NewConverterSet(converters.DbTypeMySQL),
	}
}

// Delete soft deletes a user from MySQL.
func (r *MySQLUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// List retrieves users with pagination from MySQL.
func (r *MySQLUserRepository) List(
	ctx context.Context,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.ListWithPagination(r, ctx, status, limit, offset, "List")
}

// Search searches users by query in MySQL using FULLTEXT.
func (r *MySQLUserRepository) Search(
	ctx context.Context,
	query string,
	status entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	return adapters.SearchWithValidation(r, ctx, query, status, limit, "Search")
}

// SearchByTags searches users by tags in MySQL using JSON operations.
func (r *MySQLUserRepository) SearchByTags(
	ctx context.Context,
	tags []string,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.SearchByTagsWithValidation(r, ctx, tags, status, limit, offset, "SearchByTags")
}

// ChangeStatus changes user status in MySQL.
func (r *MySQLUserRepository) ChangeStatus(
	ctx context.Context,
	id entities.UserID,
	status entities.UserStatus,
) error {
	return adapters.ChangeStatusWithValidation(r, ctx, id, status, "ChangeStatus")
}

// Activate activates a user in MySQL.
func (r *MySQLUserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user in MySQL.
func (r *MySQLUserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user in MySQL.
func (r *MySQLUserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role in MySQL.
func (r *MySQLUserRepository) ChangeRole(
	ctx context.Context,
	id entities.UserID,
	role entities.UserRole,
) error {
	return adapters.ChangeRoleWithValidation(r, ctx, id, role, "ChangeRole")
}
