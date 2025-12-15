-- name: CreateUser :one
INSERT INTO users (
    email, username, password_hash, 
    first_name, last_name, status, role 
    profile_metadata, tags
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByUUID :one
SELECT * FROM users WHERE uuid = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: UpdateUser :one
UPDATE users 
SET 
    email = COALESCE(sqlc.narg('email'), email),
    username = COALESCE(sqlc.narg('username'), username),
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    status = COALESCE(sqlc.narg('status'), status),
    role = COALESCE(sqlc.narg('role'), role),
    is_verified = COALESCE(sqlc.narg('is_verified'), is_verified),
    profile_metadata = COALESCE(sqlc.narg('profile_metadata'), profile_metadata),
    tags = COALESCE(sqlc.narg('tags'), tags)
WHERE id = ?
RETURNING *;

-- name: UpdatePassword :exec
UPDATE users 
SET password_hash = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateUserStatus :exec
UPDATE users 
SET status = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: AddUserTag :exec
UPDATE users 
SET tags = array_append(tags, ?), updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND NOT (? = ANY(tags));

-- name: RemoveUserTag :exec
UPDATE users 
SET tags = array_remove(tags, ?), updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users 
WHERE status = sqlc.narg('status')
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: SearchUsers :many
SELECT * FROM users 
WHERE search_vector @@ plainto_tsquery('english', ?)
    AND status = sqlc.narg('status', 'active')
ORDER BY ts_rank(search_vector, plainto_tsquery('english', ?)) DESC
LIMIT ? OFFSET ?;

-- name: SearchUsersByTags :many
SELECT * FROM users 
WHERE tags && ?::TEXT[]
    AND status = sqlc.narg('status', 'active')
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountUsersByStatus :one
SELECT status, COUNT(*) as count 
FROM users 
GROUP BY status;

-- name: GetUserStats :one
WITH stats AS (
    SELECT 
        COUNT(*) as total_users,
        COUNT(CASE WHEN status = 'active' THEN 1 END) as active_users,
        COUNT(CASE WHEN status = 'inactive' THEN 1 END) as inactive_users,
        COUNT(CASE WHEN status = 'suspended' THEN 1 END) as suspended_users,
        COUNT(CASE WHEN is_verified = TRUE THEN 1 END) as verified_users,
        COUNT(CASE WHEN last_login_at IS NOT NULL THEN 1 END) as users_with_logins,
        COUNT(CASE WHEN created_at >= CURRENT_TIMESTAMP - INTERVAL '30 days' THEN 1 END) as new_users_30d,
        COUNT(CASE WHEN created_at >= CURRENT_TIMESTAMP - INTERVAL '7 days' THEN 1 END) as new_users_7d
    FROM users
)
SELECT 
    total_users,
    active_users,
    inactive_users,
    suspended_users,
    verified_users,
    users_with_logins,
    new_users_30d,
    new_users_7d,
    ROUND((active_users::NUMERIC / NULLIF(total_users, 0)) * 100, 2) as active_percentage,
    ROUND((verified_users::NUMERIC / NULLIF(active_users, 0)) * 100, 2) as verification_rate
FROM stats;

-- name: CreateUserSession :one
INSERT INTO user_sessions (
    user_id, session_token, device_info, 
    ip_address, user_agent, expires_at
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetUserSession :one
SELECT s.*, u.email, u.username, u.first_name, u.last_name
FROM user_sessions s
JOIN users u ON s.user_id = u.id
WHERE s.session_token = ? AND s.is_active = TRUE AND s.expires_at > CURRENT_TIMESTAMP;

-- name: ExpireUserSession :exec
UPDATE user_sessions 
SET is_active = FALSE 
WHERE session_token = ? OR (user_id = ? AND is_active = TRUE);

-- name: CleanupExpiredSessions :exec
UPDATE user_sessions 
SET is_active = FALSE 
WHERE expires_at <= CURRENT_TIMESTAMP AND is_active = TRUE;

-- name: GetUserActiveSessions :many
SELECT * FROM user_sessions 
WHERE user_id = ? AND is_active = TRUE AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;