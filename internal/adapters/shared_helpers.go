// Package adapters provides database adapter implementations.
// This package contains shared utilities and base implementations for database adapters.
package adapters

import (
	"context"
	"fmt"

	"github.com/LarsArtmann/template-sqlc/internal/adapters/validation"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
)

// NotImplementedMethods interface for repositories with NotImplemented method.
type NotImplementedMethods interface {
	NotImplemented(method string) error
}

// notImplementedError creates a standardized not implemented error with proper wrapping.
func notImplementedError[T NotImplementedMethods](repo T, methodName string) error {
	return fmt.Errorf(
		"method %s not implemented: %w",
		methodName,
		repo.NotImplemented(methodName),
	)
}

// NotImplementedListResult is the return type for List methods.
type NotImplementedListResult = []*entities.User

// ListWithPagination handles common pagination validation for List methods.
// Returns validation error or calls notImplementedStub if validation passes.
func ListWithPagination[T NotImplementedMethods](
	_ context.Context,
	repo T,
	_ entities.UserStatus,
	limit, offset int,
	methodName string,
) (NotImplementedListResult, error) {
	err := validation.ValidatePagination(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("limit=%v: %w", limit, err)
	}

	return nil, fmt.Errorf(
		"method %s not implemented for limit=%v: %w",
		methodName,
		limit,
		repo.NotImplemented(methodName),
	)
}

// SearchWithValidation handles common validation for Search methods.
// Returns validation error or calls notImplementedStub if validation passes.
func SearchWithValidation[T NotImplementedMethods](
	_ context.Context,
	repo T,
	query string,
	_ entities.UserStatus,
	limit int,
	methodName string,
) (NotImplementedListResult, error) {
	err := validation.ValidateSearchQuery(query, limit)
	if err != nil {
		return nil, fmt.Errorf("query=%v limit=%v: %w", query, limit, err)
	}

	return nil, fmt.Errorf(
		"method %s not implemented for query=%v limit=%v: %w",
		methodName,
		query,
		limit,
		repo.NotImplemented(methodName),
	)
}

// SearchByTagsWithValidation handles common validation for SearchByTags methods.
// Returns validation error or calls notImplementedStub if validation passes.
func SearchByTagsWithValidation[T NotImplementedMethods](
	_ context.Context,
	repo T,
	tags []string,
	_ entities.UserStatus,
	_, offset int,
	methodName string,
) (NotImplementedListResult, error) {
	err := validation.ValidateTags(tags)
	if err != nil {
		return nil, err
	}

	if offset < 0 {
		offset = 0
	}

	_ = offset

	return nil, fmt.Errorf(
		"method %s not implemented: %w",
		methodName,
		repo.NotImplemented(methodName),
	)
}

// ChangeStatusWithValidation handles common validation for ChangeStatus methods.
func ChangeStatusWithValidation[T NotImplementedMethods](
	_ context.Context,
	repo T,
	_ entities.UserID,
	status entities.UserStatus,
	methodName string,
) error {
	return validation.ValidateAndExecute(status, "status", "invalid user status", func() error {
		return notImplementedError(repo, methodName)
	})
}

// ChangeRoleWithValidation handles common validation for ChangeRole methods.
func ChangeRoleWithValidation[T NotImplementedMethods](
	_ context.Context,
	repo T,
	_ entities.UserID,
	role entities.UserRole,
	methodName string,
) error {
	return validation.ValidateAndExecute(role, "role", "invalid user role", func() error {
		return notImplementedError(repo, methodName)
	})
}

// ListUsers handles the common List implementation for user repositories.
func ListUsers[T NotImplementedMethods](
	ctx context.Context,
	repo T,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return ListWithPagination(ctx, repo, status, limit, offset, "List")
}

// SearchUsers handles the common Search implementation for user repositories.
func SearchUsers[T NotImplementedMethods](
	ctx context.Context,
	repo T,
	query string,
	status entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	return SearchWithValidation(ctx, repo, query, status, limit, "Search")
}

// SearchUsersByTags handles the common SearchByTags implementation for user repositories.
func SearchUsersByTags[T NotImplementedMethods](
	ctx context.Context,
	repo T,
	tags []string,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return SearchByTagsWithValidation(ctx, repo, tags, status, limit, offset, "SearchByTags")
}

// ChangeUserStatus handles the common ChangeStatus implementation for user repositories.
func ChangeUserStatus[T NotImplementedMethods](
	ctx context.Context,
	repo T,
	id entities.UserID,
	status entities.UserStatus,
) error {
	return ChangeStatusWithValidation(ctx, repo, id, status, "ChangeStatus")
}

// BaseUserRepository provides common implementations for user repository methods.
// Embed this struct in database-specific repositories to avoid duplicating method implementations.
type BaseUserRepository struct {
	*NotImplementedUserRepository
}

// NewBaseUserRepository creates a new BaseUserRepository with the given database name.
func NewBaseUserRepository(dbName string) *BaseUserRepository {
	return &BaseUserRepository{
		NotImplementedUserRepository: NewNotImplementedUserRepository(dbName),
	}
}

// List retrieves users with pagination.
func (r *BaseUserRepository) List(
	ctx context.Context,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return ListUsers(ctx, r, status, limit, offset)
}

// Search searches users by query.
func (r *BaseUserRepository) Search(
	ctx context.Context,
	query string,
	status entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	return SearchUsers(ctx, r, query, status, limit)
}

// SearchByTags searches users by tags.
func (r *BaseUserRepository) SearchByTags(
	ctx context.Context,
	tags []string,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	return SearchUsersByTags(ctx, r, tags, status, limit, offset)
}

// ChangeStatus changes user status.
func (r *BaseUserRepository) ChangeStatus(
	ctx context.Context,
	id entities.UserID,
	status entities.UserStatus,
) error {
	return ChangeUserStatus(ctx, r, id, status)
}

// Activate activates a user.
func (r *BaseUserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user.
func (r *BaseUserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user.
func (r *BaseUserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role.
func (r *BaseUserRepository) ChangeRole(
	ctx context.Context,
	id entities.UserID,
	role entities.UserRole,
) error {
	return ChangeRoleWithValidation(ctx, r, id, role, "ChangeRole")
}

// Delete soft deletes a user.
func (r *BaseUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}
