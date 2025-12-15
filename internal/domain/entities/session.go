package entities

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
)

// UserSession represents a user session entity
type UserSession struct {
	id         SessionID
	userID     UserID
	token      SessionToken
	deviceInfo SessionDeviceInfo
	ipAddress  net.IP
	userAgent  string
	createdAt  time.Time
	expiresAt  time.Time
	isActive   bool
}

// SessionID is a strongly-typed session identifier
type SessionID int64

func (id SessionID) Int64() int64   { return int64(id) }
func (id SessionID) String() string { return fmt.Sprintf("session:%d", id) }

// SessionToken represents a secure session token
type SessionToken uuid.UUID

func NewSessionToken() SessionToken {
	return SessionToken(uuid.New())
}

func (t SessionToken) UUID() uuid.UUID { return uuid.UUID(t) }
func (t SessionToken) String() string  { return uuid.UUID(t).String() }

// SessionDeviceInfo contains device information for a session
type SessionDeviceInfo struct {
	Platform string                 `json:"platform"`
	Device   string                 `json:"device"`
	Browser  string                 `json:"browser"`
	Version  string                 `json:"version"`
	Metadata map[string]interface{} `json:"metadata"`
}

func NewSessionDeviceInfo() SessionDeviceInfo {
	return SessionDeviceInfo{
		Metadata: make(map[string]interface{}),
	}
}

// NewUserSession creates a new user session
func NewUserSession(
	userID UserID,
	ipAddress net.IP,
	userAgent string,
	deviceInfo SessionDeviceInfo,
	duration time.Duration,
) *UserSession {
	now := time.Now()
	return &UserSession{
		userID:     userID,
		token:      NewSessionToken(),
		deviceInfo: deviceInfo,
		ipAddress:  ipAddress,
		userAgent:  userAgent,
		createdAt:  now,
		expiresAt:  now.Add(duration),
		isActive:   true,
	}
}

// Session methods
func (s *UserSession) ID() SessionID                 { return s.id }
func (s *UserSession) UserID() UserID                { return s.userID }
func (s *UserSession) Token() SessionToken           { return s.token }
func (s *UserSession) DeviceInfo() SessionDeviceInfo { return s.deviceInfo }
func (s *UserSession) IPAddress() net.IP             { return s.ipAddress }
func (s *UserSession) UserAgent() string             { return s.userAgent }
func (s *UserSession) CreatedAt() time.Time          { return s.createdAt }
func (s *UserSession) ExpiresAt() time.Time          { return s.expiresAt }
func (s *UserSession) IsActive() bool                { return s.isActive }

// IsExpired returns true if the session has expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

// IsValid returns true if the session is active and not expired
func (s *UserSession) IsValid() bool {
	return s.isActive && !s.IsExpired()
}

// Deactivate marks the session as inactive
func (s *UserSession) Deactivate() {
	s.isActive = false
}

// Extend extends the session expiration time
func (s *UserSession) Extend(duration time.Duration) {
	s.expiresAt = time.Now().Add(duration)
}

// GetMetadata returns device metadata
func (d SessionDeviceInfo) GetMetadata(key string) (interface{}, bool) {
	val, ok := d.Metadata[key]
	return val, ok
}

// SetMetadata sets device metadata
func (d SessionDeviceInfo) SetMetadata(key string, value interface{}) {
	if d.Metadata == nil {
		d.Metadata = make(map[string]interface{})
	}
	d.Metadata[key] = value
}

// Common session durations
const (
	SessionDurationShort    = 24 * time.Hour      // 1 day
	SessionDurationMedium   = 7 * 24 * time.Hour  // 1 week
	SessionDurationLong     = 30 * 24 * time.Hour // 1 month
	SessionDurationRemember = 90 * 24 * time.Hour // 3 months (remember me)
)
