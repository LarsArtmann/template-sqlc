package integration

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
)

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:                 make(map[entities.UserID]*entities.User),
		passwordVerifications: make(map[string]string),
		idCounter:             1,
	}
}

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

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	userID := m.idCounter
	m.idCounter++

	user.SetID(userID)
	m.users[userID] = user

	return nil
}

func (m *MockUserRepository) GetByID(
	ctx context.Context,
	id entities.UserID,
) (*entities.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, entities.ErrUserNotFound
	}

	return user, nil
}

func (m *MockUserRepository) SetPasswordVerification(email, password string) {
	m.passwordVerifications[email] = password
}

func (m *MockUserRepository) GetByUUID(
	ctx context.Context,
	uuid entities.UuID,
) (*entities.User, error) {
	for _, user := range m.users {
		if user.UUID().String() == string(uuid) {
			return user, nil
		}
	}

	return nil, entities.ErrUserNotFound
}

func (m *MockUserRepository) GetByEmail(
	ctx context.Context,
	email entities.Email,
) (*entities.User, error) {
	return findUserBy(m.users, func(u *entities.User) bool {
		return u.Email() == email
	})
}

func (m *MockUserRepository) GetByUsername(
	ctx context.Context,
	username entities.Username,
) (*entities.User, error) {
	return findUserBy(m.users, func(u *entities.User) bool {
		return u.Username() == username
	})
}

func (m *MockUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	delete(m.users, id)

	return nil
}

func (m *MockUserRepository) List(
	ctx context.Context,
	status entities.UserStatus,
	limit, offset int,
) ([]*entities.User, error) {
	result := make([]*entities.User, 0)

	for _, user := range m.users {
		if user.Status() == status {
			result = append(result, user)
		}
	}

	return result, nil
}

func (m *MockUserRepository) CountByStatus(
	ctx context.Context,
) (map[entities.UserStatus]int64, error) {
	counts := make(map[entities.UserStatus]int64)
	for _, user := range m.users {
		counts[user.Status()]++
	}

	return counts, nil
}

func (m *MockUserRepository) GetStats(ctx context.Context) (*entities.UserStats, error) {
	stats := &entities.UserStats{}
	for _, user := range m.users {
		stats.TotalUsers++
		if user.Status() == entities.UserStatusActive {
			stats.ActiveUsers++
		}
	}

	return stats, nil
}

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

func (m *MockSessionRepository) Create(ctx context.Context, session *entities.UserSession) error {
	sessionID := m.idCounter
	m.idCounter++

	m.sessions[sessionID] = session

	return nil
}

func (m *MockSessionRepository) GetByToken(
	ctx context.Context,
	token entities.SessionToken,
) (*entities.UserSession, error) {
	return findSessionBy(m.sessions, func(s *entities.UserSession) bool {
		return s.Token() == token
	})
}

func (m *MockSessionRepository) GetByUserID(
	ctx context.Context,
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

func (m *MockSessionRepository) Delete(ctx context.Context, id entities.SessionID) error {
	delete(m.sessions, id)

	return nil
}

func (m *MockSessionRepository) GetActiveSessions(
	ctx context.Context,
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

func (m *MockSessionRepository) GetSessionStats(
	ctx context.Context,
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
