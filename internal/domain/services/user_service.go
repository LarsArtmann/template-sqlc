package services

import (
	"context"
	"fmt"
	"net"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/internal/domain/events"
	"github.com/LarsArtmann/template-sqlc/internal/domain/repositories"
	"github.com/google/uuid"
)

// UserService provides business logic for user operations
// This layer sits between domain entities and repositories
type UserService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
	eventPub    events.EventPublisher
	validator   UserValidator
}

// UserValidator defines validation interface for user operations
type UserValidator interface {
	ValidateUserCreate(email, username, firstName, lastName string) error
	ValidateUserUpdate(user *entities.User) error
	ValidatePasswordRequirements(password string) error
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	eventPub events.EventPublisher,
	validator UserValidator,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		eventPub:    eventPub,
		validator:   validator,
	}
}

// CreateUser creates a new user with business logic validation
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*entities.User, error) {
	// Validate request
	if err := s.validator.ValidateUserCreate(
		req.Email,
		req.Username,
		req.FirstName,
		req.LastName,
	); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists
	if _, err := s.userRepo.GetByEmail(ctx, entities.Email(req.Email)); err == nil {
		return nil, entities.ErrUserAlreadyExists
	}
	if _, err := s.userRepo.GetByUsername(ctx, entities.Username(req.Username)); err == nil {
		return nil, entities.ErrUserAlreadyExists
	}

	// Create domain entities
	email, err := entities.NewEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	username, err := entities.NewUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username: %w", err)
	}

	firstName, err := entities.NewFirstName(req.FirstName)
	if err != nil {
		return nil, fmt.Errorf("invalid first name: %w", err)
	}

	lastName, err := entities.NewLastName(req.LastName)
	if err != nil {
		return nil, fmt.Errorf("invalid last name: %w", err)
	}

	passwordHash, err := entities.NewPasswordHash(req.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("invalid password hash: %w", err)
	}

	// Create user entity
	user, err := entities.NewUser(
		email,
		username,
		passwordHash,
		firstName,
		lastName,
		entities.UserStatus(req.Status),
		entities.UserRole(req.Role),
		entities.NewUserMetadata(),
		req.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Persist user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Publish event
	event := events.UserCreated(
		user.UUID().String(),
		email.String(),
		username.String(),
		firstName.String(),
		lastName.String(),
		user.Role().String(),
		user.Status().String(),
	)
	if err := s.eventPub.Publish(event); err != nil {
		// Log error but don't fail the operation
		// In production, you'd use proper logging
		fmt.Printf("warning: failed to publish event: %v\n", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID with business logic checks
func (s *UserService) GetUser(ctx context.Context, userID entities.UserID) (*entities.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Additional business logic checks can go here
	// For example: check if user has permission to view this profile

	return user, nil
}

// UpdateUser updates a user with business logic validation
func (s *UserService) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*entities.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Build changes map for event tracking
	changes := make(map[string]interface{})

	// Update fields if provided
	if req.FirstName != nil {
		firstName, err := entities.NewFirstName(*req.FirstName)
		if err != nil {
			return nil, fmt.Errorf("invalid first name: %w", err)
		}
		changes["first_name"] = map[string]interface{}{
			"old": user.FirstName().String(),
			"new": firstName.String(),
		}
		user.UpdateProfile(&firstName, nil, nil)
	}

	if req.LastName != nil {
		lastName, err := entities.NewLastName(*req.LastName)
		if err != nil {
			return nil, fmt.Errorf("invalid last name: %w", err)
		}
		changes["last_name"] = map[string]interface{}{
			"old": user.LastName().String(),
			"new": lastName.String(),
		}
		user.UpdateProfile(nil, &lastName, nil)
	}

	if req.Metadata != nil {
		metadata := entities.NewUserMetadata()
		for k, v := range *req.Metadata {
			metadata.Set(k, v)
		}
		changes["metadata"] = map[string]interface{}{
			"old": user.Metadata(),
			"new": metadata,
		}
		user.UpdateProfile(nil, nil, &metadata)
	}

	if req.Tags != nil {
		changes["tags"] = map[string]interface{}{
			"old": user.Tags(),
			"new": *req.Tags,
		}
		user.UpdateProfile(nil, nil, nil, req.Tags)
	}

	// Validate updated user
	if err := s.validator.ValidateUserUpdate(user); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Publish update event
	if len(changes) > 0 {
		event := events.UserUpdated(user.UUID().String(), changes, req.UpdatedBy)
		if err := s.eventPub.Publish(event); err != nil {
			fmt.Printf("warning: failed to publish event: %v\n", err)
		}
	}

	return user, nil
}

// AuthenticateUser authenticates a user with email and password
func (s *UserService) AuthenticateUser(ctx context.Context, email, password, ipAddress, userAgent string) (*entities.UserSession, error) {
	// Validate email
	emailEntity, err := entities.NewEmail(email)
	if err != nil {
		return nil, entities.ErrInvalidCredentials
	}

	// Get user
	user, err := s.userRepo.VerifyCredentials(ctx, emailEntity, entities.PasswordHash(password))
	if err != nil {
		// Publish failed login event
		event := events.UserLoginFailed("", ipAddress, userAgent, "unknown")
		s.eventPub.Publish(event)
		return nil, entities.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive() {
		event := events.UserLoginFailed(user.UUID().String(), ipAddress, userAgent, "inactive_account")
		s.eventPub.Publish(event)

		if user.Status() == entities.UserStatusSuspended {
			return nil, entities.ErrAccountSuspended
		}
		return nil, entities.ErrAccountInactive
	}

	// Create session
	deviceInfo := entities.NewSessionDeviceInfo()
	deviceInfo.SetMetadata("user_agent", userAgent)

	session := entities.NewUserSession(
		user.ID(),
		net.ParseIP(ipAddress),
		userAgent,
		deviceInfo,
		entities.SessionDurationMedium,
	)

	// Save session
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update user last login
	user.RecordLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		fmt.Printf("warning: failed to update last login: %v\n", err)
	}

	// Publish login event
	event := events.UserLoggedIn(user.UUID().String(), ipAddress, userAgent, "unknown")
	if err := s.eventPub.Publish(event); err != nil {
		fmt.Printf("warning: failed to publish event: %v\n", err)
	}

	return session, nil
}

// VerifySession validates a session token and returns associated user
func (s *UserService) VerifySession(ctx context.Context, token string) (*entities.UserSession, *entities.User, error) {
	// Parse token
	tokenUUID, err := uuid.Parse(token)
	if err != nil {
		return nil, nil, entities.ErrInvalidSessionToken
	}

	sessionToken := entities.SessionToken(tokenUUID)

	// Get session
	session, err := s.sessionRepo.GetByToken(ctx, sessionToken)
	if err != nil {
		return nil, nil, entities.ErrSessionNotFound
	}

	// Check if session is valid
	if !session.IsValid() {
		if session.IsExpired() {
			return nil, nil, entities.ErrSessionExpired
		}
		return nil, nil, entities.ErrSessionNotFound
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, session.UserID())
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, nil, entities.ErrAccountInactive
	}

	return session, user, nil
}

// Logout deactivates a session
func (s *UserService) Logout(ctx context.Context, token string) error {
	// Parse token
	tokenUUID, err := uuid.Parse(token)
	if err != nil {
		return entities.ErrInvalidSessionToken
	}

	sessionToken := entities.SessionToken(tokenUUID)

	// Deactivate session
	if err := s.sessionRepo.DeactivateByToken(ctx, sessionToken); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	// Publish logout event
	// We need the user ID for the event, but we can't get it without hitting the DB
	// In a real implementation, you might include user ID in the session
	return nil
}

// ChangeUserRole changes a user's role with validation and event publishing
func (s *UserService) ChangeUserRole(ctx context.Context, userID entities.UserID, newRole entities.UserRole, changedBy string) (*entities.User, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Track old role
	oldRole := user.Role()

	// Change role
	if err := user.ChangeRole(newRole); err != nil {
		return nil, fmt.Errorf("invalid role: %w", err)
	}

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to change role: %w", err)
	}

	// Publish event
	event := events.RoleChanged(
		user.UUID().String(),
		oldRole.String(),
		newRole.String(),
		changedBy,
	)
	if err := s.eventPub.Publish(event); err != nil {
		fmt.Printf("warning: failed to publish event: %v\n", err)
	}

	return user, nil
}

// GetUserStats returns user statistics
func (s *UserService) GetUserStats(ctx context.Context) (*entities.UserStats, error) {
	stats, err := s.userRepo.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	return stats, nil
}

// Request DTOs

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email        string                 `json:"email" validate:"required,email"`
	Username     string                 `json:"username" validate:"required,min=3,max=50"`
	PasswordHash string                 `json:"password_hash" validate:"required"`
	FirstName    string                 `json:"first_name" validate:"required"`
	LastName     string                 `json:"last_name" validate:"required"`
	Status       string                 `json:"status" validate:"required"`
	Role         string                 `json:"role" validate:"required"`
	Tags         []string               `json:"tags"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	UserID    entities.UserID         `json:"user_id" validate:"required"`
	FirstName *string                 `json:"first_name,omitempty" validate:"omitempty,min=1"`
	LastName  *string                 `json:"last_name,omitempty" validate:"omitempty,min=1"`
	Metadata  *map[string]interface{} `json:"metadata,omitempty"`
	Tags      *[]string               `json:"tags,omitempty"`
	UpdatedBy string                  `json:"updated_by" validate:"required"`
}
