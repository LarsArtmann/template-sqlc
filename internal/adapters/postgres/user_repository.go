package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/mappers"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
	"github.com/LarsArtmann/template-sqlc/pkg/errors"
)

// PostgresUserRepository implements UserRepository for PostgreSQL
// This adapts PostgreSQL-specific types to domain interfaces
type PostgresUserRepository struct {
	pool       *pgxpool.Pool
	mapper     mappers.UserMapper
	converters *PostgresConverterSet
}

// PostgresConverterSet holds all type converters for PostgreSQL
type PostgresConverterSet struct {
	UUID     converters.PostgresUUIDConverter
	Time     converters.TimeConverter
	Bool     converters.BoolConverter
	Email    converters.DefaultEmailConverter
	Username converters.DefaultUsernameConverter
	Password converters.DefaultPasswordHashConverter
	Status   converters.DefaultUserStatusConverter
	Role     converters.DefaultUserRoleConverter
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(pool *pgxpool.Pool) repositories.UserRepository {
	return &PostgresUserRepository{
		pool: pool,
		converters: &PostgresConverterSet{
			UUID:     converters.NewPostgresUUIDConverter(),
			Time:     converters.NewTimeConverter("postgres"),
			Bool:     converters.NewBoolConverter("postgres"),
			Email:    converters.NewDefaultEmailConverter(),
			Username: converters.NewDefaultUsernameConverter(),
			Password: converters.NewDefaultPasswordHashConverter(),
			Status:   converters.NewDefaultUserStatusConverter(),
			Role:     converters.NewDefaultUserRoleConverter(),
		},
	}
}

// Create saves a new user to PostgreSQL
func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	// Convert domain entity to PostgreSQL model
	postgresUser, err := mappers.PostgresUserFromDomain(user)
	if err != nil {
		return fmt.Errorf("failed to convert user: %w", err)
	}

	// This would use actual generated sqlc code for PostgreSQL
	// Example:
	// _, err := r.queries.CreateUser(ctx, postgresUser.(postgres.CreateUserParams))
	// return errors.NewDatabaseError("failed to create user", err)

	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// GetByID retrieves a user by ID from PostgreSQL
func (r *PostgresUserRepository) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	// This would use actual generated sqlc code for PostgreSQL
	// Example:
	// postgresUser, err := r.queries.GetUserByID(ctx, int64(id))
	// if err != nil {
	//     if err == sql.ErrNoRows {
	//         return nil, entities.ErrUserNotFound
	//     }
	//     return nil, errors.NewDatabaseError("failed to get user", err)
	// }
	// return mappers.DomainUserFromPostgres(postgresUser)

	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// GetByUUID retrieves a user by UUID from PostgreSQL
func (r *PostgresUserRepository) GetByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	// Convert string to UUID type
	uuidObj, err := r.converters.UUID.DBToDomain(uuid)
	if err != nil {
		return nil, errors.NewValidationError("uuid", "invalid UUID format")
	}

	// Query using UUID type
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// GetByEmail retrieves a user by email from PostgreSQL
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email entities.Email) (*entities.User, error) {
	// Convert to database format
	dbEmail := r.converters.Email.DomainToDB(email)

	// Query using case-insensitive search (CITEXT)
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// GetByUsername retrieves a user by username from PostgreSQL
func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username entities.Username) (*entities.User, error) {
	// Convert to database format
	dbUsername := r.converters.Username.DomainToDB(username)

	// Query using case-insensitive search (CITEXT)
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// Update updates an existing user in PostgreSQL
func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	// Convert domain entity to PostgreSQL model
	postgresUser, err := mappers.PostgresUserFromDomain(user)
	if err != nil {
		return fmt.Errorf("failed to convert user: %w", err)
	}

	// Update in database
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// Delete soft deletes a user from PostgreSQL
func (r *PostgresUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	// Soft delete by changing status
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// List retrieves users with pagination from PostgreSQL
func (r *PostgresUserRepository) List(ctx context.Context, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	// Validate pagination parameters
	if limit <= 0 || limit > 1000 {
		return nil, errors.NewValidationError("limit", "must be between 1 and 1000")
	}
	if offset < 0 {
		return nil, errors.NewValidationError("offset", "must be non-negative")
	}

	// Convert status to database format
	dbStatus := r.converters.Status.DomainToDB(status)

	// Query database
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// Search searches users by query in PostgreSQL using FTS
func (r *PostgresUserRepository) Search(ctx context.Context, query string, status entities.UserStatus, limit int) ([]*entities.User, error) {
	// Validate search query
	if len(query) == 0 {
		return nil, errors.NewValidationError("query", "cannot be empty")
	}
	if len(query) > 500 {
		return nil, errors.NewValidationError("query", "cannot exceed 500 characters")
	}
	if limit <= 0 || limit > 100 {
		return nil, errors.NewValidationError("limit", "must be between 1 and 100")
	}

	// Use PostgreSQL's tsvector search
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// SearchByTags searches users by tags in PostgreSQL using GIN index
func (r *PostgresUserRepository) SearchByTags(ctx context.Context, tags []string, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	// Validate tags
	if len(tags) == 0 {
		return nil, errors.NewValidationError("tags", "cannot be empty")
	}
	if len(tags) > 10 {
		return nil, errors.NewValidationError("tags", "cannot exceed 10 tags")
	}

	// Use PostgreSQL's array operations with GIN index
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// CountByStatus counts users by status in PostgreSQL
func (r *PostgresUserRepository) CountByStatus(ctx context.Context) (map[entities.UserStatus]int64, error) {
	// Query counts by status using PostgreSQL's GROUP BY
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// GetStats retrieves user statistics from PostgreSQL
func (r *PostgresUserRepository) GetStats(ctx context.Context) (*entities.UserStats, error) {
	// Query stats using PostgreSQL's aggregate functions
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// VerifyCredentials verifies user credentials in PostgreSQL
func (r *PostgresUserRepository) VerifyCredentials(ctx context.Context, email entities.Email, password entities.PasswordHash) (*entities.User, error) {
	// Convert to database format
	dbEmail := r.converters.Email.DomainToDB(email)
	dbPassword := r.converters.Password.DomainToDB(password)

	// Query user by email and verify password
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// UpdatePassword updates user password in PostgreSQL
func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, id entities.UserID, password entities.PasswordHash) error {
	// Convert to database format
	dbPassword := r.converters.Password.DomainToDB(password)

	// Update password
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// MarkVerified marks user as verified in PostgreSQL
func (r *PostgresUserRepository) MarkVerified(ctx context.Context, id entities.UserID) error {
	// Mark user as verified using PostgreSQL's UPDATE
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// ChangeStatus changes user status in PostgreSQL
func (r *PostgresUserRepository) ChangeStatus(ctx context.Context, id entities.UserID, status entities.UserStatus) error {
	// Validate status
	if !status.IsValid() {
		return errors.NewValidationError("status", "invalid user status")
	}

	// Convert to database format
	dbStatus := r.converters.Status.DomainToDB(status)

	// Update status
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// Activate activates a user in PostgreSQL
func (r *PostgresUserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user in PostgreSQL
func (r *PostgresUserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user in PostgreSQL
func (r *PostgresUserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role in PostgreSQL
func (r *PostgresUserRepository) ChangeRole(ctx context.Context, id entities.UserID, role entities.UserRole) error {
	// Validate role
	if !role.IsValid() {
		return errors.NewValidationError("role", "invalid user role")
	}

	// Convert to database format
	dbRole := r.converters.Role.DomainToDB(role)

	// Update role
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// Helper methods

// handlePostgresError converts PostgreSQL errors to domain errors
func (r *PostgresUserRepository) handlePostgresError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Check for common PostgreSQL error types
	switch {
	case err == sql.ErrNoRows:
		return entities.ErrUserNotFound
	case isUniqueViolationError(err):
		return entities.ErrUserAlreadyExists
	case isForeignKeyViolationError(err):
		return errors.NewValidationError("foreign_key", "referenced entity does not exist")
	case isCheckViolationError(err):
		return errors.NewValidationError("check_constraint", "check constraint violated")
	default:
		return errors.NewDatabaseError(fmt.Sprintf("%s failed", operation), err)
	}
}

// isUniqueViolationError checks for PostgreSQL unique constraint violation
func isUniqueViolationError(err error) bool {
	// PostgreSQL error code 23505 for unique violation
	if pgErr, ok := err.(interface{ Code() string }); ok {
		return pgErr.Code() == "23505"
	}
	return false
}

// isForeignKeyViolationError checks for PostgreSQL foreign key violation
func isForeignKeyViolationError(err error) bool {
	// PostgreSQL error code 23503 for foreign key violation
	if pgErr, ok := err.(interface{ Code() string }); ok {
		return pgErr.Code() == "23503"
	}
	return false
}

// isCheckViolationError checks for PostgreSQL check constraint violation
func isCheckViolationError(err error) bool {
	// PostgreSQL error code 23514 for check constraint violation
	if pgErr, ok := err.(interface{ Code() string }); ok {
		return pgErr.Code() == "23514"
	}
	return false
}
