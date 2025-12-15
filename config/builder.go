package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ConfigBuilder builds sqlc configurations from components
type ConfigBuilder struct {
	baseDir   string
	outputDir string
}

// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder(baseDir, outputDir string) *ConfigBuilder {
	return &ConfigBuilder{
		baseDir:   baseDir,
		outputDir: outputDir,
	}
}

// BuildConfig builds a complete sqlc configuration for the specified databases
func (cb *ConfigBuilder) BuildConfig(databases []string) error {
	// Load base configuration
	baseConfig, err := cb.loadBaseConfig()
	if err != nil {
		return fmt.Errorf("failed to load base config: %w", err)
	}

	// Build configurations for each database
	var configs []map[string]interface{}
	for _, db := range databases {
		dbConfig, err := cb.buildDatabaseConfig(db, baseConfig)
		if err != nil {
			return fmt.Errorf("failed to build %s config: %w", db, err)
		}
		configs = append(configs, dbConfig)
	}

	// Combine configurations
	finalConfig := map[string]interface{}{
		"version": "2",
		"rules":   baseConfig["rules"],
		"plugins": baseConfig["plugins"],
		"sql":     configs,
	}

	// Write final configuration
	return cb.writeConfig(finalConfig)
}

// loadBaseConfig loads the base configuration
func (cb *ConfigBuilder) loadBaseConfig() (map[string]interface{}, error) {
	basePath := filepath.Join(cb.baseDir, "base", "common.yaml")
	data, err := os.ReadFile(basePath)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

// buildDatabaseConfig builds configuration for a specific database
func (cb *ConfigBuilder) buildDatabaseConfig(db string, baseConfig map[string]interface{}) (map[string]interface{}, error) {
	// Load database-specific configuration
	dbPath := filepath.Join(cb.baseDir, "databases", fmt.Sprintf("%s.yaml", db))
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, err
	}

	var dbConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &dbConfig); err != nil {
		return nil, err
	}

	// Extract the SQL configuration
	sqlConfigs, ok := dbConfig["sql"].([]interface{})
	if !ok || len(sqlConfigs) == 0 {
		return nil, fmt.Errorf("no sql configuration found for %s", db)
	}

	// Return the first SQL configuration
	sqlConfig, ok := sqlConfigs[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid sql configuration format for %s", db)
	}

	return sqlConfig, nil
}

// writeConfig writes the final configuration to file
func (cb *ConfigBuilder) writeConfig(config map[string]interface{}) error {
	// Ensure output directory exists
	if err := os.MkdirAll(cb.outputDir, 0o755); err != nil {
		return err
	}

	// Write configuration
	outputPath := filepath.Join(cb.outputDir, "sqlc.yaml")
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0o644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run builder.go <databases...>")
		fmt.Println("Example: go run builder.go sqlite postgres mysql")
		os.Exit(1)
	}

	// Parse databases from command line
	databases := strings.Split(os.Args[1], ",")
	for i, db := range databases {
		databases[i] = strings.TrimSpace(db)
	}

	// Build configuration
	builder := NewConfigBuilder("internal", ".")
	if err := builder.BuildConfig(databases); err != nil {
		fmt.Printf("Error building configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Configuration built successfully!")
	fmt.Printf("Generated for databases: %v\n", databases)
}
