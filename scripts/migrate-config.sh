#!/bin/bash

# Migrate from monolithic sqlc.yaml to modular approach
# Preserves existing configuration while improving maintainability

set -euo pipefail

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

BACKUP_DIR="config/backup"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

echo -e "${BLUE}🔄 sqlc Configuration Migration${NC}"
echo -e "${YELLOW}📅 Timestamp: $TIMESTAMP${NC}"

# Check if monolithic sqlc.yaml exists
if [ ! -f "sqlc.yaml" ]; then
    echo -e "${RED}❌ sqlc.yaml not found${NC}"
    exit 1
fi

# Create backup
echo -e "${YELLOW}💾 Creating backup...${NC}"
mkdir -p "$BACKUP_DIR"
cp sqlc.yaml "$BACKUP_DIR/sqlc.yaml.$TIMESTAMP"
echo -e "${GREEN}✅ Backed up to: $BACKUP_DIR/sqlc.yaml.$TIMESTAMP${NC}"

# Analyze current configuration
echo -e "${YELLOW}📊 Analyzing current configuration...${NC}"
DATABASE_COUNT=$(yq e '.sql | length' sqlc.yaml 2>/dev/null || echo "0")
RULE_COUNT=$(yq e '.rules | length' sqlc.yaml 2>/dev/null || echo "0")
PLUGIN_COUNT=$(yq e '.plugins | length' sqlc.yaml 2>/dev/null || echo "0")

echo -e "  🗃️  Databases: $DATABASE_COUNT"
echo -e "  📏 Rules: $RULE_COUNT"
echo -e "  🔌 Plugins: $PLUGIN_COUNT"

# Extract database configurations
echo -e "${YELLOW}🔧 Extracting database configurations...${NC}"
for i in $(seq 0 $((DATABASE_COUNT - 1))); do
    DB_NAME=$(yq e ".sql[$i].name" sqlc.yaml 2>/dev/null || echo "database$i")
    echo -e "  📦 Extracting: $DB_NAME"
    
    # Extract to separate file
    yq e ".sql[$i]" sqlc.yaml > "config/extracted-${DB_NAME}.yaml"
done

# Create modular configurations
echo -e "${YELLOW}🏗️  Creating modular configurations...${NC}"
for extracted in config/extracted-*.yaml; do
    if [ -f "$extracted" ]; then
        DB_NAME=$(basename "$extracted" .yaml | sed 's/extracted-//')
        ./scripts/build-database-config.sh "$DB_NAME" "config/modular"
    fi
done

# Create usage guide
echo -e "${YELLOW}📚 Creating usage guide...${NC}"
cat > "config/MIGRATION_GUIDE.md" << 'EOF'
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
EOF

echo -e "${GREEN}✅ Migration complete!${NC}"
echo -e "${YELLOW}📚 See config/MIGRATION_GUIDE.md for usage instructions${NC}"

# Clean up extraction files
rm -f config/extracted-*.yaml

# Show before/after comparison
echo -e "${BLUE}📊 Configuration Size Comparison:${NC}"
BACKUP_SIZE=$(wc -l < "$BACKUP_DIR/sqlc.yaml.$TIMESTAMP")
echo -e "  📄 Original: $BACKUP_SIZE lines"

for config in config/modular/sqlc-*.yaml; do
    if [ -f "$config" ]; then
        LINES=$(wc -l < "$config")
        NAME=$(basename "$config")
        echo -e "  📄 $NAME: $LINES lines"
    fi
done