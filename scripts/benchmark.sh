#!/bin/bash

# Performance benchmarking script for sqlc operations
# Measures code generation time and database query performance

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Default values
WARMUP_RUNS=3
BENCHMARK_RUNS=10
DATABASES="sqlite,postgres,mysql"
OUTPUT_DIR="benchmark-results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --warmup)
            WARMUP_RUNS="$2"
            shift 2
            ;;
        --runs)
            BENCHMARK_RUNS="$2"
            shift 2
            ;;
        --databases)
            DATABASES="$2"
            shift 2
            ;;
        --output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  --warmup N     Number of warmup runs (default: 3)"
            echo "  --runs N        Number of benchmark runs (default: 10)"
            echo "  --databases    Comma-separated list of databases (default: sqlite,postgres,mysql)"
            echo "  --output DIR    Output directory for results (default: benchmark-results)"
            echo "  --help          Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🚀 sqlc Performance Benchmarking${NC}"
echo -e "${YELLOW}📅 Timestamp: $TIMESTAMP${NC}"
echo -e "${YELLOW}🔥 Warmup runs: $WARMUP_RUNS${NC}"
echo -e "${YELLOW}📊 Benchmark runs: $BENCHMARK_RUNS${NC}"
echo -e "${YELLOW}🗃️  Databases: $DATABASES${NC}"

# Create output directory
OUTPUT_PATH="$OUTPUT_DIR/$TIMESTAMP"
mkdir -p "$OUTPUT_PATH"

# Dependencies check
echo -e "${CYAN}🔍 Checking dependencies...${NC}"

if ! command -v hyperfine &> /dev/null; then
    echo -e "${RED}❌ hyperfine not found. Install with: brew install hyperfine${NC}"
    exit 1
fi

if ! command -v sqlc &> /dev/null; then
    echo -e "${RED}❌ sqlc not found. Please install sqlc${NC}"
    exit 1
fi

if ! command -v yq &> /dev/null; then
    echo -e "${RED}❌ yq not found. Install with: brew install yq${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All dependencies found${NC}"

# Setup test data
echo -e "${CYAN}📋 Setting up test data...${NC}"
make examples > /dev/null 2>&1

# Global results
GLOBAL_RESULTS="$OUTPUT_PATH/global_results.json"
echo "{" > "$GLOBAL_RESULTS"
echo "\"timestamp\": \"$TIMESTAMP\"," >> "$GLOBAL_RESULTS"
echo "\"warmup_runs\": $WARMUP_RUNS," >> "$GLOBAL_RESULTS"
echo "\"benchmark_runs\": $BENCHMARK_RUNS," >> "$GLOBAL_RESULTS"
echo "\"databases\": [$(echo $DATABASES | sed 's/,/","/g')]," >> "$GLOBAL_RESULTS"
echo "\"results\": {" >> "$GLOBAL_RESULTS"

# Benchmark each database
IFS=',' read -ra DB_ARRAY <<< "$DATABASES"
for db in "${DB_ARRAY[@]}"; do
    db=$(echo "$db" | xargs) # Trim whitespace
    echo -e "\n${PURPLE}🗃️  Benchmarking $db...${NC}"
    
    DB_RESULTS="$OUTPUT_PATH/${db}_results.json"
    DB_CONFIG="config/modular/sqlc-${db}.yaml"
    
    if [ ! -f "$DB_CONFIG" ]; then
        echo -e "${RED}❌ Configuration not found: $DB_CONFIG${NC}"
        continue
    fi
    
    # Warmup runs
    echo -e "${YELLOW}🔥 Warming up ($WARMUP_RUNS runs)...${NC}"
    for ((i=1; i<=$WARMUP_RUNS; i++)); do
        echo -n "."
        sqlc -f "$DB_CONFIG" generate > /dev/null 2>&1
    done
    echo ""
    
    # Benchmark runs
    echo -e "${CYAN}📊 Running benchmark ($BENCHMARK_RUNS runs)...${NC}"
    
    hyperfine \
        --warmup $WARMUP_RUNS \
        --runs $BENCHMARK_RUNS \
        --shell none \
        --export-json "$DB_RESULTS.json" \
        --output none \
        "sqlc -f $DB_CONFIG generate" \
        --command-name "sqlc-$db-generate"
    
    # Extract key metrics
    if [ -f "$DB_RESULTS.json" ]; then
        MEAN_TIME=$(yq e '.results[0].mean' "$DB_RESULTS.json" 2>/dev/null || echo "0")
        MIN_TIME=$(yq e '.results[0].min' "$DB_RESULTS.json" 2>/dev/null || echo "0")
        MAX_TIME=$(yq e '.results[0].max' "$DB_RESULTS.json" 2>/dev/null || echo "0")
        STDDEV=$(yq e '.results[0].stddev' "$DB_RESULTS.json" 2>/dev/null || echo "0")
        
        echo -e "${GREEN}📈 $db Results:${NC}"
        echo -e "  ⏱️  Mean: $(printf "%.3f" $MEAN_TIME)s"
        echo -e "  ⚡ Min: $(printf "%.3f" $MIN_TIME)s"
        echo -e "  🐌 Max: $(printf "%.3f" $MAX_TIME)s"
        echo -e "  📊 StdDev: $(printf "%.3f" $STDDEV)s"
        
        # Add to global results
        echo "\"$db\": {" >> "$GLOBAL_RESULTS"
        echo "\"mean_seconds\": $MEAN_TIME," >> "$GLOBAL_RESULTS"
        echo "\"min_seconds\": $MIN_TIME," >> "$GLOBAL_RESULTS"
        echo "\"max_seconds\": $MAX_TIME," >> "$GLOBAL_RESULTS"
        echo "\"stddev_seconds\": $STDDEV" >> "$GLOBAL_RESULTS"
        echo "}," >> "$GLOBAL_RESULTS"
        
        # Copy detailed results
        mv "$DB_RESULTS.json" "$OUTPUT_PATH/${db}_hyperfine.json"
    else
        echo -e "${RED}❌ Failed to get results for $db${NC}"
    fi
done

# Close global results
echo "\"}}" >> "$GLOBAL_RESULTS"

# Generate summary report
echo -e "\n${BLUE}📋 Generating summary report...${NC}"
cat > "$OUTPUT_PATH/benchmark_report.md" << EOF
# sqlc Performance Benchmark Report

**Timestamp:** $TIMESTAMP  
**Warmup Runs:** $WARMUP_RUNS  
**Benchmark Runs:** $BENCHMARK_RUNS  
**Databases Tested:** $(echo $DATABASES | tr ',' ' ')

## Results Summary

| Database | Mean (s) | Min (s) | Max (s) | StdDev (s) |
|----------|------------|-----------|-----------|-------------|
EOF

# Extract results for table
if [ -f "$GLOBAL_RESULTS" ]; then
    for db in "${DB_ARRAY[@]}"; do
        db=$(echo "$db" | xargs)
        MEAN_TIME=$(yq e ".results.$db.mean_seconds" "$GLOBAL_RESULTS" 2>/dev/null || echo "N/A")
        MIN_TIME=$(yq e ".results.$db.min_seconds" "$GLOBAL_RESULTS" 2>/dev/null || echo "N/A")
        MAX_TIME=$(yq e ".results.$db.max_seconds" "$GLOBAL_RESULTS" 2>/dev/null || echo "N/A")
        STDDEV=$(yq e ".results.$db.stddev_seconds" "$GLOBAL_RESULTS" 2>/dev/null || echo "N/A")
        
        echo "| $db | $MEAN_TIME | $MIN_TIME | $MAX_TIME | $STDDEV |" >> "$OUTPUT_PATH/benchmark_report.md"
    done
fi

cat >> "$OUTPUT_PATH/benchmark_report.md" << EOF

## Environment

- **sqlc Version:** $(sqlc version 2>/dev/null || echo "Unknown")
- **Go Version:** $(go version 2>/dev/null || echo "Unknown")
- **OS:** $(uname -s) $(uname -r)
- **Architecture:** $(uname -m)

## Files Generated

EOF

# List generated files
find "$OUTPUT_PATH" -name "*.json" -o -name "*.md" | while read file; do
    filename=$(basename "$file")
    echo "- \`${filename}\`" >> "$OUTPUT_PATH/benchmark_report.md"
done

cat >> "$OUTPUT_PATH/benchmark_report.md" << EOF

## Performance Analysis

### Fastest Generation
$(echo "$DATABASES" | tr ',' '\n' | while read db; do
    db=$(echo "$db" | xargs)
    MEAN_TIME=$(yq e ".results.$db.mean_seconds" "$GLOBAL_RESULTS" 2>/dev/null || echo "999")
    echo "$MEAN_TIME $db"
done | sort -n | head -1 | cut -d' ' -f2-)

### Slowest Generation
$(echo "$DATABASES" | tr ',' '\n' | while read db; do
    db=$(echo "$db" | xargs)
    MEAN_TIME=$(yq e ".results.$db.mean_seconds" "$GLOBAL_RESULTS" 2>/dev/null || echo "0")
    echo "$MEAN_TIME $db"
done | sort -n | tail -1 | cut -d' ' -f2-)

### Most Consistent
$(echo "$DATABASES" | tr ',' '\n' | while read db; do
    db=$(echo "$db" | xargs)
    STDDEV=$(yq e ".results.$db.stddev_seconds" "$GLOBAL_RESULTS" 2>/dev/null || echo "999")
    echo "$STDDEV $db"
done | sort -n | head -1 | cut -d' ' -f2-)

## Recommendations

1. **Configuration Size:** Monitor the size of each configuration file (smaller = faster)
2. **Database-Specific Optimizations:** Consider database-specific type overrides
3. **Code Generation Patterns:** Optimize SQL schemas for faster code generation
4. **Parallel Generation:** For large projects, consider generating code in parallel

---

*This report was generated automatically on $(date)*
EOF

# Generate performance graph
echo -e "${CYAN}📈 Generating performance graph...${NC}"
if command -v python3 &> /dev/null; then
    python3 -c "
import json
import matplotlib.pyplot as plt
import sys

# Load results
try:
    with open('$GLOBAL_RESULTS', 'r') as f:
        data = json.load(f)
except Exception as e:
    print(f'Error loading results: {e}')
    sys.exit(1)

# Extract data for plotting
databases = []
mean_times = []

if 'results' in data:
    for db, metrics in data['results'].items():
        if 'mean_seconds' in metrics:
            databases.append(db)
            mean_times.append(metrics['mean_seconds'])

# Create bar chart
if databases:
    plt.figure(figsize=(10, 6))
    bars = plt.bar(databases, mean_times, color=['#FF6B6B', '#4ECDC4', '#45B7D1'])
    
    plt.title('sqlc Code Generation Performance', fontsize=16, fontweight='bold')
    plt.xlabel('Database', fontsize=12)
    plt.ylabel('Mean Generation Time (seconds)', fontsize=12)
    plt.grid(axis='y', alpha=0.3)
    
    # Add value labels on bars
    for bar, time in zip(bars, mean_times):
        plt.text(bar.get_x() + bar.get_width()/2, bar.get_height() + 0.01,
                f'{time:.3f}s', ha='center', va='bottom', fontweight='bold')
    
    plt.tight_layout()
    plt.savefig('$OUTPUT_PATH/performance_graph.png', dpi=300, bbox_inches='tight')
    print('✅ Performance graph generated: performance_graph.png')
else:
    print('⚠️  No data available for plotting')
"
else
    echo -e "${YELLOW}⚠️  Python3 not available, skipping graph generation${NC}"
fi

# Cleanup
echo -e "${CYAN}🧹 Cleaning up generated code...${NC}"
find internal/db -name "*.go" -delete 2>/dev/null || true

echo -e "\n${GREEN}✨ Benchmark complete!${NC}"
echo -e "${BLUE}📁 Results saved to: $OUTPUT_PATH${NC}"
echo -e "${YELLOW}📄 Report: $OUTPUT_PATH/benchmark_report.md${NC}"
echo -e "${YELLOW}📊 Raw data: $GLOBAL_RESULTS${NC}"

# Show summary
echo -e "\n${PURPLE}📊 Performance Summary:${NC}"
if [ -f "$GLOBAL_RESULTS" ]; then
    for db in "${DB_ARRAY[@]}"; do
        db=$(echo "$db" | xargs)
        MEAN_TIME=$(yq e ".results.$db.mean_seconds" "$GLOBAL_RESULTS" 2>/dev/null || echo "N/A")
        echo -e "  🗃️  ${db}: $(printf "%-8s" "$MEAN_TIME")s"
    done
fi

echo -e "\n${CYAN}💡 Recommendations:${NC}"
echo -e "  📋 View the full report: $OUTPUT_PATH/benchmark_report.md"
echo -e "  📈 Check performance trends by comparing with previous runs"
echo -e "  🚀 Optimize configurations based on fastest database patterns"