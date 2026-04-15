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
	// Validation errors.
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField     ErrorCode = "MISSING_FIELD"
	ErrCodeInvalidFormat    ErrorCode = "INVALID_FORMAT"
	ErrCodeConstraintFailed ErrorCode = "CONSTRAINT_FAILED"

	// Authentication errors.
	ErrCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid       ErrorCode = "TOKEN_INVALID"

	// Authorization errors.
	ErrCodeForbidden              ErrorCode = "FORBIDDEN"
	ErrCodeInsufficientPrivileges ErrorCode = "INSUFFICIENT_PRIVILEGES"
	ErrCodeAccountSuspended       ErrorCode = "ACCOUNT_SUSPENDED"
	ErrCodeAccountInactive        ErrorCode = "ACCOUNT_INACTIVE"

	// Resource errors.
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	ErrCodeAlreadyExists    ErrorCode = "ALREADY_EXISTS"
	ErrCodeResourceConflict ErrorCode = "RESOURCE_CONFLICT"

	// System errors.
	ErrCodeInternal    ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase    ErrorCode = "DATABASE_ERROR"
	ErrCodeNetwork     ErrorCode = "NETWORK_ERROR"
	ErrCodeTimeout     ErrorCode = "TIMEOUT"
	ErrCodeUnavailable ErrorCode = "UNAVAILABLE"

	// Business logic errors.
	ErrCodeBusinessLogic    ErrorCode = "BUSINESS_LOGIC_ERROR"
	ErrCodeInvalidState     ErrorCode = "INVALID_STATE"
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

// Validation error constructors.
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

func NewMissingFieldError(field string) *AppError {
	return newFieldError(ErrCodeMissingField, "Required field is missing", field)
}

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

// Authentication error constructors.
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

func NewInvalidCredentialsError() *AppError {
	return NewAppError(ErrCodeInvalidCredentials, "Invalid credentials", http.StatusUnauthorized)
}

func NewTokenExpiredError() *AppError {
	return NewAppError(ErrCodeTokenExpired, "Token has expired", http.StatusUnauthorized)
}

func NewTokenInvalidError() *AppError {
	return NewAppError(ErrCodeTokenInvalid, "Invalid token", http.StatusUnauthorized)
}

// Authorization error constructors.
func NewForbiddenError(message string) *AppError {
	return NewAppError(ErrCodeForbidden, message, http.StatusForbidden)
}

func NewInsufficientPrivilegesError() *AppError {
	return NewAppError(
		ErrCodeInsufficientPrivileges,
		"Insufficient privileges",
		http.StatusForbidden,
	)
}

func NewAccountSuspendedError() *AppError {
	return NewAppError(ErrCodeAccountSuspended, "Account suspended", http.StatusForbidden)
}

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

func NewNotFoundError(resource string) *AppError {
	return newResourceError(ErrCodeResourceNotFound, "Resource not found", resource)
}

func NewUserNotFoundError() *AppError {
	return NewNotFoundError("user")
}

func NewSessionNotFoundError() *AppError {
	return NewNotFoundError("session")
}

func NewAlreadyExistsError(resource string) *AppError {
	return newResourceError(ErrCodeAlreadyExists, "Resource already exists", resource)
}

func NewResourceConflictError(resource, message string) *AppError {
	return NewAppErrorWithDetails(
		ErrCodeResourceConflict,
		message,
		http.StatusConflict,
		map[string]any{"resource": resource},
	)
}

// System error constructors.
func NewInternalError(message string, cause error) *AppError {
	return NewAppErrorWithCause(ErrCodeInternal, message, http.StatusInternalServerError, cause)
}

func NewDatabaseError(message string, cause error) *AppError {
	return NewAppErrorWithCause(ErrCodeDatabase, message, http.StatusInternalServerError, cause)
}

func NewNetworkError(message string, cause error) *AppError {
	return NewAppErrorWithCause(ErrCodeNetwork, message, http.StatusServiceUnavailable, cause)
}

func NewTimeoutError(operation string) *AppError {
	return NewAppError(ErrCodeTimeout, "Operation timed out: "+operation, http.StatusRequestTimeout)
}

func NewUnavailableError(service string) *AppError {
	return NewAppError(
		ErrCodeUnavailable,
		"Service unavailable: "+service,
		http.StatusServiceUnavailable,
	)
}

// Business logic error constructors.
func NewBusinessLogicError(message string) *AppError {
	return NewAppError(ErrCodeBusinessLogic, message, http.StatusBadRequest)
}

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

func NewPermissionDeniedError(operation string) *AppError {
	return NewAppErrorWithDetails(
		ErrCodePermissionDenied,
		"Permission denied",
		http.StatusForbidden,
		map[string]any{"operation": operation},
	)
}

// Error checking functions.
func IsAppError(err error) bool {
	appError := &AppError{}
	ok := errors.As(err, &appError)

	return ok
}

func IsValidationError(err error) bool {
	return hasErrorCode(err,
		ErrCodeValidationFailed,
		ErrCodeInvalidInput,
		ErrCodeMissingField,
		ErrCodeInvalidFormat,
		ErrCodeConstraintFailed,
	)
}

func IsNotFoundError(err error) bool {
	appErr := &AppError{}
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeResourceNotFound || appErr.Code == ErrCodeNotFound
	}

	return false
}

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

func IsForbiddenError(err error) bool {
	return hasErrorCode(err,
		ErrCodeForbidden,
		ErrCodeInsufficientPrivileges,
		ErrCodeAccountSuspended,
		ErrCodeAccountInactive,
	)
}

func IsInternalServerError(err error) bool {
	return hasErrorCode(err,
		ErrCodeInternal,
		ErrCodeDatabase,
		ErrCodeNetwork,
	)
}

// HTTP Status Code mapping.
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

func (e *AppError) StatusCode() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}

	if status, ok := errorCodeToHTTPStatus[e.Code]; ok {
		return status
	}

	return http.StatusInternalServerError
}
