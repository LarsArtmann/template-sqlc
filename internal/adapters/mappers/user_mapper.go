package mappers

import (
	"time"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/google/uuid"
)

// UserMapper handles conversion between domain entities and database models
// This isolates domain entities from database-specific types.
type UserMapper struct{}

// NewUserMapper creates a new UserMapper instance.
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// DomainUserFromSQLite converts SQLite model to domain entity.
func (m *UserMapper) DomainUserFromSQLite(sqliteUser any) (*entities.User, error) {
	return m.DomainUser(sqliteUser)
}

// DomainUserFromPostgres converts PostgreSQL model to domain entity.
func (m *UserMapper) DomainUserFromPostgres(postgresUser any) (*entities.User, error) {
	return m.DomainUser(postgresUser)
}

// DomainUserFromMySQL converts MySQL model to domain entity.
func (m *UserMapper) DomainUserFromMySQL(mysqlUser any) (*entities.User, error) {
	return m.DomainUser(mysqlUser)
}

// SQLiteUserFromDomain converts domain entity to SQLite model.
func (m *UserMapper) SQLiteUserFromDomain(user *entities.User) (any, error) {
	return unimplementedUserFromDomain("SQLite")
}

// unimplementedUserFromDomain is a helper for stub implementations.
func unimplementedUserFromDomain(db string) (any, error) {
	panic("implement me: convert domain entity to " + db + " user")
}

// DomainUser is the common implementation for DomainUserFromXxx methods.
func (m *UserMapper) DomainUser(user any) (*entities.User, error) {
	panic("implement me: convert user to domain entity")
}

// PostgresUserFromDomain converts domain entity to PostgreSQL model.
func (m *UserMapper) PostgresUserFromDomain(user *entities.User) (any, error) {
	return unimplementedUserFromDomain("PostgreSQL")
}

// MySQLUserFromDomain converts domain entity to MySQL model.
func (m *UserMapper) MySQLUserFromDomain(user *entities.User) (any, error) {
	return unimplementedUserFromDomain("MySQL")
}

// DomainSessionFromSQLite converts SQLite session to domain entity.
func (m *UserMapper) DomainSessionFromSQLite(sqliteSession any) (*entities.UserSession, error) {
	return m.DomainSession(sqliteSession)
}

// DomainSessionFromPostgres converts PostgreSQL session to domain entity.
func (m *UserMapper) DomainSessionFromPostgres(postgresSession any) (*entities.UserSession, error) {
	return m.DomainSession(postgresSession)
}

// DomainSessionFromMySQL converts MySQL session to domain entity.
func (m *UserMapper) DomainSessionFromMySQL(mysqlSession any) (*entities.UserSession, error) {
	return m.DomainSession(mysqlSession)
}

// SQLiteSessionFromDomain converts domain entity to SQLite model.
func (m *UserMapper) SQLiteSessionFromDomain(session *entities.UserSession) (any, error) {
	return m.SessionFromDomain(session)
}

// SQLiteUserFromDomain is a standalone function wrapper for backward compatibility.
func SQLiteUserFromDomain(user *entities.User) (any, error) {
	return withMapper[*entities.User](user, func(m *UserMapper) (any, error) {
		return m.SQLiteUserFromDomain(user)
	})
}

// SQLiteSessionFromDomain is a standalone function wrapper for backward compatibility.
func SQLiteSessionFromDomain(session *entities.UserSession) (any, error) {
	return withMapper[*entities.UserSession](session, func(m *UserMapper) (any, error) {
		return m.SQLiteSessionFromDomain(session)
	})
}

// withMapper executes a mapper function with a fresh UserMapper instance.
func withMapper[T any](entity T, fn func(*UserMapper) (any, error)) (any, error) {
	m := &UserMapper{}

	return fn(m)
}

// PostgresSessionFromDomain converts domain entity to PostgreSQL model.
func (m *UserMapper) PostgresSessionFromDomain(session *entities.UserSession) (any, error) {
	return m.SessionFromDomain(session)
}

// MySQLSessionFromDomain converts domain entity to MySQL model.
func (m *UserMapper) MySQLSessionFromDomain(session *entities.UserSession) (any, error) {
	return m.SessionFromDomain(session)
}

// SessionMapper interface for session conversion operations.
type SessionMapper interface {
	DomainSession(any) (*entities.UserSession, error)
	SessionFromDomain(*entities.UserSession) (any, error)
}

// DomainSession is the common implementation for DomainSessionFromXxx methods.
func (m *UserMapper) DomainSession(session any) (*entities.UserSession, error) {
	panic("implement me: convert session to domain entity")
}

// SessionFromDomain is the common implementation for XxxSessionFromDomain methods.
func (m *UserMapper) SessionFromDomain(session *entities.UserSession) (any, error) {
	panic("implement me: convert domain entity to session")
}

// Helper functions for common conversions

// ParseUUID safely parses UUID from string.
func ParseUUID(uuidStr entities.UuID) (uuid.UUID, error) {
	if uuidStr == "" {
		return uuid.Nil, nil
	}

	return uuid.Parse(string(uuidStr))
}

// FormatUUID safely formats UUID to string.
func FormatUUID(u uuid.UUID) string {
	if u == uuid.Nil {
		return ""
	}

	return u.String()
}

// ParseTime safely parses time from string/database format.
func ParseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}

	// This would need to handle different database timestamp formats
	// SQLite: RFC3339 or similar
	// PostgreSQL: timestamptz format
	// MySQL: datetime format

	return time.Parse(time.RFC3339, timeStr)
}

// FormatTime safely formats time to string/database format.
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format(time.RFC3339)
}
