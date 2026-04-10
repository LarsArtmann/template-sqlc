package unit

import (
	"testing"
	"time"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Validatable defines the IsValid method.
type Validatable interface {
	IsValid() bool
}

// ValidatableTestCase is a test case for validatable types.
type ValidatableTestCase[T Validatable] struct {
	Name     string
	Value    T
	Expected bool
}

// runValidatableTests runs validation tests for a validatable type.
func runValidatableTests[T Validatable](t *testing.T, tests []ValidatableTestCase[T]) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Value.IsValid())
		})
	}
}

// UserCreationTestCase is a test case for user creation validation.
type UserCreationTestCase struct {
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
}

// validUserTestCase creates a valid user test case with the given field values.
func validUserTestCase(
	name, email, username, password, firstName, lastName, status, role string,
) UserCreationTestCase {
	return UserCreationTestCase{
		name:        name,
		email:       email,
		username:    username,
		password:    password,
		firstName:   firstName,
		lastName:    lastName,
		status:      status,
		role:        role,
		expectError: false,
	}
}

// valueObject is an interface for entities that have a String() method.
type valueObject interface {
	String() string
}

// testEntityValidation is a helper function to test value object creation.
func testEntityValidation[T valueObject](
	t *testing.T,
	name string,
	value string,
	constructor func(string) (T, error),
	expectedSuccess bool,
) {
	entity, err := constructor(value)
	if expectedSuccess {
		assert.NoError(t, err)
		assert.Equal(t, value, entity.String())
	} else {
		assert.Error(t, err)
	}
}

func TestUserCreation(t *testing.T) {
	tests := []UserCreationTestCase{
		validUserTestCase(
			"valid user",
			"test@example.com",
			"testuser",
			"$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsQ5pPjZ5yVlWK5WAe",
			"John",
			"Doe",
			"active",
			"user",
		),
		validUserTestCase(
			"valid user with all fields",
			"user@example.com",
			"newuser",
			"$2a$10$abcdefghijklmnopqrstuvwx1234567890ABCDEFGHIJKLMNOP",
			"Jane",
			"Smith",
			"pending",
			"admin",
		),
		generateInvalidStatusTestCase("invalid", entities.ErrInvalidUserStatus),
		generateInvalidRoleTestCase("invalid", entities.ErrInvalidUserRole),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create domain values - validation failures should cause test to fail
			email, err := entities.NewEmail(tt.email)
			require.NoError(t, err, "NewEmail should not fail for test case")

			username, err := entities.NewUsername(tt.username)
			require.NoError(t, err, "NewUsername should not fail for test case")

			passwordHash, err := entities.NewPasswordHash(tt.password)
			require.NoError(t, err, "NewPasswordHash should not fail for test case")

			firstName, err := entities.NewFirstName(tt.firstName)
			require.NoError(t, err, "NewFirstName should not fail for test case")

			lastName, err := entities.NewLastName(tt.lastName)
			require.NoError(t, err, "NewLastName should not fail for test case")

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
	newFirstName, _ := entities.NewFirstName("Jane")
	err = user.UpdateProfile(&newFirstName, nil, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, newFirstName, user.FirstName())
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
			testEntityValidation(t, tt.name, tt.email, entities.NewEmail, tt.expected)
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
			testEntityValidation(t, tt.name, tt.username, entities.NewUsername, tt.expected)
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
			testEntityValidation(t, tt.name, tt.password, entities.NewPasswordHash, tt.expected)
		})
	}
}

func TestUserStatusValidation(t *testing.T) {
	runValidatableTests(t, []ValidatableTestCase[entities.UserStatus]{
		{"active status", entities.UserStatusActive, true},
		{"inactive status", entities.UserStatusInactive, true},
		{"suspended status", entities.UserStatusSuspended, true},
		{"pending status", entities.UserStatusPending, true},
		{"invalid status", entities.UserStatus("invalid"), false},
	})
}

func TestUserRoleValidation(t *testing.T) {
	runValidatableTests(t, []ValidatableTestCase[entities.UserRole]{
		{"user role", entities.UserRoleUser, true},
		{"admin role", entities.UserRoleAdmin, true},
		{"moderator role", entities.UserRoleModerator, true},
		{"invalid role", entities.UserRole("invalid"), false},
	})
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

// generateInvalidStatusTestCase creates a test case for invalid status validation.
func generateInvalidStatusTestCase(invalidStatus string, expectedError error) UserCreationTestCase {
	base := validUserTestCase(
		"invalid status",
		"test@example.com",
		"testuser",
		"$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsQ5pPjZ5yVlWK5WAe",
		"John",
		"Doe",
		invalidStatus,
		"user",
	)
	base.expectError = true
	base.errorType = expectedError

	return base
}

// generateInvalidRoleTestCase creates a test case for invalid role validation.
func generateInvalidRoleTestCase(invalidRole string, expectedError error) UserCreationTestCase {
	base := validUserTestCase(
		"invalid role",
		"test@example.com",
		"testuser",
		"$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsQ5pPjZ5yVlWK5WAe",
		"John",
		"Doe",
		"active",
		invalidRole,
	)
	base.expectError = true
	base.errorType = expectedError

	return base
}

func BenchmarkUserCreation(b *testing.B) {
	email, _ := entities.NewEmail("test@example.com")
	username, _ := entities.NewUsername("testuser")
	passwordHash, _ := entities.NewPasswordHash("hashedpassword")
	firstName, _ := entities.NewFirstName("John")
	lastName, _ := entities.NewLastName("Doe")
	metadata := entities.NewUserMetadata()
	tags := []string{"tag1", "tag2"}

	for b.Loop() {
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
