-- name: CreateUser :one
INSERT INTO users (
    email, username, password_hash, 
    first_name, last_name, status, role, 
    profile_metadata
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetUserByID :one
SELECT 
    id,
    BIN_TO_UUID(uuid, 1) as uuid,
    email,
    username,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at,
    last_login_at,
    status,
    role,
    is_verified,
    profile_metadata
FROM users 
WHERE id = ?;

-- name: GetUserByUUID :one
SELECT 
    id,
    BIN_TO_UUID(uuid, 1) as uuid,
    email,
    username,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at,
    last_login_at,
    status,
    role,
    is_verified,
    profile_metadata
FROM users 
WHERE uuid = UUID_TO_BIN(?, 1);

-- name: GetUserByEmail :one
SELECT 
    id,
    BIN_TO_UUID(uuid, 1) as uuid,
    email,
    username,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at,
    last_login_at,
    status,
    role,
    is_verified,
    profile_metadata
FROM users 
WHERE email = ?;

-- name: GetUserByUsername :one
SELECT 
    id,
    BIN_TO_UUID(uuid, 1) as uuid,
    email,
    username,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at,
    last_login_at,
    status,
    role,
    is_verified,
    profile_metadata
FROM users 
WHERE username = ?;

-- name: UpdateUser :one
UPDATE users 
SET 
    email = IFNULL(?, email),
    username = IFNULL(?, username),
    first_name = IFNULL(?, first_name),
    last_name = IFNULL(?, last_name),
    status = IFNULL(?, status),
    role = IFNULL(?, role),
    is_verified = IFNULL(?, is_verified),
    profile_metadata = IFNULL(?, profile_metadata)
WHERE id = ?
RETURNING 
    id,
    BIN_TO_UUID(uuid, 1) as uuid,
    email,
    username,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at,
    last_login_at,
    status,
    role,
    is_verified,
    profile_metadata;

-- name: UpdatePassword :exec
UPDATE users 
SET password_hash = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateUserStatus :exec
UPDATE users 
SET status = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ListUsers :many
SELECT 
    id,
    BIN_TO_UUID(uuid, 1) as uuid,
    email,
    username,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at,
    last_login_at,
    status,
    role,
    is_verified,
    profile_metadata
FROM users 
WHERE status = IFNULL(?, 'active')
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: SearchUsers :many
SELECT 
    id,
    BIN_TO_UUID(uuid, 1) as uuid,
    email,
    username,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at,
    last_login_at,
    status,
    role,
    is_verified,
    profile_metadata,
    MATCH(searchable_content) AGAINST(? IN NATURAL LANGUAGE MODE) as relevance
FROM users 
WHERE MATCH(searchable_content) AGAINST(? IN NATURAL LANGUAGE MODE)
    AND status = IFNULL(?, 'active')
ORDER BY relevance DESC, created_at DESC
LIMIT ? OFFSET ?;

-- name: CountUsersByStatus :many
SELECT status, COUNT(*) as count 
FROM users 
GROUP BY status;

-- name: GetUserStats :one
SELECT 
    COUNT(*) as total_users,
    SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active_users,
    SUM(CASE WHEN status = 'inactive' THEN 1 ELSE 0 END) as inactive_users,
    SUM(CASE WHEN status = 'suspended' THEN 1 ELSE 0 END) as suspended_users,
    SUM(CASE WHEN is_verified = TRUE THEN 1 ELSE 0 END) as verified_users,
    SUM(CASE WHEN last_login_at IS NOT NULL THEN 1 ELSE 0 END) as users_with_logins,
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY) THEN 1 ELSE 0 END) as new_users_30d,
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN 1 ELSE 0 END) as new_users_7d,
    ROUND(
        (SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) / COUNT(*)) * 100, 2
    ) as active_percentage,
    ROUND(
        (SUM(CASE WHEN is_verified = TRUE THEN 1 ELSE 0 END) / 
         NULLIF(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0)) * 100, 2
    ) as verification_rate
FROM users;

-- name: CreateUserSession :one
INSERT INTO user_sessions (
    user_id, session_token, device_info, 
    ip_address, user_agent, expires_at
) VALUES (
    ?, UUID_TO_BIN(?, 1), ?, ?, ?, ?
)
RETURNING *;

-- name: GetUserSession :one
SELECT 
    s.id,
    s.user_id,
    BIN_TO_UUID(s.session_token, 1) as session_token,
    s.device_info,
    s.ip_address,
    s.user_agent,
    s.created_at,
    s.expires_at,
    s.is_active,
    u.email,
    u.username,
    u.first_name,
    u.last_name
FROM user_sessions s
JOIN users u ON s.user_id = u.id
WHERE s.session_token = UUID_TO_BIN(?, 1) 
    AND s.is_active = TRUE 
    AND s.expires_at > NOW();

-- name: ExpireUserSession :exec
UPDATE user_sessions 
SET is_active = FALSE 
WHERE session_token = UUID_TO_BIN(?, 1);

-- name: ExpireAllUserSessions :exec
UPDATE user_sessions 
SET is_active = FALSE 
WHERE user_id = ? AND is_active = TRUE;

-- name: CleanupExpiredSessions :exec
CALL CleanupExpiredSessions();

-- name: GetUserActiveSessions :many
SELECT 
    id,
    user_id,
    BIN_TO_UUID(session_token, 1) as session_token,
    device_info,
    ip_address,
    user_agent,
    created_at,
    expires_at,
    is_active
FROM user_sessions 
WHERE user_id = ? AND is_active = TRUE AND expires_at > NOW()
ORDER BY created_at DESC;