package mappers

import (
	"time"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/google/uuid"
)

// UserMapper handles conversion between domain entities and database models
// This isolates domain entities from database-specific types

// DomainUserFromSQLite converts SQLite model to domain entity
func DomainUserFromSQLite(sqliteUser interface{}) (*entities.User, error) {
	// This would be implemented based on actual generated SQLite types
	// Example implementation - adapt to your actual generated types

	// You would typically do something like:
	// dbUser := sqliteUser.(sqlite.Users)
	// return &entities.User{
	//     id: entities.UserID(dbUser.ID),
	//     // ... field mappings
	// }, nil

	panic("implement me: convert SQLite user to domain entity")
}

// DomainUserFromPostgres converts PostgreSQL model to domain entity
func DomainUserFromPostgres(postgresUser interface{}) (*entities.User, error) {
	// Similar implementation for PostgreSQL types
	panic("implement me: convert PostgreSQL user to domain entity")
}

// DomainUserFromMySQL converts MySQL model to domain entity
func DomainUserFromMySQL(mysqlUser interface{}) (*entities.User, error) {
	// Similar implementation for MySQL types
	panic("implement me: convert MySQL user to domain entity")
}

// SQLiteUserFromDomain converts domain entity to SQLite model
func SQLiteUserFromDomain(user *entities.User) (interface{}, error) {
	// Convert domain entity to SQLite-specific model
	// This would be implemented based on your actual generated types

	// Example:
	// return &sqlite.Users{
	//     ID:           int64(user.ID()),
	//     UUID:         user.UUID().String(),
	//     Email:        user.Email().String(),
	//     // ... field mappings
	// }, nil

	panic("implement me: convert domain entity to SQLite user")
}

// PostgresUserFromDomain converts domain entity to PostgreSQL model
func PostgresUserFromDomain(user *entities.User) (interface{}, error) {
	// Similar implementation for PostgreSQL
	panic("implement me: convert domain entity to PostgreSQL user")
}

// MySQLUserFromDomain converts domain entity to MySQL model
func MySQLUserFromDomain(user *entities.User) (interface{}, error) {
	// Similar implementation for MySQL
	panic("implement me: convert domain entity to MySQL user")
}

// DomainSessionFromSQLite converts SQLite session to domain entity
func DomainSessionFromSQLite(sqliteSession interface{}) (*entities.UserSession, error) {
	panic("implement me: convert SQLite session to domain entity")
}

// DomainSessionFromPostgres converts PostgreSQL session to domain entity
func DomainSessionFromPostgres(postgresSession interface{}) (*entities.UserSession, error) {
	panic("implement me: convert PostgreSQL session to domain entity")
}

// DomainSessionFromMySQL converts MySQL session to domain entity
func DomainSessionFromMySQL(mysqlSession interface{}) (*entities.UserSession, error) {
	panic("implement me: convert MySQL session to domain entity")
}

// SQLiteSessionFromDomain converts domain entity to SQLite model
func SQLiteSessionFromDomain(session *entities.UserSession) (interface{}, error) {
	panic("implement me: convert domain entity to SQLite session")
}

// PostgresSessionFromDomain converts domain entity to PostgreSQL model
func PostgresSessionFromDomain(session *entities.UserSession) (interface{}, error) {
	panic("implement me: convert domain entity to PostgreSQL session")
}

// MySQLSessionFromDomain converts domain entity to MySQL model
func MySQLSessionFromDomain(session *entities.UserSession) (interface{}, error) {
	panic("implement me: convert domain entity to MySQL session")
}

// Helper functions for common conversions

// ParseUUID safely parses UUID from string
func ParseUUID(uuidStr string) (uuid.UUID, error) {
	if uuidStr == "" {
		return uuid.Nil, nil
	}
	return uuid.Parse(uuidStr)
}

// FormatUUID safely formats UUID to string
func FormatUUID(u uuid.UUID) string {
	if u == uuid.Nil {
		return ""
	}
	return u.String()
}

// ParseTime safely parses time from string/database format
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

// FormatTime safely formats time to string/database format
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// ParseBool safely parses boolean from various database formats
func ParseBool(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int64:
		return v != 0
	case int:
		return v != 0
	case string:
		return v == "true" || v == "1"
	case nil:
		return false
	default:
		return false
	}
}

// FormatBool safely formats boolean to database format
func FormatBool(b bool) interface{} {
	return b
}

// ParseInterface safely parses interface{} to string
func ParseInterface(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}
