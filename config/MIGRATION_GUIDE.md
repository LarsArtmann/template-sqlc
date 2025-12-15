# sqlc Configuration Migration Guide

## What Happened

Your monolithic `sqlc.yaml` has been split into maintainable, modular configurations.

## New Structure

```
config/
├── internal/
│   ├── base/
│   │   └── common.yaml          # Base rules and plugins
│   └── databases/
│       ├── sqlite.yaml           # SQLite-specific config
│       ├── postgres.yaml         # PostgreSQL-specific config
│       └── mysql.yaml           # MySQL-specific config
├── generated/                   # Auto-generated configs
│   ├── sqlc-sqlite.yaml
│   ├── sqlc-postgres.yaml
│   └── sqlc-mysql.yaml
└── backup/                      # Your original configs
    └── sqlc.yaml.<timestamp>
```

## Usage

### Option 1: Use Individual Databases

```bash
# SQLite only
sqlc -f config/generated/sqlc-sqlite.yaml generate

# PostgreSQL only
sqlc -f config/generated/sqlc-postgres.yaml generate

# MySQL only
sqlc -f config/generated/sqlc-mysql.yaml generate
```

### Option 2: Rebuild Complete Configuration

```bash
# Build complete config from components
./scripts/build-config.sh "sqlite,postgres,mysql"

# Use generated config
sqlc generate
```

### Option 3: Custom Database Selection

```bash
# Only build PostgreSQL and MySQL
./scripts/build-config.sh "postgres,mysql"
```

## Benefits

- ✅ Maintainable: Each database < 150 lines
- ✅ Reusable: Shared base configuration
- ✅ Focused: Database-specific optimizations
- ✅ Testable: Individual database validation
- ✅ Versioned: Backup of original configuration

## Migration Steps

1. Test with new configurations
2. Update CI/CD to use new approach
3. Delete old `sqlc.yaml` when confident
