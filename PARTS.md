# Project Components Analysis

> Analysis of extractable libraries/SDKs from template-sqlc
> Last Updated: 2026-02-28

## Executive Summary

This project contains **7 potential reusable components** that could be extracted as standalone libraries. Three are **high-value** (unique differentiation), three are **medium-value** (incremental improvements), and one is **low-value** (commodity).

---

## High-Value Extractions

### 1. `sqlc-converter` — Database-Agnostic Type Conversion Library

**Location:** `internal/adapters/converters/`

**What it does:** Provides database-agnostic type conversion interfaces and implementations for common Go types (UUID, time, bool) and domain types (Email, Username, PasswordHash, UserStatus, etc.).

**Key Features:**

- Interface-based converter pattern (`DomainToDB`, `DBToDomain`)
- Database-specific implementations:
  - **SQLite:** UUID as string, bool as int/string
  - **PostgreSQL:** Native UUID type, native bool
  - **MySQL:** UUID as binary, bool as TINYINT
- Domain value object conversion with validation
- Factory functions for database-specific converter selection

**Alternatives:**
| Library | Limitation |
|---------|------------|
| `github.com/lib/pq` | PostgreSQL-specific, no interface abstraction |
| `github.com/go-sql-driver/mysql` | MySQL-specific, no domain types |
| `github.com/mattn/go-sqlite3` | SQLite-specific, no abstraction |
| `gorm.io/gorm/schema` | Tied to GORM, no domain value objects |
| Manual conversion | Boilerplate, error-prone, inconsistent |

**Our Unique Value:**

1. **Database portability:** Switch databases without changing domain layer
2. **Domain value objects:** Built-in validation for Email, Username, PasswordHash
3. **Interface-first:** Enables mocking, testing, and clean architecture
4. **Zero dependencies:** No ORM lock-in, works with any sql driver

**Extraction Path:**

```
sqlc-converter/
├── converter.go          # Core interfaces
├── uuid.go               # UUID converters (sqlite, postgres, mysql)
├── time.go               # Time converters
├── bool.go               # Boolean converters
├── domain/
│   ├── email.go          # Email value object + converter
│   ├── username.go       # Username value object + converter
│   └── password.go       # PasswordHash value object + converter
├── factory.go            # NewUUIDConverter(db), NewTimeConverter(db), etc.
└── errors.go             # ConversionError type
```

**Recommended Package:** `github.com/LarsArtmann/sqlc-converter`

---

### 2. `sqlc-config` — Comprehensive Configuration Template Library

**Location:** `sqlc.yaml`, `config/builder.go`, `config/modular/`, `config/internal/`

**What it does:** Provides production-ready sqlc configuration templates and a builder for generating custom configurations.

**Key Features:**

- **850+ line comprehensive template** with:
  - CEL validation rules (no-select-star, no-delete-without-where, no-drop-table, require-limit)
  - Multi-database support (SQLite/Turso, PostgreSQL, MySQL)
  - All emit\_\* options documented with trade-offs
  - Type overrides for 30+ common patterns
  - Verified WASM plugins (Python, Kotlin, TypeScript)
  - Project type templates (hobby, enterprise, microservices, multi-tenant, etc.)
- **Modular builder** for composing custom configurations
- **Database-specific presets** with isolated configs

**Alternatives:**
| Source | Limitation |
|--------|------------|
| sqlc docs examples | Fragmented, not comprehensive |
| `github.com/sqlc-dev/sqlc` examples | Basic only, no validation rules |
| Community templates | Inconsistent quality, unmaintained |
| Manual config | Time-consuming, error-prone, incomplete |

**Our Unique Value:**

1. **Exhaustive coverage:** Every sqlc feature documented and explained
2. **Production-hardened:** CEL rules prevent dangerous queries
3. **Zero-copy adoption:** Copy-paste ready for any project type
4. **Builder pattern:** Generate configs programmatically for CI/CD

**Extraction Path:**

```
sqlc-config/
├── preset/
│   ├── sqlite.yaml       # SQLite production config
│   ├── postgres.yaml     # PostgreSQL production config
│   ├── mysql.yaml        # MySQL production config
│   └── minimal.yaml      # Quick-start config
├── rules/
│   ├── safety.yaml       # no-select-star, no-delete-without-where
│   └── performance.yaml  # require-limit-on-select
├── templates/
│   ├── hobby.yaml        # Simple single-db
│   ├── enterprise.yaml   # Multi-env, full validation
│   ├── microservices.yaml # Service-specific configs
│   └── multi-tenant.yaml # Schema-per-tenant
├── builder/
│   ├── builder.go        # Config composition
│   └── validator.go      # Config validation
└── overrides/
    ├── sqlite.yaml       # SQLite type mappings
    ├── postgres.yaml     # PostgreSQL type mappings
    └── mysql.yaml        # MySQL type mappings
```

**Recommended Package:** `github.com/LarsArtmann/sqlc-config`

---

### 3. `domainkit` — Strong-Typed Domain Entity Patterns

**Location:** `internal/domain/entities/`

**What it does:** Provides patterns for creating strongly-typed domain entities with value objects, validation, and immutable accessors.

**Key Features:**

- **Value objects with validation:**
  - `Email` — validated email with normalization
  - `Username` — length constraints, trimming
  - `PasswordHash` — minimum length validation
  - `FirstName`/`LastName` — non-empty validation
  - `UserStatus`/`UserRole` — enum with `IsValid()` check
- **Entity pattern:**
  - Private fields with getter methods
  - Factory constructors (`NewUser`, `NewEmail`)
  - Behavior methods (`ChangeStatus`, `Verify`, `RecordLogin`)
- **Domain errors:** Sentinel errors for type-safe error handling

**Alternatives:**
| Library | Limitation |
|---------|------------|
| `go-playground/validator` | Struct tags only, no value object encapsulation |
| `github.com/go-ozzo/ozzo-validation` | Validation focused, no entity pattern |
| Manual structs | No encapsulation, validation scattered |
| DDD frameworks | Heavy, opinionated, overkill for simple cases |

**Our Unique Value:**

1. **Encapsulation:** Private fields prevent invalid state
2. **Self-validating:** Value objects validate at construction
3. **Zero dependencies:** Pure Go, no reflection
4. **Composable:** Mix and match value objects as needed
5. **sqlc-friendly:** Works seamlessly with generated code

**Extraction Path:**

```
domainkit/
├── value/
│   ├── email.go          # Email value object
│   ├── username.go       # Username value object
│   ├── password.go       # PasswordHash value object
│   ├── name.go           # FirstName, LastName
│   └── id.go             # Generic ID type
├── enum/
│   ├── status.go         # UserStatus enum pattern
│   └── role.go           # UserRole enum pattern
├── entity/
│   ├── base.go           # Common entity patterns
│   └── timestamps.go     # CreatedAt, UpdatedAt handling
├── errors/
│   └── domain.go         # Domain error types
└── meta/
    └── metadata.go       # Flexible metadata map
```

**Recommended Package:** `github.com/LarsArtmann/domainkit`

---

## Medium-Value Extractions

### 4. `sqlc-repo` — Repository Interface & Adapter Pattern

**Location:** `internal/domain/repositories/`, `internal/adapters/`

**What it does:** Provides repository interfaces and database-specific adapter implementations for sqlc-generated code.

**Key Features:**

- **Database-agnostic interfaces:**
  - CRUD operations
  - List/Search with pagination
  - Aggregations (Count, Stats)
  - Authentication operations
  - Transaction support
- **Database-specific adapters:**
  - SQLite, PostgreSQL, MySQL implementations
  - Mapper pattern for entity ↔ model conversion
  - Error translation (SQL errors → domain errors)

**Alternatives:**
| Library | Limitation |
|---------|------------|
| `gorm.io/gorm` | ORM, not sqlc-compatible |
| `github.com/jmoiron/sqlx` | No interface abstraction |
| Manual repositories | Repetitive, error-prone |
| Clean architecture templates | Generic, not sqlc-optimized |

**Our Unique Value:**

1. **sqlc-native:** Designed for sqlc-generated code
2. **Interface-first:** Testable, mockable repositories
3. **Error translation:** Consistent domain errors
4. **Transaction abstraction:** Clean transaction handling

**Recommendation:** Medium priority. The pattern is valuable but project-specific. Extract as **example templates** rather than a library.

---

### 5. `sqlc-metrics` — Prometheus Metrics for sqlc Applications

**Location:** `internal/monitoring/metrics.go`

**What it does:** Pre-configured Prometheus metrics for monitoring sqlc-based applications.

**Key Features:**

- Code generation metrics (duration, errors, total)
- Query metrics (duration, errors, total, connections)
- Domain operation metrics (user ops, sessions)
- Build metrics (duration, success, failures)
- Built-in HTTP server for `/metrics` endpoint

**Alternatives:**
| Library | Limitation |
|---------|------------|
| `prometheus/client_golang` | Requires manual setup |
| OpenTelemetry | More complex, different paradigm |
| Custom metrics | Reinventing the wheel |

**Our Unique Value:**

1. **sqlc-specific:** Metrics tailored to sqlc workflow
2. **Ready-to-use:** Pre-configured histograms with sensible buckets
3. **Domain-aware:** User/session operation metrics included

**Recommendation:** Medium priority. Good for teams using Prometheus + sqlc, but niche.

---

### 6. `sqlc-testing` — Testing Utilities for sqlc Projects

**Location:** `test/`, `internal/tests/`

**What it does:** Testing utilities and patterns for sqlc-based applications.

**Key Features:**

- Test data management (`test/testdata/`)
- BDD test structure (`test/features/`)
- Integration test patterns (`internal/tests/integration/`)
- E2E test patterns (`internal/tests/e2e/`)
- Database fixture management

**Alternatives:**
| Library | Limitation |
|---------|------------|
| `testcontainers/testcontainers-go` | Heavy, Docker required |
| `DATA-DOG/go-sqlmock` | Mock only, no real DB testing |
| `onsi/ginkgo` | BDD framework only, no sqlc utilities |

**Our Unique Value:**

1. **sqlc-specific:** Patterns for testing generated code
2. **Multi-database:** SQLite, PostgreSQL, MySQL test fixtures
3. **BDD-ready:** Ginkgo/Gomega integration

**Recommendation:** Medium priority. Valuable patterns but better as documentation/examples.

---

## Low-Value Extractions

### 7. `sqlc-examples` — SQL Schema & Query Examples

**Location:** `examples/`

**What it does:** Example SQL schemas and queries for different databases.

**Key Features:**

- User management schemas (SQLite, PostgreSQL, MySQL)
- Session management patterns
- Database-specific features (FTS, JSON, Enums)

**Alternatives:**
| Source | Quality |
|--------|---------|
| sqlc docs | Official, maintained |
| `github.com/sqlc-dev/sqlc/tree/main/examples` | Official examples |
| Real-world projects | More comprehensive |

**Our Unique Value:**

- Minimal. These are educational examples, not a library.

**Recommendation:** Low priority. Keep as project documentation, not a separate library.

---

## Summary Matrix

| Component        | Value  | Uniqueness | Extraction Effort | Recommendation   |
| ---------------- | ------ | ---------- | ----------------- | ---------------- |
| `sqlc-converter` | HIGH   | HIGH       | Medium            | **Extract now**  |
| `sqlc-config`    | HIGH   | HIGH       | Low               | **Extract now**  |
| `domainkit`      | HIGH   | MEDIUM     | Medium            | **Extract soon** |
| `sqlc-repo`      | MEDIUM | MEDIUM     | High              | Templates only   |
| `sqlc-metrics`   | MEDIUM | LOW        | Low               | Niche library    |
| `sqlc-testing`   | MEDIUM | LOW        | Medium            | Documentation    |
| `sqlc-examples`  | LOW    | LOW        | N/A               | Keep in project  |

---

## Extraction Roadmap

### Phase 1: Immediate (High ROI)

1. **`sqlc-config`** — Copy existing config, add presets, publish
2. **`sqlc-converter`** — Extract interfaces, add factory functions, publish

### Phase 2: Near-term (Medium ROI)

3. **`domainkit`** — Generalize value objects, add generic ID types, publish

### Phase 3: Future (Optional)

4. **`sqlc-repo`** — Publish as example templates, not a library
5. **`sqlc-metrics`** — Publish if demand exists

---

## Implementation Notes

### For `sqlc-converter`:

```go
// Usage example
import "github.com/LarsArtmann/sqlc-converter"

func main() {
    // Database-specific converters
    uuidConv := sqlc_converter.NewUUIDConverter("postgres")
    timeConv := sqlc_converter.NewTimeConverter("postgres")

    // Domain converters (database-agnostic)
    emailConv := sqlc_converter.NewEmailConverter()

    // Use in repository layer
    dbUUID := uuidConv.DomainToDB(user.UUID())
    email, _ := emailConv.DBToDomain(dbEmail)
}
```

### For `sqlc-config`:

```bash
# Usage example
# Copy preset
cp sqlc-config/preset/postgres.yaml myproject/sqlc.yaml

# Or use builder
go run github.com/LarsArtmann/sqlc-config/cmd/builder \
    --database postgres \
    --rules safety,performance \
    --output sqlc.yaml
```

### For `domainkit`:

```go
// Usage example
import "github.com/LarsArtmann/domainkit/value"
import "github.com/LarsArtmann/domainkit/enum"

type User struct {
    id       int64
    email    value.Email
    username value.Username
    status   enum.UserStatus
}

func NewUser(email, username string) (*User, error) {
    e, err := value.NewEmail(email)
    if err != nil {
        return nil, err
    }
    u, err := value.NewUsername(username)
    if err != nil {
        return nil, err
    }
    return &User{email: e, username: u, status: enum.StatusActive}, nil
}
```

---

## Related Files

- `PROJECT_SPLIT_EXECUTIVE_REPORT.md` — Original split proposal (focuses on project structure)
- `config/MIGRATION_GUIDE.md` — Configuration migration documentation
- `HOW_TO_GOLANG.md` (external) — Go development standards

---

## Conclusion

**Three libraries are worth extracting:**

1. **`sqlc-converter`** — Fills a real gap in the ecosystem for database-agnostic type conversion
2. **`sqlc-config`** — Production-ready configs that save significant setup time
3. **`domainkit`** — Clean value object patterns that complement sqlc's generated code

The remaining components are better served as **documentation, examples, or templates** rather than standalone libraries.
