package events

import (
	"time"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
)

// UserEvent represents a domain event related to users
type UserEvent struct {
	ID        entities.IDID   `json:"id"`
	Type      EventType       `json:"type"`
	UserID    entities.UserID `json:"user_id"`
	Data      any             `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
	Version   string          `json:"version"`
}

// EventType represents the type of domain event
type EventType string

const (
	// User lifecycle events
	EventUserCreated     EventType = "user.created"
	EventUserUpdated     EventType = "user.updated"
	EventUserDeleted     EventType = "user.deleted"
	EventUserActivated   EventType = "user.activated"
	EventUserDeactivated EventType = "user.deactivated"
	EventUserSuspended   EventType = "user.suspended"

	// Authentication events
	EventUserLogin     EventType = "user.login"
	EventUserLogout    EventType = "user.logout"
	EventUserLoginFail EventType = "user.login.failed"

	// Verification events
	EventUserVerified              EventType = "user.verified"
	EventUserVerificationRequested EventType = "user.verification.requested"

	// Password events
	EventPasswordChanged        EventType = "password.changed"
	EventPasswordReset          EventType = "password.reset"
	EventPasswordResetRequested EventType = "password.reset.requested"

	// Profile events
	EventProfileUpdated EventType = "profile.updated"
	EventRoleChanged    EventType = "role.changed"
)

// UserCreatedEvent data for user creation
type UserCreatedEvent struct {
	UserID    entities.UserID `json:"user_id"`
	Email     string          `json:"email"`
	Username  string          `json:"username"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Role      string          `json:"role"`
	Status    string          `json:"status"`
}

// UserUpdatedEvent data for user updates
type UserUpdatedEvent struct {
	UserID    entities.UserID `json:"user_id"`
	Changes   map[string]any  `json:"changes"`
	UpdatedBy entities.UserID `json:"updated_by"`
}

// UserLoginEvent data for user login
type UserLoginEvent struct {
	UserID    entities.UserID `json:"user_id"`
	IPAddress string          `json:"ip_address"`
	UserAgent string          `json:"user_agent"`
	Device    string          `json:"device"`
	Success   bool            `json:"success"`
}

// UserVerifiedEvent data for user verification
type UserVerifiedEvent struct {
	UserID    entities.UserID `json:"user_id"`
	Method    string          `json:"method"`
	Timestamp time.Time       `json:"timestamp"`
}

// RoleChangedEvent data for role changes
type RoleChangedEvent struct {
	UserID    entities.UserID `json:"user_id"`
	OldRole   string          `json:"old_role"`
	NewRole   string          `json:"new_role"`
	ChangedBy entities.UserID `json:"changed_by"`
}

// NewUserEvent creates a new user domain event
func NewUserEvent(eventType EventType, userID entities.UserID, data any) *UserEvent {
	return &UserEvent{
		ID:        entities.AsIDID(time.Now().UnixNano()),
		Type:      eventType,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// UserCreated creates a user created event
func UserCreated(
	userID entities.UserID,
	email, username, firstName, lastName, role, status string,
) *UserEvent {
	data := UserCreatedEvent{
		UserID:    userID,
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		Status:    status,
	}
	return NewUserEvent(EventUserCreated, userID, data)
}

// UserUpdated creates a user updated event
func UserUpdated(
	userID entities.UserID,
	changes map[string]any,
	updatedBy entities.UserID,
) *UserEvent {
	data := UserUpdatedEvent{
		UserID:    userID,
		Changes:   changes,
		UpdatedBy: updatedBy,
	}
	return NewUserEvent(EventUserUpdated, userID, data)
}

// UserLoginAttempt creates a user login attempt event
func UserLoginAttempt(userID entities.UserID, ipAddress, userAgent, device string, success bool, eventType EventType) *UserEvent {
	data := UserLoginEvent{
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Device:    device,
		Success:   success,
	}
	return NewUserEvent(eventType, userID, data)
}

// UserLoggedIn creates a user login event
func UserLoggedIn(userID entities.UserID, ipAddress, userAgent, device string) *UserEvent {
	return UserLoginAttempt(userID, ipAddress, userAgent, device, true, EventUserLogin)
}

// UserLoginFailed creates a user login failed event
func UserLoginFailed(userID entities.UserID, ipAddress, userAgent, device string) *UserEvent {
	return UserLoginAttempt(userID, ipAddress, userAgent, device, false, EventUserLoginFail)
}

// UserVerified creates a user verified event
func UserVerified(userID entities.UserID, method string) *UserEvent {
	data := UserVerifiedEvent{
		UserID:    userID,
		Method:    method,
		Timestamp: time.Now(),
	}
	return NewUserEvent(EventUserVerified, userID, data)
}

// RoleChanged creates a role changed event
func RoleChanged(
	userID entities.UserID,
	oldRole, newRole string,
	changedBy entities.UserID,
) *UserEvent {
	data := RoleChangedEvent{
		UserID:    userID,
		OldRole:   oldRole,
		NewRole:   newRole,
		ChangedBy: changedBy,
	}
	return NewUserEvent(EventRoleChanged, userID, data)
}

// EventPublisher interface for publishing domain events
type EventPublisher interface {
	Publish(event *UserEvent) error
	PublishBatch(events []*UserEvent) error
}

// InMemoryEventPublisher is a simple in-memory event publisher
type InMemoryEventPublisher struct {
	events []*UserEvent
}

func NewInMemoryEventPublisher() *InMemoryEventPublisher {
	return &InMemoryEventPublisher{
		events: make([]*UserEvent, 0),
	}
}

func (p *InMemoryEventPublisher) Publish(event *UserEvent) error {
	p.events = append(p.events, event)
	return nil
}

func (p *InMemoryEventPublisher) PublishBatch(events []*UserEvent) error {
	p.events = append(p.events, events...)
	return nil
}

func (p *InMemoryEventPublisher) Events() []*UserEvent {
	return p.events
}

func (p *InMemoryEventPublisher) Clear() {
	p.events = make([]*UserEvent, 0)
}

// EventType returns the string representation of EventType
func (e EventType) String() string {
	return string(e)
}

// IsValidEventType checks if the event type is valid
func (e EventType) IsValid() bool {
	validTypes := map[EventType]bool{
		EventUserCreated:               true,
		EventUserUpdated:               true,
		EventUserDeleted:               true,
		EventUserActivated:             true,
		EventUserDeactivated:           true,
		EventUserSuspended:             true,
		EventUserLogin:                 true,
		EventUserLogout:                true,
		EventUserLoginFail:             true,
		EventUserVerified:              true,
		EventUserVerificationRequested: true,
		EventPasswordChanged:           true,
		EventPasswordReset:             true,
		EventPasswordResetRequested:    true,
		EventProfileUpdated:            true,
		EventRoleChanged:               true,
	}
	return validTypes[e]
}
