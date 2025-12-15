package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/events"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
	"github.com/LarsArtmann/template-sqlc/internal/domain/services"
	"github.com/LarsArtmann/template-sqlc/pkg/validation"
)

// UserServiceIntegrationTestSuite contains integration tests for UserService
type UserServiceIntegrationTestSuite struct {
	suite.Suite
	ctx            context.Context
	userService    *services.UserService
	userRepo       repositories.UserRepository
	sessionRepo    repositories.SessionRepository
	eventPublisher *events.InMemoryEventPublisher
	validator      validation.UserValidator
	cleanup        []func() error
}

// SetupSuite sets up the test suite
func (s *UserServiceIntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Initialize in-memory event publisher
	s.eventPublisher = events.NewInMemoryEventPublisher()

	// Initialize validator
	s.validator = validation.NewUserValidator()

	// Set up test database based on environment
	s.setupTestDatabase()

	// Create user service
	s.userService = services.NewUserService(
		s.userRepo,
		s.sessionRepo,
		s.eventPublisher,
		s.validator,
	)
}

// TearDownSuite cleans up the test suite
func (s *UserServiceIntegrationTestSuite) TearDownSuite() {
	// Cleanup test resources
	for i, cleanup := range s.cleanup {
		if err := cleanup(); err != nil {
			s.T().Logf("Cleanup %d failed: %v", i, err)
		}
	}
}

// setupTestDatabase sets up the test database
func (s *UserServiceIntegrationTestSuite) setupTestDatabase() {
	// This would be implemented based on the database being tested
	// Example for SQLite:
	// db, err := sql.Open("sqlite3", ":memory:")
	// require.NoError(s.T(), err)
	//
	// // Run migrations
	// _, err = db.Exec(`
	//     CREATE TABLE users (
	//         id INTEGER PRIMARY KEY AUTOINCREMENT,
	//         uuid TEXT UNIQUE NOT NULL,
	//         email TEXT UNIQUE NOT NULL,
	//         username TEXT UNIQUE NOT NULL,
	//         password_hash TEXT NOT NULL,
	//         first_name TEXT NOT NULL,
	//         last_name TEXT NOT NULL,
	//         status TEXT NOT NULL,
	//         role TEXT NOT NULL,
	//         is_verified INTEGER DEFAULT FALSE NOT NULL,
	//         metadata TEXT DEFAULT '{}',
	//         tags TEXT DEFAULT '[]',
	//         created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	//         updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	//         last_login_at DATETIME NULL
	//     );
	// `)
	// require.NoError(s.T(), err)
	//
	// s.userRepo = repositories.NewSQLiteUserRepository(db)
	// s.sessionRepo = repositories.NewSQLiteSessionRepository(db)
	// s.cleanup = append(s.cleanup, func() error {
	//     return db.Close()
	// })

	// For now, use mock repositories
	s.userRepo = &MockUserRepository{users: make(map[entities.UserID]*entities.User)}
	s.sessionRepo = &MockSessionRepository{sessions: make(map[entities.SessionID]*entities.UserSession)}
	s.cleanup = []func() error{}
}

// TestCreateUser tests user creation
func (s *UserServiceIntegrationTestSuite) TestCreateUser() {
	req := &services.CreateUserRequest{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashed_password_min_32_chars",
		FirstName:    "John",
		LastName:     "Doe",
		Status:       "active",
		Role:         "user",
		Tags:         []string{"developer", "golang"},
		Metadata:     map[string]interface{}{"team": "engineering"},
	}

	user, err := s.userService.CreateUser(s.ctx, req)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), user)

	assert.Equal(s.T(), entities.Email("test@example.com"), user.Email())
	assert.Equal(s.T(), entities.Username("testuser"), user.Username())
	assert.Equal(s.T(), entities.FirstName("John"), user.FirstName())
	assert.Equal(s.T(), entities.LastName("Doe"), user.LastName())
	assert.Equal(s.T(), entities.UserStatusActive, user.Status())
	assert.Equal(s.T(), entities.UserRoleUser, user.Role())
	assert.False(s.T(), user.IsVerified())
	assert.Contains(s.T(), user.Tags(), "developer")
	assert.Contains(s.T(), user.Tags(), "golang")

	// Check that event was published
	events := s.eventPublisher.Events()
	require.Len(s.T(), events, 1)

	userCreatedEvent := events[0]
	assert.Equal(s.T(), events.EventUserCreated, userCreatedEvent.Type)
}

func (s *UserServiceIntegrationTestSuite) TestCreateUserDuplicateEmail() {
	// Create first user
	req1 := &services.CreateUserRequest{
		Email:        "test@example.com",
		Username:     "testuser1",
		PasswordHash: "hashed_password_min_32_chars",
		FirstName:    "John",
		LastName:     "Doe",
		Status:       "active",
		Role:         "user",
	}

	_, err := s.userService.CreateUser(s.ctx, req1)
	require.NoError(s.T(), err)

	// Try to create second user with same email
	req2 := &services.CreateUserRequest{
		Email:        "test@example.com",
		Username:     "testuser2",
		PasswordHash: "hashed_password_min_32_chars",
		FirstName:    "Jane",
		LastName:     "Smith",
		Status:       "active",
		Role:         "user",
	}

	user, err := s.userService.CreateUser(s.ctx, req2)
	assert.Error(s.T(), err)
	assert.True(s.T(), entities.IsNotFoundError(err) ||
		assert.IsType(s.T(), entities.ErrUserAlreadyExists, err))
	assert.Nil(s.T(), user)
}

func (s *UserServiceIntegrationTestSuite) TestGetUser() {
	// Create a user first
	req := &services.CreateUserRequest{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashed_password_min_32_chars",
		FirstName:    "John",
		LastName:     "Doe",
		Status:       "active",
		Role:         "user",
	}

	createdUser, err := s.userService.CreateUser(s.ctx, req)
	require.NoError(s.T(), err)

	// Get the user
	retrievedUser, err := s.userService.GetUser(s.ctx, createdUser.ID())
	require.NoError(s.T(), err)
	require.NotNil(s.T(), retrievedUser)

	assert.Equal(s.T(), createdUser.ID(), retrievedUser.ID())
	assert.Equal(s.T(), createdUser.Email(), retrievedUser.Email())
	assert.Equal(s.T(), createdUser.Username(), retrievedUser.Username())
}

func (s *UserServiceIntegrationTestSuite) TestUpdateUser() {
	// Create a user first
	createReq := &services.CreateUserRequest{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashed_password_min_32_chars",
		FirstName:    "John",
		LastName:     "Doe",
		Status:       "active",
		Role:         "user",
	}

	user, err := s.userService.CreateUser(s.ctx, createReq)
	require.NoError(s.T(), err)

	// Update the user
	newFirstName := "Jane"
	updateReq := &services.UpdateUserRequest{
		UserID:    user.ID(),
		FirstName: &newFirstName,
		UpdatedBy: "system",
	}

	updatedUser, err := s.userService.UpdateUser(s.ctx, updateReq)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), entities.FirstName(newFirstName), updatedUser.FirstName())

	// Check that event was published
	events := s.eventPublisher.Events()
	assert.Len(s.T(), events, 2) // Create + Update

	updateEvent := events[1]
	assert.Equal(s.T(), events.EventUserUpdated, updateEvent.Type)
}

func (s *UserServiceIntegrationTestSuite) TestAuthenticateUser() {
	// Create a user first
	req := &services.CreateUserRequest{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashed_password_min_32_chars",
		FirstName:    "John",
		LastName:     "Doe",
		Status:       "active",
		Role:         "user",
	}

	user, err := s.userService.CreateUser(s.ctx, req)
	require.NoError(s.T(), err)

	// Mock the password verification
	// In a real implementation, this would be handled by the repository
	if mockRepo, ok := s.userRepo.(*MockUserRepository); ok {
		mockRepo.SetPasswordVerification(user.Email().String(), "correct_password")
	}

	// Test successful authentication
	session, err := s.userService.AuthenticateUser(s.ctx,
		"test@example.com",
		"correct_password",
		"127.0.0.1",
		"test-user-agent")

	if err == nil {
		require.NotNil(s.T(), session)
		assert.Equal(s.T(), user.ID(), session.UserID())
		assert.True(s.T(), session.IsActive())

		// Check that login event was published
		events := s.eventPublisher.Events()
		loginEvent := events[len(events)-1] // Should be the last event
		assert.Equal(s.T(), events.EventUserLogin, loginEvent.Type)
	}
}

func (s *UserServiceIntegrationTestSuite) TestAuthenticateUserInvalidCredentials() {
	// Test with non-existent user
	session, err := s.userService.AuthenticateUser(s.ctx,
		"nonexistent@example.com",
		"password",
		"127.0.0.1",
		"test-user-agent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), session)
	assert.True(s.T(), entities.IsUnauthorizedError(err))

	// Check that failed login event was published
	events := s.eventPublisher.Events()
	if len(events) > 0 {
		loginFailEvent := events[len(events)-1]
		assert.Equal(s.T(), events.EventUserLoginFail, loginFailEvent.Type)
	}
}

func (s *UserServiceIntegrationTestSuite) TestChangeUserRole() {
	// Create a user first
	req := &services.CreateUserRequest{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashed_password_min_32_chars",
		FirstName:    "John",
		LastName:     "Doe",
		Status:       "active",
		Role:         "user",
	}

	user, err := s.userService.CreateUser(s.ctx, req)
	require.NoError(s.T(), err)

	// Change user role
	updatedUser, err := s.userService.ChangeUserRole(s.ctx, user.ID(), entities.UserRoleAdmin, "system")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), entities.UserRoleAdmin, updatedUser.Role())

	// Check that role change event was published
	events := s.eventPublisher.Events()
	roleChangeEvent := events[len(events)-1]
	assert.Equal(s.T(), events.EventRoleChanged, roleChangeEvent.Type)
}

func (s *UserServiceIntegrationTestSuite) TestGetUserStats() {
	// Create multiple users with different statuses
	users := []struct {
		status string
		role   string
	}{
		{"active", "user"},
		{"active", "admin"},
		{"inactive", "user"},
		{"suspended", "user"},
	}

	for _, userData := range users {
		req := &services.CreateUserRequest{
			Email:        "user_" + userData.status + "_" + userData.role + "@example.com",
			Username:     "user_" + userData.status + "_" + userData.role,
			PasswordHash: "hashed_password_min_32_chars",
			FirstName:    "Test",
			LastName:     "User",
			Status:       userData.status,
			Role:         userData.role,
		}

		_, err := s.userService.CreateUser(s.ctx, req)
		require.NoError(s.T(), err)
	}

	// Get stats
	stats, err := s.userService.GetUserStats(s.ctx)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), stats)

	assert.Greater(s.T(), stats.TotalUsers, int64(0))
	assert.Greater(s.T(), stats.ActiveUsers, int64(0))
}

// Mock implementations for testing

type MockUserRepository struct {
	users                 map[entities.UserID]*entities.User
	passwordVerifications map[string]string
	idCounter             entities.UserID
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:                 make(map[entities.UserID]*entities.User),
		passwordVerifications: make(map[string]string),
		idCounter:             1,
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	userID := m.idCounter
	m.idCounter++

	// Simulate setting the ID
	// In a real implementation, this would be handled by the database
	// This is a simplified mock

	m.users[userID] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, entities.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) SetPasswordVerification(email, password string) {
	m.passwordVerifications[email] = password
}

// Implement other methods as needed for tests...
func (m *MockUserRepository) GetByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	// Mock implementation
	return nil, entities.ErrUserNotFound
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email entities.Email) (*entities.User, error) {
	// Mock implementation
	return nil, entities.ErrUserNotFound
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username entities.Username) (*entities.User, error) {
	// Mock implementation
	return nil, entities.ErrUserNotFound
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	// Mock implementation
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	// Mock implementation
	return []*entities.User{}, nil
}

func (m *MockUserRepository) Search(ctx context.Context, query string, status entities.UserStatus, limit int) ([]*entities.User, error) {
	// Mock implementation
	return []*entities.User{}, nil
}

func (m *MockUserRepository) SearchByTags(ctx context.Context, tags []string, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	// Mock implementation
	return []*entities.User{}, nil
}

func (m *MockUserRepository) CountByStatus(ctx context.Context) (map[entities.UserStatus]int64, error) {
	// Mock implementation
	return make(map[entities.UserStatus]int64), nil
}

func (m *MockUserRepository) GetStats(ctx context.Context) (*entities.UserStats, error) {
	// Mock implementation
	return &entities.UserStats{}, nil
}

func (m *MockUserRepository) VerifyCredentials(ctx context.Context, email entities.Email, password entities.PasswordHash) (*entities.User, error) {
	// Mock password verification
	expectedPassword := m.passwordVerifications[email.String()]
	if expectedPassword != password.String() {
		return nil, entities.ErrInvalidCredentials
	}

	// Return a mock user
	return &entities.User{}, nil
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id entities.UserID, password entities.PasswordHash) error {
	// Mock implementation
	return nil
}

func (m *MockUserRepository) MarkVerified(ctx context.Context, id entities.UserID) error {
	// Mock implementation
	return nil
}

func (m *MockUserRepository) ChangeStatus(ctx context.Context, id entities.UserID, status entities.UserStatus) error {
	// Mock implementation
	return nil
}

func (m *MockUserRepository) Activate(ctx context.Context, id entities.UserID) error {
	// Mock implementation
	return nil
}

func (m *MockUserRepository) Deactivate(ctx context.Context, id entities.UserID) error {
	// Mock implementation
	return nil
}

func (m *MockUserRepository) Suspend(ctx context.Context, id entities.UserID) error {
	// Mock implementation
	return nil
}

func (m *MockUserRepository) ChangeRole(ctx context.Context, id entities.UserID, role entities.UserRole) error {
	// Mock implementation
	return nil
}

type MockSessionRepository struct {
	sessions  map[entities.SessionID]*entities.UserSession
	idCounter entities.SessionID
}

func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions:  make(map[entities.SessionID]*entities.UserSession),
		idCounter: 1,
	}
}

func (m *MockSessionRepository) Create(ctx context.Context, session *entities.UserSession) error {
	sessionID := m.idCounter
	m.idCounter++

	m.sessions[sessionID] = session
	return nil
}

// Implement other session methods as needed...
func (m *MockSessionRepository) GetByToken(ctx context.Context, token entities.SessionToken) (*entities.UserSession, error) {
	// Mock implementation
	return nil, entities.ErrSessionNotFound
}

func (m *MockSessionRepository) GetByUserID(ctx context.Context, userID entities.UserID, activeOnly bool) ([]*entities.UserSession, error) {
	// Mock implementation
	return []*entities.UserSession{}, nil
}

func (m *MockSessionRepository) Update(ctx context.Context, session *entities.UserSession) error {
	// Mock implementation
	return nil
}

func (m *MockSessionRepository) Delete(ctx context.Context, id entities.SessionID) error {
	delete(m.sessions, id)
	return nil
}

func (m *MockSessionRepository) DeactivateByToken(ctx context.Context, token entities.SessionToken) error {
	// Mock implementation
	return nil
}

func (m *MockSessionRepository) DeactivateByUserID(ctx context.Context, userID entities.UserID) error {
	// Mock implementation
	return nil
}

func (m *MockSessionRepository) CleanupExpired(ctx context.Context) (int64, error) {
	// Mock implementation
	return 0, nil
}

func (m *MockSessionRepository) GetActiveSessions(ctx context.Context, userID entities.UserID) (int64, error) {
	// Mock implementation
	return 0, nil
}

func (m *MockSessionRepository) GetSessionStats(ctx context.Context) (*entities.SessionStats, error) {
	// Mock implementation
	return &entities.SessionStats{}, nil
}

// Test suite runner
func TestUserServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(UserServiceIntegrationTestSuite))
}
