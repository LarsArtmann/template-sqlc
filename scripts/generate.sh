#!/bin/bash
set -e

# Run sqlc generate
sqlc generate

# Deduplicate by using shared package for Queries type
# This removes duplicate db.go files and uses shared.Queries directly

SQLITE_DB="internal/db/sqlite/db.go"
MYSQL_DB="internal/db/mysql/db.go"

# Remove duplicate db.go files - they will use shared.Queries directly
rm -f "$SQLITE_DB" "$MYSQL_DB"

# Update sqlite/querier.go to use shared.Queries
SQLITE_QUERIER="internal/db/sqlite/querier.go"
if [ -f "$SQLITE_QUERIER" ]; then
	# Add import for shared package and change compile check to use shared.Queries
	sed -i '' 's/import (/import (\n\t"github.com\/LarsArtmann\/template-sqlc\/internal\/db\/shared"/' "$SQLITE_QUERIER"
	sed -i '' 's/var _ Querier = (\*Queries)(nil)/var _ Querier = (*shared.Queries)(nil)/' "$SQLITE_QUERIER"
fi

# Update mysql/querier.go to use shared.Queries
MYSQL_QUERIER="internal/db/mysql/querier.go"
if [ -f "$MYSQL_QUERIER" ]; then
	# Add import for shared package and change compile check to use shared.Queries
	sed -i '' 's/import (/import (\n\t"github.com\/LarsArtmann\/template-sqlc\/internal\/db\/shared"/' "$MYSQL_QUERIER"
	sed -i '' 's/var _ Querier = (\*Queries)(nil)/var _ Querier = (*shared.Queries)(nil)/' "$MYSQL_QUERIER"
fi

# Fix MySQL duplicate json imports (sqlc bug)
for file in internal/db/mysql/models.go internal/db/mysql/user.sql.go; do
	if [ -f "$file" ]; then
		grep -v '^\s*"json"$' "$file" > "$file.tmp" || true
		mv "$file.tmp" "$file"
	fi
done

echo "Code generation and deduplication complete"
