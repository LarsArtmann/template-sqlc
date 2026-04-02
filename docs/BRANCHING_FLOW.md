# Branching Flow & Git Workflow Documentation

**Date:** 2026-03-29
**Project:** template-sqlc
**Status:** Active

---

## Overview

This document outlines the branching strategy, git workflow, and development processes for the `template-sqlc` project. It serves as a single source of truth for contributors.

---

## Branching Strategy

### Current State

| Branch   | Purpose                          | Protection   |
| -------- | -------------------------------- | ------------ |
| `master` | Main development branch          | ✅ Protected |
| `main`   | Alias for master (CI configured) | ✅ Protected |

### Branch Naming Convention

```
<type>/<short-description>
```

**Types:**
| Type | Purpose | Example |
|------|---------|---------|
| `feat/` | New features | `feat/add-user-validation` |
| `fix/` | Bug fixes | `fix/session-expiry-bug` |
| `refactor/` | Code refactoring | `refactor/user-service` |
| `docs/` | Documentation | `docs/update-readme` |
| `chore/` | Maintenance tasks | `chore/update-deps` |
| `test/` | Test improvements | `test/add-bdd-scenarios` |

---

## Git Workflow

### Feature Branch Flow

```bash
# 1. Ensure master is up to date
git checkout master
git pull origin master

# 2. Create feature branch
git checkout -b feat/your-feature-name

# 3. Make changes, commit frequently
git add <files>
git commit -m "type: descriptive message"

# 4. Push branch
git push origin feat/your-feature-name

# 5. Create PR on GitHub
# CI runs automatically on PR

# 6. After approval, merge via PR
```

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Example:**

```
feat(user): add email validation to user creation

- Add RFC 5322 compliant email regex
- Update UserValidator with new validation
- Add test cases for edge cases

Closes #123
```

### Types

| Type       | Description                          |
| ---------- | ------------------------------------ |
| `feat`     | New feature                          |
| `fix`      | Bug fix                              |
| `docs`     | Documentation changes                |
| `style`    | Formatting, missing semicolons, etc. |
| `refactor` | Code refactoring                     |
| `test`     | Adding tests                         |
| `chore`    | Maintenance                          |

---

## CI/CD Pipeline

### GitHub Actions Workflow

**File:** `.github/workflows/validate.yml`

**Triggers:**

- Push to `master` or `main`
- Pull requests to `master` or `main`

**Jobs:**

1. **validate** - Matrix testing with sqlc versions 1.29.0 and latest

**Steps:**

```yaml
- Checkout code
- Install sqlc
- Validate configuration: ./scripts/validate-config.sh
- Check YAML syntax: yq
```

### Validation Scripts

| Script                       | Purpose                      |
| ---------------------------- | ---------------------------- |
| `scripts/validate-config.sh` | Validates sqlc configuration |
| `scripts/build-config.sh`    | Builds modular configuration |
| `scripts/migrate-config.sh`  | Migrates configurations      |

---

## Development Commands

### Build Commands (via justfile)

```bash
just build          # Build the project
just test            # Run tests
just lint           # Run golangci-lint
just generate       # Run sqlc generate
just validate       # Validate sqlc configuration
just generate-db    # Generate for specific DB
```

### Full Validation Suite

```bash
# Run all checks
sqlc compile && sqlc generate && just lint && just test
```

---

## Quality Gates

Before merging any PR, ensure:

- [ ] `sqlc compile` passes
- [ ] `sqlc generate` succeeds
- [ ] `golangci-lint` passes
- [ ] All tests pass
- [ ] No hardcoded secrets
- [ ] Documentation updated (if needed)

---

## Database-Specific Development

### Build Tags

Each database has its own build tag:

| Database   | Build Tag             | Path                          |
| ---------- | --------------------- | ----------------------------- |
| SQLite     | `//go:build sqlite`   | `internal/adapters/sqlite/`   |
| PostgreSQL | `//go:build postgres` | `internal/adapters/postgres/` |
| MySQL      | `//go:build mysql`    | `internal/adapters/mysql/`    |

### Testing Specific Database

```bash
# SQLite
go test -tags=sqlite ./...

# PostgreSQL
go test -tags=postgres ./...

# MySQL
go test -tags=mysql ./...
```

---

## Project Architecture

```
template-sqlc/
├── internal/
│   ├── domain/          # Core business logic
│   │   ├── entities/    # User, Session entities
│   │   ├── services/    # UserService
│   │   └── repositories/ # Repository interfaces
│   ├── adapters/       # Infrastructure
│   │   ├── sqlite/     # SQLite implementations
│   │   ├── postgres/   # PostgreSQL implementations
│   │   └── mysql/      # MySQL implementations
│   └── db/             # sqlc generated code
├── sql/
│   ├── sqlite/         # SQLite schemas & queries
│   ├── postgres/       # PostgreSQL schemas & queries
│   └── mysql/          # MySQL schemas & queries
└── test/
    ├── features/       # BDD feature files
    └── testdata/       # Test fixtures
```

---

## BDD Testing Workflow

### Feature File Location

```
test/features/user/user_management.feature
```

### Step Definition Location

```
internal/tests/bdd/user_features_test.go
```

### Running BDD Tests

```bash
# Run godog tests
go test -v ./internal/tests/bdd/...

# Run with tags
go test -tags=bdd ./...
```

---

## Recent Commit History Analysis

### Commit Types (Last 30)

| Type     | Count |
| -------- | ----- |
| feat     | 12    |
| refactor | 5     |
| chore    | 4     |
| docs     | 4     |
| fix      | 1     |

### Active Development Areas

1. **BDD Testing** - Comprehensive BDD test review and enhancement
2. **Repository Patterns** - Cache-specific patterns with error handling
3. **Linting** - Enhanced golangci-lint configuration
4. **Domain Services** - User repository implementations

---

## Action Items from Latest Analysis

### Priority P0 (Critical)

| Task                               | Status  | Effort |
| ---------------------------------- | ------- | ------ |
| Implement missing BDD steps        | Pending | 2h     |
| Fix background steps               | Pending | 30m    |
| Extract hardcoded test credentials | Pending | 15m    |

### Priority P1 (High)

| Task                                    | Status  | Effort |
| --------------------------------------- | ------- | ------ |
| Add t.Parallel() to tests               | Pending | 5m     |
| Rewrite scenarios from user perspective | Pending | 4h     |
| Add edge case tests                     | Pending | 3h     |

### Priority P2 (Medium)

| Task                      | Status  | Effort |
| ------------------------- | ------- | ------ |
| Evaluate Ginkgo migration | Pending | 1h     |

---

## Related Documentation

- [CONTRIBUTING.md](./CONTRIBUTING.md) - Contribution guidelines
- [BDD_TESTS_REVIEW.md](./BDD_TESTS_REVIEW.md) - BDD test analysis
- [docs/status/](./docs/status/) - Status reports

---

**Last Updated:** 2026-03-29
