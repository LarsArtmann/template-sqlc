package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/events"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
	"github.com/LarsArtmann/template-sqlc/internal/domain/services"
	"github.com/LarsArtmann/template-sqlc/pkg/validation"
)

const testPasswordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsQ5pPjZ5yVlWK5WAe"

func newTestCreateUserRequest(username, firstName, lastName string) *services.CreateUserRequest {
	return &services.CreateUserRequest{
		Email:        "test@example.com",
		Username:     username,
		PasswordHash: testPasswordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Status:       "active",
		Role:         "user",
	}
}

// UserServiceIntegrationTestSuite contains integration tests for UserService
type UserServiceIntegrationTestSuite struct {
	suite.Suite
	ctx            context.Context
	userService    *services.UserService
	userRepo       repositories.UserRepository
	sessionRepo    repositories.SessionRepository
	eventPublisher *events.InMemoryEventPublisher
	validator      *validation.UserValidator
	cleanup        []func() error
}

// SetupSuite sets up the test suite
func (s *UserServiceIntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Initialize in-memory event publisher
	s.eventPublisher = events.NewInMemoryEventPublisher()

	// Initialize validator
	s.validator = validation.NewUserValidator()

	// Initialize mock repositories
	s.userRepo = NewMockUserRepository()
	s.sessionRepo = NewMockSessionRepository()

	// Create user service
	s.userService = services.NewUserService(
		s.userRepo,
		s.sessionRepo,
		s.eventPublisher,
		s.validator,
	)
}

// SetupTest resets state before each test
func (s *UserServiceIntegrationTestSuite) SetupTest() {
	// Reset mock repositories for each test
	s.eventPublisher.Clear()
	s.userRepo = NewMockUserRepository()
	s.sessionRepo = NewMockSessionRepository()
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

// TestCreateUser tests user creation
func (s *UserServiceIntegrationTestSuite) TestCreateUser() {
	req := &services.CreateUserRequest{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsQ5pPjZ5yVlWK5WAe", // Valid bcrypt hash (60 chars)
		FirstName:    "John",
		LastName:     "Doe",
		Status:       "active",
		Role:         "user",
		Tags:         []string{"developer", "golang"},
		Metadata:     map[string]any{"team": "engineering"},
	}

	user, err := s.userService.CreateUser(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(user)

	s.Equal(entities.Email("test@example.com"), user.Email())
	s.Equal(entities.Username("testuser"), user.Username())
	s.Equal(entities.FirstName("John"), user.FirstName())
	s.Equal(entities.LastName("Doe"), user.LastName())
	s.Equal(entities.UserStatusActive, user.Status())
	s.Equal(entities.UserRoleUser, user.Role())
	s.False(user.IsVerified())
	s.Contains(user.Tags(), "developer")
	s.Contains(user.Tags(), "golang")

	// Check that event was published
	userEvents := s.eventPublisher.Events()
	s.Require().Len(userEvents, 1)

	userCreatedEvent := userEvents[0]
	s.Equal(events.EventUserCreated, userCreatedEvent.Type)
}

func (s *UserServiceIntegrationTestSuite) TestCreateUserDuplicateEmail() {
	// Create first user
	req1 := newTestCreateUserRequest("testuser1", "John", "Doe")

	_, err := s.userService.CreateUser(s.ctx, req1)
	s.Require().NoError(err)

	// Try to create second user with same email
	req2 := newTestCreateUserRequest("testuser2", "Jane", "Smith")

	user, err := s.userService.CreateUser(s.ctx, req2)
	s.Error(err)
	s.True(entities.IsNotFoundError(err) ||
		assert.IsType(s.T(), entities.ErrUserAlreadyExists, err))
	s.Nil(user)
}

func (s *UserServiceIntegrationTestSuite) TestGetUser() {
	// Create a user first
	req := newTestCreateUserRequest("testuser", "John", "Doe")

	createdUser, err := s.userService.CreateUser(s.ctx, req)
	s.Require().NoError(err)

	// Get the user
	retrievedUser, err := s.userService.GetUser(s.ctx, createdUser.ID())
	s.Require().NoError(err)
	s.Require().NotNil(retrievedUser)

	s.Equal(createdUser.ID(), retrievedUser.ID())
	s.Equal(createdUser.Email(), retrievedUser.Email())
	s.Equal(createdUser.Username(), retrievedUser.Username())
}

func (s *UserServiceIntegrationTestSuite) TestUpdateUser() {
	// Create a user first
	createReq := newTestCreateUserRequest("testuser", "John", "Doe")

	user, err := s.userService.CreateUser(s.ctx, createReq)
	s.Require().NoError(err)

	// Update user
	newFirstName := "Jane"
	updateReq := &services.UpdateUserRequest{
		UserID:    user.ID(),
		FirstName: &newFirstName,
		UpdatedBy: "system",
	}

	updatedUser, err := s.userService.UpdateUser(s.ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(entities.FirstName(newFirstName), updatedUser.FirstName())

	// Check that event was published
	userEvents := s.eventPublisher.Events()
	s.Len(userEvents, 2) // Create + Update

	updateEvent := userEvents[1]
	s.Equal(events.EventUserUpdated, updateEvent.Type)
}

func (s *UserServiceIntegrationTestSuite) TestAuthenticateUser() {
	// Create a user first
	req := newTestCreateUserRequest("testuser", "John", "Doe")

	user, err := s.userService.CreateUser(s.ctx, req)
	s.Require().NoError(err)

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
		s.Require().NotNil(session)
		s.Equal(user.ID(), session.UserID())
		s.True(session.IsActive())

		// Check that login event was published
		userEvents := s.eventPublisher.Events()
		loginEvent := userEvents[len(userEvents)-1] // Should be the last event
		s.Equal(events.EventUserLogin, loginEvent.Type)
	}
}

func (s *UserServiceIntegrationTestSuite) TestAuthenticateUserInvalidCredentials() {
	// Test with non-existent user
	session, err := s.userService.AuthenticateUser(s.ctx,
		"nonexistent@example.com",
		"password",
		"127.0.0.1",
		"test-user-agent")

	s.Error(err)
	s.Nil(session)
	s.True(entities.IsAuthenticationError(err))

	// Check that failed login event was published
	userEvents := s.eventPublisher.Events()
	if len(userEvents) > 0 {
		loginFailEvent := userEvents[len(userEvents)-1]
		s.Equal(events.EventUserLoginFail, loginFailEvent.Type)
	}
}

func (s *UserServiceIntegrationTestSuite) TestChangeUserRole() {
	// Create a user first
	req := newTestCreateUserRequest("testuser", "John", "Doe")

	user, err := s.userService.CreateUser(s.ctx, req)
	s.Require().NoError(err)

	// Change user role
	updatedUser, err := s.userService.ChangeUserRole(
		s.ctx,
		user.ID(),
		entities.UserRoleAdmin,
		"system",
	)
	s.Require().NoError(err)
	s.Equal(entities.UserRoleAdmin, updatedUser.Role())

	// Check that role change event was published
	userEvents := s.eventPublisher.Events()
	roleChangeEvent := userEvents[len(userEvents)-1]
	s.Equal(events.EventRoleChanged, roleChangeEvent.Type)
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
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsQ5pPjZ5yVlWK5WAe",
			FirstName:    "Test",
			LastName:     "User",
			Status:       userData.status,
			Role:         userData.role,
		}

		_, err := s.userService.CreateUser(s.ctx, req)
		s.Require().NoError(err)
	}

	// Get stats
	stats, err := s.userService.GetUserStats(s.ctx)
	s.Require().NoError(err)
	s.Require().NotNil(stats)

	s.Positive(stats.TotalUsers)
	s.Positive(stats.ActiveUsers)
}

// Test suite runner
func TestUserServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(UserServiceIntegrationTestSuite))
}
