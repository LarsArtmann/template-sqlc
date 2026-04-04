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
// This isolates domain entities from database-specific type handling

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

// Default implementations

// SQLiteUUIDConverter handles UUID conversion for SQLite (stores as string).
type SQLiteUUIDConverter struct{}

func NewSQLiteUUIDConverter() *SQLiteUUIDConverter {
	return &SQLiteUUIDConverter{}
}

func (c *SQLiteUUIDConverter) DomainToDB(domain uuid.UUID) any {
	if domain == uuid.Nil {
		return nil
	}

	return domain.String()
}

func (c *SQLiteUUIDConverter) DBToDomain(db any) (uuid.UUID, error) {
	if db == nil {
		return uuid.Nil, nil
	}

	str, ok := db.(string)
	if !ok {
		return uuid.Nil, NewConversionError("expected string for UUID", db)
	}

	return uuid.Parse(str)
}

// PostgresUUIDConverter handles UUID conversion for PostgreSQL (stores as UUID type).
type PostgresUUIDConverter struct{}

func NewPostgresUUIDConverter() *PostgresUUIDConverter {
	return &PostgresUUIDConverter{}
}

func (c *PostgresUUIDConverter) DomainToDB(domain uuid.UUID) any {
	return domain // PostgreSQL handles UUID natively
}

func (c *PostgresUUIDConverter) DBToDomain(db any) (uuid.UUID, error) {
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

// MySQLUUIDConverter handles UUID conversion for MySQL (stores as binary).
type MySQLUUIDConverter struct{}

func NewMySQLUUIDConverter() *MySQLUUIDConverter {
	return &MySQLUUIDConverter{}
}

func (c *MySQLUUIDConverter) DomainToDB(domain uuid.UUID) any {
	return domain[:] // Convert to byte slice
}

func (c *MySQLUUIDConverter) DBToDomain(db any) (uuid.UUID, error) {
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

// SQLiteTimeConverter handles time conversion for SQLite.
type SQLiteTimeConverter struct{}

func NewSQLiteTimeConverter() *SQLiteTimeConverter {
	return &SQLiteTimeConverter{}
}

func (c *SQLiteTimeConverter) DomainToDB(domain time.Time) any {
	if domain.IsZero() {
		return nil
	}

	return domain
}

func (c *SQLiteTimeConverter) DBToDomain(db any) (time.Time, error) {
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

// SQLiteBoolConverter handles boolean conversion for SQLite.
type SQLiteBoolConverter struct{}

func NewSQLiteBoolConverter() *SQLiteBoolConverter {
	return &SQLiteBoolConverter{}
}

func (c *SQLiteBoolConverter) DomainToDB(domain bool) any {
	return domain // SQLite supports boolean natively in recent versions
}

func (c *SQLiteBoolConverter) DBToDomain(db any) (bool, error) {
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

// Conversion errors.
type ConversionError struct {
	Message string
	Value   any
}

func NewConversionError(message string, value any) *ConversionError {
	return &ConversionError{
		Message: message,
		Value:   value,
	}
}

func (e *ConversionError) Error() string {
	return e.Message
}

// Factory functions to create converters for different databases

// NewUUIDConverter creates UUID converter for specified database.
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

// NewTimeConverter creates time converter for specified database.
func NewTimeConverter(_ string) TimeConverter {
	return NewSQLiteTimeConverter()
}

// NewBoolConverter creates boolean converter for specified database.
func NewBoolConverter(_ string) BoolConverter {
	return NewSQLiteBoolConverter()
}

// Default converters

// DefaultEmailConverter handles email conversion.
type DefaultEmailConverter struct{}

func NewDefaultEmailConverter() *DefaultEmailConverter {
	return &DefaultEmailConverter{}
}

func (c *DefaultEmailConverter) DomainToDB(domain entities.Email) string {
	return domain.String()
}

func (c *DefaultEmailConverter) DBToDomain(db string) (entities.Email, error) {
	return convertSimpleValue(db, entities.NewEmail)
}

// DefaultUsernameConverter handles username conversion.
type DefaultUsernameConverter struct{}

func NewDefaultUsernameConverter() *DefaultUsernameConverter {
	return &DefaultUsernameConverter{}
}

func (c *DefaultUsernameConverter) DomainToDB(domain entities.Username) string {
	return domain.String()
}

func (c *DefaultUsernameConverter) DBToDomain(db string) (entities.Username, error) {
	return convertSimpleValue(db, entities.NewUsername)
}

// DefaultPasswordHashConverter handles password hash conversion.
type DefaultPasswordHashConverter struct{}

func NewDefaultPasswordHashConverter() *DefaultPasswordHashConverter {
	return &DefaultPasswordHashConverter{}
}

func (c *DefaultPasswordHashConverter) DomainToDB(domain entities.PasswordHash) string {
	return domain.String()
}

func (c *DefaultPasswordHashConverter) DBToDomain(db string) (entities.PasswordHash, error) {
	return convertSimpleValue(db, entities.NewPasswordHash)
}

// DefaultUserStatusConverter handles user status conversion.
type DefaultUserStatusConverter struct{}

func NewDefaultUserStatusConverter() *DefaultUserStatusConverter {
	return &DefaultUserStatusConverter{}
}

func (c *DefaultUserStatusConverter) DomainToDB(domain entities.UserStatus) string {
	return domain.String()
}

func (c *DefaultUserStatusConverter) DBToDomain(db string) (entities.UserStatus, error) {
	return convertEnumString(db, entities.UserStatusActive, "user status")
}

// convertEnumString converts a string to an enum type with validation.
type enum interface {
	~string
	IsValid() bool
}

func convertEnumString[T enum](db string, defaultVal T, typeName string) (T, error) {
	val := T(db)
	if !val.IsValid() {
		return defaultVal, NewConversionError("invalid "+typeName, db)
	}

	return val, nil
}

// stringConstructor is a function type for simple string-to-value conversion.
type stringConstructor[T any] func(string) (T, error)

// convertSimpleValue converts a database string using a simple constructor function.
func convertSimpleValue[T any](db string, constructor stringConstructor[T]) (T, error) {
	return constructor(db)
}

// DefaultUserRoleConverter handles user role conversion.
type DefaultUserRoleConverter struct{}

func NewDefaultUserRoleConverter() *DefaultUserRoleConverter {
	return &DefaultUserRoleConverter{}
}

func (c *DefaultUserRoleConverter) DomainToDB(domain entities.UserRole) string {
	return domain.String()
}

func (c *DefaultUserRoleConverter) DBToDomain(db string) (entities.UserRole, error) {
	return convertEnumString(db, entities.UserRoleUser, "user role")
}

// DefaultSessionTokenConverter handles session token conversion.
type DefaultSessionTokenConverter struct{}

func NewDefaultSessionTokenConverter() *DefaultSessionTokenConverter {
	return &DefaultSessionTokenConverter{}
}

func (c *DefaultSessionTokenConverter) DomainToDB(domain entities.SessionToken) any {
	return domain.UUID().String()
}

func (c *DefaultSessionTokenConverter) DBToDomain(db any) (entities.SessionToken, error) {
	var tokenUUID uuid.UUID

	switch v := db.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return entities.SessionToken{}, NewConversionError(
				fmt.Sprintf("invalid UUID format for token %s: %v", tokenUUID, err),
				db,
			)
		}

		tokenUUID = parsed
	case uuid.UUID:
		tokenUUID = v
	case []byte:
		parsed, err := uuid.FromBytes(v)
		if err != nil {
			return entities.SessionToken{}, NewConversionError(
				fmt.Sprintf("invalid UUID bytes for token %s: %v", tokenUUID, err),
				db,
			)
		}

		tokenUUID = parsed
	default:
		return entities.SessionToken{}, NewConversionError(
			fmt.Sprintf("expected string, UUID, or bytes for token %s", tokenUUID),
			db,
		)
	}

	return entities.SessionToken(tokenUUID), nil
}

// ConverterSet holds all type converters for user repository operations.
type ConverterSet struct {
	UUID     UUIDConverter
	Time     TimeConverter
	Bool     BoolConverter
	Email    EmailConverter
	Username UsernameConverter
	Password PasswordHashConverter
	Status   UserStatusConverter
	Role     UserRoleConverter
}

// Helper functions

// ConvertToGeneric converts any type to interface{} for database operations.
func ConvertToGeneric(value any) any {
	if value == nil {
		return nil
	}

	return value
}

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

// SafeTime safely converts interface{} to time.Time.
func SafeTime(value any) time.Time {
	if value == nil {
		return time.Time{}
	}

	if t, ok := value.(time.Time); ok {
		return t
	}

	return time.Time{}
}
