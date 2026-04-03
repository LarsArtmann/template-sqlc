package postgres

import (
	"context"
	"database/sql"
	stderrors "errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/mappers"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/validation"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
	"github.com/LarsArtmann/template-sqlc/pkg/errors"
)

// PostgresUserRepository implements UserRepository for PostgreSQL
// This adapts PostgreSQL-specific types to domain interfaces
type PostgresUserRepository struct {
	pool       *pgxpool.Pool
	mapper     *mappers.UserMapper
	converters *converters.ConverterSet
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(pool *pgxpool.Pool) repositories.UserRepository {
	return &PostgresUserRepository{
		pool:   pool,
		mapper: mappers.NewUserMapper(),
		converters: &converters.ConverterSet{
			UUID:     converters.NewUUIDConverter("postgres"),
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

// convertToModel converts a domain user to the database model and returns an error on failure
func (r *PostgresUserRepository) convertToModel(user *entities.User) (interface{}, error) {
	model, err := r.mapper.PostgresUserFromDomain(user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user %s: %w", user.ID(), err)
	}
	return model, nil
}

// Create saves a new user to PostgreSQL
func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	_, err := r.convertToModel(user)
	if err != nil {
		return err
	}

	// This would use actual generated sqlc code for PostgreSQL
	// Example:
	// _, err := r.queries.CreateUser(ctx, postgresUser.(postgres.CreateUserParams))
	// return errors.NewDatabaseError("failed to create user", err)

	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// GetByID retrieves a user by ID from PostgreSQL
func (r *PostgresUserRepository) GetByID(
	ctx context.Context,
	id entities.UserID,
) (*entities.User, error) {
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
func (r *PostgresUserRepository) GetByUUID(
	ctx context.Context,
	uuid entities.UuID,
) (*entities.User, error) {
	_, err := r.getByUUID(ctx, uuid)
	return nil, err
}

// GetByEmail retrieves a user by email from PostgreSQL
func (r *PostgresUserRepository) GetByEmail(
	ctx context.Context,
	email entities.Email,
) (*entities.User, error) {
	_, err := r.getByEmail(ctx, email)
	return nil, err
}

// GetByUsername retrieves a user by username from PostgreSQL
func (r *PostgresUserRepository) GetByUsername(
	ctx context.Context,
	username entities.Username,
) (*entities.User, error) {
	_, err := r.getByUsername(ctx, username)
	return nil, err
}

func (r *PostgresUserRepository) getByUUID(ctx context.Context, uuid entities.UuID) (struct{}, error) {
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

func (r *PostgresUserRepository) getByEmail(ctx context.Context, email entities.Email) (struct{}, error) {
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

func (r *PostgresUserRepository) getByUsername(ctx context.Context, username entities.Username) (struct{}, error) {
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// Update updates an existing user in PostgreSQL
func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	_, err := r.convertToModel(user)
	if err != nil {
		return err
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
func (r *PostgresUserRepository) List(
	ctx context.Context,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	if err := validation.ValidatePagination(limit, offset); err != nil {
		return nil, err
	}

	// Query database
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// Search searches users by query in PostgreSQL using FTS
func (r *PostgresUserRepository) Search(
	ctx context.Context,
	query string,
	status entities.UserStatus,
	limit int,
) ([]*entities.User, error) {
	if err := validation.ValidateSearchQuery(query, limit); err != nil {
		return nil, err
	}

	// Use PostgreSQL's tsvector search
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// SearchByTags searches users by tags in PostgreSQL using GIN index
func (r *PostgresUserRepository) SearchByTags(
	ctx context.Context,
	tags []string,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	if err := validation.ValidateTags(tags); err != nil {
		return nil, err
	}

	// Use PostgreSQL's array operations with GIN index
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// CountByStatus counts users by status in PostgreSQL
func (r *PostgresUserRepository) CountByStatus(
	ctx context.Context,
) (map[entities.UserStatus]int64, error) {
	// Query counts by status using PostgreSQL's GROUP BY
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// GetStats retrieves user statistics from PostgreSQL
func (r *PostgresUserRepository) GetStats(ctx context.Context) (*entities.UserStats, error) {
	// Query stats using PostgreSQL's aggregate functions
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// VerifyCredentials verifies user credentials in PostgreSQL
func (r *PostgresUserRepository) VerifyCredentials(
	ctx context.Context,
	email entities.Email,
	password entities.PasswordHash,
) (*entities.User, error) {
	// Query user by email and verify password
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// UpdatePassword updates user password in PostgreSQL
func (r *PostgresUserRepository) UpdatePassword(
	ctx context.Context,
	id entities.UserID,
	password entities.PasswordHash,
) error {
	// Update password
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// MarkVerified marks user as verified in PostgreSQL
func (r *PostgresUserRepository) MarkVerified(ctx context.Context, id entities.UserID) error {
	// Mark user as verified using PostgreSQL's UPDATE
	panic("implement me: use actual sqlc generated code for PostgreSQL")
}

// validateAndUpdateStatus validates and updates user status
func (r *PostgresUserRepository) validateAndUpdateStatus(
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
func (r *PostgresUserRepository) validateAndUpdateRole(
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

// ChangeStatus changes user status in PostgreSQL
func (r *PostgresUserRepository) ChangeStatus(
	ctx context.Context,
	id entities.UserID,
	status entities.UserStatus,
) error {
	return r.validateAndUpdateStatus(ctx, id, status, func() error {
		panic("implement me: use actual sqlc generated code for PostgreSQL")
	})
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
func (r *PostgresUserRepository) ChangeRole(
	ctx context.Context,
	id entities.UserID,
	role entities.UserRole,
) error {
	return r.validateAndUpdateRole(ctx, id, role, func() error {
		panic("implement me: use actual sqlc generated code for PostgreSQL")
	})
}

// Helper methods

// handlePostgresError converts PostgreSQL errors to domain errors
func (r *PostgresUserRepository) handlePostgresError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Check for common PostgreSQL error types
	switch {
	case stderrors.Is(err, sql.ErrNoRows):
		return entities.ErrUserNotFound
	case isUniqueViolationError(err):
		return entities.ErrUserAlreadyExists
	case isForeignKeyViolationError(err):
		return errors.NewValidationError("foreign_key", "referenced entity does not exist")
	case isCheckViolationError(err):
		return errors.NewValidationError("check_constraint", "check constraint violated")
	default:
		return errors.NewDatabaseError(operation+" failed", err)
	}
}
