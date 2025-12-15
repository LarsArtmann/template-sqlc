package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/mappers"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
	"github.com/LarsArtmann/template-sqlc/pkg/errors"
)

// SQLiteSessionRepository implements SessionRepository for SQLite
type SQLiteSessionRepository struct {
	db         *sql.DB
	mapper     mappers.UserMapper
	converters *ConverterSet
}

// NewSQLiteSessionRepository creates a new SQLite session repository
func NewSQLiteSessionRepository(db *sql.DB) repositories.SessionRepository {
	return &SQLiteSessionRepository{
		db: db,
		converters: &ConverterSet{
			SessionToken: converters.NewDefaultSessionTokenConverter(),
			Time:         converters.NewTimeConverter("sqlite"),
			Bool:         converters.NewBoolConverter("sqlite"),
		},
	}
}

// Create saves a new session to SQLite
func (r *SQLiteSessionRepository) Create(ctx context.Context, session *entities.UserSession) error {
	// Convert domain entity to SQLite model
	sqliteSession, err := mappers.SQLiteSessionFromDomain(session)
	if err != nil {
		return fmt.Errorf("failed to convert session: %w", err)
	}

	// This would use actual generated sqlc code
	// _, err := r.queries.CreateSession(ctx, sqliteSession.(sqlite.CreateSessionParams))
	// return errors.NewDatabaseError("failed to create session", err)

	panic("implement me: use actual sqlc generated code")
}

// GetByToken retrieves a session by token from SQLite
func (r *SQLiteSessionRepository) GetByToken(ctx context.Context, token entities.SessionToken) (*entities.UserSession, error) {
	// Convert token to database format
	dbToken := r.converters.SessionToken.DomainToDB(token)

	// Query database
	// sqliteSession, err := r.queries.GetSessionByToken(ctx, dbToken)
	// if err != nil {
	//     if err == sql.ErrNoRows {
	//         return nil, entities.ErrSessionNotFound
	//     }
	//     return nil, errors.NewDatabaseError("failed to get session", err)
	// }
	// return mappers.DomainSessionFromSQLite(sqliteSession)

	panic("implement me: use actual sqlc generated code")
}

// GetByUserID retrieves sessions by user ID from SQLite
func (r *SQLiteSessionRepository) GetByUserID(ctx context.Context, userID entities.UserID, activeOnly bool) ([]*entities.UserSession, error) {
	// Query sessions by user ID
	// sqliteSessions, err := r.queries.GetSessionsByUserID(ctx, int64(userID), activeOnly)
	// if err != nil {
	//     return nil, errors.NewDatabaseError("failed to get sessions", err)
	// }

	// Convert to domain entities
	// sessions := make([]*entities.UserSession, len(sqliteSessions))
	// for i, sqliteSession := range sqliteSessions {
	//     session, err := mappers.DomainSessionFromSQLite(sqliteSession)
	//     if err != nil {
	//         return nil, fmt.Errorf("failed to convert session: %w", err)
	//     }
	//     sessions[i] = session
	// }
	// return sessions, nil

	panic("implement me: use actual sqlc generated code")
}

// Update updates a session in SQLite
func (r *SQLiteSessionRepository) Update(ctx context.Context, session *entities.UserSession) error {
	// Convert domain entity to SQLite model
	sqliteSession, err := mappers.SQLiteSessionFromDomain(session)
	if err != nil {
		return fmt.Errorf("failed to convert session: %w", err)
	}

	// Update in database
	panic("implement me: use actual sqlc generated code")
}

// Delete removes a session from SQLite
func (r *SQLiteSessionRepository) Delete(ctx context.Context, id entities.SessionID) error {
	// Delete session
	panic("implement me: use actual sqlc generated code")
}

// DeactivateByToken deactivates a session by token in SQLite
func (r *SQLiteSessionRepository) DeactivateByToken(ctx context.Context, token entities.SessionToken) error {
	// Convert token to database format
	dbToken := r.converters.SessionToken.DomainToDB(token)

	// Deactivate session
	// _, err := r.queries.DeactivateSessionByToken(ctx, dbToken)
	// return errors.NewDatabaseError("failed to deactivate session", err)

	panic("implement me: use actual sqlc generated code")
}

// DeactivateByUserID deactivates all sessions for a user in SQLite
func (r *SQLiteSessionRepository) DeactivateByUserID(ctx context.Context, userID entities.UserID) error {
	// Deactivate all user sessions
	// _, err := r.queries.DeactivateSessionsByUserID(ctx, int64(userID))
	// return errors.NewDatabaseError("failed to deactivate user sessions", err)

	panic("implement me: use actual sqlc generated code")
}

// CleanupExpired removes expired sessions from SQLite
func (r *SQLiteSessionRepository) CleanupExpired(ctx context.Context) (int64, error) {
	// Clean up expired sessions
	// count, err := r.queries.CleanupExpiredSessions(ctx)
	// return count, errors.NewDatabaseError("failed to cleanup expired sessions", err)

	panic("implement me: use actual sqlc generated code")
}

// GetActiveSessions returns count of active sessions for a user in SQLite
func (r *SQLiteSessionRepository) GetActiveSessions(ctx context.Context, userID entities.UserID) (int64, error) {
	// Count active sessions
	// count, err := r.queries.GetActiveSessionCount(ctx, int64(userID))
	// return count, errors.NewDatabaseError("failed to get active session count", err)

	panic("implement me: use actual sqlc generated code")
}

// GetSessionStats returns session statistics from SQLite
func (r *SQLiteSessionRepository) GetSessionStats(ctx context.Context) (*entities.SessionStats, error) {
	// Query session statistics
	panic("implement me: use actual sqlc generated code")
}

// Helper methods

// handleSessionError converts database errors to domain errors
func (r *SQLiteSessionRepository) handleSessionError(err error, operation string) error {
	if err == nil {
		return nil
	}

	switch {
	case err == sql.ErrNoRows:
		return entities.ErrSessionNotFound
	case isSessionUniqueConstraintError(err):
		return entities.ErrUserAlreadyExists // or session-specific error
	default:
		return errors.NewDatabaseError(fmt.Sprintf("%s failed", operation), err)
	}
}

// isSessionUniqueConstraintError checks if error is a session-related unique constraint
func isSessionUniqueConstraintError(err error) bool {
	// This would check for SQLite-specific session constraint errors
	return err != nil &&
		(fmt.Sprintf("%s", err) == "UNIQUE constraint failed: sessions.token" ||
			fmt.Sprintf("%s", err) == "session token already exists")
}
