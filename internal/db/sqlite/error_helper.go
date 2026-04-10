package sqlite

import (
	"database/sql"
	stderrors "errors"
	"strings"

	apperrors "github.com/LarsArtmann/template-sqlc/pkg/errors"
)

// HandleDBError converts database errors to domain errors.
// Takes entity-specific notFoundErr for sql.ErrNoRows and a custom
// uniqueConstraintChecker for constraint violations.
func HandleDBError(
	err error,
	operation string,
	notFoundErr, conflictErr error,
	uniqueConstraintChecker func(error) bool,
) error {
	if err == nil {
		return nil
	}

	switch {
	case stderrors.Is(err, sql.ErrNoRows):
		return notFoundErr
	case uniqueConstraintChecker(err):
		return conflictErr
	default:
		return apperrors.NewDatabaseError(operation+" failed", err)
	}
}

// isNilOrErrorMsgContains checks if the error message contains any of the given substrings.
func isNilOrErrorMsgContains(err error, substrs ...string) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	for _, substr := range substrs {
		if strings.Contains(msg, substr) {
			return true
		}
	}

	return false
}

// IsSQLiteUniqueConstraintError checks if error is a SQLite UNIQUE constraint violation.
func IsSQLiteUniqueConstraintError(err error) bool {
	return isNilOrErrorMsgContains(err, "UNIQUE constraint failed", "is not unique")
}

// IsSQLiteSessionTokenConstraintError checks if error is a session token constraint violation.
func IsSQLiteSessionTokenConstraintError(err error) bool {
	return isNilOrErrorMsgContains(
		err,
		"UNIQUE constraint failed: sessions.token",
		"session token already exists",
	)
}
