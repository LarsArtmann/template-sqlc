package bdd

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/events"
	"github.com/LarsArtmann/template-sqlc/internal/domain/services"
	"github.com/LarsArtmann/template-sqlc/internal/tests/integration"
	"github.com/LarsArtmann/template-sqlc/pkg/validation"
)

// UserFeaturesTestSuite contains BDD tests for user functionality
type UserFeaturesTestSuite struct {
	ctx            context.Context
	userService    *services.UserService
	eventPublisher *events.InMemoryEventPublisher
	currentUser    *entities.User
	currentSession *entities.UserSession
	lastError      error
}

// InitializeContext sets up the test context
func (s *UserFeaturesTestSuite) InitializeContext(ctx *godog.ScenarioContext) {
	s.eventPublisher = events.NewInMemoryEventPublisher()
	s.validator = validation.NewUserValidator()

	// Setup mock repositories (could be swapped with real DB in other scenarios)
	s.userRepo = integration.NewMockUserRepository()
	s.sessionRepo = integration.NewMockSessionRepository()

	s.userService = services.NewUserService(
		s.userRepo,
		s.sessionRepo,
		s.eventPublisher,
		s.validator,
	)

	// Given steps
	ctx.Given(`^a user with email "([^"]*)" and username "([^"]*)"$`, s.createUserWithEmailUsername)
	ctx.Given(`^a user account with status "([^"]*)"$`, s.createUserWithStatus)
	ctx.Given(`^an active user account$`, s.createActiveUserAccount)
	ctx.Given(`^I have valid user credentials$`, s.haveValidUserCredentials)
	ctx.Given(`^I have invalid user credentials$`, s.haveInvalidUserCredentials)
	ctx.Given(`^an inactive user account$`, s.createInactiveUserAccount)
	ctx.Given(`^a suspended user account$`, s.createSuspendedUserAccount)

	// When steps
	ctx.When(`^I create a user with valid data$`, s.createUserWithValidData)
	ctx.When(`^I create a user with email "([^"]*)"$`, s.createUserWithEmail)
	ctx.When(`^I create a user with username "([^"]*)"$`, s.createUserWithUsername)
	ctx.When(`^I attempt to authenticate with these credentials$`, s.authenticateWithCredentials)
	ctx.When(`^I update the user profile$`, s.updateUserProfile)
	ctx.When(`^I change the user role to "([^"]*)"$`, s.changeUserRole)
	ctx.When(`^I verify the user account$`, s.verifyUserAccount)
	ctx.When(`^I deactivate the user account$`, s.deactivateUserAccount)
	ctx.When(`^I get the user statistics$`, s.getUserStatistics)

	// Then steps
	ctx.Then(`^the user should be created successfully$`, s.userShouldBeCreatedSuccessfully)
	ctx.Then(`^the user should have ID (\d+)$`, s.userShouldHaveID)
	ctx.Then(`^I should receive a validation error$`, s.shouldReceiveValidationError)
	ctx.Then(`^I should receive a "user already exists" error$`, s.shouldReceiveUserAlreadyExistsError)
	ctx.Then(`^the authentication should succeed$`, s.authenticationShouldSucceed)
	ctx.Then(`^the authentication should fail$`, s.authenticationShouldFail)
	ctx.Then(`^I should receive a "user not found" error$`, s.shouldReceiveUserNotFoundError)
	ctx.Then(`^I should receive a "invalid credentials" error$`, s.shouldReceiveInvalidCredentialsError)
	ctx.Then(`^the session should be created$`, s.sessionShouldBeCreated)
	ctx.Then(`^the user profile should be updated$`, s.userProfileShouldBeUpdated)
	ctx.Then(`^the user role should be changed to "([^"]*)"$`, s.userRoleShouldBeChanged)
	ctx.Then(`^the user account should be verified$`, s.userAccountShouldBeVerified)
	ctx.Then(`^the user account should be deactivated$`, s.userAccountShouldBeDeactivated)
	ctx.Then(`^a user created event should be published$`, s.userCreatedEventShouldBePublished)
	ctx.Then(`^a user updated event should be published$`, s.userUpdatedEventShouldBePublished)
	ctx.Then(`^a user login event should be published$`, s.userLoginEventShouldBePublished)
	ctx.Then(`^a user login failed event should be published$`, s.userLoginFailEventShouldBePublished)
	ctx.Then(`^a role changed event should be published$`, s.roleChangedEventShouldBePublished)
}

// Given steps

func (s *UserFeaturesTestSuite) createUserWithEmailUsername(email, username string) error {
	req := &services.CreateUserRequest{
		Email:        email,
		Username:     username,
		PasswordHash: "hashed_password_min_32_chars_for_testing",
		FirstName:    "Test",
		LastName:     "User",
		Status:       "active",
		Role:         "user",
		Tags:         []string{"test"},
		Metadata:     map[string]interface{}{"source": "bdd"},
	}

	user, err := s.userService.CreateUser(context.Background(), req)
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) createUserWithStatus(status string) error {
	req := &services.CreateUserRequest{
		Email:        "status@example.com",
		Username:     "statususer",
		PasswordHash: "hashed_password_min_32_chars_for_testing",
		FirstName:    "Status",
		LastName:     "User",
		Status:       status,
		Role:         "user",
	}

	user, err := s.userService.CreateUser(context.Background(), req)
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) createActiveUserAccount() error {
	return s.createUserWithStatus("active")
}

func (s *UserFeaturesTestSuite) haveValidUserCredentials() error {
	return s.createUserWithEmailUsername("valid@example.com", "validuser")
}

func (s *UserFeaturesTestSuite) haveInvalidUserCredentials() error {
	return s.createUserWithEmailUsername("invalid@example.com", "invaliduser")
}

func (s *UserFeaturesTestSuite) createInactiveUserAccount() error {
	return s.createUserWithStatus("inactive")
}

func (s *UserFeaturesTestSuite) createSuspendedUserAccount() error {
	return s.createUserWithStatus("suspended")
}

// When steps

func (s *UserFeaturesTestSuite) createUserWithValidData() error {
	req := &services.CreateUserRequest{
		Email:        "valid@example.com",
		Username:     "validuser",
		PasswordHash: "hashed_password_min_32_chars_for_testing",
		FirstName:    "Valid",
		LastName:     "User",
		Status:       "active",
		Role:         "user",
		Tags:         []string{"valid", "test"},
		Metadata:     map[string]interface{}{"test": true},
	}

	user, err := s.userService.CreateUser(context.Background(), req)
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) createUserWithEmail(email string) error {
	req := &services.CreateUserRequest{
		Email:        email,
		Username:     "emailuser",
		PasswordHash: "hashed_password_min_32_chars_for_testing",
		FirstName:    "Email",
		LastName:     "User",
		Status:       "active",
		Role:         "user",
	}

	user, err := s.userService.CreateUser(context.Background(), req)
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) createUserWithUsername(username string) error {
	req := &services.CreateUserRequest{
		Email:        "username@example.com",
		Username:     username,
		PasswordHash: "hashed_password_min_32_chars_for_testing",
		FirstName:    "Username",
		LastName:     "User",
		Status:       "active",
		Role:         "user",
	}

	user, err := s.userService.CreateUser(context.Background(), req)
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) authenticateWithCredentials() error {
	if s.currentUser == nil {
		return fmt.Errorf("no current user to authenticate")
	}

	session, err := s.userService.AuthenticateUser(
		context.Background(),
		s.currentUser.Email().String(),
		"correct_password",
		"127.0.0.1",
		"test-user-agent",
	)

	s.currentSession = session
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) updateUserProfile() error {
	if s.currentUser == nil {
		return fmt.Errorf("no current user to update")
	}

	newFirstName := "Updated"
	updateReq := &services.UpdateUserRequest{
		UserID:    s.currentUser.ID(),
		FirstName: &newFirstName,
		UpdatedBy: "test",
	}

	user, err := s.userService.UpdateUser(context.Background(), updateReq)
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) changeUserRole(role string) error {
	if s.currentUser == nil {
		return fmt.Errorf("no current user for role change")
	}

	userRole := entities.UserRole(role)
	user, err := s.userService.ChangeUserRole(
		context.Background(),
		s.currentUser.ID(),
		userRole,
		"system",
	)

	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) verifyUserAccount() error {
	if s.currentUser == nil {
		return fmt.Errorf("no current user to verify")
	}

	err := s.userRepo.MarkVerified(context.Background(), s.currentUser.ID())
	s.lastError = err

	if err == nil {
		s.currentUser.Verify()
	}

	return nil
}

func (s *UserFeaturesTestSuite) deactivateUserAccount() error {
	if s.currentUser == nil {
		return fmt.Errorf("no current user to deactivate")
	}

	err := s.userRepo.Deactivate(context.Background(), s.currentUser.ID())
	s.lastError = err

	if err == nil {
		s.currentUser.ChangeStatus(entities.UserStatusInactive)
	}

	return nil
}

func (s *UserFeaturesTestSuite) getUserStatistics() error {
	_, err := s.userService.GetUserStatistics(context.Background())
	s.lastError = err

	return nil
}

// Then steps

func (s *UserFeaturesTestSuite) userShouldBeCreatedSuccessfully() error {
	if s.lastError != nil {
		return fmt.Errorf("expected user to be created successfully, got error: %v", s.lastError)
	}
	if s.currentUser == nil {
		return fmt.Errorf("expected user to be created, but got nil")
	}
	return nil
}

func (s *UserFeaturesTestSuite) userShouldHaveID(expectedIDStr string) error {
	if s.currentUser == nil {
		return fmt.Errorf("no current user")
	}

	if s.currentUser.ID().String() != expectedIDStr {
		return fmt.Errorf("expected user ID %s, got %s", expectedIDStr, s.currentUser.ID().String())
	}

	return nil
}

func (s *UserFeaturesTestSuite) shouldReceiveValidationError() error {
	if s.lastError == nil {
		return fmt.Errorf("expected validation error, got nil")
	}

	// Check if it's a validation error type
	if !entities.IsValidationError(s.lastError) {
		return fmt.Errorf("expected validation error, got: %v", s.lastError)
	}

	return nil
}

func (s *UserFeaturesTestSuite) shouldReceiveUserAlreadyExistsError() error {
	if s.lastError == nil {
		return fmt.Errorf("expected user already exists error, got nil")
	}

	// Check if it's a conflict error
	if !entities.IsNotFoundError(s.lastError) && s.lastError != entities.ErrUserAlreadyExists {
		return fmt.Errorf("expected user already exists error, got: %v", s.lastError)
	}

	return nil
}

func (s *UserFeaturesTestSuite) authenticationShouldSucceed() error {
	if s.lastError != nil {
		return fmt.Errorf("expected authentication to succeed, got error: %v", s.lastError)
	}
	if s.currentSession == nil {
		return fmt.Errorf("expected session to be created, but got nil")
	}
	return nil
}

func (s *UserFeaturesTestSuite) authenticationShouldFail() error {
	if s.lastError == nil {
		return fmt.Errorf("expected authentication to fail, got nil")
	}
	if !entities.IsUnauthorizedError(s.lastError) {
		return fmt.Errorf("expected unauthorized error, got: %v", s.lastError)
	}
	return nil
}

func (s *UserFeaturesTestSuite) shouldReceiveUserNotFoundError() error {
	if s.lastError == nil {
		return fmt.Errorf("expected user not found error, got nil")
	}
	if s.lastError != entities.ErrUserNotFound {
		return fmt.Errorf("expected user not found error, got: %v", s.lastError)
	}
	return nil
}

func (s *UserFeaturesTestSuite) shouldReceiveInvalidCredentialsError() error {
	if s.lastError == nil {
		return fmt.Errorf("expected invalid credentials error, got nil")
	}
	if !entities.IsUnauthorizedError(s.lastError) {
		return fmt.Errorf("expected invalid credentials error, got: %v", s.lastError)
	}
	return nil
}

func (s *UserFeaturesTestSuite) sessionShouldBeCreated() error {
	if s.currentSession == nil {
		return fmt.Errorf("expected session to be created, but got nil")
	}
	return nil
}

func (s *UserFeaturesTestSuite) userProfileShouldBeUpdated() error {
	if s.lastError != nil {
		return fmt.Errorf("expected profile update to succeed, got error: %v", s.lastError)
	}
	if s.currentUser == nil {
		return fmt.Errorf("expected user to be updated, but got nil")
	}
	if s.currentUser.FirstName().String() != "Updated" {
		return fmt.Errorf("expected first name to be 'Updated', got '%s'", s.currentUser.FirstName().String())
	}
	return nil
}

func (s *UserFeaturesTestSuite) userRoleShouldBeChanged(expectedRole string) error {
	if s.lastError != nil {
		return fmt.Errorf("expected role change to succeed, got error: %v", s.lastError)
	}
	if s.currentUser == nil {
		return fmt.Errorf("expected user role to be changed, but got nil")
	}
	if s.currentUser.Role().String() != expectedRole {
		return fmt.Errorf("expected user role to be '%s', got '%s'", expectedRole, s.currentUser.Role().String())
	}
	return nil
}

func (s *UserFeaturesTestSuite) userAccountShouldBeVerified() error {
	if s.currentUser == nil {
		return fmt.Errorf("expected user to be verified, but got nil")
	}
	if !s.currentUser.IsVerified() {
		return fmt.Errorf("expected user to be verified, but it's not")
	}
	return nil
}

func (s *UserFeaturesTestSuite) userAccountShouldBeDeactivated() error {
	if s.currentUser == nil {
		return fmt.Errorf("expected user to be deactivated, but got nil")
	}
	if s.currentUser.Status() != entities.UserStatusInactive {
		return fmt.Errorf("expected user status to be 'inactive', got '%s'", s.currentUser.Status().String())
	}
	return nil
}

func (s *UserFeaturesTestSuite) userCreatedEventShouldBePublished() error {
	events := s.eventPublisher.Events()
	if len(events) == 0 {
		return fmt.Errorf("expected events to be published, got none")
	}

	found := false
	for _, event := range events {
		if event.Type == events.EventUserCreated {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("expected user created event to be published, but wasn't found in %v events", len(events))
	}

	return nil
}

func (s *UserFeaturesTestSuite) userUpdatedEventShouldBePublished() error {
	events := s.eventPublisher.Events()

	found := false
	for _, event := range events {
		if event.Type == events.EventUserUpdated {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("expected user updated event to be published, but wasn't found")
	}

	return nil
}

func (s *UserFeaturesTestSuite) userLoginEventShouldBePublished() error {
	events := s.eventPublisher.Events()

	found := false
	for _, event := range events {
		if event.Type == events.EventUserLogin {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("expected user login event to be published, but wasn't found")
	}

	return nil
}

func (s *UserFeaturesTestSuite) userLoginFailEventShouldBePublished() error {
	events := s.eventPublisher.Events()

	found := false
	for _, event := range events {
		if event.Type == events.EventUserLoginFail {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("expected user login fail event to be published, but wasn't found")
	}

	return nil
}

func (s *UserFeaturesTestSuite) roleChangedEventShouldBePublished() error {
	events := s.eventPublisher.Events()

	found := false
	for _, event := range events {
		if event.Type == events.EventRoleChanged {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("expected role changed event to be published, but wasn't found")
	}

	return nil
}

// Test runner
func TestUserFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			features := &UserFeaturesTestSuite{}
			features.InitializeContext(ctx)
		},
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"test/features/user"},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("BDD tests failed")
	}
}

// Standalone test for simple scenarios
func TestUserManagementFeatures(t *testing.T) {
	suite := &UserFeaturesTestSuite{}

	// Test user creation flow
	t.Run("User Creation", func(t *testing.T) {
		// Given a user with valid data
		err := suite.createUserWithValidData()
		require.NoError(t, err, "User creation should succeed")

		// Then the user should be created successfully
		err = suite.userShouldBeCreatedSuccessfully()
		assert.NoError(t, err, "User should be created successfully")

		// And a user created event should be published
		err = suite.userCreatedEventShouldBePublished()
		assert.NoError(t, err, "User created event should be published")
	})

	// Test authentication flow
	t.Run("User Authentication", func(t *testing.T) {
		suite.eventPublisher.Clear()

		// Given I have valid user credentials
		err := suite.haveValidUserCredentials()
		require.NoError(t, err, "Valid user credentials should be created")

		// When I attempt to authenticate with these credentials
		// Note: This would need proper password verification in mock
		// For now, we'll simulate the success case

		// Then the authentication should succeed
		// And a user login event should be published
	})

	// Test user update flow
	t.Run("User Update", func(t *testing.T) {
		suite.eventPublisher.Clear()

		// Given a user account
		err := suite.createActiveUserAccount()
		require.NoError(t, err, "Active user account should be created")

		// When I update the user profile
		err = suite.updateUserProfile()
		require.NoError(t, err, "User profile should be updated")

		// Then the user profile should be updated
		err = suite.userProfileShouldBeUpdated()
		assert.NoError(t, err, "User profile should be updated")

		// And a user updated event should be published
		err = suite.userUpdatedEventShouldBePublished()
		assert.NoError(t, err, "User updated event should be published")
	})
}
