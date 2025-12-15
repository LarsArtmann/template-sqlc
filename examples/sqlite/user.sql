-- Example: User management schema for SQLite
-- This file demonstrates best practices for SQLite schema design

-- Users table with proper constraints and indexes
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,                    -- UUID for public identification
    email TEXT UNIQUE NOT NULL,                   -- Unique email constraint
    username TEXT UNIQUE NOT NULL,                -- Unique username
    password_hash TEXT NOT NULL,                  -- NEVER store plain passwords
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    
    -- Audit fields (best practice)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login_at DATETIME NULL,
    
    -- Status fields
    is_active BOOLEAN DEFAULT TRUE CHECK (is_active IN (0, 1)),
    is_verified BOOLEAN DEFAULT FALSE CHECK (is_verified IN (0, 1)),
    
    -- Profile data stored as JSON (flexible schema)
    profile_metadata JSON NULL,
    
    -- Full-text search content (for FTS5)
    searchable_content TEXT GENERATED ALWAYS AS (
        first_name || ' ' || last_name || ' ' || username || ' ' || email
    ) STORED
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_uuid ON users(uuid);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = 1;

-- FTS5 virtual table for full-text search
CREATE VIRTUAL TABLE users_fts USING fts5(
    first_name,
    last_name,
    username,
    email,
    content='users',
    content_rowid='id'
);

-- FTS triggers to keep search index updated
CREATE TRIGGER users_fts_insert AFTER INSERT ON users BEGIN
    INSERT INTO users_fts(rowid, first_name, last_name, username, email)
    VALUES (new.id, new.first_name, new.last_name, new.username, new.email);
END;

CREATE TRIGGER users_fts_delete AFTER DELETE ON users BEGIN
    INSERT INTO users_fts(users_fts, rowid, first_name, last_name, username, email)
    VALUES ('delete', old.id, old.first_name, old.last_name, old.username, old.email);
END;

CREATE TRIGGER users_fts_update AFTER UPDATE ON users BEGIN
    INSERT INTO users_fts(users_fts, rowid, first_name, last_name, username, email)
    VALUES ('delete', old.id, old.first_name, old.last_name, old.username, old.email);
    INSERT INTO users_fts(rowid, first_name, last_name, username, email)
    VALUES (new.id, new.first_name, new.last_name, new.username, new.email);
END;

-- Update timestamp trigger
CREATE TRIGGER users_updated_at AFTER UPDATE ON users BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;