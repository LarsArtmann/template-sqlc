// Package validation provides validation utilities for database adapters
package validation

import (
	"github.com/LarsArtmann/template-sqlc/pkg/errors"
)

// ValidatePagination validates pagination parameters.
func ValidatePagination(limit, offset int) error {
	if limit <= 0 || limit > 1000 {
		return errors.NewValidationError("limit", "must be between 1 and 1000")
	}

	if offset < 0 {
		return errors.NewValidationError("offset", "must be non-negative")
	}

	return nil
}

// ValidateTags validates tags parameter.
func ValidateTags(tags []string) error {
	if len(tags) == 0 {
		return errors.NewValidationError("tags", "cannot be empty")
	}

	if len(tags) > 10 {
		return errors.NewValidationError("tags", "cannot exceed 10 tags")
	}

	return nil
}

// ValidateSearchQuery validates search query and limit.
func ValidateSearchQuery(query string, limit int) error {
	if len(query) == 0 {
		return errors.NewValidationError("query", "cannot be empty")
	}

	if len(query) > 500 {
		return errors.NewValidationError("query", "cannot exceed 500 characters")
	}

	if limit <= 0 || limit > 100 {
		return errors.NewValidationError("limit", "must be between 1 and 100")
	}

	return nil
}

// Validatable is an interface for entities that can validate themselves.
type Validatable interface {
	IsValid() bool
}

// Validate validates an entity using its IsValid method.
func Validate[T interface{ IsValid() bool }](entity T, fieldName, invalidMessage string) error {
	if !entity.IsValid() {
		return errors.NewValidationError(fieldName, invalidMessage)
	}

	return nil
}
