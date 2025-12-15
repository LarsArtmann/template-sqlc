package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
)

func TestUserCreation(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		username    string
		password    string
		firstName   string
		lastName    string
		status      string
		role        string
		expectError bool
		errorType   error
	}{
		{
			name:      "valid user",
			email:     "test@example.com",
			username:  "testuser",
			password:  "SecurePass123!",
			firstName: "John",
			lastName:  "Doe",
			status:    "active",
			role:      "user",
		},
		{
			name:        "invalid email",
			email:       "invalid-email",
			username:    "testuser",
			password:    "SecurePass123!",
			firstName:   "John",
			lastName:    "Doe",
			status:      "active",
			role:        "user",
			expectError: true,
			errorType:   entities.ErrInvalidEmail,
		},
		{
			name:        "short username",
			email:       "test@example.com",
			username:    "ab",
			password:    "SecurePass123!",
			firstName:   "John",
			lastName:    "Doe",
			status:      "active",
			role:        "user",
			expectError: true,
			errorType:   entities.ErrInvalidUsername,
		},
		{
			name:        "empty first name",
			email:       "test@example.com",
			username:    "testuser",
			password:    "SecurePass123!",
			firstName:   "",
			lastName:    "Doe",
			status:      "active",
			role:        "user",
			expectError: true,
			errorType:   entities.ErrInvalidFirstName,
		},
		{
			name:        "invalid status",
			email:       "test@example.com",
			username:    "testuser",
			password:    "SecurePass123!",
			firstName:   "John",
			lastName:    "Doe",
			status:      "invalid",
			role:        "user",
			expectError: true,
			errorType:   entities.ErrInvalidUserStatus,
		},
		{
			name:        "invalid role",
			email:       "test@example.com",
			username:    "testuser",
			password:    "SecurePass123!",
			firstName:   "John",
			lastName:    "Doe",
			status:      "active",
			role:        "invalid",
			expectError: true,
			errorType:   entities.ErrInvalidUserRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create domain values
			email, err := entities.NewEmail(tt.email)
			if err != nil && !tt.expectError {
				t.Fatalf("Failed to create email: %v", err)
			}

			username, err := entities.NewUsername(tt.username)
			if err != nil && !tt.expectError {
				t.Fatalf("Failed to create username: %v", err)
			}

			passwordHash, err := entities.NewPasswordHash(tt.password)
			if err != nil && !tt.expectError {
				t.Fatalf("Failed to create password hash: %v", err)
			}

			firstName, err := entities.NewFirstName(tt.firstName)
			if err != nil && !tt.expectError {
				t.Fatalf("Failed to create first name: %v", err)
			}

			lastName, err := entities.NewLastName(tt.lastName)
			if err != nil && !tt.expectError {
				t.Fatalf("Failed to create last name: %v", err)
			}

			status := entities.UserStatus(tt.status)
			role := entities.UserRole(tt.role)

			// Create user entity
			user, err := entities.NewUser(
				email,
				username,
				passwordHash,
				firstName,
				lastName,
				status,
				role,
				entities.NewUserMetadata(),
				[]string{},
			)

			if tt.expectError {
				assert.Error(t, err)
				assert.IsType(t, tt.errorType, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, email, user.Email())
				assert.Equal(t, username, user.Username())
				assert.Equal(t, firstName, user.FirstName())
				assert.Equal(t, lastName, user.LastName())
				assert.Equal(t, status, user.Status())
				assert.Equal(t, role, user.Role())
				assert.False(t, user.IsVerified())
				assert.WithinDuration(t, time.Now(), user.CreatedAt(), time.Second)
			}
		})
	}
}

func TestUserMethods(t *testing.T) {
	// Create a valid user
	email, _ := entities.NewEmail("test@example.com")
	username, _ := entities.NewUsername("testuser")
	passwordHash, _ := entities.NewPasswordHash("hashedpassword")
	firstName, _ := entities.NewFirstName("John")
	lastName, _ := entities.NewLastName("Doe")

	user, err := entities.NewUser(
		email,
		username,
		passwordHash,
		firstName,
		lastName,
		entities.UserStatusActive,
		entities.UserRoleUser,
		entities.NewUserMetadata(),
		[]string{},
	)
	require.NoError(t, err)

	// Test IsActive method
	assert.True(t, user.IsActive())

	// Test status changes
	err = user.ChangeStatus(entities.UserStatusInactive)
	assert.NoError(t, err)
	assert.False(t, user.IsActive())

	err = user.ChangeStatus(entities.UserStatusSuspended)
	assert.NoError(t, err)
	assert.False(t, user.IsActive())

	// Test role changes
	err = user.ChangeRole(entities.UserRoleAdmin)
	assert.NoError(t, err)
	assert.Equal(t, entities.UserRoleAdmin, user.Role())

	// Test verification
	assert.False(t, user.IsVerified())
	user.Verify()
	assert.True(t, user.IsVerified())

	// Test last login
	assert.Nil(t, user.LastLoginAt())
	user.RecordLogin()
	assert.NotNil(t, user.LastLoginAt())
	assert.WithinDuration(t, time.Now(), *user.LastLoginAt(), time.Second)

	// Test tags
	assert.Empty(t, user.Tags())
	user.AddTag("test")
	assert.Contains(t, user.Tags(), "test")
	user.RemoveTag("test")
	assert.NotContains(t, user.Tags(), "test")

	// Test profile update
	newFirstName := "Jane"
	firstNamePtr := &newFirstName
	err = user.UpdateProfile(firstNamePtr, nil, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, entities.FirstName(newFirstName), user.FirstName())
}

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"valid email", "test@example.com", true},
		{"invalid email - no @", "testexample.com", false},
		{"invalid email - no domain", "test@", false},
		{"invalid email - no local", "@example.com", false},
		{"valid email with subdomain", "test@sub.example.com", true},
		{"valid email with numbers", "test123@example456.com", true},
		{"invalid email with special chars", "test@exa!mple.com", false},
		{"valid email with +", "test+alias@example.com", true},
		{"empty email", "", false},
		{"email too long", string(make([]byte, 255)) + "@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := entities.NewEmail(tt.email)
			if tt.expected {
				assert.NoError(t, err)
				assert.Equal(t, tt.email, email.String())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestUsernameValidation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected bool
	}{
		{"valid username", "testuser", true},
		{"valid username with underscore", "test_user", true},
		{"valid username with hyphen", "test-user", true},
		{"valid username with numbers", "test123", true},
		{"too short username", "ab", false},
		{"too long username", string(make([]byte, 51)), false},
		{"username with space", "test user", false},
		{"username with special char", "test@user", false},
		{"username with period", "test.user", false},
		{"reserved username", "admin", false},
		{"reserved username uppercase", "ADMIN", false},
		{"empty username", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, err := entities.NewUsername(tt.username)
			if tt.expected {
				assert.NoError(t, err)
				assert.Equal(t, tt.username, username.String())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestPasswordHashValidation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"valid password hash", "hash_with_sufficient_length_for_bcrypt", true},
		{"valid bcrypt hash", "$2a$10$abcdefghijklmnopqrstuvwx1234567890", true},
		{"too short hash", "short", false},
		{"empty hash", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passwordHash, err := entities.NewPasswordHash(tt.password)
			if tt.expected {
				assert.NoError(t, err)
				assert.Equal(t, tt.password, passwordHash.String())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestUserStatusValidation(t *testing.T) {
	tests := []struct {
		name     string
		status   entities.UserStatus
		expected bool
	}{
		{"active status", entities.UserStatusActive, true},
		{"inactive status", entities.UserStatusInactive, true},
		{"suspended status", entities.UserStatusSuspended, true},
		{"pending status", entities.UserStatusPending, true},
		{"invalid status", entities.UserStatus("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestUserRoleValidation(t *testing.T) {
	tests := []struct {
		name     string
		role     entities.UserRole
		expected bool
	}{
		{"user role", entities.UserRoleUser, true},
		{"admin role", entities.UserRoleAdmin, true},
		{"moderator role", entities.UserRoleModerator, true},
		{"invalid role", entities.UserRole("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.IsValid())
		})
	}
}

func TestUserMetadata(t *testing.T) {
	metadata := entities.NewUserMetadata()

	// Test empty metadata
	assert.Empty(t, metadata)

	// Test setting and getting values
	metadata.Set("key1", "value1")
	metadata.Set("key2", 123)
	metadata.Set("key3", true)

	val, ok := metadata.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	val, ok = metadata.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, 123, val)

	val, ok = metadata.Get("key3")
	assert.True(t, ok)
	assert.Equal(t, true, val)

	// Test getting non-existent key
	val, ok = metadata.Get("nonexistent")
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestUserID(t *testing.T) {
	userID := entities.UserID(123)

	assert.Equal(t, int64(123), userID.Int64())
	assert.Equal(t, "user:123", userID.String())
}

func BenchmarkUserCreation(b *testing.B) {
	email, _ := entities.NewEmail("test@example.com")
	username, _ := entities.NewUsername("testuser")
	passwordHash, _ := entities.NewPasswordHash("hashedpassword")
	firstName, _ := entities.NewFirstName("John")
	lastName, _ := entities.NewLastName("Doe")
	metadata := entities.NewUserMetadata()
	tags := []string{"tag1", "tag2"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = entities.NewUser(
			email,
			username,
			passwordHash,
			firstName,
			lastName,
			entities.UserStatusActive,
			entities.UserRoleUser,
			metadata,
			tags,
		)
	}
}
