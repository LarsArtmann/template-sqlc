package entities

import (
	"fmt"
)

// Domain errors for user entity
var (
	// Validation errors
	ErrInvalidEmail        = NewValidationError("email", "must be a valid email address")
	ErrInvalidUsername     = NewValidationError("username", "must be 3-50 characters")
	ErrInvalidPasswordHash = NewValidationError("password_hash", "must be a valid hash")
	ErrInvalidFirstName    = NewValidationError("first_name", "must not be empty")
	ErrInvalidLastName     = NewValidationError("last_name", "must not be empty")
	ErrInvalidUserStatus   = NewValidationError("status", "must be a valid user status")
	ErrInvalidUserRole     = NewValidationError("role", "must be a valid user role")

	// Business logic errors
	ErrUserNotFound           = NewNotFoundError("user", "user not found")
	ErrUserAlreadyExists      = NewConflictError("user", "user already exists")
	ErrInvalidCredentials     = NewAuthenticationError("invalid credentials")
	ErrAccountSuspended       = NewAuthorizationError("account suspended")
	ErrAccountInactive        = NewAuthorizationError("account inactive")
	ErrInsufficientPrivileges = NewAuthorizationError("insufficient privileges")

	// Session errors
	ErrSessionNotFound     = NewNotFoundError("session", "session not found")
	ErrSessionExpired      = NewAuthenticationError("session expired")
	ErrInvalidSessionToken = NewAuthenticationError("invalid session token")
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string `json:"resource"`
	Message  string `json:"message"`
}

func NewNotFoundError(resource, message string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		Message:  message,
	}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Resource, e.Message)
}

// ConflictError represents a resource conflict error
type ConflictError struct {
	Resource string `json:"resource"`
	Message  string `json:"message"`
}

func NewConflictError(resource, message string) *ConflictError {
	return &ConflictError{
		Resource: resource,
		Message:  message,
	}
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s conflict: %s", e.Resource, e.Message)
}

// AuthenticationError represents an authentication failure
type AuthenticationError struct {
	Message string `json:"message"`
}

func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{
		Message: message,
	}
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error: %s", e.Message)
}

// AuthorizationError represents an authorization failure
type AuthorizationError struct {
	Message string `json:"message"`
}

func NewAuthorizationError(message string) *AuthorizationError {
	return &AuthorizationError{
		Message: message,
	}
}

func (e *AuthorizationError) Error() string {
	return fmt.Sprintf("authorization error: %s", e.Message)
}

// InternalError represents an internal server error
type InternalError struct {
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func NewInternalError(message string, cause error) *InternalError {
	return &InternalError{
		Message: message,
		Cause:   cause,
	}
}

func (e *InternalError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("internal error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("internal error: %s", e.Message)
}

func (e *InternalError) Unwrap() error {
	return e.Cause
}
