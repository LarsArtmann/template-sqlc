package sqlite

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// SQLiteSessionRepository implements SessionRepository for SQLite.
type SQLiteSessionRepository struct {
	db         any
	converters *converters.ConverterSet
}

// NewSQLiteSessionRepository creates a new SQLite session repository.
func NewSQLiteSessionRepository(db any) repositories.SessionRepository {
	return &SQLiteSessionRepository{
		db:         db,
		converters: converters.NewConverterSet(converters.DbTypeSQLite),
	}
}

// Create saves a new session to SQLite.
func (r *SQLiteSessionRepository) Create(_ context.Context, _ *entities.UserSession) error {
	panic("implement me: use actual sqlc generated code")
}

// GetByToken retrieves a session by token from SQLite.
func (r *SQLiteSessionRepository) GetByToken(
	_ context.Context,
	_ entities.SessionToken,
) (*entities.UserSession, error) {
	panic("implement me: use actual sqlc generated code")
}

// GetByUserID retrieves sessions by user ID from SQLite.
func (r *SQLiteSessionRepository) GetByUserID(
	_ context.Context,
	_ entities.UserID,
	_ bool,
) ([]*entities.UserSession, error) {
	panic("implement me: use actual sqlc generated code")
}

// Update updates a session in SQLite.
func (r *SQLiteSessionRepository) Update(_ context.Context, _ *entities.UserSession) error {
	panic("implement me: use actual sqlc generated code")
}

// Delete removes a session from SQLite.
func (r *SQLiteSessionRepository) Delete(_ context.Context, _ entities.SessionID) error {
	return entities.StubNotImplemented("Delete", "SQLite")
}

// DeactivateByToken deactivates a session by token in SQLite.
func (r *SQLiteSessionRepository) DeactivateByToken(
	_ context.Context,
	_ entities.SessionToken,
) error {
	panic("implement me: use actual sqlc generated code")
}

// DeactivateByUserID deactivates all sessions for a user in SQLite.
func (r *SQLiteSessionRepository) DeactivateByUserID(
	_ context.Context,
	_ entities.UserID,
) error {
	return entities.StubNotImplemented("DeactivateByUserID", "SQLite")
}

// CleanupExpired removes expired sessions from SQLite.
func (r *SQLiteSessionRepository) CleanupExpired(_ context.Context) (int64, error) {
	panic("implement me: use actual sqlc generated code")
}

// GetActiveSessions returns count of active sessions for a user in SQLite.
func (r *SQLiteSessionRepository) GetActiveSessions(
	_ context.Context,
	_ entities.UserID,
) (int64, error) {
	panic("implement me: use actual sqlc generated code")
}

// GetSessionStats returns session statistics from SQLite.
func (r *SQLiteSessionRepository) GetSessionStats(
	_ context.Context,
) (*entities.SessionStats, error) {
	panic("implement me: use actual sqlc generated code")
}
