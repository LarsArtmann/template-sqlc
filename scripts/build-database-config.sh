#!/bin/bash

# Build individual database configurations
# Creates separate, maintainable configuration files

set -euo pipefail

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Parse arguments
DATABASE="${1:-sqlite}"
OUTPUT_DIR="${2:-config/generated}"

# Validate database
case "$DATABASE" in
    "sqlite"|"postgres"|"mysql")
        echo -e "${GREEN}🗃️  Building $DATABASE configuration...${NC}"
        ;;
    *)
        echo -e "${RED}❌ Invalid database: $DATABASE${NC}"
        echo "Valid options: sqlite, postgres, mysql"
        exit 1
        ;;
esac

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Copy and customize base configuration
BASE_CONFIG="config/internal/base/common.yaml"
DB_CONFIG="config/internal/databases/${DATABASE}.yaml"
OUTPUT_FILE="$OUTPUT_DIR/sqlc-${DATABASE}.yaml"

echo -e "${YELLOW}📋 Copying base configuration...${NC}"
cp "$BASE_CONFIG" "$OUTPUT_FILE"

# Append database-specific configuration
echo -e "${YELLOW}🔗 Adding $DATABASE-specific configuration...${NC}"
cat >> "$OUTPUT_FILE" << EOF

# === $DATABASE SPECIFIC CONFIGURATION ===
sql:
EOF

# Extract SQL section from database config
yq e '.sql' "$DB_CONFIG" >> "$OUTPUT_FILE"

# Add database-specific validation rules
case "$DATABASE" in
    "postgres")
        yq e '.sql[0].rules // []' "$DB_CONFIG" >> "$OUTPUT_FILE"
        ;;
esac

echo -e "${GREEN}✅ Generated: $OUTPUT_FILE${NC}"

# Show stats
LINES=$(wc -l < "$OUTPUT_FILE")
echo -e "${YELLOW}📊 Configuration: $LINES lines${NC}"

echo -e "${GREEN}💡 Use with: sqlc -f $OUTPUT_FILE generate${NC}"