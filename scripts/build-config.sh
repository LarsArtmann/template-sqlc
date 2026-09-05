#!/bin/bash

# Configuration build script
# Builds sqlc.yaml from modular components

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default databases
DEFAULT_DATABASES="sqlite,postgres,mysql"

# Parse command line arguments
DATABASES="${1:-$DEFAULT_DATABASES}"
OUTPUT_DIR="${2:-.}"

echo -e "${BLUE}🏗️  Building sqlc configuration...${NC}"
echo -e "${YELLOW}📋 Databases: ${DATABASES}${NC}"
echo -e "${YELLOW}📁 Output: ${OUTPUT_DIR}${NC}"

# Check if config directory exists
if [ ! -d "config/internal" ]; then
	echo -e "${RED}❌ config/internal directory not found${NC}"
	exit 1
fi

# Build configuration
cd config
go run builder.go "$DATABASES"
cd ..

# Validate generated configuration
echo -e "${YELLOW}🔍 Validating generated configuration...${NC}"
if command -v yq &>/dev/null; then
	if yq eval . sqlc.yaml >/dev/null 2>&1; then
		echo -e "${GREEN}✅ Configuration is valid YAML${NC}"
	else
		echo -e "${RED}❌ Generated configuration is invalid YAML${NC}"
		exit 1
	fi
else
	echo -e "${YELLOW}⚠️  yq not found, skipping YAML validation${NC}"
fi

# Show configuration stats
echo -e "${BLUE}📊 Configuration statistics:${NC}"
if command -v yq &>/dev/null; then
	DATABASE_COUNT=$(yq e '.sql | length' sqlc.yaml)
	RULE_COUNT=$(yq e '.rules | length' sqlc.yaml)
	PLUGIN_COUNT=$(yq e '.plugins | length' sqlc.yaml)
	LINES=$(wc -l <sqlc.yaml)

	echo -e "  🗃️  Databases: ${DATABASE_COUNT}"
	echo -e "  📏 Rules: ${RULE_COUNT}"
	echo -e "  🔌 Plugins: ${PLUGIN_COUNT}"
	echo -e "  📄 Total lines: ${LINES}"
else
	LINES=$(wc -l <sqlc.yaml)
	echo -e "  📄 Total lines: ${LINES}"
fi

echo -e "${GREEN}✨ Configuration build complete!${NC}"
echo -e "${YELLOW}💡 Run 'sqlc compile' to validate the configuration${NC}"
