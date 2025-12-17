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

// WorkingSQLiteUserRepository is a simplified implementation that works without generated sqlc code
// This demonstrates the pattern while we fix the SQL syntax issues
type WorkingSQLiteUserRepository struct {
	db         *sql.DB
	mapper     mappers.UserMapper
	converters *converters.SQLiteConverterSet
}

// NewWorkingSQLiteUserRepository creates a new working SQLite user repository
func NewWorkingSQLiteUserRepository(db *sql.DB) repositories.UserRepository {
	return &WorkingSQLiteUserRepository{
		db:     db,
		mapper: mappers.UserMapper{},
		converters: &converters.SQLiteConverterSet{
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

// Create creates a new user in SQLite
func (r *WorkingSQLiteUserRepository) Create(ctx context.Context, user *entities.User) error {
	// For now, implement using raw SQL to avoid generated code dependency
	query := `
		INSERT INTO users (email, username, password_hash, first_name, last_name, status, role, is_verified, metadata, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Convert domain values to database-compatible types
	email := r.converters.Email.DomainToDB(user.Email())
	username := r.converters.Username.DomainToDB(user.Username())
	passwordHash := r.converters.Password.DomainToDB(user.PasswordHash())
	firstName := user.FirstName().String()
	lastName := user.LastName().String()
	status := r.converters.Status.DomainToDB(user.Status())
	role := r.converters.Role.DomainToDB(user.Role())
	isVerified := r.converters.Bool.DomainToDB(user.IsVerified())

	// Convert metadata and tags to JSON
	metadataJSON := "{}"
	tagsJSON := "[]"

	result, err := r.db.ExecContext(ctx, query,
		email, username, passwordHash, firstName, lastName, status, role, isVerified, metadataJSON, tagsJSON)
	if err != nil {
		return errors.NewDatabaseError("failed to create user", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError("failed to check affected rows", err)
	}

	if rows == 0 {
		return errors.NewDatabaseError("no rows affected", fmt.Errorf("user creation failed"))
	}

	return nil
}

// GetByID retrieves a user by ID from SQLite
func (r *WorkingSQLiteUserRepository) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, status, role, 
		       is_verified, metadata, tags, created_at, updated_at, last_login_at
		FROM users 
		WHERE id = ?
	`

	user := &entities.User{} // This is wrong - need proper constructor
	var email, username, passwordHash, firstName, lastName, status, role string
	var isVerified bool
	var metadataJSON, tagsJSON string
	var createdAt, updatedAt, lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&id, &email, &username, &passwordHash, &firstName, &lastName, &status, &role,
		&isVerified, &metadataJSON, &tagsJSON, &createdAt, &updatedAt, &lastLoginAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, errors.NewDatabaseError("failed to get user by ID", err)
	}

	// This is a simplified example - proper implementation would:
	// 1. Convert database types to domain entities
	// 2. Use proper entity constructors
	// 3. Handle all field conversions

	// For now, return nil to show pattern
	return nil, fmt.Errorf("implementation in progress - user found with ID %d", id)
}

// GetByUUID retrieves a user by UUID from SQLite
func (r *WorkingSQLiteUserRepository) GetByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	// Implementation following same pattern as GetByID
	return nil, fmt.Errorf("implementation in progress - get by UUID %s", uuid)
}

// GetByEmail retrieves a user by email from SQLite
func (r *WorkingSQLiteUserRepository) GetByEmail(ctx context.Context, email entities.Email) (*entities.User, error) {
	// Implementation following same pattern as GetByID
	return nil, fmt.Errorf("implementation in progress - get by email %s", email.String())
}

// GetByUsername retrieves a user by username from SQLite
func (r *WorkingSQLiteUserRepository) GetByUsername(ctx context.Context, username entities.Username) (*entities.User, error) {
	// Implementation following same pattern as GetByID
	return nil, fmt.Errorf("implementation in progress - get by username %s", username.String())
}

// Update updates an existing user in SQLite
func (r *WorkingSQLiteUserRepository) Update(ctx context.Context, user *entities.User) error {
	// Implementation with UPDATE query
	return fmt.Errorf("implementation in progress - update user ID %d", user.ID())
}

// Delete soft deletes a user from SQLite
func (r *WorkingSQLiteUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	// Implementation with soft delete (UPDATE status)
	return fmt.Errorf("implementation in progress - delete user ID %d", id)
}

// List retrieves users with pagination from SQLite
func (r *WorkingSQLiteUserRepository) List(ctx context.Context, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	// Implementation with LIMIT and OFFSET
	return nil, fmt.Errorf("implementation in progress - list users with status %s, limit %d, offset %d", status.String(), limit, offset)
}

// Search searches users by query in SQLite
func (r *WorkingSQLiteUserRepository) Search(ctx context.Context, query string, status entities.UserStatus, limit int) ([]*entities.User, error) {
	// Implementation with LIKE or FTS5
	return nil, fmt.Errorf("implementation in progress - search users with query '%s', status %s, limit %d", query, status.String(), limit)
}

// SearchByTags searches users by tags in SQLite
func (r *WorkingSQLiteUserRepository) SearchByTags(ctx context.Context, tags []string, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	// Implementation with JSON operations
	return nil, fmt.Errorf("implementation in progress - search users by tags %v, status %s, limit %d, offset %d", tags, status.String(), limit, offset)
}

// CountByStatus counts users by status in SQLite
func (r *WorkingSQLiteUserRepository) CountByStatus(ctx context.Context) (map[entities.UserStatus]int64, error) {
	// Implementation with GROUP BY
	return nil, fmt.Errorf("implementation in progress - count users by status")
}

// GetStats retrieves user statistics from SQLite
func (r *WorkingSQLiteUserRepository) GetStats(ctx context.Context) (*entities.UserStats, error) {
	// Implementation with aggregate functions
	return nil, fmt.Errorf("implementation in progress - get user stats")
}

// VerifyCredentials verifies user credentials in SQLite
func (r *WorkingSQLiteUserRepository) VerifyCredentials(ctx context.Context, email entities.Email, password entities.PasswordHash) (*entities.User, error) {
	// Implementation with WHERE clause
	return nil, fmt.Errorf("implementation in progress - verify credentials for email %s", email.String())
}

// UpdatePassword updates user password in SQLite
func (r *WorkingSQLiteUserRepository) UpdatePassword(ctx context.Context, id entities.UserID, password entities.PasswordHash) error {
	// Implementation with UPDATE password_hash
	return fmt.Errorf("implementation in progress - update password for user ID %d", id)
}

// MarkVerified marks user as verified in SQLite
func (r *WorkingSQLiteUserRepository) MarkVerified(ctx context.Context, id entities.UserID) error {
	// Implementation with UPDATE is_verified
	return fmt.Errorf("implementation in progress - mark user verified for ID %d", id)
}

// ChangeStatus changes user status in SQLite
func (r *WorkingSQLiteUserRepository) ChangeStatus(ctx context.Context, id entities.UserID, status entities.UserStatus) error {
	// Implementation with UPDATE status
	return fmt.Errorf("implementation in progress - change status to %s for user ID %d", status.String(), id)
}

// Activate activates a user in SQLite
func (r *WorkingSQLiteUserRepository) Activate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusActive)
}

// Deactivate deactivates a user in SQLite
func (r *WorkingSQLiteUserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusInactive)
}

// Suspend suspends a user in SQLite
func (r *WorkingSQLiteUserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	return r.ChangeStatus(ctx, id, entities.UserStatusSuspended)
}

// ChangeRole changes user role in SQLite
func (r *WorkingSQLiteUserRepository) ChangeRole(ctx context.Context, id entities.UserID, role entities.UserRole) error {
	// Implementation with UPDATE role
	return fmt.Errorf("implementation in progress - change role to %s for user ID %d", role.String(), id)
}
