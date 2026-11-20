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
	CheckInterval int    `yaml:"check_interval"` // seconds for the watcher
}

// ProjectConfig represents the project-specific configuration stored in .lanup.yaml
type ProjectConfig struct {
	Vars       map[string]string `yaml:"vars"`
	Output     string            `yaml:"output,omitempty"`
	Outputs    []OutputConfig    `yaml:"outputs,omitempty"`
	AutoDetect AutoDetectConfig  `yaml:"auto_detect"`
}

// OutputConfig is one env file updated by lanup, with the variables it receives
// on top of the shared ones. Paths are relative to the directory lanup runs in.
type OutputConfig struct {
	Path string            `yaml:"path"`
	Vars map[string]string `yaml:"vars"`
}

// Target is a resolved env file to write: its path, every variable it gets,
// and whether it is updated in place (outputs) or regenerated (legacy output).
type Target struct {
	Path    string
	Vars    map[string]string
	InPlace bool
}

// Targets lists the env files to write. The legacy `output` keeps its
// regenerate behavior; each entry of `outputs` gets the shared vars plus its own
// (its own win) and is updated in place, leaving every other line untouched.
func (c *ProjectConfig) Targets() []Target {
	var targets []Target
	if c.Output != "" {
		targets = append(targets, Target{Path: c.Output, Vars: copyVars(c.Vars)})
	}
	for _, o := range c.Outputs {
		vars := copyVars(c.Vars)
		for k, v := range o.Vars {
			vars[k] = v
		}
		targets = append(targets, Target{Path: o.Path, Vars: vars, InPlace: true})
	}
	return targets
}

func copyVars(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
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

	if c.CheckInterval < 1 {
		return fmt.Errorf("check_interval must be at least 1 second, got %d", c.CheckInterval)
	}

	return nil
}

// Validate checks if the ProjectConfig has valid values
func (c *ProjectConfig) Validate() error {
	if c.Output == "" && len(c.Outputs) == 0 {
		return fmt.Errorf("output file path cannot be empty (set output or outputs)")
	}

	if c.Vars == nil {
		c.Vars = make(map[string]string)
	}
	if err := validateVars(c.Vars); err != nil {
		return err
	}

	seen := make(map[string]bool)
	if c.Output != "" {
		seen[filepath.Clean(c.Output)] = true
	}
	for i, o := range c.Outputs {
		if o.Path == "" {
			return fmt.Errorf("outputs[%d]: path cannot be empty", i)
		}
		clean := filepath.Clean(o.Path)
		if seen[clean] {
			return fmt.Errorf("outputs[%d]: %s is listed twice", i, o.Path)
		}
		seen[clean] = true
		if len(o.Vars) == 0 && len(c.Vars) == 0 {
			return fmt.Errorf("outputs[%d] (%s): no variables to write", i, o.Path)
		}
		if err := validateVars(o.Vars); err != nil {
			return fmt.Errorf("outputs[%d] (%s): %w", i, o.Path, err)
		}
	}

	return nil
}

// validateVars checks that no variable key or value is empty
func validateVars(vars map[string]string) error {
	for key, value := range vars {
		if key == "" {
			return fmt.Errorf("variable key cannot be empty")
		}
		if value == "" {
			return fmt.Errorf("variable %s has empty value", key)
		}
	}
	return nil
}
