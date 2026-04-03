package postgres

import "errors"

// pgCodeChecker is an interface for PostgreSQL errors that have error codes.
type pgCodeChecker interface {
	Code() string
}

// isPostgresErrorCode checks if the error has the specified PostgreSQL error code.
// See https://www.postgresql.org/docs/current/errcodes-appendix.html for error codes.
func isPostgresErrorCode(err error, code string) bool {
	var pgErr pgCodeChecker
	if errors.As(err, &pgErr) {
		return pgErr.Code() == code
	}
	return false
}

// isUniqueViolationError checks for PostgreSQL unique constraint violation (23505)
func isUniqueViolationError(err error) bool {
	return isPostgresErrorCode(err, "23505")
}

// isForeignKeyViolationError checks for PostgreSQL foreign key violation (23503)
func isForeignKeyViolationError(err error) bool {
	return isPostgresErrorCode(err, "23503")
}

// isCheckViolationError checks for PostgreSQL check constraint violation (23514)
func isCheckViolationError(err error) bool {
	return isPostgresErrorCode(err, "23514")
}
