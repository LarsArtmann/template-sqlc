# sqlc Template Examples

This directory contains real-world examples demonstrating how to use the sqlc template with different databases.

## Structure

```
examples/
├── sqlite/
│   ├── user.sql              # User schema with SQLite-specific features
│   └── queries/
│       └── user.sql          # User queries with SQLite patterns
├── postgres/
│   ├── user.sql              # User schema with PostgreSQL features
│   └── queries/
│       └── user.sql          # User queries with PostgreSQL patterns
├── mysql/
│   ├── user.sql              # User schema with MySQL features
│   └── queries/
│       └── user.sql          # User queries with MySQL patterns
└── README.md                 # This file
```

## Database-Specific Features

### SQLite

- **FTS5** full-text search
- Generated columns for computed fields
- JSON metadata storage
- Trigger-based audit fields

### PostgreSQL

- **UUID** with built-in extension
- **CITEXT** for case-insensitive fields
- **ENUM** types for status fields
- **JSONB** with GIN indexes
- **tsvector** for full-text search
- Array fields and GIN indexes

### MySQL

- **UUID** stored as binary
- **ENUM** types
- **JSON** fields with generated columns
- **FULLTEXT** search indexes
- Stored procedures for cleanup

## Quick Start

1. Choose your database engine
2. Copy the relevant schema to your project
3. Set up the directory structure:
   ```bash
   mkdir -p sql/{sqlite,postgres,mysql}/{queries,schema}
   ```
4. Copy schema and query files
5. Run sqlc to generate code:
   ```bash
   sqlc generate
   ```

## Best Practices Demonstrated

- **Never store plain passwords** - always use password_hash
- **Proper constraints** for data integrity
- **Audit fields** (created_at, updated_at)
- **Status enums** instead of magic strings
- **Appropriate indexes** for performance
- **Full-text search** where applicable
- **Session management** with proper expiration

## Generated Code Features

All examples are designed to generate code with:

- Type-safe queries
- Proper NULL handling
- JSON tag generation
- Interface generation for testing
- Optimized prepared statements

## Testing the Examples

```bash
# SQLite
sqlite3 example.db < examples/sqlite/user.sql
sqlc generate

# PostgreSQL (requires running instance)
psql -d example_db < examples/postgres/user.sql
sqlc generate

# MySQL (requires running instance)
mysql example_db < examples/mysql/user.sql
sqlc generate
```
