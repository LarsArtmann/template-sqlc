package adapters

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// NotImplementedRepository provides default implementations for repository methods
// that are not yet implemented. Embed this struct in database-specific repositories
// to avoid duplicating stub code.
type NotImplementedRepository struct{}

// NotImplementedUserRepository provides stub implementations for UserRepository methods.
type NotImplementedUserRepository struct {
	NotImplementedRepository

	dbName string
}

// NewNotImplementedUserRepository creates a new NotImplementedUserRepository.
func NewNotImplementedUserRepository(dbName string) *NotImplementedUserRepository {
	return &NotImplementedUserRepository{dbName: dbName}
}

func (r *NotImplementedUserRepository) NotImplemented(method string) error {
	return entities.StubNotImplemented(method, r.dbName)
}

// Create is a stub implementation.
func (r *NotImplementedUserRepository) Create(_ context.Context, _ *entities.User) error {
	return r.NotImplemented("Create")
}

// GetByID is a stub implementation.
func (r *NotImplementedUserRepository) GetByID(
	_ context.Context,
	_ entities.UserID,
) (*entities.User, error) {
	return nil, r.NotImplemented("GetByID")
}

// GetByUUID is a stub implementation.
func (r *NotImplementedUserRepository) GetByUUID(
	_ context.Context,
	_ entities.UuID,
) (*entities.User, error) {
	return nil, r.NotImplemented("GetByUUID")
}

// GetByEmail is a stub implementation.
func (r *NotImplementedUserRepository) GetByEmail(
	_ context.Context,
	_ entities.Email,
) (*entities.User, error) {
	return nil, r.NotImplemented("GetByEmail")
}

// GetByUsername is a stub implementation.
func (r *NotImplementedUserRepository) GetByUsername(
	_ context.Context,
	_ entities.Username,
) (*entities.User, error) {
	return nil, r.NotImplemented("GetByUsername")
}

// Update is a stub implementation.
func (r *NotImplementedUserRepository) Update(_ context.Context, _ *entities.User) error {
	return r.NotImplemented("Update")
}

// Delete is a stub implementation.
func (r *NotImplementedUserRepository) Delete(_ context.Context, _ entities.UserID) error {
	return r.NotImplemented("Delete")
}

// List is a stub implementation.
func (r *NotImplementedUserRepository) List(
	_ context.Context,
	_ entities.UserStatus,
	_, _ int,
) ([]*entities.User, error) {
	return nil, r.NotImplemented("List")
}

// Search is a stub implementation.
func (r *NotImplementedUserRepository) Search(
	_ context.Context,
	_ string,
	_ entities.UserStatus,
	_ int,
) ([]*entities.User, error) {
	return nil, r.NotImplemented("Search")
}

// SearchByTags is a stub implementation.
func (r *NotImplementedUserRepository) SearchByTags(
	_ context.Context,
	_ []string,
	_ entities.UserStatus,
	_, _ int,
) ([]*entities.User, error) {
	return nil, r.NotImplemented("SearchByTags")
}

// CountByStatus is a stub implementation.
func (r *NotImplementedUserRepository) CountByStatus(
	_ context.Context,
) (map[entities.UserStatus]int64, error) {
	return nil, r.NotImplemented("CountByStatus")
}

// GetStats is a stub implementation.
func (r *NotImplementedUserRepository) GetStats(_ context.Context) (*entities.UserStats, error) {
	return nil, r.NotImplemented("GetStats")
}

// VerifyCredentials is a stub implementation.
func (r *NotImplementedUserRepository) VerifyCredentials(
	_ context.Context,
	_ entities.Email,
	_ entities.PasswordHash,
) (*entities.User, error) {
	return nil, r.NotImplemented("VerifyCredentials")
}

// UpdatePassword is a stub implementation.
func (r *NotImplementedUserRepository) UpdatePassword(
	_ context.Context,
	_ entities.UserID,
	_ entities.PasswordHash,
) error {
	return r.NotImplemented("UpdatePassword")
}

// MarkVerified is a stub implementation.
func (r *NotImplementedUserRepository) MarkVerified(_ context.Context, _ entities.UserID) error {
	return r.NotImplemented("MarkVerified")
}

// ChangeStatus is a stub implementation.
func (r *NotImplementedUserRepository) ChangeStatus(
	_ context.Context,
	_ entities.UserID,
	_ entities.UserStatus,
) error {
	return r.NotImplemented("ChangeStatus")
}

// Activate is a stub implementation.
func (r *NotImplementedUserRepository) Activate(_ context.Context, _ entities.UserID) error {
	return r.NotImplemented("Activate")
}

// Deactivate is a stub implementation.
func (r *NotImplementedUserRepository) Deactivate(_ context.Context, _ entities.UserID) error {
	return r.NotImplemented("Deactivate")
}

// Suspend is a stub implementation.
func (r *NotImplementedUserRepository) Suspend(_ context.Context, _ entities.UserID) error {
	return r.NotImplemented("Suspend")
}

// ChangeRole is a stub implementation.
func (r *NotImplementedUserRepository) ChangeRole(
	_ context.Context,
	_ entities.UserID,
	_ entities.UserRole,
) error {
	return r.NotImplemented("ChangeRole")
}

// Ensure NotImplementedUserRepository implements UserRepository.
var _ repositories.UserRepository = (*NotImplementedUserRepository)(nil)

// NotImplementedSessionRepository provides stub implementations for SessionRepository methods.
type NotImplementedSessionRepository struct {
	NotImplementedRepository

	dbName string
}

// NewNotImplementedSessionRepository creates a new NotImplementedSessionRepository.
func NewNotImplementedSessionRepository(dbName string) *NotImplementedSessionRepository {
	return &NotImplementedSessionRepository{dbName: dbName}
}

func (r *NotImplementedSessionRepository) NotImplemented(method string) error {
	return entities.StubNotImplemented(method, r.dbName)
}

// Create is a stub implementation.
func (r *NotImplementedSessionRepository) Create(_ context.Context, _ *entities.UserSession) error {
	return r.NotImplemented("Create")
}

// GetByToken is a stub implementation.
func (r *NotImplementedSessionRepository) GetByToken(
	_ context.Context,
	_ entities.SessionToken,
) (*entities.UserSession, error) {
	return nil, r.NotImplemented("GetByToken")
}

// GetByUserID is a stub implementation.
func (r *NotImplementedSessionRepository) GetByUserID(
	_ context.Context,
	_ entities.UserID,
	_ bool,
) ([]*entities.UserSession, error) {
	return nil, r.NotImplemented("GetByUserID")
}

// Update is a stub implementation.
func (r *NotImplementedSessionRepository) Update(_ context.Context, _ *entities.UserSession) error {
	return r.NotImplemented("Update")
}

// Delete is a stub implementation.
func (r *NotImplementedSessionRepository) Delete(_ context.Context, _ entities.SessionID) error {
	return r.NotImplemented("Delete")
}

// DeactivateByToken is a stub implementation.
func (r *NotImplementedSessionRepository) DeactivateByToken(
	_ context.Context,
	_ entities.SessionToken,
) error {
	return r.NotImplemented("DeactivateByToken")
}

// DeactivateByUserID is a stub implementation.
func (r *NotImplementedSessionRepository) DeactivateByUserID(
	_ context.Context,
	_ entities.UserID,
) error {
	return r.NotImplemented("DeactivateByUserID")
}

// CleanupExpired is a stub implementation.
func (r *NotImplementedSessionRepository) CleanupExpired(_ context.Context) (int64, error) {
	return 0, r.NotImplemented("CleanupExpired")
}

// GetActiveSessions is a stub implementation.
func (r *NotImplementedSessionRepository) GetActiveSessions(
	_ context.Context,
	_ entities.UserID,
) (int64, error) {
	return 0, r.NotImplemented("GetActiveSessions")
}

// GetSessionStats is a stub implementation.
func (r *NotImplementedSessionRepository) GetSessionStats(
	_ context.Context,
) (*entities.SessionStats, error) {
	return nil, r.NotImplemented("GetSessionStats")
}

// Ensure NotImplementedSessionRepository implements SessionRepository.
var _ repositories.SessionRepository = (*NotImplementedSessionRepository)(nil)
