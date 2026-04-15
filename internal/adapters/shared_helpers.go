// Package adapters provides database adapter implementations.
// This package contains shared utilities and base implementations for database adapters.
package adapters

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/adapters/validation"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
)

// NotImplementedMethods interface for repositories with NotImplemented method.
type NotImplementedMethods interface {
	NotImplemented(method string) error
}

// NotImplementedListResult is the return type for List methods.
type NotImplementedListResult = []*entities.User

// ListWithPagination handles common pagination validation for List methods.
// Returns validation error or calls notImplementedStub if validation passes.
func ListWithPagination[T NotImplementedMethods](
	repo T,
	_ context.Context,
	_ entities.UserStatus,
	limit, offset int,
	methodName string,
) (NotImplementedListResult, error) {
	err := validation.ValidatePagination(limit, offset)
	if err != nil {
		return nil, err
	}

	return nil, repo.NotImplemented(methodName)
}

// SearchWithValidation handles common validation for Search methods.
// Returns validation error or calls notImplementedStub if validation passes.
func SearchWithValidation[T NotImplementedMethods](
	repo T,
	_ context.Context,
	query string,
	_ entities.UserStatus,
	limit int,
	methodName string,
) (NotImplementedListResult, error) {
	err := validation.ValidateSearchQuery(query, limit)
	if err != nil {
		return nil, err
	}

	return nil, repo.NotImplemented(methodName)
}

// SearchByTagsWithValidation handles common validation for SearchByTags methods.
// Returns validation error or calls notImplementedStub if validation passes.
func SearchByTagsWithValidation[T NotImplementedMethods](
	repo T,
	_ context.Context,
	tags []string,
	_ entities.UserStatus,
	limit, offset int,
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

	return nil, repo.NotImplemented(methodName)
}

// ChangeStatusWithValidation handles common validation for ChangeStatus methods.
func ChangeStatusWithValidation[T NotImplementedMethods](
	repo T,
	_ context.Context,
	_ entities.UserID,
	status entities.UserStatus,
	methodName string,
) error {
	return validation.ValidateAndExecute(status, "status", "invalid user status", func() error {
		return repo.NotImplemented(methodName)
	})
}

// ChangeRoleWithValidation handles common validation for ChangeRole methods.
func ChangeRoleWithValidation[T NotImplementedMethods](
	repo T,
	_ context.Context,
	_ entities.UserID,
	role entities.UserRole,
	methodName string,
) error {
	return validation.ValidateAndExecute(role, "role", "invalid user role", func() error {
		return repo.NotImplemented(methodName)
	})
}
