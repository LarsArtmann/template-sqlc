-- Example: User management schema for MySQL
-- This file demonstrates MySQL-specific features and best practices

-- Users table with MySQL-specific features
CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    uuid BINARY(16) NOT NULL UNIQUE,                    -- UUID stored as binary for efficiency
    email VARCHAR(255) NOT NULL UNIQUE,               -- Email with unique constraint
    username VARCHAR(50) NOT NULL UNIQUE,              -- Username with unique constraint
    password_hash VARCHAR(255) NOT NULL,               -- NEVER store plain passwords
    
    -- Profile fields
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    
    -- MySQL timestamp fields with automatic updates
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP NULL,
    
    -- Enum fields for type safety (MySQL 8.0+)
    status ENUM('active', 'inactive', 'suspended', 'pending') DEFAULT 'active' NOT NULL,
    role ENUM('user', 'admin', 'moderator') DEFAULT 'user' NOT NULL,
    is_verified TINYINT(1) DEFAULT FALSE NOT NULL,    -- MySQL uses TINYINT for booleans
    
    -- JSON field with virtual generated column for search
    profile_metadata JSON DEFAULT (JSON_OBJECT()),
    
    -- Full-text search content
    searchable_content VARCHAR(512) GENERATED ALWAYS AS (
        CONCAT_WS(' ', 
            COALESCE(first_name, ''), 
            COALESCE(last_name, ''), 
            COALESCE(username, ''), 
            COALESCE(email, '')
        )
    ) STORED,
    
    -- Constraints
    CONSTRAINT valid_email CHECK (email REGEXP '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$'),
    CONSTRAINT valid_username CHECK (LENGTH(username) >= 3 AND LENGTH(username) <= 50),
    CONSTRAINT valid_name CHECK (LENGTH(first_name) > 0 AND LENGTH(last_name) > 0),
    
    -- Indexes
    INDEX idx_users_email (email),
    INDEX idx_users_username (username),
    INDEX idx_users_status (status),
    INDEX idx_users_role (role),
    INDEX idx_users_created_at (created_at),
    INDEX idx_users_last_login (last_login_at),
    
    -- Full-text search index (MySQL 5.7+)
    FULLTEXT INDEX idx_users_search (searchable_content),
    
    -- JSON index for efficient queries
    INDEX idx_profile_metadata ((CAST(profile_metadata AS CHAR(255) ARRAY)))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- User sessions table
CREATE TABLE user_sessions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    session_token BINARY(16) NOT NULL UNIQUE,
    device_info JSON DEFAULT (JSON_OBJECT()),
    ip_address VARCHAR(45) NULL,                     -- IPv6 compatible
    user_agent TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    is_active TINYINT(1) DEFAULT TRUE NOT NULL,
    
    -- Foreign key with CASCADE delete
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_sessions_user_id (user_id),
    INDEX idx_sessions_token (session_token),
    INDEX idx_sessions_expires (expires_at),
    
    -- Constraints
    CONSTRAINT valid_expires_at CHECK (expires_at > created_at)
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Function to generate UUID (MySQL 8.0+)
DELIMITER $$
CREATE FUNCTION IF NOT EXISTS generate_uuid()
RETURNS BINARY(16)
DETERMINISTIC
BEGIN
    RETURN UUID_TO_BIN(UUID(), 1);
END$$
DELIMITER ;

-- Trigger to set UUID on insert
DELIMITER $$
CREATE TRIGGER before_user_insert
BEFORE INSERT ON users
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = X'00000000000000000000000000000000' THEN
        SET NEW.uuid = generate_uuid();
    END IF;
END$$
DELIMITER ;

-- Trigger to generate session token
DELIMITER $$
CREATE TRIGGER before_session_insert
BEFORE INSERT ON user_sessions
FOR EACH ROW
BEGIN
    IF NEW.session_token IS NULL OR NEW.session_token = X'00000000000000000000000000000000' THEN
        SET NEW.session_token = generate_uuid();
    END IF;
END$$
DELIMITER ;

-- View for active users with session count (MySQL 8.0+)
CREATE OR REPLACE VIEW active_users_with_sessions AS
SELECT 
    u.*,
    (SELECT COUNT(*) FROM user_sessions s 
     WHERE s.user_id = u.id AND s.is_active = TRUE AND s.expires_at > NOW()) as active_sessions,
    (SELECT MAX(s.created_at) FROM user_sessions s 
     WHERE s.user_id = u.id) as last_session_created
FROM users u
WHERE u.status = 'active';

-- Procedure to cleanup expired sessions
DELIMITER $$
CREATE PROCEDURE CleanupExpiredSessions()
BEGIN
    UPDATE user_sessions 
    SET is_active = FALSE 
    WHERE expires_at <= NOW() AND is_active = TRUE;
    
    SELECT ROW_COUNT() as sessions_cleaned;
END$$
DELIMITER ;