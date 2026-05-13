// Package converters provides type conversion utilities for database operations.
// It handles conversion between domain entities and database-specific types.
package converters

import (
	"fmt"
	"time"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/google/uuid"
)

// Database type constants for consistent database identification.
const (
	DbTypeSQLite   = "sqlite"
	DbTypePostgres = "postgres"
	DbTypeMySQL    = "mysql"
)

// TypeConverter handles database-specific type conversions
// This isolates domain entities from database-specific type handling.

// TypeConverter is a generic interface for domain <-> database type conversions.
type TypeConverter[Domain any, DB any] interface {
	DomainToDB(domain Domain) DB
	DBToDomain(db DB) (Domain, error)
}

// UUIDConverter handles UUID conversions between domain and database.
type UUIDConverter interface {
	TypeConverter[uuid.UUID, any]
}

// TimeConverter handles time conversions between domain and database.
type TimeConverter interface {
	TypeConverter[time.Time, any]
}

// BoolConverter handles boolean conversions between domain and database.
type BoolConverter interface {
	TypeConverter[bool, any]
}

// EmailConverter handles email conversions between domain and database.
type EmailConverter interface {
	TypeConverter[entities.Email, string]
}

// UsernameConverter handles username conversions between domain and database.
type UsernameConverter interface {
	TypeConverter[entities.Username, string]
}

// PasswordHashConverter handles password hash conversions between domain and database.
type PasswordHashConverter interface {
	TypeConverter[entities.PasswordHash, string]
}

// UserStatusConverter handles user status conversions between domain and database.
type UserStatusConverter interface {
	TypeConverter[entities.UserStatus, string]
}

// UserRoleConverter handles user role conversions between domain and database.
type UserRoleConverter interface {
	TypeConverter[entities.UserRole, string]
}

// SessionTokenConverter handles session token conversions between domain and database.
type SessionTokenConverter interface {
	TypeConverter[entities.SessionToken, any]
}

// SQLiteUUIDConverter handles UUID conversion for SQLite (stores as string).
type SQLiteUUIDConverter struct{}

// NewSQLiteUUIDConverter creates a new SQLiteUUIDConverter.
func NewSQLiteUUIDConverter() *SQLiteUUIDConverter { return &SQLiteUUIDConverter{} }

// DomainToDB converts a domain UUID to a SQLite-compatible string representation.
func (c *SQLiteUUIDConverter) DomainToDB(domain uuid.UUID) any {
	if domain == uuid.Nil {
		return nil
	}

	return domain.String()
}

// DBToDomain converts a SQLite UUID (string) to a domain UUID.
func (c *SQLiteUUIDConverter) DBToDomain(db any) (uuid.UUID, error) {
	if db == nil {
		return uuid.Nil, nil
	}

	return parseUUIDFromAny(db)
}

// PostgresUUIDConverter handles UUID conversion for PostgreSQL (stores as UUID type).
type PostgresUUIDConverter struct{}

// NewPostgresUUIDConverter creates a new PostgresUUIDConverter.
func NewPostgresUUIDConverter() *PostgresUUIDConverter { return &PostgresUUIDConverter{} }

// DomainToDB converts a domain UUID to a PostgreSQL UUID.
func (c *PostgresUUIDConverter) DomainToDB(domain uuid.UUID) any {
	return domain
}

// DBToDomain converts a PostgreSQL UUID to a domain UUID.
func (c *PostgresUUIDConverter) DBToDomain(db any) (uuid.UUID, error) {
	if db == nil {
		return uuid.Nil, nil
	}

	return parseUUIDFromAny(db)
}

// MySQLUUIDConverter handles UUID conversion for MySQL (stores as binary).
type MySQLUUIDConverter struct{}

// NewMySQLUUIDConverter creates a new MySQLUUIDConverter.
func NewMySQLUUIDConverter() *MySQLUUIDConverter { return &MySQLUUIDConverter{} }

// DomainToDB converts a domain UUID to MySQL binary format.
func (c *MySQLUUIDConverter) DomainToDB(domain uuid.UUID) any {
	return domain[:]
}

// DBToDomain converts MySQL binary UUID to a domain UUID.
func (c *MySQLUUIDConverter) DBToDomain(db any) (uuid.UUID, error) {
	if db == nil {
		return uuid.Nil, nil
	}

	return parseUUIDFromAny(db)
}

// SQLiteTimeConverter handles time conversion for SQLite.
type SQLiteTimeConverter struct{}

// NewSQLiteTimeConverter creates a new SQLiteTimeConverter.
func NewSQLiteTimeConverter() *SQLiteTimeConverter { return &SQLiteTimeConverter{} }

// DomainToDB converts a domain time.Time to a SQLite-compatible format.
func (c *SQLiteTimeConverter) DomainToDB(domain time.Time) any {
	if domain.IsZero() {
		return nil
	}

	return domain
}

// DBToDomain converts a SQLite time value to a domain time.Time.
func (c *SQLiteTimeConverter) DBToDomain(db any) (time.Time, error) {
	if db == nil {
		return time.Time{}, nil
	}

	if t, ok := db.(time.Time); ok {
		return t, nil
	}

	if str, ok := db.(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid time format: %w", err)
		}

		return parsedTime, nil
	}

	return time.Time{}, NewConversionError("expected time or string", db)
}

// SQLiteBoolConverter handles boolean conversion for SQLite.
type SQLiteBoolConverter struct{}

// NewSQLiteBoolConverter creates a new SQLiteBoolConverter.
func NewSQLiteBoolConverter() *SQLiteBoolConverter { return &SQLiteBoolConverter{} }

// DomainToDB converts a domain bool to a SQLite-compatible format.
func (c *SQLiteBoolConverter) DomainToDB(domain bool) any { return domain }

// DBToDomain converts a SQLite bool value to a domain bool.
func (c *SQLiteBoolConverter) DBToDomain(db any) (bool, error) {
	if db == nil {
		return false, nil
	}

	switch dbValue := db.(type) {
	case bool:
		return dbValue, nil
	case int64:
		return dbValue != 0, nil
	case int:
		return dbValue != 0, nil
	case string:
		return dbValue == "true" || dbValue == "1", nil
	default:
		return false, NewConversionError("expected bool, int, or string", db)
	}
}

// ConversionError represents a conversion error.
type ConversionError struct {
	Message string
	Value   any
}

// NewConversionError creates a new ConversionError with a message and value.
func NewConversionError(message string, value any) *ConversionError {
	return &ConversionError{Message: message, Value: value}
}

func (e *ConversionError) Error() string { return e.Message }

// Factory functions

// NewUUIDConverter creates a new UUIDConverter for the specified database type.
func NewUUIDConverter(database string) UUIDConverter {
	switch database {
	case DbTypeSQLite:
		return NewSQLiteUUIDConverter()
	case DbTypePostgres:
		return NewPostgresUUIDConverter()
	case DbTypeMySQL:
		return NewMySQLUUIDConverter()
	default:
		return NewSQLiteUUIDConverter()
	}
}

// NewTimeConverter creates a new TimeConverter for the specified database type.
func NewTimeConverter(_ string) TimeConverter { return NewSQLiteTimeConverter() }

// NewBoolConverter creates a new BoolConverter for the specified database type.
func NewBoolConverter(_ string) BoolConverter { return NewSQLiteBoolConverter() }

// Default converters

// DefaultEmailConverter handles email conversions.
type DefaultEmailConverter struct{}

// NewDefaultEmailConverter creates a new DefaultEmailConverter.
func NewDefaultEmailConverter() *DefaultEmailConverter { return &DefaultEmailConverter{} }

// DomainToDB converts a domain Email to a database string.
func (c *DefaultEmailConverter) DomainToDB(domain entities.Email) string { return domain.String() }

// DBToDomain converts a database string to a domain Email.
func (c *DefaultEmailConverter) DBToDomain(db string) (entities.Email, error) {
	return convertSimpleValue(db, entities.NewEmail)
}

// DefaultUsernameConverter handles username conversions.
type DefaultUsernameConverter struct{}

// NewDefaultUsernameConverter creates a new DefaultUsernameConverter.
func NewDefaultUsernameConverter() *DefaultUsernameConverter { return &DefaultUsernameConverter{} }

// DomainToDB converts a domain Username to a database string.
func (c *DefaultUsernameConverter) DomainToDB(domain entities.Username) string {
	return domain.String()
}

// DBToDomain converts a database string to a domain Username.
func (c *DefaultUsernameConverter) DBToDomain(db string) (entities.Username, error) {
	return convertSimpleValue(db, entities.NewUsername)
}

// DefaultPasswordHashConverter handles password hash conversions.
type DefaultPasswordHashConverter struct{}

// NewDefaultPasswordHashConverter creates a new DefaultPasswordHashConverter.
func NewDefaultPasswordHashConverter() *DefaultPasswordHashConverter {
	return &DefaultPasswordHashConverter{}
}

// DomainToDB converts a domain PasswordHash to a database string.
func (c *DefaultPasswordHashConverter) DomainToDB(domain entities.PasswordHash) string {
	return domain.String()
}

// DBToDomain converts a database string to a domain PasswordHash.
func (c *DefaultPasswordHashConverter) DBToDomain(db string) (entities.PasswordHash, error) {
	return convertSimpleValue(db, entities.NewPasswordHash)
}

// DefaultUserStatusConverter handles user status conversions.
type DefaultUserStatusConverter struct{}

// NewDefaultUserStatusConverter creates a new DefaultUserStatusConverter.
func NewDefaultUserStatusConverter() *DefaultUserStatusConverter {
	return &DefaultUserStatusConverter{}
}

// DomainToDB converts a domain UserStatus to a database string.
func (c *DefaultUserStatusConverter) DomainToDB(domain entities.UserStatus) string {
	return domain.String()
}

// DBToDomain converts a database string to a domain UserStatus.
func (c *DefaultUserStatusConverter) DBToDomain(db string) (entities.UserStatus, error) {
	return convertEnumString(db, entities.UserStatusActive, "user status")
}

// convertEnumString converts a string to an enum type with validation.
type enum interface {
	~string
	IsValid() bool
}

// convertEnumString converts enum values between domain and database representations.
//
//nolint:ireturn // Generic converters intentionally return type parameters
func convertEnumString[T enum](db string, defaultVal T, typeName string) (T, error) {
	val := T(db)
	if !val.IsValid() {
		return defaultVal, NewConversionError("invalid "+typeName, db)
	}

	return val, nil
}

type stringConstructor[T any] func(string) (T, error)

// convertSimpleValue converts simple values using a constructor function.
//
//nolint:ireturn // Generic converters intentionally return type parameters
func convertSimpleValue[T any](db string, constructor stringConstructor[T]) (T, error) {
	return constructor(db)
}

// DefaultUserRoleConverter handles user role conversions.
type DefaultUserRoleConverter struct{}

// NewDefaultUserRoleConverter creates a new DefaultUserRoleConverter.
func NewDefaultUserRoleConverter() *DefaultUserRoleConverter { return &DefaultUserRoleConverter{} }

// DomainToDB converts a domain UserRole to a database string.
func (c *DefaultUserRoleConverter) DomainToDB(domain entities.UserRole) string {
	return domain.String()
}

// DBToDomain converts a database string to a domain UserRole.
func (c *DefaultUserRoleConverter) DBToDomain(db string) (entities.UserRole, error) {
	return convertEnumString(db, entities.UserRoleUser, "user role")
}

// DefaultSessionTokenConverter handles session token conversions.
type DefaultSessionTokenConverter struct{}

// NewDefaultSessionTokenConverter creates a new DefaultSessionTokenConverter.
func NewDefaultSessionTokenConverter() *DefaultSessionTokenConverter {
	return &DefaultSessionTokenConverter{}
}

// DomainToDB converts a domain SessionToken to a database-compatible format.
func (c *DefaultSessionTokenConverter) DomainToDB(domain entities.SessionToken) any {
	return domain.UUID().String()
}

// DBToDomain converts a database token to a domain SessionToken.
func (c *DefaultSessionTokenConverter) DBToDomain(db any) (entities.SessionToken, error) {
	tokenUUID, err := parseUUIDFromAny(db)
	if err != nil {
		return entities.SessionToken{}, err
	}

	return entities.SessionToken(tokenUUID), nil
}

// parseUUIDFromAny parses a UUID from various database types.
func parseUUIDFromAny(dbValue any) (uuid.UUID, error) {
	switch parsed := dbValue.(type) {
	case uuid.UUID:
		return parsed, nil
	case string:
		parsedUUID, err := uuid.Parse(parsed)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid UUID string: %w", err)
		}

		return parsedUUID, nil
	case []byte:
		parsedUUID, err := uuid.FromBytes(parsed)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid UUID bytes: %w", err)
		}

		return parsedUUID, nil
	default:
		return uuid.Nil, NewConversionError("expected UUID, string, or bytes", dbValue)
	}
}

// ConverterSet holds all type converters for user repository operations.
type ConverterSet struct {
	UUID         UUIDConverter
	Time         TimeConverter
	Bool         BoolConverter
	Email        EmailConverter
	Username     UsernameConverter
	Password     PasswordHashConverter
	Status       UserStatusConverter
	Role         UserRoleConverter
	SessionToken SessionTokenConverter
}

// NewConverterSet creates a new ConverterSet for the specified database type.
func NewConverterSet(database string) *ConverterSet {
	return &ConverterSet{
		UUID:         NewUUIDConverter(database),
		Time:         NewTimeConverter(database),
		Bool:         NewBoolConverter(database),
		Email:        NewDefaultEmailConverter(),
		Username:     NewDefaultUsernameConverter(),
		Password:     NewDefaultPasswordHashConverter(),
		Status:       NewDefaultUserStatusConverter(),
		Role:         NewDefaultUserRoleConverter(),
		SessionToken: NewDefaultSessionTokenConverter(),
	}
}

// Helper functions

// SafeString safely converts interface{} to string.
func SafeString(value any) string {
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
