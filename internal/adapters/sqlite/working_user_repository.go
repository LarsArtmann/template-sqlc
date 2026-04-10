package sqlite

import (
	"context"
	"database/sql"
	stderrors "errors"
	"fmt"

	"github.com/LarsArtmann/template-sqlc/internal/adapters/converters"
	"github.com/LarsArtmann/template-sqlc/internal/adapters/mappers"
	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
	"github.com/LarsArtmann/template-sqlc/pkg/errors"
)

// WorkingSQLiteUserRepository is a simplified implementation that works without generated sqlc code
// This demonstrates the pattern while we fix the SQL syntax issues.
type WorkingSQLiteUserRepository struct {
	db         *sql.DB
	mapper     mappers.UserMapper
	converters *converters.ConverterSet
}

// NewWorkingSQLiteUserRepository creates a new working SQLite user repository.
func NewWorkingSQLiteUserRepository(db *sql.DB) repositories.UserRepository {
	return &WorkingSQLiteUserRepository{
		db:         db,
		mapper:     mappers.UserMapper{},
		converters: converters.NewConverterSet(converters.DbTypeSQLite),
	}
}

// Create creates a new user in SQLite.
func (r *WorkingSQLiteUserRepository) Create(ctx context.Context, user *entities.User) error {
	// For now, implement using raw SQL to avoid generated code dependency
	query := `
		INSERT INTO users (email, username, password_hash, first_name, last_name, status, role, is_verified, metadata, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Convert domain values to database-compatible types
	email := r.converters.Email.DomainToDB(user.Email())
	username := r.converters.Username.DomainToDB(user.Username())
	passwordHash := r.converters.Password.DomainToDB(
		entities.PasswordHash("placeholder"),
	) // Password should come from user
	firstName := user.FirstName().String()
	lastName := user.LastName().String()
	status := r.converters.Status.DomainToDB(user.Status())
	role := r.converters.Role.DomainToDB(user.Role())
	isVerified := r.converters.Bool.DomainToDB(user.IsVerified())

	// Convert metadata and tags to JSON
	metadataJSON := "{}"
	tagsJSON := "[]"

	result, err := r.db.ExecContext(
		ctx,
		query,
		email,
		username,
		passwordHash,
		firstName,
		lastName,
		status,
		role,
		isVerified,
		metadataJSON,
		tagsJSON,
	)
	if err != nil {
		return errors.NewDatabaseError(
			fmt.Sprintf("failed to create user %s: %v", username, err),
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(
			fmt.Sprintf("failed to check affected rows for user %s: %v", username, err),
			err,
		)
	}

	if rows == 0 {
		return errors.NewDatabaseError("no rows affected", stderrors.New("user creation failed"))
	}

	return nil
}

// GetByID retrieves a user by ID from SQLite.
func (r *WorkingSQLiteUserRepository) GetByID(
	ctx context.Context,
	id entities.UserID,
) (*entities.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, status, role, 
		       is_verified, metadata, tags, created_at, updated_at, last_login_at
		FROM users 
		WHERE id = ?
	`

	_ = &entities.User{} // This is wrong - need proper constructor

	var (
		email, username, passwordHash, firstName, lastName, status, role string
		isVerified                                                       bool
		metadataJSON, tagsJSON                                           string
		createdAt, updatedAt, lastLoginAt                                sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&id, &email, &username, &passwordHash, &firstName, &lastName, &status, &role,
		&isVerified, &metadataJSON, &tagsJSON, &createdAt, &updatedAt, &lastLoginAt,
	)
	if err != nil {
		if stderrors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user %s not found: %w", id, entities.ErrUserNotFound)
		}

		return nil, fmt.Errorf("failed to get user by ID %s: %w", id, err)
	}

	// This is a simplified example - proper implementation would:
	// 1. Convert database types to domain entities
	// 2. Use proper entity constructors
	// 3. Handle all field conversions

	// For now, return nil to show pattern
	return nil, fmt.Errorf("implementation in progress - user found with ID %s", id)
}

// GetByUUID retrieves a user by UUID from SQLite.
func (r *WorkingSQLiteUserRepository) GetByUUID(
	ctx context.Context,
	uuid entities.UuID,
) (*entities.User, error) {
	// Implementation following same pattern as GetByID
	stubPanic()

	return nil, nil
}

// GetByEmail retrieves a user by email from SQLite.
func (r *WorkingSQLiteUserRepository) GetByEmail(
	_ context.Context,
	_ entities.Email,
) (*entities.User, error) {
	stubPanic()

	return nil, nil
}

// GetByUsername retrieves a user by username from SQLite.
func (r *WorkingSQLiteUserRepository) GetByUsername(
	_ context.Context,
	_ entities.Username,
) (*entities.User, error) {
	stubPanic()

	return nil, nil
}

// Update updates an existing user in SQLite.
func (r *WorkingSQLiteUserRepository) Update(_ context.Context, _ *entities.User) error {
	stubPanic()

	return nil
}

// Delete soft deletes a user from SQLite.
func (r *WorkingSQLiteUserRepository) Delete(_ context.Context, _ entities.UserID) error {
	stubPanic()

	return nil
}

// List retrieves users with pagination from SQLite.
func (r *WorkingSQLiteUserRepository) List(
	_ context.Context,
	_ entities.UserStatus,
	_, _ int,
) ([]*entities.User, error) {
	stubPanic()

	return nil, nil
}

// Search searches users by query in SQLite.
func (r *WorkingSQLiteUserRepository) Search(
	_ context.Context,
	_ string,
	_ entities.UserStatus,
	_ int,
) ([]*entities.User, error) {
	stubPanic()

	return nil, nil
}

// SearchByTags searches users by tags in SQLite.
func (r *WorkingSQLiteUserRepository) SearchByTags(
	_ context.Context,
	_ []string,
	_ entities.UserStatus,
	_, _ int,
) ([]*entities.User, error) {
	stubPanic()

	return nil, nil
}

// CountByStatus counts users by status in SQLite.
func (r *WorkingSQLiteUserRepository) CountByStatus(
	_ context.Context,
) (map[entities.UserStatus]int64, error) {
	stubPanic()

	return nil, nil
}

// GetStats retrieves user statistics from SQLite.
func (r *WorkingSQLiteUserRepository) GetStats(_ context.Context) (*entities.UserStats, error) {
	stubPanic()

	return nil, nil
}

// VerifyCredentials verifies user credentials in SQLite.
func (r *WorkingSQLiteUserRepository) VerifyCredentials(
	_ context.Context,
	_ entities.Email,
	_ entities.PasswordHash,
) (*entities.User, error) {
	stubPanic()

	return nil, nil
}

// UpdatePassword updates user password in SQLite.
func (r *WorkingSQLiteUserRepository) UpdatePassword(
	_ context.Context,
	_ entities.UserID,
	_ entities.PasswordHash,
) error {
	stubPanic()

	return nil
}

// MarkVerified marks user as verified in SQLite.
func (r *WorkingSQLiteUserRepository) MarkVerified(_ context.Context, _ entities.UserID) error {
	stubPanic()

	return nil
}

// ChangeStatus changes user status in SQLite.
func (r *WorkingSQLiteUserRepository) ChangeStatus(
	_ context.Context,
	_ entities.UserID,
	_ entities.UserStatus,
) error {
	stubPanic()

	return nil
}

// Activate activates a user in SQLite.
func (r *WorkingSQLiteUserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user in SQLite.
func (r *WorkingSQLiteUserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user in SQLite.
func (r *WorkingSQLiteUserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role in SQLite.
func (r *WorkingSQLiteUserRepository) ChangeRole(
	_ context.Context,
	_ entities.UserID,
	_ entities.UserRole,
) error {
	stubPanic()

	return nil
}
