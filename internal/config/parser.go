package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadGlobalConfig reads the global configuration from ~/.lanup/config.yaml
func LoadGlobalConfig() (*GlobalConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(home, ".lanup", "config.yaml")

	// If config doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := GetDefaultGlobalConfig()
		if err := ensureGlobalConfigDir(); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}
		if err := saveGlobalConfig(configPath, defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config GlobalConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// LoadProjectConfig reads the project configuration from .lanup.yaml in the current directory
func LoadProjectConfig(path string) (*ProjectConfig, error) {
	if path == "" {
		path = ".lanup.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("project config file not found: %s (run 'lanup init' to create one)", path)
		}
		return nil, fmt.Errorf("failed to read project config: %w", err)
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse project config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid project configuration: %w", err)
	}

	return &config, nil
}

// SaveProjectConfig writes the project configuration to a file in YAML format
func SaveProjectConfig(path string, config *ProjectConfig) error {
	if path == "" {
		path = ".lanup.yaml"
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultGlobalConfig returns a GlobalConfig with default values
func GetDefaultGlobalConfig() *GlobalConfig {
	home, _ := os.UserHomeDir()
	logPath := filepath.Join(home, ".lanup", "logs", "lanup.log")

	return &GlobalConfig{
		LogPath:       logPath,
		LogLevel:      "info",
		DefaultPort:   8080,
		CheckInterval: 5,
	}
}

// GetDefaultProjectConfig returns a ProjectConfig with default values
func GetDefaultProjectConfig() *ProjectConfig {
	return &ProjectConfig{
		Vars: map[string]string{
			"SUPABASE_URL":      "http://localhost:54321",
			"SUPABASE_ANON_KEY": "your-anon-key",
			"API_URL":           "http://localhost:8000",
			"DASHBOARD_URL":     "http://localhost:3000",
		},
		Output: ".env.local",
		AutoDetect: AutoDetectConfig{
			Docker:   true,
			Supabase: true,
		},
	}
}

// ensureGlobalConfigDir creates the ~/.lanup directory structure if it doesn't exist
func ensureGlobalConfigDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".lanup")
	logsDir := filepath.Join(configDir, "logs")

	// Create config directory with 0755 permissions
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create logs directory with 0755 permissions
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	return nil
}

// saveGlobalConfig writes the global configuration to a file
func saveGlobalConfig(path string, config *GlobalConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with 0600 permissions for security
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
