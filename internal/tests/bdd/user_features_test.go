package bdd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/events"
	"github.com/LarsArtmann/template-sqlc/internal/domain/services"
	"github.com/LarsArtmann/template-sqlc/internal/tests/integration"
	apperrors "github.com/LarsArtmann/template-sqlc/pkg/errors"
	"github.com/LarsArtmann/template-sqlc/pkg/validation"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPasswordHash is a bcrypt hash for "test_password" used in tests.
//
//nolint:gosec,G101 // This is a test constant, not a production secret.
const TestPasswordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsQ5pPjZ5yVlWK5WAe"

// UserFeaturesTestSuite contains BDD tests for user functionality.
type UserFeaturesTestSuite struct {
	userService    *services.UserService
	userRepo       *integration.MockUserRepository
	sessionRepo    *integration.MockSessionRepository
	eventPublisher *events.InMemoryEventPublisher
	validator      *validation.UserValidator
	currentUser    *entities.User
	currentSession *entities.UserSession
	lastError      error
	stats          *entities.UserStats
}

// InitializeContext sets up the test context.
func (s *UserFeaturesTestSuite) InitializeContext(ctx *godog.ScenarioContext) {
	s.eventPublisher = events.NewInMemoryEventPublisher()
	s.validator = validation.NewUserValidator()
	s.userRepo = integration.NewMockUserRepository()
	s.sessionRepo = integration.NewMockSessionRepository()

	s.userService = services.NewUserService(
		s.userRepo,
		s.sessionRepo,
		s.eventPublisher,
		s.validator,
	)

	// Background steps
	ctx.Given(`^a clean user system$`, s.cleanUserSystem)
	ctx.Given(`^the event publisher is cleared$`, s.clearEventPublisher)

	// Given steps
	ctx.Given(`^a user with email "([^"]*)" and username "([^"]*)"$`, s.createUserWithEmailUsername)
	ctx.Given(`^a user account with status "([^"]*)"$`, s.createUserWithStatus)
	ctx.Given(`^an active user account$`, s.createActiveUserAccount)
	ctx.Given(`^I have valid user credentials$`, s.haveValidUserCredentials)
	ctx.Given(`^I have invalid user credentials$`, s.haveInvalidUserCredentials)
	ctx.Given(`^an inactive user account$`, s.createInactiveUserAccount)
	ctx.Given(`^a suspended user account$`, s.createSuspendedUserAccount)
	ctx.Given(
		`^I have valid user credentials for this account$`,
		s.haveValidCredentialsForCurrentAccount,
	)
	ctx.Given(`^multiple user accounts with different statuses$`, s.createMultipleStatusAccounts)
	ctx.Given(`^a user account with status "([^"]*)"$`, s.createUserWithStatus)

	// When steps
	ctx.When(`^I create a user with valid data$`, s.createUserWithValidData)
	ctx.When(`^I create a user with email "([^"]*)"$`, s.createUserWithEmail)
	ctx.When(`^I create a user with username "([^"]*)"$`, s.createUserWithUsername)
	ctx.When(`^I attempt to authenticate with these credentials$`, s.authenticateWithCredentials)
	ctx.When(`^I update the user profile$`, s.updateUserProfile)
	ctx.When(`^I change the user role to "([^"]*)"$`, s.changeUserRole)
	ctx.When(`^I verify the user account$`, s.verifyUserAccount)
	ctx.When(`^I deactivate the user account$`, s.deactivateUserAccount)
	ctx.When(`^I get user statistics$`, s.getUserStatistics)
	ctx.When(`^the user role is "([^"]*)"$`, s.setUserRole)
	ctx.When(`^the user status is "([^"]*)"$`, s.setUserStatus)
	ctx.When(`^the session expires$`, s.expireSession)
	ctx.When(`^I attempt to authenticate from multiple devices$`, s.authenticateFromMultipleDevices)

	// Then steps
	ctx.Then(`^the user should be created successfully$`, s.userShouldBeCreatedSuccessfully)
	ctx.Then(`^the user should have ID (\d+)$`, s.userShouldHaveID)
	ctx.Then(`^I should receive a validation error$`, s.shouldReceiveValidationError)
	ctx.Then(
		`^I should receive a "user already exists" error$`,
		s.shouldReceiveUserAlreadyExistsError,
	)
	ctx.Then(`^the authentication should succeed$`, s.authenticationShouldSucceed)
	ctx.Then(`^the authentication should fail$`, s.authenticationShouldFail)
	ctx.Then(`^I should receive a "user not found" error$`, s.shouldReceiveUserNotFoundError)
	ctx.Then(
		`^I should receive a "invalid credentials" error$`,
		s.shouldReceiveInvalidCredentialsError,
	)
	ctx.Then(`^the session should be created$`, s.sessionShouldBeCreated)
	ctx.Then(`^the user profile should be updated$`, s.userProfileShouldBeUpdated)
	ctx.Then(`^the user role should be changed to "([^"]*)"$`, s.userRoleShouldBeChanged)
	ctx.Then(`^the user account should be verified$`, s.userAccountShouldBeVerified)
	ctx.Then(`^the user account should be deactivated$`, s.userAccountShouldBeDeactivated)
	ctx.Then(`^the user account should be in pending state$`, s.userAccountShouldBePending)
	ctx.Then(`^the user account should be suspended$`, s.userAccountShouldBeSuspended)
	ctx.Then(`^the user should not be able to authenticate$`, s.userShouldNotAuthenticate)
	ctx.Then(`^the user should have the specified metadata$`, s.userShouldHaveSpecifiedMetadata)
	ctx.Then(`^the user should have the specified tags$`, s.userShouldHaveSpecifiedTags)
	ctx.Then(`^the user should have admin privileges$`, s.userShouldHaveAdminPrivileges)
	ctx.Then(`^the user should have moderator privileges$`, s.userShouldHaveModeratorPrivileges)
	ctx.Then(
		`^the statistics should include counts for each status$`,
		s.statisticsShouldIncludeCounts,
	)
	ctx.Then(`^the session should no longer be valid$`, s.sessionShouldNotBeValid)
	ctx.Then(`^multiple sessions should be created$`, s.multipleSessionsShouldBeCreated)
	ctx.Then(`^all sessions should be active$`, s.allSessionsShouldBeActive)
	ctx.Then(`^a user created event should be published$`, s.userCreatedEventShouldBePublished)
	ctx.Then(`^a user updated event should be published$`, s.userUpdatedEventShouldBePublished)
	ctx.Then(`^a user login event should be published$`, s.userLoginEventShouldBePublished)
	ctx.Then(
		`^a user login failed event should be published$`,
		s.userLoginFailEventShouldBePublished,
	)
	ctx.Then(`^a role changed event should be published$`, s.roleChangedEventShouldBePublished)
	ctx.Then(`^a user verified event should be published$`, s.userVerifiedEventShouldBePublished)
}

// Background steps

func (s *UserFeaturesTestSuite) cleanUserSystem() error {
	s.userRepo = integration.NewMockUserRepository()
	s.sessionRepo = integration.NewMockSessionRepository()
	s.userService = services.NewUserService(
		s.userRepo,
		s.sessionRepo,
		s.eventPublisher,
		s.validator,
	)
	s.currentUser = nil
	s.currentSession = nil
	s.lastError = nil

	return nil
}

func (s *UserFeaturesTestSuite) clearEventPublisher() error {
	s.eventPublisher.Clear()

	return nil
}

// Given steps

func (s *UserFeaturesTestSuite) createUserWithEmailUsername(email, username string) error {
	req := &services.CreateUserRequest{
		Email:        email,
		Username:     username,
		PasswordHash: TestPasswordHash,
		FirstName:    "Test",
		LastName:     "User",
		Status:       "active",
		Role:         "user",
		Tags:         []string{"test"},
		Metadata:     map[string]any{"source": "bdd"},
	}

	user, err := s.userService.CreateUser(context.Background(), req)
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) createUserWithStatus(status string) error {
	req := &services.CreateUserRequest{
		Email:        "status@example.com",
		Username:     fmt.Sprintf("statususer%d", time.Now().UnixNano()),
		PasswordHash: TestPasswordHash,
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

func (s *UserFeaturesTestSuite) createUserWithPrefix(prefix string) (*entities.User, error) {
	req := &services.CreateUserRequest{
		Email:        prefix + "@example.com",
		Username:     prefix,
		PasswordHash: TestPasswordHash,
		FirstName:    strings.Title(prefix),
		LastName:     "User",
		Status:       "active",
		Role:         "user",
	}

	return s.userService.CreateUser(context.Background(), req)
}

func (s *UserFeaturesTestSuite) haveValidUserCredentials() error {
	user, err := s.createUserWithPrefix("valid")
	s.currentUser = user
	s.lastError = err

	if user != nil {
		s.userRepo.SetPasswordVerification(user.Email().String(), "correct_password")
	}

	return nil
}

func (s *UserFeaturesTestSuite) haveInvalidUserCredentials() error {
	user, err := s.createUserWithPrefix("invalid")
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) createInactiveUserAccount() error {
	return s.createUserWithStatus("inactive")
}

func (s *UserFeaturesTestSuite) createSuspendedUserAccount() error {
	return s.createUserWithStatus("suspended")
}

func (s *UserFeaturesTestSuite) haveValidCredentialsForCurrentAccount() error {
	if s.currentUser == nil {
		return errors.New("no current user to set credentials for")
	}

	s.userRepo.SetPasswordVerification(s.currentUser.Email().String(), "correct_password")

	return nil
}

func (s *UserFeaturesTestSuite) createMultipleStatusAccounts() error {
	statuses := []string{"active", "inactive", "suspended", "pending"}
	for i, status := range statuses {
		req := &services.CreateUserRequest{
			Email:        fmt.Sprintf("user%d@example.com", i),
			Username:     fmt.Sprintf("user%d", i),
			PasswordHash: TestPasswordHash,
			FirstName:    "Multi",
			LastName:     "User",
			Status:       status,
			Role:         "user",
		}

		_, err := s.userService.CreateUser(context.Background(), req)
		if err != nil {
			return err
		}
	}

	return nil
}

// When steps

func (s *UserFeaturesTestSuite) createUserWithValidData() error {
	req := &services.CreateUserRequest{
		Email:        "valid@example.com",
		Username:     "validuser",
		PasswordHash: TestPasswordHash,
		FirstName:    "Valid",
		LastName:     "User",
		Status:       "active",
		Role:         "user",
		Tags:         []string{"valid", "test"},
		Metadata:     map[string]any{"test": true},
	}

	user, err := s.userService.CreateUser(context.Background(), req)
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) createUserWithEmail(email string) error {
	req := &services.CreateUserRequest{
		Email:        email,
		Username:     fmt.Sprintf("emailuser%d", time.Now().UnixNano()),
		PasswordHash: TestPasswordHash,
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
		Email:        fmt.Sprintf("username%d@example.com", time.Now().UnixNano()),
		Username:     username,
		PasswordHash: TestPasswordHash,
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

func (s *UserFeaturesTestSuite) authenticateCurrentUser() (*entities.UserSession, error) {
	return s.userService.AuthenticateUser(
		context.Background(),
		s.currentUser.Email().String(),
		"correct_password",
		"127.0.0.1",
		"test-user-agent",
	)
}

func (s *UserFeaturesTestSuite) authenticateWithCredentials() error {
	if s.currentUser == nil {
		return errors.New("no current user to authenticate")
	}

	session, err := s.authenticateCurrentUser()

	s.currentSession = session
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) updateUserProfile() error {
	if s.currentUser == nil {
		return errors.New("no current user to update")
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
		return errors.New("no current user for role change")
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
		return errors.New("no current user to verify")
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
		return errors.New("no current user to deactivate")
	}

	err := s.userRepo.Deactivate(context.Background(), s.currentUser.ID())
	s.lastError = err

	if err == nil {
		_ = s.currentUser.ChangeStatus(entities.UserStatusInactive)
	}

	return nil
}

func (s *UserFeaturesTestSuite) getUserStatistics() error {
	stats, err := s.userService.GetUserStats(context.Background())
	s.stats = stats
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) setUserRole(role string) error {
	if s.currentUser == nil {
		return errors.New("no current user to set role")
	}

	req := &services.CreateUserRequest{
		Email:        fmt.Sprintf("roleuser%d@example.com", time.Now().UnixNano()),
		Username:     fmt.Sprintf("roleuser%d", time.Now().UnixNano()),
		PasswordHash: TestPasswordHash,
		FirstName:    "Role",
		LastName:     "User",
		Status:       "active",
		Role:         role,
	}

	user, err := s.userService.CreateUser(context.Background(), req)
	s.currentUser = user
	s.lastError = err

	return nil
}

func (s *UserFeaturesTestSuite) setUserStatus(status string) error {
	req := &services.CreateUserRequest{
		Email:        fmt.Sprintf("statususer%d@example.com", time.Now().UnixNano()),
		Username:     fmt.Sprintf("statususer%d", time.Now().UnixNano()),
		PasswordHash: TestPasswordHash,
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

func (s *UserFeaturesTestSuite) expireSession() error {
	if s.currentSession == nil {
		return errors.New("no current session to expire")
	}

	s.currentSession.Deactivate()

	return nil
}

func (s *UserFeaturesTestSuite) authenticateFromMultipleDevices() error {
	if s.currentUser == nil {
		return errors.New("no current user for multiple device auth")
	}

	devices := []string{"desktop", "mobile", "tablet"}
	for _, device := range devices {
		session, err := s.userService.AuthenticateUser(
			context.Background(),
			s.currentUser.Email().String(),
			"correct_password",
			"127.0.0.1",
			device,
		)
		if err != nil {
			s.lastError = err

			return err
		}

		if session != nil && s.currentSession == nil {
			s.currentSession = session
		}
	}

	return nil
}

// Then steps

func (s *UserFeaturesTestSuite) userShouldBeCreatedSuccessfully() error {
	return s.expectSuccessfulOperationWithSession(func() error {
		if s.currentUser == nil {
			return errors.New("expected user to be created, but got nil")
		}

		return nil
	}, "expected user to be created successfully")
}

func (s *UserFeaturesTestSuite) expectSuccessfulOperationWithSession(
	check func() error,
	msg string,
) error {
	if s.lastError != nil {
		return fmt.Errorf("%s, got error: %w", msg, s.lastError)
	}

	return check()
}

func (s *UserFeaturesTestSuite) requireCurrentUserStringProperty(
	nilError string,
	getProperty func() string,
	expected string,
	propertyName string,
) error {
	if s.currentUser == nil {
		return errors.New(nilError)
	}

	actual := getProperty()
	if actual != expected {
		return fmt.Errorf("expected user %s to be '%s', got '%s'", propertyName, expected, actual)
	}

	return nil
}

func (s *UserFeaturesTestSuite) userShouldHaveID(expectedIDStr string) error {
	return s.requireCurrentUserStringProperty(
		"no current user",
		func() string { return s.currentUser.ID().String() },
		expectedIDStr,
		"ID",
	)
}

func (s *UserFeaturesTestSuite) shouldReceiveValidationError() error {
	if s.lastError == nil {
		return errors.New("expected validation error, got nil")
	}

	// Check for validation errors from pkg/errors
	if !apperrors.IsValidationError(s.lastError) {
		return fmt.Errorf("expected validation error, got: %w", s.lastError)
	}

	return nil
}

func (s *UserFeaturesTestSuite) shouldReceiveUserAlreadyExistsError() error {
	if s.lastError == nil {
		return errors.New("expected user already exists error, got nil")
	}

	if !errors.Is(s.lastError, entities.ErrUserNotFound) &&
		!errors.Is(s.lastError, entities.ErrUserAlreadyExists) {
		return fmt.Errorf("expected user already exists error, got: %w", s.lastError)
	}

	return nil
}

func (s *UserFeaturesTestSuite) authenticationShouldSucceed() error {
	return s.expectSuccessfulOperationWithSession(func() error {
		if s.currentSession == nil {
			return errors.New("expected session to be created, but got nil")
		}

		return nil
	}, "expected authentication to succeed")
}

func (s *UserFeaturesTestSuite) authenticationShouldFail() error {
	return s.expectUnauthorizedOrAuthError("expected authentication to fail")
}

func (s *UserFeaturesTestSuite) expectUnauthorizedOrAuthError(msg string) error {
	if s.lastError == nil {
		return fmt.Errorf("%s, got nil", msg)
	}

	if !entities.IsUnauthorizedError(s.lastError) &&
		!entities.IsAuthenticationError(s.lastError) {
		return fmt.Errorf("%s, got: %w", msg, s.lastError)
	}

	return nil
}

func (s *UserFeaturesTestSuite) shouldReceiveUserNotFoundError() error {
	if s.lastError == nil {
		return errors.New("expected user not found error, got nil")
	}

	if !entities.IsNotFoundError(s.lastError) &&
		!errors.Is(s.lastError, entities.ErrUserNotFound) {
		return fmt.Errorf("expected user not found error, got: %w", s.lastError)
	}

	return nil
}

func (s *UserFeaturesTestSuite) shouldReceiveInvalidCredentialsError() error {
	return s.expectUnauthorizedOrAuthError("expected invalid credentials error")
}

func (s *UserFeaturesTestSuite) sessionShouldBeCreated() error {
	if s.currentSession == nil {
		return errors.New("expected session to be created, but got nil")
	}

	return nil
}

func (s *UserFeaturesTestSuite) userProfileShouldBeUpdated() error {
	if s.lastError != nil {
		return fmt.Errorf("expected profile update to succeed, got error: %w", s.lastError)
	}

	if s.currentUser == nil {
		return errors.New("expected user to be updated, but got nil")
	}

	if s.currentUser.FirstName().String() != "Updated" {
		return fmt.Errorf(
			"expected first name to be 'Updated', got '%s'",
			s.currentUser.FirstName().String(),
		)
	}

	return nil
}

func (s *UserFeaturesTestSuite) userRoleShouldBeChanged(expectedRole string) error {
	if s.lastError != nil {
		return fmt.Errorf("expected role change to succeed, got error: %w", s.lastError)
	}

	return s.requireCurrentUserStringProperty(
		"expected user role to be changed, but got nil",
		func() string { return s.currentUser.Role().String() },
		expectedRole,
		"role",
	)
}

func (s *UserFeaturesTestSuite) userAccountShouldBeVerified() error {
	if s.currentUser == nil {
		return errors.New("expected user to be verified, but got nil")
	}

	if !s.currentUser.IsVerified() {
		return errors.New("expected user to be verified, but it's not")
	}

	return nil
}

func (s *UserFeaturesTestSuite) userAccountShouldHaveStatus(
	expectedStatus entities.UserStatus,
	statusName string,
) error {
	if s.currentUser == nil {
		return fmt.Errorf("expected user to be %s, but got nil", statusName)
	}

	if s.currentUser.Status() != expectedStatus {
		return fmt.Errorf(
			"expected user status to be '%s', got '%s'",
			statusName,
			s.currentUser.Status().String(),
		)
	}

	return nil
}

func (s *UserFeaturesTestSuite) userAccountShouldBeDeactivated() error {
	return s.userAccountShouldHaveStatus(entities.UserStatusInactive, "inactive")
}

func (s *UserFeaturesTestSuite) userAccountShouldBePending() error {
	return s.userAccountShouldHaveStatus(entities.UserStatusPending, "pending")
}

func (s *UserFeaturesTestSuite) userAccountShouldBeSuspended() error {
	return s.userAccountShouldHaveStatus(entities.UserStatusSuspended, "suspended")
}

func (s *UserFeaturesTestSuite) userShouldNotAuthenticate() error {
	if s.currentUser == nil {
		return errors.New("no current user to test authentication")
	}

	_, err := s.authenticateCurrentUser()
	if err == nil {
		return errors.New("expected authentication to fail for non-active user, but it succeeded")
	}

	// Check for authentication or authorization errors
	if !entities.IsAuthenticationError(err) &&
		!entities.IsUnauthorizedError(err) {
		return fmt.Errorf("expected authentication/authorization error, got: %w", err)
	}

	return nil
}

func (s *UserFeaturesTestSuite) userShouldHaveSpecifiedMetadata() error {
	if s.currentUser == nil {
		return errors.New("no current user to check metadata")
	}

	if s.currentUser.Metadata() == nil {
		return errors.New("expected user to have metadata, but got nil")
	}

	return nil
}

func (s *UserFeaturesTestSuite) userShouldHaveSpecifiedTags() error {
	if s.currentUser == nil {
		return errors.New("no current user to check tags")
	}

	if len(s.currentUser.Tags()) == 0 {
		return errors.New("expected user to have tags, but got none")
	}

	return nil
}

func (s *UserFeaturesTestSuite) userShouldHaveRole(
	expectedRole entities.UserRole,
	roleName string,
) error {
	if s.currentUser == nil {
		return errors.New("no current user to check privileges")
	}

	if s.currentUser.Role() != expectedRole {
		return fmt.Errorf(
			"expected user role to be '%s', got '%s'",
			roleName,
			s.currentUser.Role().String(),
		)
	}

	return nil
}

func (s *UserFeaturesTestSuite) userShouldHaveAdminPrivileges() error {
	return s.userShouldHaveRole(entities.UserRoleAdmin, "admin")
}

func (s *UserFeaturesTestSuite) userShouldHaveModeratorPrivileges() error {
	return s.userShouldHaveRole(entities.UserRoleModerator, "moderator")
}

func (s *UserFeaturesTestSuite) statisticsShouldIncludeCounts() error {
	if s.lastError != nil {
		return fmt.Errorf("expected stats to be retrieved, got error: %w", s.lastError)
	}

	if s.stats == nil {
		return errors.New("expected stats to be populated, got nil")
	}

	if s.stats.TotalUsers == 0 {
		return errors.New("expected total users count to be greater than 0")
	}

	return nil
}

func (s *UserFeaturesTestSuite) sessionShouldNotBeValid() error {
	if s.currentSession == nil {
		return errors.New("no current session to check validity")
	}

	if s.currentSession.IsActive() {
		return errors.New("expected session to be expired, but it's still active")
	}

	return nil
}

func (s *UserFeaturesTestSuite) multipleSessionsShouldBeCreated() error {
	if s.currentUser == nil {
		return errors.New("no current user to check sessions")
	}

	sessions, err := s.sessionRepo.GetByUserID(context.Background(), s.currentUser.ID(), false)
	if err != nil {
		return fmt.Errorf("failed to get sessions: %w", err)
	}

	if len(sessions) < 2 {
		return fmt.Errorf("expected at least 2 sessions, got %d", len(sessions))
	}

	return nil
}

func (s *UserFeaturesTestSuite) allSessionsShouldBeActive() error {
	if s.currentUser == nil {
		return errors.New("no current user to check sessions")
	}

	sessions, err := s.sessionRepo.GetByUserID(context.Background(), s.currentUser.ID(), true)
	if err != nil {
		return fmt.Errorf("failed to get sessions: %w", err)
	}

	for _, session := range sessions {
		if !session.IsActive() {
			return fmt.Errorf(
				"expected all sessions to be active, but session %s is not",
				session.Token().String(),
			)
		}
	}

	return nil
}

func (s *UserFeaturesTestSuite) userCreatedEventShouldBePublished() error {
	return s.assertEventPublished(events.EventUserCreated, "user created")
}

func (s *UserFeaturesTestSuite) assertEventPublished(
	eventType events.EventType,
	eventName string,
) error {
	userEvents := s.eventPublisher.Events()

	for _, event := range userEvents {
		if event.Type == eventType {
			return nil
		}
	}

	return fmt.Errorf("expected %s event to be published, but wasn't found", eventName)
}

func (s *UserFeaturesTestSuite) userUpdatedEventShouldBePublished() error {
	return s.assertEventPublished(events.EventUserUpdated, "user updated")
}

func (s *UserFeaturesTestSuite) userLoginEventShouldBePublished() error {
	return s.assertEventPublished(events.EventUserLogin, "user login")
}

func (s *UserFeaturesTestSuite) userLoginFailEventShouldBePublished() error {
	return s.assertEventPublished(events.EventUserLoginFail, "user login fail")
}

func (s *UserFeaturesTestSuite) roleChangedEventShouldBePublished() error {
	return s.assertEventPublished(events.EventRoleChanged, "role changed")
}

func (s *UserFeaturesTestSuite) userVerifiedEventShouldBePublished() error {
	return s.assertEventPublished(events.EventUserVerified, "user verified")
}

// Test runner.
func TestUserFeatures(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	featurePath := filepath.Join(wd, "..", "..", "..", "test", "features", "user")

	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			features := &UserFeaturesTestSuite{}
			features.InitializeContext(ctx)
		},
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{featurePath},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("BDD tests failed")
	}
}

// Standalone test for simple scenarios.
func TestUserManagementFeatures(t *testing.T) {
	t.Parallel()

	suite := &UserFeaturesTestSuite{}
	suite.eventPublisher = events.NewInMemoryEventPublisher()
	suite.validator = validation.NewUserValidator()
	suite.userRepo = integration.NewMockUserRepository()
	suite.sessionRepo = integration.NewMockSessionRepository()
	suite.userService = services.NewUserService(
		suite.userRepo,
		suite.sessionRepo,
		suite.eventPublisher,
		suite.validator,
	)

	t.Run("User Creation", func(t *testing.T) {
		err := suite.createUserWithValidData()
		require.NoError(t, err, "User creation should succeed")

		err = suite.userShouldBeCreatedSuccessfully()
		assert.NoError(t, err, "User should be created successfully")

		err = suite.userCreatedEventShouldBePublished()
		assert.NoError(t, err, "User created event should be published")
	})

	t.Run("User Authentication", func(t *testing.T) {
		suite.eventPublisher.Clear()

		err := suite.haveValidUserCredentials()
		require.NoError(t, err, "Valid user credentials should be created")
	})

	t.Run("User Update", func(t *testing.T) {
		suite.eventPublisher.Clear()

		err := suite.createActiveUserAccount()
		require.NoError(t, err, "Active user account should be created")

		err = suite.updateUserProfile()
		require.NoError(t, err, "User profile should be updated")

		err = suite.userProfileShouldBeUpdated()
		assert.NoError(t, err, "User profile should be updated")

		err = suite.userUpdatedEventShouldBePublished()
		assert.NoError(t, err, "User updated event should be published")
	})
}
