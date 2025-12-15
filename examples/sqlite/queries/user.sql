-- name: CreateUser :one
INSERT INTO users (
    uuid, email, username, password_hash, 
    first_name, last_name, profile_metadata, is_active
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? AND is_active = TRUE;

-- name: GetUserByUUID :one
SELECT * FROM users WHERE uuid = ? AND is_active = TRUE;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? AND is_active = TRUE;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? AND is_active = TRUE;

-- name: UpdateUser :one
UPDATE users 
SET 
    email = COALESCE(?, email),
    username = COALESCE(?, username),
    first_name = COALESCE(?, first_name),
    last_name = COALESCE(?, last_name),
    profile_metadata = COALESCE(?, profile_metadata),
    is_active = COALESCE(?, is_active),
    is_verified = COALESCE(?, is_verified)
WHERE id = ?
RETURNING *;

-- name: UpdatePassword :exec
UPDATE users 
SET password_hash = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SoftDeleteUser :exec
UPDATE users 
SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users 
WHERE is_active = TRUE 
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: SearchUsers :many
SELECT u.* FROM users u
JOIN users_fts fts ON u.id = fts.rowid
WHERE users_fts MATCH ? AND u.is_active = TRUE
ORDER BY u.created_at DESC
LIMIT ?;

-- name: CountActiveUsers :one
SELECT COUNT(*) FROM users WHERE is_active = TRUE;

-- name: UpdateLastLogin :exec
UPDATE users 
SET last_login_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: VerifyUser :exec
UPDATE users 
SET is_verified = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE uuid = ?;

-- name: GetUserStats :one
SELECT 
    COUNT(*) as total_users,
    COUNT(CASE WHEN is_active = TRUE THEN 1 END) as active_users,
    COUNT(CASE WHEN is_verified = TRUE THEN 1 END) as verified_users,
    COUNT(CASE WHEN last_login_at IS NOT NULL THEN 1 END) as users_with_logins,
    COUNT(CASE WHEN created_at >= datetime('now', '-30 days') THEN 1 END) as new_users_30d
FROM users;