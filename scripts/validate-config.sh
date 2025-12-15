#!/bin/bash

# Configuration validation script for sqlc templates
# Ensures type safety and configuration correctness

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔍 Validating sqlc configuration...${NC}"

# Check if sqlc is installed
if ! command -v sqlc &> /dev/null; then
    echo -e "${RED}❌ sqlc is not installed. Please install it first.${NC}"
    echo "Visit: https://docs.sqlc.dev/en/stable/overview/install.html"
    exit 1
fi

# Validate main configuration
echo -e "${YELLOW}📋 Checking sqlc.yaml syntax...${NC}"
if sqlc compile 2>/dev/null; then
    echo -e "${GREEN}✅ sqlc.yaml is valid${NC}"
else
    echo -e "${RED}❌ sqlc.yaml has errors${NC}"
    exit 1
fi

# Check for required directories
echo -e "${YELLOW}📁 Checking directory structure...${NC}"
required_dirs=("sql/sqlite/queries" "sql/sqlite/schema" "sql/postgres/queries" "sql/postgres/schema" "sql/mysql/queries" "sql/mysql/schema")
for dir in "${required_dirs[@]}"; do
    if [ ! -d "$dir" ]; then
        echo -e "${YELLOW}⚠️  Directory $dir does not exist (expected for template)${NC}"
    fi
done

# Validate YAML structure
echo -e "${YELLOW}🏗️  Checking YAML structure...${NC}"
if command -v yq &> /dev/null; then
    # Check if version is specified
    version=$(yq e '.version' sqlc.yaml 2>/dev/null || echo "")
    if [ -z "$version" ]; then
        echo -e "${RED}❌ Missing version in sqlc.yaml${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ YAML structure is valid${NC}"
else
    echo -e "${YELLOW}⚠️  yq not found, skipping detailed YAML validation${NC}"
fi

echo -e "${GREEN}✨ Configuration validation complete!${NC}"