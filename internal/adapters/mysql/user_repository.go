package mysql

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/adapters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/validation"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// MySQLUserRepository implements UserRepository for MySQL
// This adapts MySQL-specific types to domain interfaces.
type MySQLUserRepository struct {
	*adapters.NotImplementedUserRepository
	db         any
	converters any
}

// NewMySQLUserRepository creates a new MySQL user repository.
func NewMySQLUserRepository(db any) repositories.UserRepository {
	return &MySQLUserRepository{
		NotImplementedUserRepository: adapters.NewNotImplementedUserRepository("MySQL"),
		db:         db,
		converters: nil,
	}
}

// Delete soft deletes a user from MySQL.
func (r *MySQLUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// List retrieves users with pagination from MySQL.
func (r *MySQLUserRepository) List(
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

// Search searches users by query in MySQL using FULLTEXT.
func (r *MySQLUserRepository) Search(
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

// SearchByTags searches users by tags in MySQL using JSON operations.
func (r *MySQLUserRepository) SearchByTags(
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

// ChangeStatus changes user status in MySQL.
func (r *MySQLUserRepository) ChangeStatus(
	_ context.Context,
	_ entities.UserID,
	status entities.UserStatus,
) error {
	return validation.ValidateAndExecute(status, "status", "invalid user status", func() error {
		return r.NotImplemented("ChangeStatus")
	})
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
	_ context.Context,
	_ entities.UserID,
	role entities.UserRole,
) error {
	return validation.ValidateAndExecute(role, "role", "invalid user role", func() error {
		return r.NotImplemented("ChangeRole")
	})
}
