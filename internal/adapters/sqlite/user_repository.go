package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/mappers"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/validation"
	sqlitedb "github.com/LarsArtmann/template-sqlc/internal/db/sqlite"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
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
	UUID         converters.UUIDConverter
	Time         converters.TimeConverter
	Bool         converters.BoolConverter
	Email        *converters.DefaultEmailConverter
	Username     *converters.DefaultUsernameConverter
	Password     *converters.DefaultPasswordHashConverter
	Status       *converters.DefaultUserStatusConverter
	Role         *converters.DefaultUserRoleConverter
	SessionToken *converters.DefaultSessionTokenConverter
}

// NewSQLiteUserRepository creates a new SQLite user repository
func NewSQLiteUserRepository(db *sql.DB) repositories.UserRepository {
	return &SQLiteUserRepository{
		db: db,
		converters: &ConverterSet{
			UUID:     converters.NewUUIDConverter("sqlite"),
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
	_, err := r.mapUserToSQLite(user)
	if err != nil {
		return err
	}

	// This would use the actual generated sqlc code
	// Example:
	// _, err := r.queries.CreateUser(ctx, sqliteUser.(sqlite.CreateUserParams))
	// return errors.NewDatabaseError("failed to create user", err)

	panic("implement me: use actual sqlc generated code")
}

// GetByID retrieves a user by ID from SQLite
func (r *SQLiteUserRepository) GetByID(
	ctx context.Context,
	id entities.UserID,
) (*entities.User, error) {
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
func (r *SQLiteUserRepository) GetByUUID(
	ctx context.Context,
	uuid entities.UuID,
) (*entities.User, error) {
	_, err := r.getByUUID(ctx, uuid)
	return nil, err
}

// GetByEmail retrieves a user by email from SQLite
func (r *SQLiteUserRepository) GetByEmail(
	ctx context.Context,
	email entities.Email,
) (*entities.User, error) {
	_, err := r.getByEmail(ctx, email)
	return nil, err
}

// GetByUsername retrieves a user by username from SQLite
func (r *SQLiteUserRepository) GetByUsername(
	ctx context.Context,
	username entities.Username,
) (*entities.User, error) {
	_, err := r.getByUsername(ctx, username)
	return nil, err
}

func (r *SQLiteUserRepository) getByUUID(ctx context.Context, uuid entities.UuID) (struct{}, error) {
	panic("implement me: use actual sqlc generated code")
}

func (r *SQLiteUserRepository) getByEmail(ctx context.Context, email entities.Email) (struct{}, error) {
	panic("implement me: use actual sqlc generated code")
}

func (r *SQLiteUserRepository) getByUsername(ctx context.Context, username entities.Username) (struct{}, error) {
	panic("implement me: use actual sqlc generated code")
}

// Update updates an existing user in SQLite
func (r *SQLiteUserRepository) Update(ctx context.Context, user *entities.User) error {
	_, err := r.mapUserToSQLite(user)
	if err != nil {
		return err
	}

	// Update in database
	panic("implement me: use actual sqlc generated code")
}

// mapUserToSQLite converts a domain User entity to a SQLite model
func (r *SQLiteUserRepository) mapUserToSQLite(user *entities.User) (any, error) {
	sqliteUser, err := mappers.SQLiteUserFromDomain(user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user %s: %w", user.ID(), err)
	}
	return sqliteUser, nil
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
func (r *SQLiteUserRepository) List(
	ctx context.Context,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	if err := validation.ValidatePagination(limit, offset); err != nil {
		return nil, err
	}

	// Query database
	panic("implement me: use actual sqlc generated code")
}

// Search searches users by query in SQLite
func (r *SQLiteUserRepository) Search(
	ctx context.Context,
	query string,
	status entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	if err := validation.ValidateSearchQuery(query, limit); err != nil {
		return nil, err
	}

	// Search database
	panic("implement me: use actual sqlc generated code")
}

// SearchByTags searches users by tags in SQLite
func (r *SQLiteUserRepository) SearchByTags(
	ctx context.Context,
	tags []string,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	if err := validation.ValidateTags(tags); err != nil {
		return nil, err
	}

	// Search by tags
	panic("implement me: use actual sqlc generated code")
}

// CountByStatus counts users by status in SQLite
func (r *SQLiteUserRepository) CountByStatus(
	ctx context.Context,
) (map[entities.UserStatus]int64, error) {
	// Query counts by status
	panic("implement me: use actual sqlc generated code")
}

// GetStats retrieves user statistics from SQLite
func (r *SQLiteUserRepository) GetStats(ctx context.Context) (*entities.UserStats, error) {
	// Query stats
	panic("implement me: use actual sqlc generated code")
}

// VerifyCredentials verifies user credentials in SQLite
func (r *SQLiteUserRepository) VerifyCredentials(
	ctx context.Context,
	email entities.Email,
	password entities.PasswordHash,
) (*entities.User, error) {
	// Query user by email and verify password
	panic("implement me: use actual sqlc generated code")
}

// UpdatePassword updates user password in SQLite
func (r *SQLiteUserRepository) UpdatePassword(
	ctx context.Context,
	id entities.UserID,
	password entities.PasswordHash,
) error {
	// Update password
	panic("implement me: use actual sqlc generated code")
}

// MarkVerified marks user as verified in SQLite
func (r *SQLiteUserRepository) MarkVerified(ctx context.Context, id entities.UserID) error {
	// Mark user as verified
	panic("implement me: use actual sqlc generated code")
}

// validateAndUpdateStatus validates and updates user status
func (r *SQLiteUserRepository) validateAndUpdateStatus(
	_ context.Context,
	_ entities.UserID,
	status entities.UserStatus,
	updateFn func() error,
) error {
	if err := validation.ValidateStatus(status); err != nil {
		return err
	}
	return updateFn()
}

// validateAndUpdateRole validates and updates user role
func (r *SQLiteUserRepository) validateAndUpdateRole(
	_ context.Context,
	_ entities.UserID,
	role entities.UserRole,
	updateFn func() error,
) error {
	if err := validation.ValidateRole(role); err != nil {
		return err
	}
	return updateFn()
}

// ChangeStatus changes user status in SQLite
func (r *SQLiteUserRepository) ChangeStatus(
	ctx context.Context,
	id entities.UserID,
	status entities.UserStatus,
) error {
	return r.validateAndUpdateStatus(ctx, id, status, func() error {
		panic("implement me: use actual sqlc generated code")
	})
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
func (r *SQLiteUserRepository) ChangeRole(
	ctx context.Context,
	id entities.UserID,
	role entities.UserRole,
) error {
	return r.validateAndUpdateRole(ctx, id, role, func() error {
		panic("implement me: use actual sqlc generated code")
	})
}

// Helper method

func (r *SQLiteUserRepository) handleError(err error, operation string) error {
	return sqlitedb.HandleDBError(err, operation, entities.ErrUserNotFound, entities.ErrUserAlreadyExists, sqlitedb.IsSQLiteUniqueConstraintError)
}
