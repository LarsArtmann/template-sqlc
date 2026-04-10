// Package validation provides validation utilities for database adapters
package validation

import (
	"github.com/LarsArtmann/template-sqlc/pkg/errors"
)

const (
	maxPaginationLimit = 1000
	maxSearchLimit     = 100
)

// ValidatePagination validates pagination parameters.
func ValidatePagination(limit, offset int) error {
	if !isValidPaginationLimit(limit) {
		return errors.NewValidationError("limit", "must be between 1 and 1000")
	}

	if offset < 0 {
		return errors.NewValidationError("offset", "must be non-negative")
	}

	return nil
}

func isValidPaginationLimit(limit int) bool {
	return limit > 0 && limit <= maxPaginationLimit
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
	if query == "" {
		return errors.NewValidationError("query", "cannot be empty")
	}

	if len(query) > 500 {
		return errors.NewValidationError("query", "cannot exceed 500 characters")
	}

	if !isValidSearchLimit(limit) {
		return errors.NewValidationError("limit", "must be between 1 and 100")
	}

	return nil
}

func isValidSearchLimit(limit int) bool {
	return limit > 0 && limit <= maxSearchLimit
}

// Validate validates an entity using its IsValid method.
func Validate[T interface{ IsValid() bool }](entity T, fieldName, invalidMessage string) error {
	if !entity.IsValid() {
		return errors.NewValidationError(fieldName, invalidMessage)
	}

	return nil
}

// ValidateAndExecute validates an entity and executes the update function if validation passes.
func ValidateAndExecute[T interface{ IsValid() bool }](
	entity T,
	fieldName, invalidMessage string,
	updateFn func() error,
) error {
	err := Validate(entity, fieldName, invalidMessage)
	if err != nil {
		return err
	}

	return updateFn()
}
