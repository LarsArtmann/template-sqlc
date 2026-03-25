# template-sqlc Project Context

Production-ready Go template using sqlc for type-safe SQL code generation with multi-database support.

## Build Commands

```bash
just build          # Build the project
just test           # Run tests
just lint           # Run golangci-lint
just generate       # Run sqlc generate
```

## Architecture

- **Hexagonal/Clean Architecture**: Domain entities, services, repositories, adapters
- **Multi-database support**: SQLite, PostgreSQL, MySQL
- **sqlc**: Type-safe SQL code generation from queries

## Project Structure

```
internal/
├── domain/         # Core business logic
│   ├── entities/   # Domain entities (User, Session)
│   ├── services/   # Business services
│   └── repositories/ # Repository interfaces
├── adapters/       # Infrastructure adapters
│   ├── sqlite/     # SQLite implementations
│   ├── postgres/   # PostgreSQL implementations
│   ├── mysql/      # MySQL implementations
│   ├── mappers/    # DTO mappers
│   └── converters/ # Type converters
├── db/             # sqlc generated code
│   ├── sqlite/
│   ├── postgres/
│   └── mysql/
├── monitoring/     # Metrics and observability
└── tests/          # Test suites
    ├── unit/
    ├── integration/
    ├── e2e/
    └── bdd/        # Godog/Cucumber tests
```

## Database Build Tags

Each database has its own build tag:

- `//go:build sqlite` - SQLite code
- `//go:build postgres` - PostgreSQL code
- `//go:build mysql` - MySQL code

## Code Style

- Functional programming patterns preferred
- Early returns over nested conditionals
- Explicit over implicit
- Descriptive names over comments
- Type-first development

## Testing

- Unit tests in `internal/tests/unit/`
- Integration tests in `internal/tests/integration/`
- BDD tests in `internal/tests/bdd/` using Godog

## SQL Schema & Queries

- Schema files: `sql/<db>/schema/*.sql`
- Query files: `sql/<db>/queries/*.sql`
- Generated code: `internal/db/<db>/`

## Configuration

- sqlc config: `sqlc.yaml`
- Linting: `.golangci.yml`
- Build: `justfile`
