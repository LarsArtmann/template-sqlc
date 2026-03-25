-- name: CreateUser :one
INSERT INTO users (
    uuid, email, username, password_hash, 
    first_name, last_name, profile_metadata, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND is_active = TRUE;

-- name: GetUserByUUID :one
SELECT * FROM users WHERE uuid = $1 AND is_active = TRUE;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND is_active = TRUE;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 AND is_active = TRUE;

-- name: UpdateUser :one
UPDATE users 
SET 
    email = COALESCE($2, email),
    username = COALESCE($3, username),
    first_name = COALESCE($4, first_name),
    last_name = COALESCE($5, last_name),
    profile_metadata = COALESCE($6, profile_metadata),
    is_active = COALESCE($7, is_active),
    is_verified = COALESCE($8, is_verified)
WHERE id = $1
RETURNING *;

-- name: UpdatePassword :exec
UPDATE users 
SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: SoftDeleteUser :exec
UPDATE users 
SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users 
WHERE is_active = TRUE 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountActiveUsers :one
SELECT COUNT(*) FROM users WHERE is_active = TRUE;

-- name: UpdateLastLogin :exec
UPDATE users 
SET last_login_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: VerifyUser :exec
UPDATE users 
SET is_verified = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE uuid = $1;

-- name: GetUserStats :one
SELECT 
    COUNT(*) as total_users,
    COUNT(*) FILTER (WHERE is_active = TRUE) as active_users,
    COUNT(*) FILTER (WHERE is_verified = TRUE) as verified_users,
    COUNT(*) FILTER (WHERE last_login_at IS NOT NULL) as users_with_logins
FROM users;
