#!/usr/bin/env just

# sqlc Template Automation with Just
# A more modern alternative to Make

# Default target - show help
default:
    @just --list

# Validation targets
validate: # Validate sqlc configuration
    #!/usr/bin/env bash
    echo "🔍 Validating sqlc configuration..."
    ./scripts/validate-config.sh

validate-all db="": # Validate all database configurations
    #!/usr/bin/env bash
    if [ -z "$db" ]; then
        echo "🔍 Validating all database configurations..."
        for db in sqlite postgres mysql; do
            echo "Validating $db configuration..."
            sqlc -f sqlc.yaml -y sql/$db/queries -y sql/$db/schema compile || exit 1
        done
    else
        echo "🔍 Validating $db configuration..."
        sqlc -f sqlc.yaml -y sql/$db/queries -y sql/$db/schema compile
    fi

# Code generation targets
generate: # Generate code from sqlc
    echo "⚡ Generating code..."
    sqlc generate

generate-verbose: # Generate code with verbose output
    echo "⚡ Generating code (verbose)..."
    sqlc generate -v

generate-db db="": # Generate code for specific database
    #!/usr/bin/env bash
    if [ -z "$db" ]; then
        echo "Usage: just generate-db <sqlite|postgres|mysql>"
        exit 1
    fi
    echo "⚡ Generating code for $db..."
    sqlc -f sqlc.yaml -y sql/$db/queries -y sql/$db/schema generate

# Example targets
examples: # Set up example schemas
    #!/usr/bin/env bash
    echo "📋 Setting up example schemas..."
    mkdir -p sql/{sqlite,postgres,mysql}/{queries,schema}
    cp -r examples/sqlite/user.sql sql/sqlite/schema/001_users.sql
    cp -r examples/sqlite/queries/user.sql sql/sqlite/queries/
    cp -r examples/postgres/user.sql sql/postgres/schema/001_users.sql
    cp -r examples/postgres/queries/user.sql sql/postgres/queries/
    cp -r examples/mysql/user.sql sql/mysql/schema/001_users.sql
    cp -r examples/mysql/queries/user.sql sql/mysql/queries/
    echo "✅ Example schemas set up"

test-examples: # Test all example configurations
    #!/usr/bin/env bash
    echo "🧪 Testing all example configurations..."
    for db in sqlite postgres mysql; do
        echo "Testing $db example..."
        sqlc -f sqlc.yaml -y examples/$db/queries -y examples/$db/schema compile || exit 1
    done
    echo "✅ All examples validated"

# Cleanup targets
clean: # Clean generated files
    echo "🧹 Cleaning generated files..."
    rm -rf internal/db/
    rm -f *.test *.out
    echo "✅ Cleaned generated files"

# Testing targets
test: validate # Run all tests
    echo "🧪 Running tests..."
    echo "✅ All tests passed"

# Development targets
install-deps: # Install development dependencies
    #!/usr/bin/env bash
    echo "📦 Installing dependencies..."
    if command -v sqlc >/dev/null 2>&1; then
        echo "✅ sqlc is already installed"
    else
        echo "Installing sqlc..."
        curl -L https://github.com/kyleconroy/sqlc/releases/latest/download/sqlc_$(uname -s)_$(uname -m).tar.gz | tar -xz -C /usr/local/bin sqlc
    fi
    if command -v yq >/dev/null 2>&1; then
        echo "✅ yq is already installed"
    else
        echo "Installing yq..."
        curl -L https://github.com/mikefarah/yq/releases/latest/download/yq_$(uname -s)_$(uname -m) -o /usr/local/bin/yq && chmod +x /usr/local/bin/yq
    fi

dev-setup: install-deps examples # Set up development environment
    echo "🛠️  Development environment setup complete"
    echo "Next steps:"
    echo "  1. Customize sqlc.yaml for your needs"
    echo "  2. Add your SQL files to appropriate directories"
    echo "  3. Run 'just generate' to create Go code"

# Documentation targets
docs: # Generate documentation
    echo "📚 Generating documentation..."
    echo "# sqlc Template Documentation" > docs/generated.md
    echo "" >> docs/generated.md
    echo "Generated on $(date)" >> docs/generated.md
    echo "" >> docs/generated.md
    echo "## Configuration" >> docs/generated.md
    echo '```yaml' >> docs/generated.md
    cat sqlc.yaml >> docs/generated.md
    echo '```' >> docs/generated.md
    echo "✅ Documentation generated"

# Release targets
version: # Show version information
    echo "📋 Version information:"
    if command -v sqlc >/dev/null 2>&1; then
        sqlc version
    else
        echo "sqlc: not installed"
    fi
    if command -v yq >/dev/null 2>&1; then
        echo "yq: $(yq --version)"
    else
        echo "yq: not installed"
    fi
    echo "Template version: $(git describe --tags --always --dirty 2>/dev/null || echo 'unknown')"

# Security targets
security-audit: # Audit configuration for security issues
    echo "🔒 Running security audit..."
    echo "Checking for potential security issues..."
    if grep -q "password" sqlc.yaml; then
        echo "⚠️  Found password-related configuration - ensure secrets are properly managed"
    fi
    if grep -q "token" sqlc.yaml; then
        echo "⚠️  Found token-related configuration - ensure secrets are properly managed"
    fi
    echo "✅ Security audit complete"

# Watch targets (requires entr)
watch: # Watch for changes and regenerate
    #!/usr/bin/env bash
    echo "👀 Watching for SQL changes..."
    if ! command -v entr >/dev/null 2>&1; then
        echo "entr not found. Install with: brew install entr or apt-get install entr"
        exit 1
    fi
    find sql/ -name "*.sql" | entr -r just generate

# Format targets
format: # Format SQL files
    #!/usr/bin/env bash
    echo "🎨 Formatting SQL files..."
    if command -v sqlfluff >/dev/null 2>&1; then
        find sql/ -name "*.sql" -exec sqlfluff format {} \;
        echo "✅ SQL files formatted"
    else
        echo "sqlfluff not found. Install with: pip install sqlfluff"
    fi

# Lint targets
lint: # Lint SQL files
    #!/usr/bin/env bash
    echo "🔍 Linting SQL files..."
    if command -v sqlfluff >/dev/null 2>&1; then
        sqlfluff lint sql/
    else
        echo "sqlfluff not found. Install with: pip install sqlfluff"
    fi

# Performance targets
benchmark: # Benchmark sqlc performance
    #!/usr/bin/env bash
    echo "⚡ Benchmarking sqlc performance..."
    echo "Time taken to generate code:"
    time sqlc generate