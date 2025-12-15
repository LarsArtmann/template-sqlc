-- Example: User management schema for PostgreSQL
-- This file demonstrates PostgreSQL-specific features and best practices

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table with PostgreSQL-specific features
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'pending');
CREATE TYPE user_role AS ENUM ('user', 'admin', 'moderator');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    
    -- Authentication fields
    email CITEXT UNIQUE NOT NULL,                    -- Case-insensitive email
    username CITEXT UNIQUE NOT NULL,                 -- Case-insensitive username
    password_hash TEXT NOT NULL,                    -- NEVER store plain passwords
    
    -- Profile fields
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    
    -- PostgreSQL-specific timestamp fields
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMPTZ NULL,
    
    -- Enum fields for type safety
    status user_status DEFAULT 'active' NOT NULL,
    role user_role DEFAULT 'user' NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE NOT NULL,
    
    -- JSONB field with GIN index for efficient queries
    profile_metadata JSONB DEFAULT '{}'::jsonb,
    
    -- Array fields
    tags TEXT[] DEFAULT '{}',
    
    -- Full-text search with tsvector
    search_vector tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(first_name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(last_name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(username, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(email, '')), 'C')
    ) STORED,
    
    -- Constraints
    CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT valid_username CHECK (length(username) >= 3 AND length(username) <= 50),
    CONSTRAINT valid_name CHECK (length(first_name) > 0 AND length(last_name) > 0)
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_uuid ON users(uuid);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_last_login ON users(last_login_at) WHERE last_login_at IS NOT NULL;

-- GIN index for JSONB field
CREATE INDEX idx_users_profile_metadata ON users USING GIN(profile_metadata);

-- GIN index for full-text search
CREATE INDEX idx_users_search_vector ON users USING GIN(search_vector);

-- Array index for tags
CREATE INDEX idx_users_tags ON users USING GIN(tags);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger to update search vector
CREATE OR REPLACE FUNCTION update_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := 
        setweight(to_tsvector('english', coalesce(NEW.first_name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.last_name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.username, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(NEW.email, '')), 'C');
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_search_vector
    BEFORE INSERT OR UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_search_vector();

-- User sessions table with proper constraints
CREATE TABLE user_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    device_info JSONB DEFAULT '{}'::jsonb,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    
    CONSTRAINT valid_expires_at CHECK (expires_at > created_at)
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX idx_user_sessions_expires ON user_sessions(expires_at);

-- View for active users with sessions
CREATE VIEW active_users_with_sessions AS
SELECT 
    u.*,
    COUNT(s.id) FILTER (WHERE s.is_active = TRUE AND s.expires_at > CURRENT_TIMESTAMP) as active_sessions,
    MAX(s.created_at) as last_session_created
FROM users u
LEFT JOIN user_sessions s ON u.id = s.user_id
WHERE u.status = 'active'
GROUP BY u.id;