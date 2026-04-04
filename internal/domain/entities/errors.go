package entities

import (
	"errors"
	"fmt"
)

// Domain errors for user entity.
var (
	// Validation errors.
	ErrInvalidEmail        = NewValidationError("email", "must be a valid email address")
	ErrInvalidUsername     = NewValidationError("username", "must be 3-50 characters")
	ErrInvalidPasswordHash = NewValidationError("password_hash", "must be a valid hash")
	ErrInvalidFirstName    = NewValidationError("first_name", "must not be empty")
	ErrInvalidLastName     = NewValidationError("last_name", "must not be empty")
	ErrInvalidUserStatus   = NewValidationError("status", "must be a valid user status")
	ErrInvalidUserRole     = NewValidationError("role", "must be a valid user role")

	// Business logic errors.
	ErrUserNotFound           = NewNotFoundError("user", "user not found")
	ErrUserAlreadyExists      = NewConflictError("user", "user already exists")
	ErrInvalidCredentials     = NewAuthenticationError("invalid credentials")
	ErrAccountSuspended       = NewAuthorizationError("account suspended")
	ErrAccountInactive        = NewAuthorizationError("account inactive")
	ErrInsufficientPrivileges = NewAuthorizationError("insufficient privileges")

	// Session errors.
	ErrSessionNotFound     = NewNotFoundError("session", "session not found")
	ErrSessionExpired      = NewAuthenticationError("session expired")
	ErrInvalidSessionToken = NewAuthenticationError("invalid session token")
)

// ValidationError represents a field validation error.
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

// ValidateSearchQuery validates search query parameters
// Returns nil if valid, otherwise returns a ValidationError.
func ValidateSearchQuery(query string, limit int) error {
	if len(query) == 0 {
		return NewValidationError("query", "cannot be empty")
	}

	if len(query) > 500 {
		return NewValidationError("query", "cannot exceed 500 characters")
	}

	if limit <= 0 || limit > 100 {
		return NewValidationError("limit", "must be between 1 and 100")
	}

	return nil
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

func NewConflictError(resource, message string) *ConflictError {
	return &ConflictError{
		ResourceError{Resource: resource, Message: message, Prefix: "conflict"},
	}
}

func (e *ConflictError) Error() string {
	return e.ResourceError.Error()
}

func newResourceError(resource, message, prefix string) ResourceError {
	return ResourceError{Resource: resource, Message: message, Prefix: prefix}
}

// AuthenticationError represents an authentication failure.
type AuthenticationError struct {
	Message string `json:"message"`
}

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

// is[T any] is a generic helper that checks if err is of type T.
func is[T any](err error, target *T) bool {
	if err == nil {
		return false
	}

	return errors.As(err, target)
}

// IsValidationError checks if an error is a ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError

	return is(err, &ve)
}

// IsNotFoundError checks if an error is a NotFoundError.
func IsNotFoundError(err error) bool {
	var ne *NotFoundError

	return is(err, &ne)
}

// IsConflictError checks if an error is a ConflictError.
func IsConflictError(err error) bool {
	var ce *ConflictError

	return is(err, &ce)
}

// IsAuthenticationError checks if an error is an AuthenticationError.
func IsAuthenticationError(err error) bool {
	var ae *AuthenticationError

	return is(err, &ae)
}

// IsUnauthorizedError checks if an error is an AuthorizationError.
func IsUnauthorizedError(err error) bool {
	var aze *AuthorizationError

	return is(err, &aze)
}

// IsInternalError checks if an error is an InternalError.
func IsInternalError(err error) bool {
	var ie *InternalError

	return is(err, &ie)
}

// StubNotImplemented returns an error for stub implementations.
func StubNotImplemented(method, db string) error {
	return fmt.Errorf("implement me: use actual sqlc generated code for %s", db)
}
