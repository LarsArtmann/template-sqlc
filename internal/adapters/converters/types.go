package converters

import (
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

func NewSQLiteUUIDConverter() *SQLiteUUIDConverter { return &SQLiteUUIDConverter{} }

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
	return parseUUIDFromAny(db)
}

// PostgresUUIDConverter handles UUID conversion for PostgreSQL (stores as UUID type).
type PostgresUUIDConverter struct{}

func NewPostgresUUIDConverter() *PostgresUUIDConverter { return &PostgresUUIDConverter{} }

func (c *PostgresUUIDConverter) DomainToDB(domain uuid.UUID) any {
	return domain
}

func (c *PostgresUUIDConverter) DBToDomain(db any) (uuid.UUID, error) {
	if db == nil {
		return uuid.Nil, nil
	}
	return parseUUIDFromAny(db)
}

// MySQLUUIDConverter handles UUID conversion for MySQL (stores as binary).
type MySQLUUIDConverter struct{}

func NewMySQLUUIDConverter() *MySQLUUIDConverter { return &MySQLUUIDConverter{} }

func (c *MySQLUUIDConverter) DomainToDB(domain uuid.UUID) any {
	return domain[:]
}

func (c *MySQLUUIDConverter) DBToDomain(db any) (uuid.UUID, error) {
	if db == nil {
		return uuid.Nil, nil
	}
	return parseUUIDFromAny(db)
}

// SQLiteTimeConverter handles time conversion for SQLite.
type SQLiteTimeConverter struct{}

func NewSQLiteTimeConverter() *SQLiteTimeConverter { return &SQLiteTimeConverter{} }

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

func NewSQLiteBoolConverter() *SQLiteBoolConverter { return &SQLiteBoolConverter{} }

func (c *SQLiteBoolConverter) DomainToDB(domain bool) any { return domain }

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

// ConversionError represents a conversion error.
type ConversionError struct {
	Message string
	Value   any
}

func NewConversionError(message string, value any) *ConversionError {
	return &ConversionError{Message: message, Value: value}
}

func (e *ConversionError) Error() string { return e.Message }

// Factory functions

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

func NewTimeConverter(_ string) TimeConverter { return NewSQLiteTimeConverter() }
func NewBoolConverter(_ string) BoolConverter { return NewSQLiteBoolConverter() }

// Default converters

type DefaultEmailConverter struct{}

func NewDefaultEmailConverter() *DefaultEmailConverter { return &DefaultEmailConverter{} }

func (c *DefaultEmailConverter) DomainToDB(domain entities.Email) string { return domain.String() }

func (c *DefaultEmailConverter) DBToDomain(db string) (entities.Email, error) {
	return convertSimpleValue(db, entities.NewEmail)
}

type DefaultUsernameConverter struct{}

func NewDefaultUsernameConverter() *DefaultUsernameConverter { return &DefaultUsernameConverter{} }

func (c *DefaultUsernameConverter) DomainToDB(domain entities.Username) string {
	return domain.String()
}

func (c *DefaultUsernameConverter) DBToDomain(db string) (entities.Username, error) {
	return convertSimpleValue(db, entities.NewUsername)
}

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

type stringConstructor[T any] func(string) (T, error)

func convertSimpleValue[T any](db string, constructor stringConstructor[T]) (T, error) {
	return constructor(db)
}

type DefaultUserRoleConverter struct{}

func NewDefaultUserRoleConverter() *DefaultUserRoleConverter { return &DefaultUserRoleConverter{} }

func (c *DefaultUserRoleConverter) DomainToDB(domain entities.UserRole) string {
	return domain.String()
}

func (c *DefaultUserRoleConverter) DBToDomain(db string) (entities.UserRole, error) {
	return convertEnumString(db, entities.UserRoleUser, "user role")
}

type DefaultSessionTokenConverter struct{}

func NewDefaultSessionTokenConverter() *DefaultSessionTokenConverter {
	return &DefaultSessionTokenConverter{}
}

func (c *DefaultSessionTokenConverter) DomainToDB(domain entities.SessionToken) any {
	return domain.UUID().String()
}

func (c *DefaultSessionTokenConverter) DBToDomain(db any) (entities.SessionToken, error) {
	tokenUUID, err := parseUUIDFromAny(db)
	if err != nil {
		return entities.SessionToken{}, err
	}
	return entities.SessionToken(tokenUUID), nil
}

// parseUUIDFromAny parses a UUID from various database types.
func parseUUIDFromAny(db any) (uuid.UUID, error) {
	switch v := db.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	case []byte:
		return uuid.FromBytes(v)
	default:
		return uuid.Nil, NewConversionError("expected UUID, string, or bytes", db)
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
