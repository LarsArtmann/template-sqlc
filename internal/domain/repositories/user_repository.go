package repositories

import (
	"context"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
)

// UserRepository defines the interface for user data access
// This abstracts away database-specific implementation details
type UserRepository interface {
	// CRUD operations
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id entities.UserID) (*entities.User, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.User, error)
	GetByEmail(ctx context.Context, email entities.Email) (*entities.User, error)
	GetByUsername(ctx context.Context, username entities.Username) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id entities.UserID) error

	// List and search operations
	List(ctx context.Context, status entities.UserStatus, limit, offset int) ([]*entities.User, error)
	Search(ctx context.Context, query string, status entities.UserStatus, limit int) ([]*entities.User, error)
	SearchByTags(ctx context.Context, tags []string, status entities.UserStatus, limit, offset int) ([]*entities.User, error)

	// Aggregate operations
	CountByStatus(ctx context.Context) (map[entities.UserStatus]int64, error)
	GetStats(ctx context.Context) (*entities.UserStats, error)

	// Authentication operations
	VerifyCredentials(ctx context.Context, email entities.Email, password entities.PasswordHash) (*entities.User, error)
	UpdatePassword(ctx context.Context, id entities.UserID, password entities.PasswordHash) error
	MarkVerified(ctx context.Context, id entities.UserID) error

	// Status operations
	ChangeStatus(ctx context.Context, id entities.UserID, status entities.UserStatus) error
	Activate(ctx context.Context, id entities.UserID) error
	Deactivate(ctx context.Context, id entities.UserID) error
	Suspend(ctx context.Context, id entities.UserID) error

	// Role operations
	ChangeRole(ctx context.Context, id entities.UserID, role entities.UserRole) error
}

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	// CRUD operations
	Create(ctx context.Context, session *entities.UserSession) error
	GetByToken(ctx context.Context, token entities.SessionToken) (*entities.UserSession, error)
	GetByUserID(ctx context.Context, userID entities.UserID, activeOnly bool) ([]*entities.UserSession, error)
	Update(ctx context.Context, session *entities.UserSession) error
	Delete(ctx context.Context, id entities.SessionID) error

	// Session management
	DeactivateByToken(ctx context.Context, token entities.SessionToken) error
	DeactivateByUserID(ctx context.Context, userID entities.UserID) error
	CleanupExpired(ctx context.Context) (int64, error)

	// Analytics
	GetActiveSessions(ctx context.Context, userID entities.UserID) (int64, error)
	GetSessionStats(ctx context.Context) (*entities.SessionStats, error)
}

// TransactionalRepository defines transaction support
type TransactionalRepository interface {
	// Transaction operations
	BeginTx(ctx context.Context) (Transaction, error)
	RunInTransaction(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}

// Transaction defines the transaction interface
type Transaction interface {
	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error

	// Repository interfaces within transaction
	UserRepository() UserRepository
	SessionRepository() SessionRepository
}

// UserStats represents user statistics
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

// SessionStats represents session statistics
type SessionStats struct {
	TotalSessions   int64 `json:"total_sessions"`
	ActiveSessions  int64 `json:"active_sessions"`
	ExpiredSessions int64 `json:"expired_sessions"`
	Sessions24h     int64 `json:"sessions_24h"`
	Sessions7d      int64 `json:"sessions_7d"`
	Sessions30d     int64 `json:"sessions_30d"`
}
