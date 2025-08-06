# template-sqlc

> 🚀 **The Ultimate sqlc Configuration Template**  
> Production-ready, comprehensive sqlc configuration template that works for ALL project types

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![sqlc](https://img.shields.io/badge/sqlc-v1.29.0-blue)](https://sqlc.dev/)
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)](https://golang.org/)

## 🎯 Overview

This repository provides a **comprehensive, production-ready sqlc configuration template** that takes advantage of **every available sqlc feature**. Based on analysis of 21+ real-world projects, this template includes:

- ✅ **Multi-database support**: SQLite/Turso, PostgreSQL, MySQL
- ✅ **All sqlc v1.29.0 features**: Every configuration option explained
- ✅ **Production best practices**: Security, performance, validation
- ✅ **Universal compatibility**: Works for ALL project types
- ✅ **Comprehensive documentation**: Every setting explained with trade-offs

## 🚀 Quick Start

### 1. Copy the Template
```bash
# Clone this repository
git clone https://github.com/LarsArtmann/template-sqlc.git
cd template-sqlc

# Copy to your project
cp sqlc.yaml /path/to/your/project/
```

### 2. Choose Your Configuration
The template includes configurations for different project types:

```yaml
# 🔥 ACTIVE: Full multi-database setup (lines 71-350)
sql:
  - name: "sqlite"    # SQLite/Turso configuration
  - name: "postgres"  # PostgreSQL configuration  
  - name: "mysql"     # MySQL configuration

# 📝 EXAMPLES: Alternative configurations (lines 635-830)
# - Hobby/Small projects
# - Enterprise/Large projects
# - Microservices
# - Testing/CI
# - Performance-critical
# - Legacy integration
# - And more...
```

### 3. Customize for Your Needs
```bash
# Edit the configuration
vim sqlc.yaml

# Uncomment your preferred setup
# Modify paths, database URLs, and domain-specific overrides
# Test the configuration
sqlc compile
```

## 📋 Table of Contents

- [Features](#features)
- [Configuration Sections](#configuration-sections)
- [Project Type Templates](#project-type-templates)
- [Database Support](#database-support)
- [Usage Examples](#usage-examples)
- [Production Deployment](#production-deployment)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## ✨ Features

### 🛡️ **Production Safety**
- **CEL validation rules**: Prevent dangerous queries (`SELECT *`, `DELETE` without `WHERE`)
- **Strict validation**: Function checks, ORDER BY validation
- **Environment variables**: Secure credential handling
- **Prepared statements**: Performance optimization with safety

### 🎨 **Code Generation Excellence**
- **All emit_* options**: 15+ code generation features fully configured
- **Smart type overrides**: 30+ real-world type mapping patterns
- **Comprehensive naming**: Go conventions, acronym handling, domain patterns
- **Build tags**: Conditional compilation support

### 🌍 **Universal Project Support**
- **Hobby projects**: Simple, minimal configuration
- **Enterprise**: Multi-database, full validation, managed databases
- **Microservices**: Service-specific database configurations
- **Testing/CI**: Fast, in-memory database setups
- **Analytics**: Read-heavy, flexible connection patterns
- **Legacy integration**: Existing database compatibility

### 🔧 **Advanced Features**
- **Multi-language plugins**: Python, Gleam, Kotlin, TypeScript
- **Cloud integration**: sqlc Cloud support for validation
- **Flexible paths**: Directories, files, glob patterns
- **Migration patterns**: Sequential migrations, schema evolution

## 📚 Configuration Sections

### 🌐 Global Configuration
```yaml
# Cloud integration for schema verification
cloud:
  project: "your-project-id"
  token: "${SQLC_TOKEN}"

# Validation rules (CEL expressions)
rules:
  - name: "no-select-star"     # Security: prevent SELECT *
  - name: "no-delete-without-where"  # Safety: require WHERE clauses
  
# Multi-language plugins  
plugins:
  - name: "py"      # Python via WASM
  - name: "gleam"   # Gleam functional language
```

### 💾 Database Configurations

#### SQLite/Turso
```yaml
sql:
  - name: "sqlite"
    engine: "sqlite"
    queries: "sql/sqlite/queries"
    schema: "sql/sqlite/schema"
    # Full configuration with all options explained
```

#### PostgreSQL
```yaml
sql:
  - name: "postgres" 
    engine: "postgresql"
    sql_package: "pgx/v5"        # High-performance driver
    emit_prepared_queries: true   # Performance optimization
    # Comprehensive PostgreSQL type mappings
```

#### MySQL
```yaml
sql:
  - name: "mysql"
    engine: "mysql"
    sql_package: "database/sql"
    # Complete MySQL type support
```

## 🏗️ Project Type Templates

### 🔬 Hobby/Small Projects
```yaml
sql:
  - engine: "sqlite"
    queries: "queries"
    schema: "schema.sql"
    gen:
      go:
        package: "db"
        emit_json_tags: true
        emit_interface: true
```

### 🏢 Enterprise/Large Projects  
```yaml
sql:
  - name: "production"
    engine: "postgresql"
    queries: "internal/storage/queries/"
    schema: "migrations/*.up.sql"
    database:
      uri: "${PRODUCTION_DATABASE_URL}"
      managed: false
    rules:
      - "sqlc/db-prepare"
```

### 🧪 Testing/CI Configurations
```yaml
sql:
  - name: "test"
    engine: "sqlite"  
    database:
      uri: ":memory:"          # Fast in-memory tests
    gen:
      go:
        emit_interface: true    # Essential for mocking
        omit_sqlc_version: true # Cleaner test code
```

### ⚡ Performance-Critical
```yaml
sql:
  - engine: "postgresql"
    gen:
      go:
        emit_prepared_queries: true      # Pre-compiled queries
        emit_result_struct_pointers: true # Memory efficiency  
        emit_params_struct_pointers: true # Large parameter efficiency
        omit_unused_structs: true        # Smaller binaries
```

[See full project type examples in the template →](sqlc.yaml#L635-L830)

## 💾 Database Support

### 🗂️ **SQLite/Turso**
- ✅ Full-text search (FTS5) support
- ✅ JSON handling with `json.RawMessage`
- ✅ Time handling with proper nullable types
- ✅ Comprehensive type mappings

### 🐘 **PostgreSQL**
- ✅ pgx/v5 driver optimization
- ✅ Advanced types: UUID, JSONB, arrays, network types
- ✅ Prepared query support
- ✅ PostGIS compatibility ready

### 🐬 **MySQL**
- ✅ All MySQL data types supported
- ✅ Proper time zone handling
- ✅ DECIMAL precision for financial data
- ✅ Legacy database integration patterns

## 🔧 Usage Examples

### Basic Setup
```bash
# 1. Copy template
cp sqlc.yaml your-project/

# 2. Create directories
mkdir -p sql/{sqlite,postgres,mysql}/{queries,schema}

# 3. Add your schema
echo "CREATE TABLE users (id INTEGER PRIMARY KEY, email TEXT);" > sql/sqlite/schema/001_users.sql

# 4. Add queries
echo "-- name: GetUser :one\nSELECT * FROM users WHERE id = ?;" > sql/sqlite/queries/users.sql

# 5. Generate code
sqlc generate
```

### Multi-Database Project
```bash
# Generate for all databases
sqlc generate

# Output structure:
# internal/db/sqlite/    - SQLite generated code
# internal/db/postgres/  - PostgreSQL generated code  
# internal/db/mysql/     - MySQL generated code
```

### Single Database Project
```bash
# Uncomment single-database configuration in sqlc.yaml
# Example: SQLite only (lines 500-520)

sqlc generate
# Output: internal/db/
```

## 🚀 Production Deployment

### Environment Variables
```bash
# Set database URLs
export DATABASE_URL="sqlite:///prod.db"
export POSTGRES_DATABASE_URL="postgresql://user:pass@localhost/db"
export MYSQL_DATABASE_URL="mysql://user:pass@localhost/db"
export SQLC_TOKEN="your-sqlc-cloud-token"
```

### Validation
```bash
# Validate configuration
sqlc compile

# Run linting rules
sqlc vet

# Verify against live database
sqlc verify
```

### CI/CD Integration
```yaml
# .github/workflows/sqlc.yml
name: sqlc
on: [push, pull_request]
jobs:
  sqlc:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup sqlc
        run: |
          curl -sSfL https://raw.githubusercontent.com/sqlc-dev/sqlc/main/scripts/install-sqlc.sh | sh
          sudo mv ./bin/sqlc /usr/local/bin/
      - name: Validate sqlc
        run: |
          sqlc compile
          sqlc vet
```

## 🛠️ Troubleshooting

### Common Issues

#### "mutually exclusive" Errors
```bash
# Error: emit_methods_with_db_argument and emit_prepared_queries options are mutually exclusive
```
**Solution**: These options have different architectures:
- `emit_prepared_queries: true` → Better performance, connection stored on struct
- `emit_methods_with_db_argument: true` → More flexible, connection passed to methods

Choose based on your needs. The template defaults to prepared queries for performance.

#### PostgreSQL Connection Errors
```bash
# Error: no PostgreSQL database server found
```
**Solutions**:
1. **Development**: Use SQLite configuration instead
2. **Testing**: Set `managed: false` and provide connection string
3. **Production**: Use sqlc Cloud with `managed: true`

#### Type Override Errors
```bash
# Error: Package override `go_type` specifier "[]byte" is not a Go basic type
```
**Solution**: Use proper type imports or built-in types:
```yaml
# ❌ Wrong
go_type: "[]byte"

# ✅ Correct  
go_type: "string"  # Use string for most cases
# OR custom types with proper imports
```

### Performance Tips

1. **Use prepared queries**: `emit_prepared_queries: true` for repeated queries
2. **Enable struct pointers**: `emit_result_struct_pointers: true` for large structs  
3. **Omit unused**: `omit_unused_structs: true` for smaller binaries
4. **Optimize builds**: Use `build_tags` for conditional compilation

## 🤝 Contributing

We welcome contributions! This template is built from real-world usage across 21+ projects.

### How to Contribute
1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Test** your changes: `sqlc compile && sqlc vet`
4. **Commit** your changes: `git commit -m 'Add amazing feature'`
5. **Push** to the branch: `git push origin feature/amazing-feature`
6. **Open** a Pull Request

### What We're Looking For
- 🐛 **Bug fixes**: Configuration errors, typos, invalid settings
- 📚 **Documentation**: Better explanations, more examples
- 🏗️ **New patterns**: Additional project type configurations
- 🔧 **Improvements**: Better defaults, performance optimizations

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **sqlc team**: For building an amazing type-safe SQL generator
- **Community**: 21+ real projects that informed this template
- **Contributors**: Everyone who helps make this template better

---

<div align="center">

**⭐ If this template helps your project, please give it a star! ⭐**

[🐛 Report Bug](https://github.com/LarsArtmann/template-sqlc/issues) • 
[✨ Request Feature](https://github.com/LarsArtmann/template-sqlc/issues) • 
[📖 Documentation](https://docs.sqlc.dev/)

</div>