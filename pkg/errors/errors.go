// Package errors provides standardized application error types with HTTP status codes.
package errors

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
)

// ErrorCode represents standardized error codes.
type ErrorCode string

const (
	// ErrCodeValidationFailed indicates a validation error occurred.
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	// ErrCodeInvalidInput indicates invalid input was provided.
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	// ErrCodeMissingField indicates a required field is missing.
	ErrCodeMissingField ErrorCode = "MISSING_FIELD"
	// ErrCodeInvalidFormat indicates the format of a value is invalid.
	ErrCodeInvalidFormat ErrorCode = "INVALID_FORMAT"
	// ErrCodeConstraintFailed indicates a database constraint was violated.
	ErrCodeConstraintFailed ErrorCode = "CONSTRAINT_FAILED"

	// ErrCodeUnauthorized indicates authentication is required.
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	// ErrCodeInvalidCredentials indicates the provided credentials are invalid.
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS" //nolint:gosec // This is an error code, not a credential
	// ErrCodeTokenExpired indicates an authentication token has expired.
	ErrCodeTokenExpired ErrorCode = "TOKEN_EXPIRED"
	// ErrCodeTokenInvalid indicates an authentication token is invalid.
	ErrCodeTokenInvalid ErrorCode = "TOKEN_INVALID"

	// ErrCodeForbidden indicates the user lacks permission for this action.
	ErrCodeForbidden ErrorCode = "FORBIDDEN"
	// ErrCodeInsufficientPrivileges indicates insufficient privileges for this operation.
	ErrCodeInsufficientPrivileges ErrorCode = "INSUFFICIENT_PRIVILEGES"
	// ErrCodeAccountSuspended indicates the account is suspended.
	ErrCodeAccountSuspended ErrorCode = "ACCOUNT_SUSPENDED"
	// ErrCodeAccountInactive indicates the account is inactive.
	ErrCodeAccountInactive ErrorCode = "ACCOUNT_INACTIVE"

	// ErrCodeNotFound indicates the requested resource was not found.
	ErrCodeNotFound ErrorCode = "NOT_FOUND"
	// ErrCodeResourceNotFound indicates a specific resource type was not found.
	ErrCodeResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	// ErrCodeAlreadyExists indicates the resource already exists.
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	// ErrCodeResourceConflict indicates a conflict with the current state of the resource.
	ErrCodeResourceConflict ErrorCode = "RESOURCE_CONFLICT"

	// ErrCodeInternal indicates an internal server error occurred.
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
	// ErrCodeDatabase indicates a database error occurred.
	ErrCodeDatabase ErrorCode = "DATABASE_ERROR"
	// ErrCodeNetwork indicates a network error occurred.
	ErrCodeNetwork ErrorCode = "NETWORK_ERROR"
	// ErrCodeTimeout indicates an operation timed out.
	ErrCodeTimeout ErrorCode = "TIMEOUT"
	// ErrCodeUnavailable indicates the service is unavailable.
	ErrCodeUnavailable ErrorCode = "UNAVAILABLE"

	// ErrCodeBusinessLogic indicates a business logic rule was violated.
	ErrCodeBusinessLogic ErrorCode = "BUSINESS_LOGIC_ERROR"
	// ErrCodeInvalidState indicates the resource is in an invalid state for this operation.
	ErrCodeInvalidState ErrorCode = "INVALID_STATE"
	// ErrCodePermissionDenied indicates permission was denied for this operation.
	ErrCodePermissionDenied ErrorCode = "PERMISSION_DENIED"
)

// AppError represents a structured application error.
type AppError struct {
	Code       ErrorCode      `json:"code"`
	Message    string         `json:"message"`
	Details    map[string]any `json:"details,omitempty"`
	HTTPStatus int            `json:"-"`
	Cause      error          `json:"-"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("%s: %s (details: %+v)", e.Code, e.Message, e.Details)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError creates a new application error.
func NewAppError(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// NewAppErrorWithCause creates a new application error with underlying cause.
func NewAppErrorWithCause(code ErrorCode, message string, httpStatus int, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Cause:      cause,
	}
}

// NewAppErrorWithDetails creates a new application error with details.
func NewAppErrorWithDetails(
	code ErrorCode,
	message string,
	httpStatus int,
	details map[string]any,
) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Details:    details,
	}
}

// newBadRequestError creates an AppError with http.StatusBadRequest.
func newBadRequestError(code ErrorCode, message string, kvPairs ...string) *AppError {
	details := make(map[string]any)
	for i := 0; i+1 < len(kvPairs); i += 2 {
		details[kvPairs[i]] = kvPairs[i+1]
	}

	return NewAppErrorWithDetails(code, message, http.StatusBadRequest, details)
}

// NewValidationError creates a validation error for a specific field.
func NewValidationError(field, message string) *AppError {
	return newBadRequestError(
		ErrCodeValidationFailed,
		"Validation failed",
		"field",
		field,
		"message",
		message,
	)
}

// NewInvalidInputError creates an invalid input error.
func NewInvalidInputError(message string) *AppError {
	return NewAppError(ErrCodeInvalidInput, message, http.StatusBadRequest)
}

// newFieldError creates an AppError with a field detail.
func newFieldError(code ErrorCode, message, field string) *AppError {
	return NewAppErrorWithDetails(
		code,
		message,
		http.StatusBadRequest,
		map[string]any{"field": field},
	)
}

// NewMissingFieldError creates a missing field error.
func NewMissingFieldError(field string) *AppError {
	return newFieldError(ErrCodeMissingField, "Required field is missing", field)
}

// NewInvalidFormatError creates an invalid format error.
func NewInvalidFormatError(field, format string) *AppError {
	return newBadRequestError(
		ErrCodeInvalidFormat,
		"Invalid format",
		"field",
		field,
		"format",
		format,
	)
}

// NewUnauthorizedError creates an unauthorized error.
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

// NewInvalidCredentialsError creates an invalid credentials error.
func NewInvalidCredentialsError() *AppError {
	return NewAppError(ErrCodeInvalidCredentials, "Invalid credentials", http.StatusUnauthorized)
}

// NewTokenExpiredError creates a token expired error.
func NewTokenExpiredError() *AppError {
	return NewAppError(ErrCodeTokenExpired, "Token has expired", http.StatusUnauthorized)
}

// NewTokenInvalidError creates a token invalid error.
func NewTokenInvalidError() *AppError {
	return NewAppError(ErrCodeTokenInvalid, "Invalid token", http.StatusUnauthorized)
}

// NewForbiddenError creates a forbidden error.
func NewForbiddenError(message string) *AppError {
	return NewAppError(ErrCodeForbidden, message, http.StatusForbidden)
}

// NewInsufficientPrivilegesError creates an insufficient privileges error.
func NewInsufficientPrivilegesError() *AppError {
	return NewAppError(
		ErrCodeInsufficientPrivileges,
		"Insufficient privileges",
		http.StatusForbidden,
	)
}

// NewAccountSuspendedError creates an account suspended error.
func NewAccountSuspendedError() *AppError {
	return NewAppError(ErrCodeAccountSuspended, "Account suspended", http.StatusForbidden)
}

// NewAccountInactiveError creates an account inactive error.
func NewAccountInactiveError() *AppError {
	return NewAppError(ErrCodeAccountInactive, "Account inactive", http.StatusForbidden)
}

// Resource error constructors.
func newResourceError(code ErrorCode, message, resource string) *AppError {
	return NewAppErrorWithDetails(
		code,
		message,
		http.StatusNotFound,
		map[string]any{"resource": resource},
	)
}

// NewNotFoundError creates a not found error for a resource.
func NewNotFoundError(resource string) *AppError {
	return newResourceError(ErrCodeResourceNotFound, "Resource not found", resource)
}

// NewUserNotFoundError creates a user not found error.
func NewUserNotFoundError() *AppError {
	return NewNotFoundError("user")
}

// NewSessionNotFoundError creates a session not found error.
func NewSessionNotFoundError() *AppError {
	return NewNotFoundError("session")
}

// NewAlreadyExistsError creates an already exists error.
func NewAlreadyExistsError(resource string) *AppError {
	return newResourceError(ErrCodeAlreadyExists, "Resource already exists", resource)
}

// NewResourceConflictError creates a resource conflict error.
func NewResourceConflictError(resource, message string) *AppError {
	return NewAppErrorWithDetails(
		ErrCodeResourceConflict,
		message,
		http.StatusConflict,
		map[string]any{"resource": resource},
	)
}

// NewInternalError creates an internal error.
func NewInternalError(message string, cause error) *AppError {
	return NewAppErrorWithCause(ErrCodeInternal, message, http.StatusInternalServerError, cause)
}

// NewDatabaseError creates a database error.
func NewDatabaseError(message string, cause error) *AppError {
	return NewAppErrorWithCause(ErrCodeDatabase, message, http.StatusInternalServerError, cause)
}

// NewNetworkError creates a network error.
func NewNetworkError(message string, cause error) *AppError {
	return NewAppErrorWithCause(ErrCodeNetwork, message, http.StatusServiceUnavailable, cause)
}

// NewTimeoutError creates a timeout error.
func NewTimeoutError(operation string) *AppError {
	return NewAppError(ErrCodeTimeout, "Operation timed out: "+operation, http.StatusRequestTimeout)
}

// NewUnavailableError creates a service unavailable error.
func NewUnavailableError(service string) *AppError {
	return NewAppError(
		ErrCodeUnavailable,
		"Service unavailable: "+service,
		http.StatusServiceUnavailable,
	)
}

// NewBusinessLogicError creates a business logic error.
func NewBusinessLogicError(message string) *AppError {
	return NewAppError(ErrCodeBusinessLogic, message, http.StatusBadRequest)
}

// NewInvalidStateError creates an invalid state error.
func NewInvalidStateError(state, operation string) *AppError {
	return newBadRequestError(
		ErrCodeInvalidState,
		"Invalid state for operation",
		"state",
		state,
		"operation",
		operation,
	)
}

// NewPermissionDeniedError creates a permission denied error.
func NewPermissionDeniedError(operation string) *AppError {
	return NewAppErrorWithDetails(
		ErrCodePermissionDenied,
		"Permission denied",
		http.StatusForbidden,
		map[string]any{"operation": operation},
	)
}

// IsAppError checks if err is an AppError.
func IsAppError(err error) bool {
	appError := &AppError{}
	ok := errors.As(err, &appError)

	return ok
}

// IsValidationError checks if err is a validation error.
func IsValidationError(err error) bool {
	return hasErrorCode(
		err,
		ErrCodeValidationFailed,
		ErrCodeInvalidInput,
		ErrCodeMissingField,
		ErrCodeInvalidFormat,
		ErrCodeConstraintFailed,
	)
}

// IsNotFoundError checks if err is a not found error.
func IsNotFoundError(err error) bool {
	appErr := &AppError{}
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeResourceNotFound || appErr.Code == ErrCodeNotFound
	}

	return false
}

// IsUnauthorizedError checks if err is an unauthorized error.
func IsUnauthorizedError(err error) bool {
	appErr := &AppError{}
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeUnauthorized ||
			appErr.Code == ErrCodeInvalidCredentials ||
			appErr.Code == ErrCodeTokenExpired ||
			appErr.Code == ErrCodeTokenInvalid
	}

	return false
}

// hasErrorCode checks if err contains any of the given error codes.
func hasErrorCode(err error, codes ...ErrorCode) bool {
	appErr := &AppError{}
	if errors.As(err, &appErr) {
		if slices.Contains(codes, appErr.Code) {
			return true
		}
	}

	return false
}

// IsForbiddenError checks if err is a forbidden error.
func IsForbiddenError(err error) bool {
	return hasErrorCode(
		err,
		ErrCodeForbidden,
		ErrCodeInsufficientPrivileges,
		ErrCodeAccountSuspended,
		ErrCodeAccountInactive,
	)
}

// IsInternalServerError checks if err is an internal server error.
func IsInternalServerError(err error) bool {
	return hasErrorCode(
		err,
		ErrCodeInternal,
		ErrCodeDatabase,
		ErrCodeNetwork,
	)
}

// HTTP Status Code mapping.
//
//nolint:gochecknoglobals // Intentional lookup table for error code to HTTP status mapping
var errorCodeToHTTPStatus = map[ErrorCode]int{
	ErrCodeValidationFailed:       http.StatusBadRequest,
	ErrCodeInvalidInput:           http.StatusBadRequest,
	ErrCodeMissingField:           http.StatusBadRequest,
	ErrCodeInvalidFormat:          http.StatusBadRequest,
	ErrCodeConstraintFailed:       http.StatusBadRequest,
	ErrCodeBusinessLogic:          http.StatusBadRequest,
	ErrCodeInvalidState:           http.StatusBadRequest,
	ErrCodeUnauthorized:           http.StatusUnauthorized,
	ErrCodeInvalidCredentials:     http.StatusUnauthorized,
	ErrCodeTokenExpired:           http.StatusUnauthorized,
	ErrCodeTokenInvalid:           http.StatusUnauthorized,
	ErrCodeForbidden:              http.StatusForbidden,
	ErrCodeInsufficientPrivileges: http.StatusForbidden,
	ErrCodeAccountSuspended:       http.StatusForbidden,
	ErrCodeAccountInactive:        http.StatusForbidden,
	ErrCodePermissionDenied:       http.StatusForbidden,
	ErrCodeNotFound:               http.StatusNotFound,
	ErrCodeResourceNotFound:       http.StatusNotFound,
	ErrCodeAlreadyExists:          http.StatusConflict,
	ErrCodeResourceConflict:       http.StatusConflict,
	ErrCodeTimeout:                http.StatusRequestTimeout,
	ErrCodeUnavailable:            http.StatusServiceUnavailable,
}

// StatusCode returns the HTTP status code for the error.
func (e *AppError) StatusCode() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}

	if status, ok := errorCodeToHTTPStatus[e.Code]; ok {
		return status
	}

	return http.StatusInternalServerError
}
