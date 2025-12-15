package converters

import (
	"time"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/google/uuid"
)

// TypeConverter handles database-specific type conversions
// This isolates domain entities from database-specific type handling

// UUIDConverter handles UUID conversions between domain and database
type UUIDConverter interface {
	DomainToDB(domain uuid.UUID) interface{}
	DBToDomain(db interface{}) (uuid.UUID, error)
}

// TimeConverter handles time conversions between domain and database
type TimeConverter interface {
	DomainToDB(domain time.Time) interface{}
	DBToDomain(db interface{}) (time.Time, error)
}

// BoolConverter handles boolean conversions between domain and database
type BoolConverter interface {
	DomainToDB(domain bool) interface{}
	DBToDomain(db interface{}) (bool, error)
}

// EmailConverter handles email conversions between domain and database
type EmailConverter interface {
	DomainToDB(domain entities.Email) string
	DBToDomain(db string) (entities.Email, error)
}

// UsernameConverter handles username conversions between domain and database
type UsernameConverter interface {
	DomainToDB(domain entities.Username) string
	DBToDomain(db string) (entities.Username, error)
}

// PasswordHashConverter handles password hash conversions between domain and database
type PasswordHashConverter interface {
	DomainToDB(domain entities.PasswordHash) string
	DBToDomain(db string) (entities.PasswordHash, error)
}

// UserStatusConverter handles user status conversions between domain and database
type UserStatusConverter interface {
	DomainToDB(domain entities.UserStatus) string
	DBToDomain(db string) (entities.UserStatus, error)
}

// UserRoleConverter handles user role conversions between domain and database
type UserRoleConverter interface {
	DomainToDB(domain entities.UserRole) string
	DBToDomain(db string) (entities.UserRole, error)
}

// SessionTokenConverter handles session token conversions between domain and database
type SessionTokenConverter interface {
	DomainToDB(domain entities.SessionToken) interface{}
	DBToDomain(db interface{}) (entities.SessionToken, error)
}

// Default implementations

// SQLiteUUIDConverter handles UUID conversion for SQLite (stores as string)
type SQLiteUUIDConverter struct{}

func NewSQLiteUUIDConverter() *SQLiteUUIDConverter {
	return &SQLiteUUIDConverter{}
}

func (c *SQLiteUUIDConverter) DomainToDB(domain uuid.UUID) interface{} {
	if domain == uuid.Nil {
		return nil
	}
	return domain.String()
}

func (c *SQLiteUUIDConverter) DBToDomain(db interface{}) (uuid.UUID, error) {
	if db == nil {
		return uuid.Nil, nil
	}

	str, ok := db.(string)
	if !ok {
		return uuid.Nil, NewConversionError("expected string for UUID", db)
	}

	return uuid.Parse(str)
}

// PostgresUUIDConverter handles UUID conversion for PostgreSQL (stores as UUID type)
type PostgresUUIDConverter struct{}

func NewPostgresUUIDConverter() *PostgresUUIDConverter {
	return &PostgresUUIDConverter{}
}

func (c *PostgresUUIDConverter) DomainToDB(domain uuid.UUID) interface{} {
	return domain // PostgreSQL handles UUID natively
}

func (c *PostgresUUIDConverter) DBToDomain(db interface{}) (uuid.UUID, error) {
	if db == nil {
		return uuid.Nil, nil
	}

	if uuidObj, ok := db.(uuid.UUID); ok {
		return uuidObj, nil
	}

	if str, ok := db.(string); ok {
		return uuid.Parse(str)
	}

	return uuid.Nil, NewConversionError("expected UUID or string", db)
}

// MySQLUUIDConverter handles UUID conversion for MySQL (stores as binary)
type MySQLUUIDConverter struct{}

func NewMySQLUUIDConverter() *MySQLUUIDConverter {
	return &MySQLUUIDConverter{}
}

func (c *MySQLUUIDConverter) DomainToDB(domain uuid.UUID) interface{} {
	return domain[:] // Convert to byte slice
}

func (c *MySQLUUIDConverter) DBToDomain(db interface{}) (uuid.UUID, error) {
	if db == nil {
		return uuid.Nil, nil
	}

	if bytes, ok := db.([]byte); ok {
		return uuid.FromBytes(bytes)
	}

	if str, ok := db.(string); ok {
		return uuid.Parse(str)
	}

	return uuid.Nil, NewConversionError("expected bytes or string for UUID", db)
}

// SQLiteTimeConverter handles time conversion for SQLite
type SQLiteTimeConverter struct{}

func NewSQLiteTimeConverter() *SQLiteTimeConverter {
	return &SQLiteTimeConverter{}
}

func (c *SQLiteTimeConverter) DomainToDB(domain time.Time) interface{} {
	if domain.IsZero() {
		return nil
	}
	return domain
}

func (c *SQLiteTimeConverter) DBToDomain(db interface{}) (time.Time, error) {
	if db == nil {
		return time.Time{}, nil
	}

	if t, ok := db.(time.Time); ok {
		return t, nil
	}

	if str, ok := db.(string); ok {
		return time.Parse(time.RFC3339, str)
	}

	return time.Time{}, NewConversionError("expected time or string", db)
}

// SQLiteBoolConverter handles boolean conversion for SQLite
type SQLiteBoolConverter struct{}

func NewSQLiteBoolConverter() *SQLiteBoolConverter {
	return &SQLiteBoolConverter{}
}

func (c *SQLiteBoolConverter) DomainToDB(domain bool) interface{} {
	return domain // SQLite supports boolean natively in recent versions
}

func (c *SQLiteBoolConverter) DBToDomain(db interface{}) (bool, error) {
	if db == nil {
		return false, nil
	}

	switch v := db.(type) {
	case bool:
		return v, nil
	case int64:
		return v != 0, nil
	case int:
		return v != 0, nil
	case string:
		return v == "true" || v == "1", nil
	default:
		return false, NewConversionError("expected bool, int, or string", db)
	}
}

// Conversion errors
type ConversionError struct {
	Message string
	Value   interface{}
}

func NewConversionError(message string, value interface{}) *ConversionError {
	return &ConversionError{
		Message: message,
		Value:   value,
	}
}

func (e *ConversionError) Error() string {
	return e.Message
}

// Factory functions to create converters for different databases

// NewUUIDConverter creates UUID converter for specified database
func NewUUIDConverter(database string) UUIDConverter {
	switch database {
	case "sqlite":
		return NewSQLiteUUIDConverter()
	case "postgres":
		return NewPostgresUUIDConverter()
	case "mysql":
		return NewMySQLUUIDConverter()
	default:
		return NewSQLiteUUIDConverter() // Default to SQLite
	}
}

// NewTimeConverter creates time converter for specified database
func NewTimeConverter(database string) TimeConverter {
	switch database {
	case "sqlite", "postgres", "mysql":
		return NewSQLiteTimeConverter()
	default:
		return NewSQLiteTimeConverter()
	}
}

// NewBoolConverter creates boolean converter for specified database
func NewBoolConverter(database string) BoolConverter {
	switch database {
	case "sqlite", "postgres", "mysql":
		return NewSQLiteBoolConverter()
	default:
		return NewSQLiteBoolConverter()
	}
}

// Default converters

// DefaultEmailConverter handles email conversion
type DefaultEmailConverter struct{}

func NewDefaultEmailConverter() *DefaultEmailConverter {
	return &DefaultEmailConverter{}
}

func (c *DefaultEmailConverter) DomainToDB(domain entities.Email) string {
	return domain.String()
}

func (c *DefaultEmailConverter) DBToDomain(db string) (entities.Email, error) {
	return entities.NewEmail(db)
}

// DefaultUsernameConverter handles username conversion
type DefaultUsernameConverter struct{}

func NewDefaultUsernameConverter() *DefaultUsernameConverter {
	return &DefaultUsernameConverter{}
}

func (c *DefaultUsernameConverter) DomainToDB(domain entities.Username) string {
	return domain.String()
}

func (c *DefaultUsernameConverter) DBToDomain(db string) (entities.Username, error) {
	return entities.NewUsername(db)
}

// DefaultPasswordHashConverter handles password hash conversion
type DefaultPasswordHashConverter struct{}

func NewDefaultPasswordHashConverter() *DefaultPasswordHashConverter {
	return &DefaultPasswordHashConverter{}
}

func (c *DefaultPasswordConverter) DomainToDB(domain entities.PasswordHash) string {
	return domain.String()
}

func (c *DefaultPasswordConverter) DBToDomain(db string) (entities.PasswordHash, error) {
	return entities.NewPasswordHash(db)
}

// DefaultUserStatusConverter handles user status conversion
type DefaultUserStatusConverter struct{}

func NewDefaultUserStatusConverter() *DefaultUserStatusConverter {
	return &DefaultUserStatusConverter{}
}

func (c *DefaultUserStatusConverter) DomainToDB(domain entities.UserStatus) string {
	return domain.String()
}

func (c *DefaultUserStatusConverter) DBToDomain(db string) (entities.UserStatus, error) {
	status := entities.UserStatus(db)
	if !status.IsValid() {
		return entities.UserStatusActive, NewConversionError("invalid user status", db)
	}
	return status, nil
}

// DefaultUserRoleConverter handles user role conversion
type DefaultUserRoleConverter struct{}

func NewDefaultUserRoleConverter() *DefaultUserRoleConverter {
	return &DefaultUserRoleConverter{}
}

func (c *DefaultUserRoleConverter) DomainToDB(domain entities.UserRole) string {
	return domain.String()
}

func (c *DefaultUserRoleConverter) DBToDomain(db string) (entities.UserRole, error) {
	role := entities.UserRole(db)
	if !role.IsValid() {
		return entities.UserRoleUser, NewConversionError("invalid user role", db)
	}
	return role, nil
}

// DefaultSessionTokenConverter handles session token conversion
type DefaultSessionTokenConverter struct{}

func NewDefaultSessionTokenConverter() *DefaultSessionTokenConverter {
	return &DefaultSessionTokenConverter{}
}

func (c *DefaultSessionTokenConverter) DomainToDB(domain entities.SessionToken) interface{} {
	return domain.UUID().String()
}

func (c *DefaultSessionTokenConverter) DBToDomain(db interface{}) (entities.SessionToken, error) {
	var tokenUUID uuid.UUID

	switch v := db.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return entities.SessionToken{}, NewConversionError("invalid UUID format", db)
		}
		tokenUUID = parsed
	case uuid.UUID:
		tokenUUID = v
	case []byte:
		parsed, err := uuid.FromBytes(v)
		if err != nil {
			return entities.SessionToken{}, NewConversionError("invalid UUID bytes", db)
		}
		tokenUUID = parsed
	default:
		return entities.SessionToken{}, NewConversionError("expected string, UUID, or bytes", db)
	}

	return entities.SessionToken(tokenUUID), nil
}

// Helper functions

// ConvertToGeneric converts any type to interface{} for database operations
func ConvertToGeneric(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	return value
}

// SafeString safely converts interface{} to string
func SafeString(value interface{}) string {
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

// SafeTime safely converts interface{} to time.Time
func SafeTime(value interface{}) time.Time {
	if value == nil {
		return time.Time{}
	}

	if t, ok := value.(time.Time); ok {
		return t
	}

	return time.Time{}
}
