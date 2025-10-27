package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GlobalConfig represents the global configuration stored in ~/.lanup/config.yaml
type GlobalConfig struct {
	LogPath       string `yaml:"log_path"`
	LogLevel      string `yaml:"log_level"`
	DefaultPort   int    `yaml:"default_port"`
	CheckInterval int    `yaml:"check_interval"` // seconds for the watcher
}

// ProjectConfig represents the project-specific configuration stored in .lanup.yaml
type ProjectConfig struct {
	Vars       map[string]string `yaml:"vars"`
	Output     string            `yaml:"output"`
	AutoDetect AutoDetectConfig  `yaml:"auto_detect"`
}

// AutoDetectConfig holds settings for automatic service detection
type AutoDetectConfig struct {
	Docker   bool `yaml:"docker"`
	Supabase bool `yaml:"supabase"`
}

// Validate checks if the GlobalConfig has valid values
func (c *GlobalConfig) Validate() error {
	if c.LogPath == "" {
		return fmt.Errorf("log_path cannot be empty")
	}

	// Expand ~ in log path
	if strings.HasPrefix(c.LogPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		c.LogPath = filepath.Join(home, c.LogPath[1:])
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[strings.ToLower(c.LogLevel)] {
		return fmt.Errorf("invalid log_level: %s (must be debug, info, warn, or error)", c.LogLevel)
	}

	if c.DefaultPort < 1 || c.DefaultPort > 65535 {
		return fmt.Errorf("default_port must be between 1 and 65535, got %d", c.DefaultPort)
	}

	if c.CheckInterval < 1 {
		return fmt.Errorf("check_interval must be at least 1 second, got %d", c.CheckInterval)
	}

	return nil
}

// Validate checks if the ProjectConfig has valid values
func (c *ProjectConfig) Validate() error {
	if c.Output == "" {
		return fmt.Errorf("output file path cannot be empty")
	}

	if c.Vars == nil {
		c.Vars = make(map[string]string)
	}

	// Validate that variable keys are not empty
	for key, value := range c.Vars {
		if key == "" {
			return fmt.Errorf("variable key cannot be empty")
		}
		if value == "" {
			return fmt.Errorf("variable %s has empty value", key)
		}
	}

	return nil
}
