package entities

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents the core domain entity for a user
// This is INDEPENDENT of database representation.
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

// UserID is a strongly-typed user identifier.
type UserID int64

func (id UserID) Int64() int64   { return int64(id) }
func (id UserID) String() string { return fmt.Sprintf("user:%d", id) }

// IDID is a strongly-typed event identifier.
type IDID int64

func (id IDID) Int64() int64   { return int64(id) }
func (id IDID) String() string { return fmt.Sprintf("event:%d", id) }

// AsIDID converts an int64 to IDID.
func AsIDID(value int64) IDID { return IDID(value) }

// UuID is a strongly-typed UUID identifier.
type UuID string

func (id UuID) String() string { return string(id) }

// NewUuID generates a new UuID from a UUID string.
func NewUuID(s string) (UuID, error) {
	if s == "" {
		return "", nil
	}

	parsed, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}

	return UuID(parsed.String()), nil
}

// NewUuIDFromUUID creates a UuID from a uuid.UUID.
func NewUuIDFromUUID(u uuid.UUID) UuID {
	return UuID(u.String())
}

// emailRegex is a simple email validation pattern.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// isValidEmail validates an email address.
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// Email represents a validated email address.
type Email string

func NewEmail(email string) (Email, error) {
	if !isValidEmail(email) {
		return "", ErrInvalidEmail
	}

	return Email(strings.ToLower(strings.TrimSpace(email))), nil
}

func (e Email) String() string { return string(e) }

// Username represents a validated username.
type Username string

var reservedUsernames = map[string]bool{
	"admin": true, "root": true, "system": true, "moderator": true,
}

var usernameValidChars = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func NewUsername(username string) (Username, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 50 {
		return "", ErrInvalidUsername
	}

	if !usernameValidChars.MatchString(username) {
		return "", ErrInvalidUsername
	}

	if reservedUsernames[strings.ToLower(username)] {
		return "", ErrInvalidUsername
	}

	return Username(username), nil
}

func (u Username) String() string { return string(u) }

// PasswordHash represents a secure password hash.
type PasswordHash string

func NewPasswordHash(hash string) (PasswordHash, error) {
	if len(hash) < 32 { // Minimum bcrypt length
		return "", ErrInvalidPasswordHash
	}

	return PasswordHash(hash), nil
}

func (p PasswordHash) String() string { return string(p) }

// FirstName represents a validated first name.
type FirstName string

func NewFirstName(name string) (FirstName, error) {
	validated, err := validateNonEmpty(name, ErrInvalidFirstName)
	if err != nil {
		return "", err
	}

	return FirstName(validated), nil
}

func (f FirstName) String() string { return string(f) }

// LastName represents a validated last name.
type LastName string

func NewLastName(name string) (LastName, error) {
	validated, err := validateNonEmpty(name, ErrInvalidLastName)
	if err != nil {
		return "", err
	}

	return LastName(validated), nil
}

func (l LastName) String() string { return string(l) }

// validateNonEmpty trims whitespace and validates the string is not empty.
func validateNonEmpty(name string, emptyErr error) (string, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return "", emptyErr
	}

	return name, nil
}

// UserStatus represents user account status.
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

// UserRole represents user role in system.
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

// UserMetadata represents flexible user metadata.
type UserMetadata map[string]any

func NewUserMetadata() UserMetadata {
	return make(UserMetadata)
}

func (m UserMetadata) Set(key string, value any) {
	m[key] = value
}

func (m UserMetadata) Get(key string) (any, bool) {
	val, ok := m[key]

	return val, ok
}

// NewUser creates a new user entity with validation.
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

// IsActive returns true if user status is active.
func (u *User) IsActive() bool {
	return u.status == UserStatusActive
}

// UpdateProfile updates user profile information.
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

// changeField updates a field with validation and timestamp using generics.
func changeField[T any](
	u *User,
	value T,
	valid func(T) bool,
	errForInvalid func() error,
	apply func(*User, T),
) error {
	if !valid(value) {
		return errForInvalid()
	}

	apply(u, value)
	u.updatedAt = time.Now()

	return nil
}

// ChangeStatus updates user status with validation.
func (u *User) ChangeStatus(status UserStatus) error {
	return changeField(
		u,
		status,
		func(s UserStatus) bool { return s.IsValid() },
		func() error { return ErrInvalidUserStatus },
		func(u *User, s UserStatus) { u.status = s },
	)
}

// ChangeRole updates user role with validation.
func (u *User) ChangeRole(role UserRole) error {
	return changeField(
		u,
		role,
		func(r UserRole) bool { return r.IsValid() },
		func() error { return ErrInvalidUserRole },
		func(u *User, r UserRole) { u.role = r },
	)
}

// Verify marks user as verified.
func (u *User) Verify() {
	u.isVerified = true
	u.updatedAt = time.Now()
}

// RecordLogin updates last login time.
func (u *User) RecordLogin() {
	now := time.Now()
	u.lastLoginAt = &now
	u.updatedAt = now
}

// AddTag adds a tag to user if not already present.
func (u *User) AddTag(tag string) {
	if slices.Contains(u.tags, tag) {
		return
	}

	u.tags = append(u.tags, tag)
	u.updatedAt = time.Now()
}

// RemoveTag removes a tag from user.
func (u *User) RemoveTag(tag string) {
	for i, existingTag := range u.tags {
		if existingTag == tag {
			u.tags = append(u.tags[:i], u.tags[i+1:]...)
			u.updatedAt = time.Now()

			return
		}
	}
}

// SetID sets the user ID (used by repository after creation)
// This is intentionally package-private to allow repository to set ID after creation.
func (u *User) SetID(id UserID) {
	u.id = id
}

// UserStats represents user statistics.
type UserStats struct {
	TotalUsers       int64   `json:"total_users"`
	ActiveUsers      int64   `json:"active_users"`
	InactiveUsers    int64   `json:"inactive_users"`
	SuspendedUsers   int64   `json:"suspended_users"`
	VerifiedUsers    int64   `json:"verified_users"`
	UsersWithLogins  int64   `json:"users_with_logins"`
	NewUsers30d      int64   `json:"new_users_30d"`
	NewUsers7d       int64   `json:"new_users_7d"`
	ActivePercentage float64 `json:"active_percentage"`
	VerificationRate float64 `json:"verification_rate"`
}

// SessionStats represents session statistics.
type SessionStats struct {
	TotalSessions   int64 `json:"total_sessions"`
	ActiveSessions  int64 `json:"active_sessions"`
	ExpiredSessions int64 `json:"expired_sessions"`
	Sessions24h     int64 `json:"sessions_24h"`
	Sessions7d      int64 `json:"sessions_7d"`
	Sessions30d     int64 `json:"sessions_30d"`
}
