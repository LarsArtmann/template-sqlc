package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/LarsArtmann/template-sqlc/internal/domain/entities"
	"github.com/LarsArtmann/template-sqlc/pkg/errors"
)

// UserValidator implements user validation logic.
type UserValidator struct {
	emailRegex    *regexp.Regexp
	usernameRegex *regexp.Regexp
}

// NewUserValidator creates a new user validator.
func NewUserValidator() *UserValidator {
	return &UserValidator{
		emailRegex:    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		usernameRegex: regexp.MustCompile(`^[a-zA-Z0-9_-]{3,50}$`),
	}
}

// ValidateUserCreate validates user creation request.
func (v *UserValidator) ValidateUserCreate(email, username, firstName, lastName string) error {
	// Validate email
	err := v.validateEmail(email)
	if err != nil {
		return err
	}

	// Validate username
	err = v.validateUsername(username)
	if err != nil {
		return err
	}

	// Validate first name
	err = v.validateName("first_name", firstName)
	if err != nil {
		return err
	}

	// Validate last name
	err = v.validateName("last_name", lastName)
	if err != nil {
		return err
	}

	return nil
}

// ValidateUserUpdate validates user update request.
func (v *UserValidator) ValidateUserUpdate(user *entities.User) error {
	// All basic validations should pass since the user entity
	// was already validated during creation
	// Additional business logic validations can go here

	// For example, check if admin users can be deactivated
	if user.Status() == entities.UserStatusSuspended && user.Role() == entities.UserRoleAdmin {
		return errors.NewBusinessLogicError("admin users cannot be suspended")
	}

	return nil
}

// ValidatePasswordRequirements validates password strength.
func (v *UserValidator) ValidatePasswordRequirements(password string) error {
	err := v.validatePasswordLength(password)
	if err != nil {
		return err
	}

	if v.countCharacterCategories(password) < 3 {
		return errors.NewValidationError(
			"password",
			"must contain at least 3 of: uppercase letters, lowercase letters, numbers, special characters",
		)
	}

	if v.isCommonPassword(password) {
		return errors.NewValidationError(
			"password",
			"password is too common, please choose a stronger one",
		)
	}

	return nil
}

// validateLength checks string length constraints.
func validateLength(field, value string, minLen, maxLen int) error {
	if minLen > 0 && len(value) < minLen {
		return errors.NewValidationError(
			field,
			fmt.Sprintf("must be at least %d characters", minLen),
		)
	}

	if maxLen > 0 && len(value) > maxLen {
		return errors.NewValidationError(
			field,
			fmt.Sprintf("must not exceed %d characters", maxLen),
		)
	}

	return nil
}

// validatePasswordLength checks password length requirements.
func (v *UserValidator) validatePasswordLength(password string) error {
	return validateLength("password", password, 8, 128)
}

// countCharacterCategories counts how many character categories are present.
func (v *UserValidator) countCharacterCategories(password string) int {
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	count := 0

	for _, has := range []bool{hasUpper, hasLower, hasNumber, hasSpecial} {
		if has {
			count++
		}
	}

	return count
}

// ValidateEmail validates email format.
func (v *UserValidator) validateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return errors.NewMissingFieldError("email")
	}

	if len(email) > 254 { //nolint:mnd // RFC 5321
		return errors.NewValidationError("email", "must not exceed 254 characters")
	}

	if !v.emailRegex.MatchString(email) {
		return errors.NewInvalidFormatError("email", "must be a valid email address")
	}

	// Additional checks
	if strings.Count(email, "@") != 1 {
		return errors.NewInvalidFormatError("email", "must contain exactly one @ symbol")
	}

	parts := strings.Split(email, "@")
	localPart := parts[0]
	domainPart := parts[1]

	if len(localPart) == 0 {
		return errors.NewInvalidFormatError("email", "local part cannot be empty")
	}

	if len(domainPart) == 0 {
		return errors.NewInvalidFormatError("email", "domain part cannot be empty")
	}

	// Check for invalid characters in local part
	if strings.Contains(localPart, "..") {
		return errors.NewInvalidFormatError("email", "cannot contain consecutive dots")
	}

	return nil
}

// ValidateUsername validates username format.
func (v *UserValidator) validateUsername(username string) error {
	username = strings.TrimSpace(username)

	if username == "" {
		return errors.NewMissingFieldError("username")
	}

	if len(username) < 3 { //nolint:mnd // minimum username length
		return errors.NewValidationError("username", "must be at least 3 characters long")
	}

	if len(username) > 50 { //nolint:mnd // maximum username length
		return errors.NewValidationError("username", "must not exceed 50 characters")
	}

	if !v.usernameRegex.MatchString(username) {
		return errors.NewInvalidFormatError(
			"username",
			"can only contain letters, numbers, underscores, and hyphens",
		)
	}

	// Check for reserved usernames
	if v.isReservedUsername(username) {
		return errors.NewValidationError("username", "username is reserved")
	}

	return nil
}

// ValidateName validates first/last name.
func (v *UserValidator) validateName(field, name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.NewMissingFieldError(field)
	}

	if len(name) > 100 { //nolint:mnd // maximum name length
		return errors.NewValidationError(field, "must not exceed 100 characters")
	}

	// Basic validation - names should contain letters and possibly spaces/hyphens
	for _, char := range name {
		if !unicode.IsLetter(char) && char != ' ' && char != '-' && char != '\'' {
			return errors.NewInvalidFormatError(
				field,
				"can only contain letters, spaces, hyphens, and apostrophes",
			)
		}
	}

	return nil
}

// isReservedUsername checks if username is reserved.
func (v *UserValidator) isReservedUsername(username string) bool {
	reserved := map[string]bool{
		"admin":         true,
		"administrator": true,
		"root":          true,
		"system":        true,
		"api":           true,
		"www":           true,
		"mail":          true,
		"support":       true,
		"nobody":        true,
		"guest":         true,
		"anonymous":     true,
		"user":          true,
		"users":         true,
		"help":          true,
		"info":          true,
		"sales":         true,
		"marketing":     true,
		"billing":       true,
		"security":      true,
		"legal":         true,
		"privacy":       true,
		"terms":         true,
		"contact":       true,
		"about":         true,
		"blog":          true,
		"news":          true,
		"press":         true,
		"careers":       true,
		"jobs":          true,
		"shop":          true,
		"store":         true,
		"cart":          true,
		"checkout":      true,
		"order":         true,
		"orders":        true,
		"account":       true,
		"profile":       true,
		"settings":      true,
		"dashboard":     true,
		"console":       true,
		"manage":        true,
		"my":            true,
		"me":            true,
		"self":          true,
	}

	lowercase := strings.ToLower(username)

	return reserved[lowercase]
}

// isCommonPassword checks against common passwords.
func (v *UserValidator) isCommonPassword(password string) bool {
	// In a real implementation, you'd use a comprehensive list
	// or integrate with a service like HaveIBeenPwned
	commonPasswords := map[string]bool{
		"password":    true,
		"123456":      true,
		"123456789":   true,
		"12345678":    true,
		"12345":       true,
		"1234567":     true,
		"1234567890":  true,
		"qwerty":      true,
		"abc123":      true,
		"password123": true,
		"admin":       true,
		"letmein":     true,
		"welcome":     true,
		"monkey":      true,
		"1234":        true,
		"dragon":      true,
		"master":      true,
		"hello":       true,
		"freedom":     true,
		"whatever":    true,
		"qazwsx":      true,
		"trustno1":    true,
		"123qwe":      true,
		"1q2w3e4r":    true,
		"zxcvbnm":     true,
		"iloveyou":    true,
		"starwars":    true,
		"football":    true,
		"baseball":    true,
		"soccer":      true,
	}

	lowercase := strings.ToLower(password)

	return commonPasswords[lowercase]
}

// ValidateUserRole validates user role.
func (v *UserValidator) ValidateUserRole(role string) error {
	validRoles := map[string]bool{
		string(entities.UserRoleUser):      true,
		string(entities.UserRoleAdmin):     true,
		string(entities.UserRoleModerator): true,
	}

	if !validRoles[role] {
		return errors.NewValidationError("role", "must be one of: user, admin, moderator")
	}

	return nil
}

// ValidateUserStatus validates user status.
func (v *UserValidator) ValidateUserStatus(status string) error {
	validStatuses := map[string]bool{
		string(entities.UserStatusActive):    true,
		string(entities.UserStatusInactive):  true,
		string(entities.UserStatusSuspended): true,
		string(entities.UserStatusPending):   true,
	}

	if !validStatuses[status] {
		return errors.NewValidationError(
			"status",
			"must be one of: active, inactive, suspended, pending",
		)
	}

	return nil
}

// ValidatePagination validates pagination parameters.
func (v *UserValidator) ValidatePagination(limit, offset int) error {
	if limit < 1 {
		return errors.NewValidationError("limit", "must be at least 1")
	}

	if limit > 1000 { //nolint:mnd // max pagination limit
		return errors.NewValidationError("limit", "must not exceed 1000")
	}

	if offset < 0 {
		return errors.NewValidationError("offset", "must be non-negative")
	}

	return nil
}

// ValidateSearchQuery validates search query.
func (v *UserValidator) ValidateSearchQuery(query string) error {
	query = strings.TrimSpace(query)

	if query == "" {
		return errors.NewValidationError("query", "search query cannot be empty")
	}

	if len(query) > 500 { //nolint:mnd // max search query length
		return errors.NewValidationError("query", "search query must not exceed 500 characters")
	}

	// Basic validation - prevent injection attempts
	invalidPatterns := []string{
		"<script",
		"</script>",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"onclick=",
		"alert(",
		"prompt(",
		"confirm(",
	}

	lowercase := strings.ToLower(query)
	for _, pattern := range invalidPatterns {
		if strings.Contains(lowercase, pattern) {
			return errors.NewValidationError("query", "search query contains invalid characters")
		}
	}

	return nil
}

// ValidateTags validates user tags.
func (v *UserValidator) ValidateTags(tags []string) error {
	if len(tags) > 50 { //nolint:mnd // max tags count
		return errors.NewValidationError("tags", "cannot have more than 50 tags")
	}

	seen := make(map[string]bool)

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)

		if tag == "" {
			continue // Skip empty tags
		}

		if len(tag) > 50 { //nolint:mnd // max tag length
			return errors.NewValidationError("tags", "tag must not exceed 50 characters")
		}

		// Basic validation for tag content
		if strings.Contains(tag, " ") || strings.Contains(tag, "\t") {
			return errors.NewValidationError("tags", "tag cannot contain whitespace")
		}

		lowercase := strings.ToLower(tag)
		if seen[lowercase] {
			return errors.NewValidationError("tags", "duplicate tag found")
		}

		seen[lowercase] = true
	}

	return nil
}

// SanitizeString sanitizes input string.
func (v *UserValidator) SanitizeString(input string) string {
	// Basic sanitization - trim whitespace and convert to lowercase if needed
	return strings.TrimSpace(input)
}

// SanitizeHTML sanitizes HTML input (basic implementation).
func (v *UserValidator) SanitizeHTML(input string) string {
	// In a real implementation, you'd use a proper HTML sanitizer library
	// like bluemonday or html-sanitizer
	input = strings.ReplaceAll(input, "<script>", "&lt;script&gt;")
	input = strings.ReplaceAll(input, "</script>", "&lt;/script&gt;")
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "&", "&amp;")

	return input
}
