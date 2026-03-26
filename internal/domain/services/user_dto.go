package services

import "github.com/LarsArtmann/template-sqlc/internal/domain/entities"

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email        string         `json:"email" validate:"required,email"`
	Username     string         `json:"username" validate:"required,min=3,max=50"`
	PasswordHash string         `json:"password_hash" validate:"required"`
	FirstName    string         `json:"first_name" validate:"required"`
	LastName     string         `json:"last_name" validate:"required"`
	Status       string         `json:"status" validate:"required"`
	Role         string         `json:"role" validate:"required"`
	Tags         []string       `json:"tags"`
	Metadata     map[string]any `json:"metadata"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	UserID    entities.UserID `json:"user_id" validate:"required"`
	FirstName *string         `json:"first_name,omitempty" validate:"omitempty,min=1"`
	LastName  *string         `json:"last_name,omitempty" validate:"omitempty,min=1"`
	Metadata  *map[string]any `json:"metadata,omitempty"`
	Tags      *[]string       `json:"tags,omitempty"`
	UpdatedBy string          `json:"updated_by" validate:"required"`
}
