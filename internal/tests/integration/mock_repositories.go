// Package integration provides integration test utilities including mock repositories.
package integration

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

// NewMockUserRepository creates a new MockUserRepository for testing.
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:                 make(map[entities.UserID]*entities.User),
		passwordVerifications: make(map[string]string),
		idCounter:             1,
	}
}

// NewMockSessionRepository creates a new MockSessionRepository for testing.
func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions:  make(map[entities.SessionID]*entities.UserSession),
		idCounter: 1,
	}
}

func findUserBy(
	m map[entities.UserID]*entities.User,
	match func(*entities.User) bool,
) (*entities.User, error) {
	for _, user := range m {
		if match(user) {
			return user, nil
		}
	}

	return nil, entities.ErrUserNotFound
}

func findSessionBy(
	m map[entities.SessionID]*entities.UserSession,
	match func(*entities.UserSession) bool,
) (*entities.UserSession, error) {
	for _, session := range m {
		if match(session) {
			return session, nil
		}
	}

	return nil, entities.ErrSessionNotFound
}

// MockUserRepositoryStub provides default stub implementations for UserRepository methods.
// Embed this in mock implementations to avoid duplicating stub code.
type MockUserRepositoryStub struct{}

// UpdatePassword stub implementation.
func (MockUserRepositoryStub) UpdatePassword(
	context.Context,
	entities.UserID,
	entities.PasswordHash,
) error {
	return nil
}

// MarkVerified stub implementation.
func (MockUserRepositoryStub) MarkVerified(context.Context, entities.UserID) error {
	return nil
}

// ChangeStatus stub implementation.
func (MockUserRepositoryStub) ChangeStatus(
	context.Context,
	entities.UserID,
	entities.UserStatus,
) error {
	return nil
}

// Activate stub implementation.
func (MockUserRepositoryStub) Activate(context.Context, entities.UserID) error {
	return nil
}

// Deactivate stub implementation.
func (MockUserRepositoryStub) Deactivate(context.Context, entities.UserID) error {
	return nil
}

// Suspend stub implementation.
func (MockUserRepositoryStub) Suspend(context.Context, entities.UserID) error {
	return nil
}

// ChangeRole stub implementation.
func (MockUserRepositoryStub) ChangeRole(
	context.Context,
	entities.UserID,
	entities.UserRole,
) error {
	return nil
}

// Update stub implementation.
func (MockUserRepositoryStub) Update(context.Context, *entities.User) error {
	return nil
}

// Search stub implementation.
func (MockUserRepositoryStub) Search(
	context.Context,
	string,
	entities.UserStatus,
	int,
) ([]*entities.User, error) {
	return []*entities.User{}, nil
}

// SearchByTags stub implementation.
func (MockUserRepositoryStub) SearchByTags(
	context.Context,
	[]string,
	entities.UserStatus,
	int,
	int,
) ([]*entities.User, error) {
	return []*entities.User{}, nil
}

// MockUserRepository implements UserRepository for testing.
type MockUserRepository struct {
	MockUserRepositoryStub

	users                 map[entities.UserID]*entities.User
	passwordVerifications map[string]string
	idCounter             entities.UserID
}

// Create stores a new user in the mock repository.
func (m *MockUserRepository) Create(_ context.Context, user *entities.User) error {
	userID := m.idCounter
	m.idCounter++

	user.SetID(userID)
	m.users[userID] = user

	return nil
}

// GetByID retrieves a user by their ID from the mock repository.
func (m *MockUserRepository) GetByID(
	_ context.Context,
	id entities.UserID,
) (*entities.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, entities.ErrUserNotFound
	}

	return user, nil
}

// SetPasswordVerification sets the expected password for an email in the mock repository.
func (m *MockUserRepository) SetPasswordVerification(email, password string) {
	m.passwordVerifications[email] = password
}

// GetByUUID retrieves a user by their UUID from the mock repository.
func (m *MockUserRepository) GetByUUID(
	_ context.Context,
	uuid entities.UuID,
) (*entities.User, error) {
	for _, user := range m.users {
		if user.UUID().String() == string(uuid) {
			return user, nil
		}
	}

	return nil, entities.ErrUserNotFound
}

// GetByEmail retrieves a user by their email from the mock repository.
func (m *MockUserRepository) GetByEmail(
	_ context.Context,
	email entities.Email,
) (*entities.User, error) {
	return findUserBy(m.users, func(u *entities.User) bool {
		return u.Email() == email
	})
}

// GetByUsername retrieves a user by their username from the mock repository.
func (m *MockUserRepository) GetByUsername(
	_ context.Context,
	username entities.Username,
) (*entities.User, error) {
	return findUserBy(m.users, func(u *entities.User) bool {
		return u.Username() == username
	})
}

// Delete removes a user from the mock repository.
func (m *MockUserRepository) Delete(_ context.Context, id entities.UserID) error {
	delete(m.users, id)

	return nil
}

// List retrieves all users with a specific status from the mock repository.
func (m *MockUserRepository) List(
	_ context.Context,
	status entities.UserStatus,
	_, _ int,
) ([]*entities.User, error) {
	result := make([]*entities.User, 0)

	for _, user := range m.users {
		if user.Status() == status {
			result = append(result, user)
		}
	}

	return result, nil
}

// CountByStatus counts users by their status in the mock repository.
func (m *MockUserRepository) CountByStatus(
	_ context.Context,
) (map[entities.UserStatus]int64, error) {
	counts := make(map[entities.UserStatus]int64)
	for _, user := range m.users {
		counts[user.Status()]++
	}

	return counts, nil
}

// GetStats retrieves user statistics from the mock repository.
func (m *MockUserRepository) GetStats(_ context.Context) (*entities.UserStats, error) {
	stats := &entities.UserStats{}
	for _, user := range m.users {
		stats.TotalUsers++
		if user.Status() == entities.UserStatusActive {
			stats.ActiveUsers++
		}
	}

	return stats, nil
}

// VerifyCredentials verifies user credentials against the mock repository.
func (m *MockUserRepository) VerifyCredentials(
	ctx context.Context,
	email entities.Email,
	password entities.PasswordHash,
) (*entities.User, error) {
	expectedPassword := m.passwordVerifications[email.String()]
	if expectedPassword != password.String() {
		return nil, entities.ErrInvalidCredentials
	}

	user, err := m.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// MockSessionRepositoryStub provides default stub implementations for SessionRepository methods.
// Embed this in mock implementations to avoid duplicating stub code.
type MockSessionRepositoryStub struct{}

// Update stub implementation.
func (MockSessionRepositoryStub) Update(context.Context, *entities.UserSession) error {
	return nil
}

// DeactivateByToken stub implementation.
func (MockSessionRepositoryStub) DeactivateByToken(context.Context, entities.SessionToken) error {
	return nil
}

// DeactivateByUserID stub implementation.
func (MockSessionRepositoryStub) DeactivateByUserID(context.Context, entities.UserID) error {
	return nil
}

// CleanupExpired stub implementation.
func (MockSessionRepositoryStub) CleanupExpired(context.Context) (int64, error) {
	return 0, nil
}

// MockSessionRepository implements SessionRepository for testing.
type MockSessionRepository struct {
	MockSessionRepositoryStub

	sessions  map[entities.SessionID]*entities.UserSession
	idCounter entities.SessionID
}

// Create stores a new session in the mock repository.
func (m *MockSessionRepository) Create(_ context.Context, session *entities.UserSession) error {
	sessionID := m.idCounter
	m.idCounter++

	m.sessions[sessionID] = session

	return nil
}

// GetByToken retrieves a session by its token from the mock repository.
func (m *MockSessionRepository) GetByToken(
	_ context.Context,
	token entities.SessionToken,
) (*entities.UserSession, error) {
	return findSessionBy(m.sessions, func(s *entities.UserSession) bool {
		return s.Token() == token
	})
}

// GetByUserID retrieves all sessions for a user from the mock repository.
func (m *MockSessionRepository) GetByUserID(
	_ context.Context,
	userID entities.UserID,
	activeOnly bool,
) ([]*entities.UserSession, error) {
	result := make([]*entities.UserSession, 0)

	for _, session := range m.sessions {
		if session.UserID() == userID {
			if !activeOnly || session.IsActive() {
				result = append(result, session)
			}
		}
	}

	return result, nil
}

// Delete removes a session from the mock repository.
func (m *MockSessionRepository) Delete(_ context.Context, id entities.SessionID) error {
	delete(m.sessions, id)

	return nil
}

// GetActiveSessions counts active sessions for a user in the mock repository.
func (m *MockSessionRepository) GetActiveSessions(
	_ context.Context,
	userID entities.UserID,
) (int64, error) {
	count := int64(0)

	for _, session := range m.sessions {
		if session.UserID() == userID && session.IsActive() {
			count++
		}
	}

	return count, nil
}

// GetSessionStats retrieves session statistics from the mock repository.
func (m *MockSessionRepository) GetSessionStats(
	_ context.Context,
) (*entities.SessionStats, error) {
	stats := &entities.SessionStats{}
	for _, session := range m.sessions {
		stats.TotalSessions++
		if session.IsActive() {
			stats.ActiveSessions++
		}
	}

	return stats, nil
}

// Ensure MockUserRepository implements UserRepository.
var _ repositories.UserRepository = (*MockUserRepository)(nil)

// Ensure MockSessionRepository implements SessionRepository.
var _ repositories.SessionRepository = (*MockSessionRepository)(nil)
