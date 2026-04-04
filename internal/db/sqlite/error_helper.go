package sqlite

import (
	"database/sql"
	stderrors "errors"

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

// IsSQLiteUniqueConstraintError checks if error is a SQLite UNIQUE constraint violation.
func IsSQLiteUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	return contains(msg, "UNIQUE constraint failed") ||
		contains(msg, "is not unique")
}

// IsSQLiteSessionTokenConstraintError checks if error is a session token constraint violation.
func IsSQLiteSessionTokenConstraintError(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	return contains(msg, "UNIQUE constraint failed: sessions.token") ||
		contains(msg, "session token already exists")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
