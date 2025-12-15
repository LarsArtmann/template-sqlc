.PHONY: help validate generate clean test docker-test install-deps examples

# Default target
help: ## Show this help message
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Validation targets
validate: ## Validate sqlc configuration
	@echo "🔍 Validating sqlc configuration..."
	@./scripts/validate-config.sh

validate-all: ## Validate configuration for all databases
	@echo "🔍 Validating all database configurations..."
	@for db in sqlite postgres mysql; do \
		echo "Validating $$db configuration..."; \
		sqlc -f sqlc.yaml -y sql/$$db/queries -y sql/$$db/schema compile || exit 1; \
	done

# Code generation targets
generate: ## Generate code from sqlc
	@echo "⚡ Generating code..."
	@sqlc generate

generate-verbose: ## Generate code with verbose output
	@echo "⚡ Generating code (verbose)..."
	@sqlc generate -v

generate-database: ## Generate code for specific database (usage: make generate-database DB=sqlite)
	@if [ -z "$(DB)" ]; then \
		echo "Usage: make generate-database DB=sqlite|postgres|mysql"; \
		exit 1; \
	fi
	@echo "⚡ Generating code for $(DB)..."
	@sqlc -f sqlc.yaml -y sql/$(DB)/queries -y sql/$(DB)/schema generate

# Example targets
examples: ## Set up example schemas
	@echo "📋 Setting up example schemas..."
	@mkdir -p sql/{sqlite,postgres,mysql}/{queries,schema}
	@cp -r examples/sqlite/user.sql sql/sqlite/schema/001_users.sql
	@cp -r examples/sqlite/queries/user.sql sql/sqlite/queries/
	@cp -r examples/postgres/user.sql sql/postgres/schema/001_users.sql
	@cp -r examples/postgres/queries/user.sql sql/postgres/queries/
	@cp -r examples/mysql/user.sql sql/mysql/schema/001_users.sql
	@cp -r examples/mysql/queries/user.sql sql/mysql/queries/
	@echo "✅ Example schemas set up"

test-examples: ## Test all example configurations
	@echo "🧪 Testing all example configurations..."
	@for db in sqlite postgres mysql; do \
		echo "Testing $$db example..."; \
		sqlc -f sqlc.yaml -y examples/$$db/queries -y examples/$$db/schema compile || exit 1; \
	done
	@echo "✅ All examples validated"

# Cleanup targets
clean: ## Clean generated files
	@echo "🧹 Cleaning generated files..."
	@rm -rf internal/db/
	@rm -f *.test *.out
	@echo "✅ Cleaned generated files"

# Testing targets
test: validate ## Run all tests
	@echo "🧪 Running tests..."
	@echo "✅ All tests passed"

docker-test: ## Run tests in Docker
	@echo "🐳 Running tests in Docker..."
	@docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --remove-orphans

# Development targets
install-deps: ## Install development dependencies
	@echo "📦 Installing dependencies..."
	@if command -v sqlc >/dev/null 2>&1; then \
		echo "✅ sqlc is already installed"; \
	else \
		echo "Installing sqlc..."; \
		curl -L https://github.com/kyleconroy/sqlc/releases/latest/download/sqlc_$(shell uname -s)_$(shell uname -m).tar.gz | tar -xz -C /usr/local/bin sqlc; \
	fi
	@if command -v yq >/dev/null 2>&1; then \
		echo "✅ yq is already installed"; \
	else \
		echo "Installing yq..."; \
		curl -L https://github.com/mikefarah/yq/releases/latest/download/yq_$(shell uname -s)_$(shell uname -m) -o /usr/local/bin/yq && chmod +x /usr/local/bin/yq; \
	fi

dev-setup: install-deps examples ## Set up development environment
	@echo "🛠️  Development environment setup complete"
	@echo "Next steps:"
	@echo "  1. Customize sqlc.yaml for your needs"
	@echo "  2. Add your SQL files to appropriate directories"
	@echo "  3. Run 'make generate' to create Go code"

# Documentation targets
docs: ## Generate documentation
	@echo "📚 Generating documentation..."
	@echo "# sqlc Template Documentation" > docs/generated.md
	@echo "" >> docs/generated.md
	@echo "Generated on $(shell date)" >> docs/generated.md
	@echo "" >> docs/generated.md
	@echo "## Configuration" >> docs/generated.md
	@echo '```yaml' >> docs/generated.md
	@cat sqlc.yaml >> docs/generated.md
	@echo '```' >> docs/generated.md
	@echo "✅ Documentation generated"

# Release targets
version: ## Show version information
	@echo "📋 Version information:"
	@if command -v sqlc >/dev/null 2>&1; then \
		sqlc version; \
	else \
		echo "sqlc: not installed"; \
	fi
	@if command -v yq >/dev/null 2>&1; then \
		echo "yq: $(shell yq --version)"; \
	else \
		echo "yq: not installed"; \
	fi
	@echo "Template version: $(shell git describe --tags --always --dirty 2>/dev/null || echo 'unknown')"

# Security targets
security-audit: ## Audit configuration for security issues
	@echo "🔒 Running security audit..."
	@echo "Checking for potential security issues..."
	@if grep -q "password" sqlc.yaml; then \
		echo "⚠️  Found password-related configuration - ensure secrets are properly managed"; \
	fi
	@if grep -q "token" sqlc.yaml; then \
		echo "⚠️  Found token-related configuration - ensure secrets are properly managed"; \
	fi
	@echo "✅ Security audit complete"