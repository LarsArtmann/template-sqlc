package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents the core domain entity for a user
// This is INDEPENDENT of database representation
type User struct {
	id          UserID
	uuid        uuid.UUID
	email       Email
	username    Username
	password    PasswordHash
	firstName   FirstName
	lastName    LastName
	status      UserStatus
	role        UserRole
	isVerified  bool
	metadata    UserMetadata
	tags        []string
	createdAt   time.Time
	updatedAt   time.Time
	lastLoginAt *time.Time
}

// UserID is a strongly-typed user identifier
type UserID int64

func (id UserID) Int64() int64   { return int64(id) }
func (id UserID) String() string { return fmt.Sprintf("user:%d", id) }

// Email represents a validated email address
type Email string

func NewEmail(email string) (Email, error) {
	if !isValidEmail(email) {
		return "", ErrInvalidEmail
	}
	return Email(strings.ToLower(strings.TrimSpace(email))), nil
}

func (e Email) String() string { return string(e) }

// Username represents a validated username
type Username string

func NewUsername(username string) (Username, error) {
	if len(username) < 3 || len(username) > 50 {
		return "", ErrInvalidUsername
	}
	return Username(strings.TrimSpace(username)), nil
}

func (u Username) String() string { return string(u) }

// PasswordHash represents a secure password hash
type PasswordHash string

func NewPasswordHash(hash string) (PasswordHash, error) {
	if len(hash) < 32 { // Minimum bcrypt length
		return "", ErrInvalidPasswordHash
	}
	return PasswordHash(hash), nil
}

func (p PasswordHash) String() string { return string(p) }

// FirstName represents a validated first name
type FirstName string

func NewFirstName(name string) (FirstName, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return "", ErrInvalidFirstName
	}
	return FirstName(name), nil
}

func (f FirstName) String() string { return string(f) }

// LastName represents a validated last name
type LastName string

func NewLastName(name string) (LastName, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return "", ErrInvalidLastName
	}
	return LastName(name), nil
}

func (l LastName) String() string { return string(l) }

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusPending   UserStatus = "pending"
)

func (s UserStatus) String() string { return string(s) }
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended, UserStatusPending:
		return true
	default:
		return false
	}
}

// UserRole represents user role in system
type UserRole string

const (
	UserRoleUser      UserRole = "user"
	UserRoleAdmin     UserRole = "admin"
	UserRoleModerator UserRole = "moderator"
)

func (r UserRole) String() string { return string(r) }
func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleUser, UserRoleAdmin, UserRoleModerator:
		return true
	default:
		return false
	}
}

// UserMetadata represents flexible user metadata
type UserMetadata map[string]interface{}

func NewUserMetadata() UserMetadata {
	return make(UserMetadata)
}

func (m UserMetadata) Set(key string, value interface{}) {
	m[key] = value
}

func (m UserMetadata) Get(key string) (interface{}, bool) {
	val, ok := m[key]
	return val, ok
}

// NewUser creates a new user entity with validation
func NewUser(
	email Email,
	username Username,
	password PasswordHash,
	firstName FirstName,
	lastName LastName,
	status UserStatus,
	role UserRole,
	metadata UserMetadata,
	tags []string,
) (*User, error) {
	if !status.IsValid() {
		return nil, ErrInvalidUserStatus
	}
	if !role.IsValid() {
		return nil, ErrInvalidUserRole
	}

	now := time.Now()
	return &User{
		uuid:       uuid.New(),
		email:      email,
		username:   username,
		password:   password,
		firstName:  firstName,
		lastName:   lastName,
		status:     status,
		role:       role,
		isVerified: false,
		metadata:   metadata,
		tags:       tags,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// Methods for the User entity

func (u *User) ID() UserID              { return u.id }
func (u *User) UUID() uuid.UUID         { return u.uuid }
func (u *User) Email() Email            { return u.email }
func (u *User) Username() Username      { return u.username }
func (u *User) FirstName() FirstName    { return u.firstName }
func (u *User) LastName() LastName      { return u.lastName }
func (u *User) Status() UserStatus      { return u.status }
func (u *User) Role() UserRole          { return u.role }
func (u *User) IsVerified() bool        { return u.isVerified }
func (u *User) Metadata() UserMetadata  { return u.metadata }
func (u *User) Tags() []string          { return u.tags }
func (u *User) CreatedAt() time.Time    { return u.createdAt }
func (u *User) UpdatedAt() time.Time    { return u.updatedAt }
func (u *User) LastLoginAt() *time.Time { return u.lastLoginAt }

// IsActive returns true if user status is active
func (u *User) IsActive() bool {
	return u.status == UserStatusActive
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(
	firstName *FirstName,
	lastName *LastName,
	metadata *UserMetadata,
	tags *[]string,
) error {
	if firstName != nil {
		u.firstName = *firstName
	}
	if lastName != nil {
		u.lastName = *lastName
	}
	if metadata != nil {
		u.metadata = *metadata
	}
	if tags != nil {
		u.tags = *tags
	}

	u.updatedAt = time.Now()
	return nil
}

// ChangeStatus updates user status with validation
func (u *User) ChangeStatus(status UserStatus) error {
	if !status.IsValid() {
		return ErrInvalidUserStatus
	}
	u.status = status
	u.updatedAt = time.Now()
	return nil
}

// ChangeRole updates user role with validation
func (u *User) ChangeRole(role UserRole) error {
	if !role.IsValid() {
		return ErrInvalidUserRole
	}
	u.role = role
	u.updatedAt = time.Now()
	return nil
}

// Verify marks user as verified
func (u *User) Verify() {
	u.isVerified = true
	u.updatedAt = time.Now()
}

// RecordLogin updates last login time
func (u *User) RecordLogin() {
	now := time.Now()
	u.lastLoginAt = &now
	u.updatedAt = now
}

// AddTag adds a tag to user if not already present
func (u *User) AddTag(tag string) {
	for _, existingTag := range u.tags {
		if existingTag == tag {
			return
		}
	}
	u.tags = append(u.tags, tag)
	u.updatedAt = time.Now()
}

// RemoveTag removes a tag from user
func (u *User) RemoveTag(tag string) {
	for i, existingTag := range u.tags {
		if existingTag == tag {
			u.tags = append(u.tags[:i], u.tags[i+1:]...)
			u.updatedAt = time.Now()
			return
		}
	}
}
