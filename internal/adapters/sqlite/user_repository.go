package sqlite

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/db/shared"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// SQLiteUserRepository implements UserRepository for SQLite
// This adapts SQLite-specific types to domain interfaces.
type SQLiteUserRepository struct {
	*adapters.NotImplementedUserRepository

	db         shared.DBTX
	converters *converters.ConverterSet
}

// NewSQLiteUserRepository creates a new SQLite user repository.
func NewSQLiteUserRepository(db shared.DBTX) repositories.UserRepository {
	return &SQLiteUserRepository{
		NotImplementedUserRepository: adapters.NewNotImplementedUserRepository("SQLite"),
		db:                           db,
		converters:                   converters.NewConverterSet(converters.DbTypeSQLite),
	}
}

// List retrieves users with pagination from SQLite.
func (r *SQLiteUserRepository) List(
	ctx context.Context,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.ListWithPagination(r, ctx, status, limit, offset, "List")
}

// Search searches users by query in SQLite.
func (r *SQLiteUserRepository) Search(
	ctx context.Context,
	query string,
	status entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	return adapters.SearchWithValidation(r, ctx, query, status, limit, "Search")
}

// SearchByTags searches users by tags in SQLite.
func (r *SQLiteUserRepository) SearchByTags(
	ctx context.Context,
	tags []string,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return adapters.SearchByTagsWithValidation(r, ctx, tags, status, limit, offset, "SearchByTags")
}

// ChangeStatus changes user status in SQLite.
func (r *SQLiteUserRepository) ChangeStatus(
	ctx context.Context,
	id entities.UserID,
	status entities.UserStatus,
) error {
	return adapters.ChangeStatusWithValidation(r, ctx, id, status, "ChangeStatus")
}

// Activate activates a user in SQLite.
func (r *SQLiteUserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user in SQLite.
func (r *SQLiteUserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user in SQLite.
func (r *SQLiteUserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role in SQLite.
func (r *SQLiteUserRepository) ChangeRole(
	ctx context.Context,
	id entities.UserID,
	role entities.UserRole,
) error {
	return adapters.ChangeRoleWithValidation(r, ctx, id, role, "ChangeRole")
}
