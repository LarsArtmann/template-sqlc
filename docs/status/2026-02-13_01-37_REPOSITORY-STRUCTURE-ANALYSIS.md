# Repository Structure Analysis

**Date:** 2026-02-13 01:37 CET
**Focus:** Complete inventory and purpose of all repository files

---

## Overview

`template-sqlc` is a production-ready sqlc configuration template supporting SQLite/Turso, PostgreSQL, and MySQL with multi-language code generation (Go, Python, Kotlin, TypeScript) via WASM plugins.

---

## Root Files

| File                      | Purpose                                                                           |
| ------------------------- | --------------------------------------------------------------------------------- |
| `sqlc.yaml`               | Main config - 853 lines with rules, plugins, overrides, 3 database configs        |
| `Makefile`                | Build targets: `validate`, `generate`, `examples`, `clean`, `test`, `docker-test` |
| `justfile`                | Modern task runner: `watch`, `format`, `lint`, `benchmark`                        |
| `docker-compose.test.yml` | PostgreSQL, MySQL, SQLite containers for integration testing                      |
| `.golangci.yml`           | Linter config: gosec, errcheck, staticcheck, govet                                |
| `.gitignore`              | Ignores generated code, databases, IDE files, .env                                |
| `LICENSE`                 | MIT License                                                                       |
| `CONTRIBUTING.md`         | Contribution guidelines                                                           |

---

## `.github/` - CI/CD

| File                     | Purpose                                                                    |
| ------------------------ | -------------------------------------------------------------------------- |
| `workflows/validate.yml` | CI workflow - runs on push/PR to main/master, tests multiple sqlc versions |

---

## `.readme/` - Auto-README Generation

| File                         | Purpose                              |
| ---------------------------- | ------------------------------------ |
| `configs/readme-config.yaml` | Config for template-readme generator |
| `assets/language-trends.svg` | Language trends graphic              |

---

## `config/` - Modular Configuration System

Go-based system for assembling `sqlc.yaml` from modular components.

| Path                               | Purpose                                  |
| ---------------------------------- | ---------------------------------------- |
| `builder.go`                       | Assembles fragments into final sqlc.yaml |
| `go.mod`, `go.sum`                 | Go module definitions                    |
| `sqlc.yaml`                        | Generated output (auto-generated)        |
| `MIGRATION_GUIDE.md`               | Migration documentation                  |
| `internal/base/common.yaml`        | Shared: rules, plugins, global settings  |
| `internal/databases/sqlite.yaml`   | SQLite-specific config                   |
| `internal/databases/postgres.yaml` | PostgreSQL-specific config               |
| `internal/databases/mysql.yaml`    | MySQL-specific config                    |
| `modular/sqlc-*.yaml`              | Standalone single-database configs       |
| `generated/sqlc-*.yaml`            | Auto-generated per-database configs      |
| `backup/`                          | Backup of original monolithic config     |

**Use Case:** Enterprise teams managing complex multi-database setups with separate concerns.

---

## `examples/` - Working Database Examples

### `sqlite/`

| File               | Purpose                                                                  |
| ------------------ | ------------------------------------------------------------------------ |
| `user.sql`         | Schema: users table with FTS5 virtual table, generated columns, triggers |
| `queries/user.sql` | Queries: CreateUser, GetUser, SearchUsers (FTS), UpdateUser, DeleteUser  |

**Features Demonstrated:** FTS5 full-text search, generated columns, AFTER triggers

### `postgres/`

| File               | Purpose                                                             |
| ------------------ | ------------------------------------------------------------------- |
| `user.sql`         | Schema: UUID primary key, CITEXT email, ENUM status, JSONB metadata |
| `queries/user.sql` | Queries: Uses `sqlc.narg`, array operations, JSONB operators        |

**Features Demonstrated:** UUID type, CITEXT, custom ENUM, JSONB, tsvector

### `mysql/`

| File               | Purpose                                                    |
| ------------------ | ---------------------------------------------------------- |
| `user.sql`         | Schema: BINARY(16) UUID, ENUM status, JSON, FULLTEXT index |
| `queries/user.sql` | Queries: UUID_TO_BIN/BIN_TO_UUID functions                 |

**Features Demonstrated:** Binary UUID storage, JSON columns, FULLTEXT search

---

## `internal/` - Domain-Driven Design Reference

Complete DDD architecture showing how to integrate sqlc-generated code.

### `domain/entities/`

| File         | Purpose                                                                     |
| ------------ | --------------------------------------------------------------------------- |
| `user.go`    | User entity with value objects (UserID, Email, Username), factory functions |
| `session.go` | UserSession entity with SessionToken, SessionDeviceInfo                     |
| `errors.go`  | Domain errors: ValidationError, NotFoundError, ConflictError, etc.          |

### `domain/services/`

| File              | Purpose                                                                 |
| ----------------- | ----------------------------------------------------------------------- |
| `user_service.go` | UserService: CreateUser, AuthenticateUser, business logic orchestration |

### `domain/repositories/`

| File                 | Purpose                                                 |
| -------------------- | ------------------------------------------------------- |
| `user_repository.go` | Interfaces: UserRepository, SessionRepository contracts |

### `domain/events/`

| File             | Purpose                                                           |
| ---------------- | ----------------------------------------------------------------- |
| `user_events.go` | Domain events: UserCreated, UserUpdated, EventPublisher interface |

### `adapters/`

| File/Dir                            | Purpose                                 | Status                         |
| ----------------------------------- | --------------------------------------- | ------------------------------ |
| `sqlite/user_repository.go`         | SQLite implementation                   | Stub (`panic("implement me")`) |
| `sqlite/session_repository.go`      | SQLite session repo                     | Stub                           |
| `sqlite/working_user_repository.go` | Working reference implementation        | **Functional**                 |
| `postgres/user_repository.go`       | PostgreSQL implementation               | Stub                           |
| `mysql/user_repository.go`          | MySQL implementation                    | Stub                           |
| `mappers/user_mapper.go`            | Domain ↔ DB conversion                  | Stub                           |
| `converters/types.go`               | UUIDConverter, TimeConverter interfaces | Defined                        |

### `monitoring/`

| File         | Purpose                                                     |
| ------------ | ----------------------------------------------------------- |
| `metrics.go` | Prometheus metrics for code generation timing, query counts |

### `tests/`

| File                               | Purpose                           |
| ---------------------------------- | --------------------------------- |
| `unit/user_test.go`                | Unit tests using testify          |
| `integration/user_service_test.go` | Integration tests with mocks      |
| `bdd/user_features_test.go`        | BDD tests using godog             |
| `e2e/`                             | End-to-end test directory (empty) |

---

## `test/` - Test Data

| Path                                    | Purpose                         |
| --------------------------------------- | ------------------------------- |
| `features/user/user_management.feature` | Gherkin scenarios for user CRUD |
| `testdata/sqlite/`                      | SQLite test fixtures            |
| `testdata/postgres/`                    | PostgreSQL test fixtures        |
| `testdata/mysql/`                       | MySQL test fixtures             |

---

## `scripts/` - Build Automation

| Script                     | Purpose                                        |
| -------------------------- | ---------------------------------------------- |
| `validate-config.sh`       | Validates sqlc.yaml syntax with `sqlc compile` |
| `build-config.sh`          | Runs Go builder to assemble config             |
| `benchmark.sh`             | Performance benchmarks using hyperfine         |
| `migrate-config.sh`        | Migrates from monolithic to modular config     |
| `build-database-config.sh` | Builds database-specific configs               |
| `sqlc/`                    | sqlc-specific scripts                          |

---

## `docs/` - Documentation

| Path                           | Purpose                           |
| ------------------------------ | --------------------------------- |
| `status/2025-12-15_18-39_*.md` | Architecture restructuring report |
| `status/2025-12-15_19-45_*.md` | Major breakthroughs day report    |
| `status/2026-02-12_19-27_*.md` | Previous status report            |

---

## Component Maturity Matrix

| Component                | Status                   | Notes                                                        |
| ------------------------ | ------------------------ | ------------------------------------------------------------ |
| `sqlc.yaml`              | Production Ready         | Fully documented, all features working                       |
| `examples/`              | Production Ready         | Working schemas and queries for all 3 DBs                    |
| `config/` modular system | Production Ready         | Builder works, generates valid configs                       |
| `internal/domain/`       | Reference Implementation | Complete architecture, good patterns                         |
| `internal/adapters/`     | Partial                  | `working_user_repository.go` is functional, others are stubs |
| `internal/tests/`        | Framework Ready          | Test structure exists, needs real DB tests                   |
| `scripts/`               | Production Ready         | All scripts functional                                       |
| CI/CD                    | Production Ready         | GitHub Actions validates on PR                               |

---

## What You Actually Need

| Use Case               | Take This                       |
| ---------------------- | ------------------------------- |
| Simple project         | Copy `sqlc.yaml` + adjust paths |
| Learn sqlc patterns    | Study `examples/`               |
| Large team/enterprise  | Use `config/` modular system    |
| Reference architecture | Study `internal/` structure     |

---

## Recommendations

1. **For most users:** Just copy `sqlc.yaml` and `examples/<your-db>/`
2. **For enterprise:** Use `config/` builder for team-managed configs
3. **For learning:** Study `internal/domain/` for DDD patterns with sqlc
4. **Adapters:** Only `working_user_repository.go` is functional - use as reference

---

## Quick Start

```bash
# Copy main config
cp sqlc.yaml your-project/

# Copy relevant examples
cp -r examples/sqlite your-project/sql/

# Generate code
cd your-project && sqlc generate
```

---

**Summary:** This is a comprehensive template repository. The core `sqlc.yaml` and `examples/` are production-ready. The `internal/` DDD structure serves as a reference implementation - not all adapters are fully implemented.
