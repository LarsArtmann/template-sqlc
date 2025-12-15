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

// SQLiteUserRepository implements UserRepository for SQLite
// This adapts SQLite-specific types to domain interfaces
type SQLiteUserRepository struct {
	db         *sql.DB
	mapper     mappers.UserMapper
	converters *ConverterSet
}

// ConverterSet holds all type converters for SQLite
type ConverterSet struct {
	UUID     converters.UUIDConverter
	Time     converters.TimeConverter
	Bool     converters.BoolConverter
	Email    converters.DefaultEmailConverter
	Username converters.DefaultUsernameConverter
	Password converters.DefaultPasswordHashConverter
	Status   converters.DefaultUserStatusConverter
	Role     converters.DefaultUserRoleConverter
}

// NewSQLiteUserRepository creates a new SQLite user repository
func NewSQLiteUserRepository(db *sql.DB) repositories.UserRepository {
	return &SQLiteUserRepository{
		db: db,
		converters: &ConverterSet{
			UUID:     converters.NewSQLiteUUIDConverter(),
			Time:     converters.NewTimeConverter("sqlite"),
			Bool:     converters.NewBoolConverter("sqlite"),
			Email:    converters.NewDefaultEmailConverter(),
			Username: converters.NewDefaultUsernameConverter(),
			Password: converters.NewDefaultPasswordHashConverter(),
			Status:   converters.NewDefaultUserStatusConverter(),
			Role:     converters.NewDefaultUserRoleConverter(),
		},
	}
}

// Create saves a new user to SQLite
func (r *SQLiteUserRepository) Create(ctx context.Context, user *entities.User) error {
	// Convert domain entity to SQLite model
	sqliteUser, err := mappers.SQLiteUserFromDomain(user)
	if err != nil {
		return fmt.Errorf("failed to convert user: %w", err)
	}

	// This would use the actual generated sqlc code
	// Example:
	// _, err := r.queries.CreateUser(ctx, sqliteUser.(sqlite.CreateUserParams))
	// return errors.NewDatabaseError("failed to create user", err)

	panic("implement me: use actual sqlc generated code")
}

// GetByID retrieves a user by ID from SQLite
func (r *SQLiteUserRepository) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	// This would use the actual generated sqlc code
	// Example:
	// sqliteUser, err := r.queries.GetUserByID(ctx, int64(id))
	// if err != nil {
	//     if err == sql.ErrNoRows {
	//         return nil, entities.ErrUserNotFound
	//     }
	//     return nil, errors.NewDatabaseError("failed to get user", err)
	// }
	// return mappers.DomainUserFromSQLite(sqliteUser)

	panic("implement me: use actual sqlc generated code")
}

// GetByUUID retrieves a user by UUID from SQLite
func (r *SQLiteUserRepository) GetByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	// Similar implementation for UUID lookup
	panic("implement me: use actual sqlc generated code")
}

// GetByEmail retrieves a user by email from SQLite
func (r *SQLiteUserRepository) GetByEmail(ctx context.Context, email entities.Email) (*entities.User, error) {
	// Similar implementation for email lookup
	panic("implement me: use actual sqlc generated code")
}

// GetByUsername retrieves a user by username from SQLite
func (r *SQLiteUserRepository) GetByUsername(ctx context.Context, username entities.Username) (*entities.User, error) {
	// Similar implementation for username lookup
	panic("implement me: use actual sqlc generated code")
}

// Update updates an existing user in SQLite
func (r *SQLiteUserRepository) Update(ctx context.Context, user *entities.User) error {
	// Convert domain entity to SQLite model
	sqliteUser, err := mappers.SQLiteUserFromDomain(user)
	if err != nil {
		return fmt.Errorf("failed to convert user: %w", err)
	}

	// Update in database
	panic("implement me: use actual sqlc generated code")
}

// Delete soft deletes a user from SQLite
func (r *SQLiteUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	// Soft delete by changing status
	// panic("implement me: use actual sqlc generated code")

	// For now, implement as hard delete
	// _, err := r.queries.DeleteUser(ctx, int64(id))
	// return errors.NewDatabaseError("failed to delete user", err)

	return nil
}

// List retrieves users with pagination from SQLite
func (r *SQLiteUserRepository) List(ctx context.Context, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	// Validate pagination parameters
	if limit <= 0 || limit > 1000 {
		return nil, errors.NewValidationError("limit", "must be between 1 and 1000")
	}
	if offset < 0 {
		return nil, errors.NewValidationError("offset", "must be non-negative")
	}

	// Query database
	panic("implement me: use actual sqlc generated code")
}

// Search searches users by query in SQLite
func (r *SQLiteUserRepository) Search(ctx context.Context, query string, status entities.UserStatus, limit int) ([]*entities.User, error) {
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

	// Search database
	panic("implement me: use actual sqlc generated code")
}

// SearchByTags searches users by tags in SQLite
func (r *SQLiteUserRepository) SearchByTags(ctx context.Context, tags []string, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	// Validate tags
	if len(tags) == 0 {
		return nil, errors.NewValidationError("tags", "cannot be empty")
	}
	if len(tags) > 10 {
		return nil, errors.NewValidationError("tags", "cannot exceed 10 tags")
	}

	// Search by tags
	panic("implement me: use actual sqlc generated code")
}

// CountByStatus counts users by status in SQLite
func (r *SQLiteUserRepository) CountByStatus(ctx context.Context) (map[entities.UserStatus]int64, error) {
	// Query counts by status
	panic("implement me: use actual sqlc generated code")
}

// GetStats retrieves user statistics from SQLite
func (r *SQLiteUserRepository) GetStats(ctx context.Context) (*entities.UserStats, error) {
	// Query stats
	panic("implement me: use actual sqlc generated code")
}

// VerifyCredentials verifies user credentials in SQLite
func (r *SQLiteUserRepository) VerifyCredentials(ctx context.Context, email entities.Email, password entities.PasswordHash) (*entities.User, error) {
	// Query user by email and verify password
	panic("implement me: use actual sqlc generated code")
}

// UpdatePassword updates user password in SQLite
func (r *SQLiteUserRepository) UpdatePassword(ctx context.Context, id entities.UserID, password entities.PasswordHash) error {
	// Update password
	panic("implement me: use actual sqlc generated code")
}

// MarkVerified marks user as verified in SQLite
func (r *SQLiteUserRepository) MarkVerified(ctx context.Context, id entities.UserID) error {
	// Mark user as verified
	panic("implement me: use actual sqlc generated code")
}

// ChangeStatus changes user status in SQLite
func (r *SQLiteUserRepository) ChangeStatus(ctx context.Context, id entities.UserID, status entities.UserStatus) error {
	// Validate status
	if !status.IsValid() {
		return errors.NewValidationError("status", "invalid user status")
	}

	// Update status
	panic("implement me: use actual sqlc generated code")
}

// Activate activates a user in SQLite
func (r *SQLiteUserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user in SQLite
func (r *SQLiteUserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user in SQLite
func (r *SQLiteUserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role in SQLite
func (r *SQLiteUserRepository) ChangeRole(ctx context.Context, id entities.UserID, role entities.UserRole) error {
	// Validate role
	if !role.IsValid() {
		return errors.NewValidationError("role", "invalid user role")
	}

	// Update role
	panic("implement me: use actual sqlc generated code")
}

// Helper methods

// handleError converts database errors to domain errors
func (r *SQLiteUserRepository) handleError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Check for common error types
	switch {
	case err == sql.ErrNoRows:
		return entities.ErrUserNotFound
	case isUniqueConstraintError(err):
		return entities.ErrUserAlreadyExists
	default:
		return errors.NewDatabaseError(fmt.Sprintf("%s failed", operation), err)
	}
}

// isUniqueConstraintError checks if error is a unique constraint violation
func isUniqueConstraintError(err error) bool {
	// This would check for SQLite-specific constraint errors
	// Example: check if error message contains "UNIQUE constraint failed"
	return err != nil &&
		(fmt.Sprintf("%s", err) == "UNIQUE constraint failed" ||
			fmt.Sprintf("%s", err) == "column ... is not unique")
}
