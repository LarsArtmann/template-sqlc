// Package entities provides domain entities and error types for the application.
// It contains core business logic objects like User and UserSession, along with
// domain-specific error types for validation, authentication, authorization, and more.
package entities

import (
	"errors"
	"fmt"
)

// Domain errors for user entity.
var (
	// ErrInvalidEmail is returned when email validation fails.
	ErrInvalidEmail        = NewValidationError("email", "must be a valid email address")
	ErrInvalidUsername     = NewValidationError("username", "must be 3-50 characters")
	ErrInvalidPasswordHash = NewValidationError("password_hash", "must be a valid hash")
	ErrInvalidFirstName    = NewValidationError("first_name", "must not be empty")
	ErrInvalidLastName     = NewValidationError("last_name", "must not be empty")
	ErrInvalidUserStatus   = NewValidationError("status", "must be a valid user status")
	ErrInvalidUserRole     = NewValidationError("role", "must be a valid user role")

	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound           = NewNotFoundError("user", "user not found")
	ErrUserAlreadyExists      = NewConflictError("user", "user already exists")
	ErrInvalidCredentials     = NewAuthenticationError("invalid credentials")
	ErrAccountSuspended       = NewAuthorizationError("account suspended")
	ErrAccountInactive        = NewAuthorizationError("account inactive")
	ErrInsufficientPrivileges = NewAuthorizationError("insufficient privileges")

	// ErrSessionNotFound is returned when a session is not found.
	ErrSessionNotFound     = NewNotFoundError("session", "session not found")
	ErrSessionExpired      = NewAuthenticationError("session expired")
	ErrInvalidSessionToken = NewAuthenticationError("invalid session token")
)

// ValidationError represents a field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationError creates a new ValidationError for the given field and message.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ResourceError represents a resource-level error with resource and message.
type ResourceError struct {
	Resource string `json:"resource"`
	Message  string `json:"message"`
	Prefix   string
}

func (e *ResourceError) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Resource, e.Prefix, e.Message)
}

// NotFoundError represents a resource not found error.
type NotFoundError struct {
	ResourceError
}

// NewNotFoundError creates a new NotFoundError for the specified resource with a message.
func NewNotFoundError(resource, message string) *NotFoundError {
	return &NotFoundError{
		ResourceError{Resource: resource, Message: message, Prefix: "not found"},
	}
}

func (e *NotFoundError) Error() string {
	return e.ResourceError.Error()
}

// ConflictError represents a resource conflict error.
type ConflictError struct {
	ResourceError
}

// NewConflictError creates a new ConflictError for the specified resource with a message.
func NewConflictError(resource, message string) *ConflictError {
	return &ConflictError{
		ResourceError{Resource: resource, Message: message, Prefix: "conflict"},
	}
}

func (e *ConflictError) Error() string {
	return e.ResourceError.Error()
}

// AuthenticationError represents an authentication failure.
type AuthenticationError struct {
	Message string `json:"message"`
}

// NewAuthenticationError creates a new AuthenticationError with the given message.
func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{
		Message: message,
	}
}

func (e *AuthenticationError) Error() string {
	return "authentication error: " + e.Message
}

// AuthorizationError represents an authorization failure.
type AuthorizationError struct {
	Message string `json:"message"`
}

// NewAuthorizationError creates a new AuthorizationError with the given message.
func NewAuthorizationError(message string) *AuthorizationError {
	return &AuthorizationError{
		Message: message,
	}
}

func (e *AuthorizationError) Error() string {
	return "authorization error: " + e.Message
}

// InternalError represents an internal server error.
type InternalError struct {
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

// NewInternalError creates a new InternalError with message and optional cause.
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

	return "internal error: " + e.Message
}

func (e *InternalError) Unwrap() error {
	return e.Cause
}

// is[E error] is a generic helper that checks if err is of type E.
func is[E error](err error) (E, bool) {
	return errors.AsType[E](err)
}

// IsValidationError checks if an error is a ValidationError.
func IsValidationError(err error) bool {
	_, ok := is[*ValidationError](err)

	return ok
}

// IsNotFoundError checks if an error is a NotFoundError.
func IsNotFoundError(err error) bool {
	_, ok := is[*NotFoundError](err)

	return ok
}

// IsConflictError checks if an error is a ConflictError.
func IsConflictError(err error) bool {
	_, ok := is[*ConflictError](err)

	return ok
}

// IsAuthenticationError checks if an error is an AuthenticationError.
func IsAuthenticationError(err error) bool {
	_, ok := is[*AuthenticationError](err)

	return ok
}

// IsUnauthorizedError checks if an error is an AuthorizationError.
func IsUnauthorizedError(err error) bool {
	_, ok := is[*AuthorizationError](err)

	return ok
}

// IsInternalError checks if an error is an InternalError.
func IsInternalError(err error) bool {
	_, ok := is[*InternalError](err)

	return ok
}

// errNotImplemented is a static error used as the base for stub not implemented errors.
var errNotImplemented = errors.New("not implemented")

// StubNotImplemented returns an error for stub implementations.
func StubNotImplemented(_, db string) error {
	return fmt.Errorf("%w: %s", errNotImplemented, db)
}
