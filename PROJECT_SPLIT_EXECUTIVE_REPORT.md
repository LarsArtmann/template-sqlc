# Project Split Executive Report: template-sqlc

## Introduction
This report outlines a strategy to refactor the `template-sqlc` project into multiple, highly focused, and independently manageable projects. The current project serves as a comprehensive template for generating type-safe Go code from SQL schemas using `sqlc`, supporting multiple database backends (MySQL, PostgreSQL, SQLite). While robust, its monolithic nature can pose challenges for maintainability, independent evolution of components, and specialized tooling.

## Proposed Project Splits

The `template-sqlc` project can be logically decomposed into the following six distinct, focused projects:

### 1. `sqlc-config-manager`
*   **Focus:** Centralized management, generation, and validation of `sqlc.yaml` configurations.
*   **Description:** This project would encapsulate all logic related to `sqlc` configuration files, including the `config/builder.go` for dynamic configuration generation and `scripts/sqlc` for configuration-related tasks. It could evolve into a powerful CLI tool or library for streamlined `sqlc` setup and maintenance.
*   **Key Components:** `config/builder.go`, `config/modular/`, `config/generated/`, `scripts/sqlc/`.

### 2. `sqlc-database-examples`
*   **Focus:** Standalone collection of SQL schemas and queries for `sqlc` usage examples.
*   **Description:** A dedicated repository for the `examples/` directory, providing clear, concise, and ready-to-use SQL examples for SQLite, MySQL, and PostgreSQL. This project would serve as an educational resource, quick-start guide, and reference for various database integrations with `sqlc`.
*   **Key Components:** `examples/sqlite/`, `examples/mysql/`, `examples/postgres/`.

### 3. `sqlc-go-adapters`
*   **Focus:** Database-specific Go interfaces and concrete implementations for `sqlc`-generated code.
*   **Description:** This Go module would house the `internal/adapters` directory, offering abstract interfaces and concrete implementations for database interaction. It would include generic converters, mappers, and specific implementations for each supported database, allowing the core domain logic to remain database-agnostic.
*   **Key Components:** `internal/adapters/converters/`, `internal/adapters/postgres/`, `internal/adapters/sqlite/`, `internal/adapters/mysql/`, `internal/adapters/mappers/`.

### 4. `sqlc-domain-model`
*   **Focus:** Core business logic, entities, and repository interfaces.
*   **Description:** This Go module would contain `internal/domain`, comprising the application's core entities, domain services, events, and abstract repository interfaces. It would be entirely decoupled from specific database implementations, relying on the interfaces defined in `sqlc-go-adapters`.
*   **Key Components:** `internal/domain/errors/`, `internal/domain/repositories/`, `internal/domain/entities/`, `internal/domain/events/`, `internal/domain/services/`.

### 5. `sqlc-testing-framework`
*   **Focus:** Standardized testing utilities and frameworks for `sqlc`-based projects.
*   **Description:** A dedicated project for all testing-related infrastructure, including `test/` and `internal/tests`. It would provide reusable components for unit, integration, BDD, and E2E testing, ensuring consistent and robust validation of `sqlc`-generated code and database interactions across projects.
*   **Key Components:** `test/testdata/`, `test/features/`, `internal/tests/bdd/`, `internal/tests/e2e/`, `internal/tests/unit/`, `internal/tests/integration/`.

### 6. `sqlc-cli-tool`
*   **Focus:** An overarching command-line interface for the entire `sqlc` development workflow.
*   **Description:** This project would provide a user-friendly CLI that integrates and orchestrates the functionalities of `sqlc-config-manager`, `sqlc-go-adapters`, and `sqlc-domain-model`. It would simplify tasks such as project initialization, code generation, migration, and local testing, abstracting away underlying complexities.
*   **Key Components:** A new top-level `cmd/` directory, integrating the other modules.

## Benefits of Splitting

*   **Improved Maintainability:** Smaller, focused codebases are easier to understand, debug, and maintain.
*   **Independent Evolution:** Each project can evolve, be versioned, and deployed independently, reducing the risk of introducing breaking changes across unrelated components.
*   **Clearer Ownership & Responsibility:** Each team or developer can have clear ownership of specific components.
*   **Enhanced Testability:** Isolated projects facilitate more targeted and efficient testing.
*   **Reusability:** Components like `sqlc-go-adapters` and `sqlc-testing-framework` become reusable in other `sqlc`-based projects.
*   **Scalability:** Allows for independent scaling of different services if the projects were to become microservices.
*   **Reduced Build Times:** Smaller codebases lead to faster build and CI/CD pipeline execution.

## Challenges and Considerations

*   **Dependency Management:** Careful management of dependencies between the new projects will be crucial (e.g., `sqlc-domain-model` depending on `sqlc-go-adapters` interfaces).
*   **CI/CD Pipeline Setup:** Each new project will require its own CI/CD pipeline, increasing initial setup overhead.
*   **Version Control:** Strategies for versioning and releasing interdependent modules need to be established (e.g., semantic versioning).
*   **Local Development Experience:** Ensuring a smooth local development experience across multiple repositories will be important, potentially requiring monorepo tools if the projects remain closely coupled.
*   **Data Migration & Schema Evolution:** Coordination across projects for database schema changes and migrations will require robust processes.

## Conclusion

Splitting the `template-sqlc` project into these focused components offers significant long-term benefits in terms of modularity, maintainability, and scalability. While there are initial challenges in setting up the new structure and managing inter-project dependencies, the advantages of clearer separation of concerns and independent evolution will ultimately lead to a more robust and adaptable system. This architectural shift aligns with best practices for building complex software systems and will empower more efficient development and deployment cycles.
