package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  GlobalConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: GlobalConfig{
				LogPath:       "/tmp/lanup.log",
				LogLevel:      "info",
				DefaultPort:   8080,
				CheckInterval: 5,
			},
			wantErr: false,
		},
		{
			name: "empty log path",
			config: GlobalConfig{
				LogPath:       "",
				LogLevel:      "info",
				DefaultPort:   8080,
				CheckInterval: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: GlobalConfig{
				LogPath:       "/tmp/lanup.log",
				LogLevel:      "invalid",
				DefaultPort:   8080,
				CheckInterval: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid port - too low",
			config: GlobalConfig{
				LogPath:       "/tmp/lanup.log",
				LogLevel:      "info",
				DefaultPort:   0,
				CheckInterval: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: GlobalConfig{
				LogPath:       "/tmp/lanup.log",
				LogLevel:      "info",
				DefaultPort:   65536,
				CheckInterval: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid check interval",
			config: GlobalConfig{
				LogPath:       "/tmp/lanup.log",
				LogLevel:      "info",
				DefaultPort:   8080,
				CheckInterval: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ProjectConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ProjectConfig{
				Vars: map[string]string{
					"API_URL": "http://localhost:8000",
				},
				Output: ".env.local",
				AutoDetect: AutoDetectConfig{
					Docker:   true,
					Supabase: true,
				},
			},
			wantErr: false,
		},
		{
			name: "empty output",
			config: ProjectConfig{
				Vars: map[string]string{
					"API_URL": "http://localhost:8000",
				},
				Output: "",
			},
			wantErr: true,
		},
		{
			name: "empty variable value",
			config: ProjectConfig{
				Vars: map[string]string{
					"API_URL": "",
				},
				Output: ".env.local",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetDefaultGlobalConfig(t *testing.T) {
	config := GetDefaultGlobalConfig()

	assert.NotNil(t, config)
	assert.NotEmpty(t, config.LogPath)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, 8080, config.DefaultPort)
	assert.Equal(t, 5, config.CheckInterval)

	// Validate the default config
	err := config.Validate()
	assert.NoError(t, err)
}

func TestGetDefaultProjectConfig(t *testing.T) {
	config := GetDefaultProjectConfig()

	assert.NotNil(t, config)
	assert.NotEmpty(t, config.Vars)
	assert.Equal(t, ".env.local", config.Output)
	assert.True(t, config.AutoDetect.Docker)
	assert.True(t, config.AutoDetect.Supabase)

	// Validate the default config
	err := config.Validate()
	assert.NoError(t, err)
}

func TestSaveAndLoadProjectConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".lanup.yaml")

	// Create a test config
	testConfig := &ProjectConfig{
		Vars: map[string]string{
			"API_URL":      "http://localhost:8000",
			"DATABASE_URL": "postgresql://localhost:5432/test",
		},
		Output: ".env.test",
		AutoDetect: AutoDetectConfig{
			Docker:   false,
			Supabase: true,
		},
	}

	// Save the config
	err := SaveProjectConfig(configPath, testConfig)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(configPath)
	require.NoError(t, err)

	// Load the config back
	loadedConfig, err := LoadProjectConfig(configPath)
	require.NoError(t, err)

	// Verify the loaded config matches
	assert.Equal(t, testConfig.Vars, loadedConfig.Vars)
	assert.Equal(t, testConfig.Output, loadedConfig.Output)
	assert.Equal(t, testConfig.AutoDetect.Docker, loadedConfig.AutoDetect.Docker)
	assert.Equal(t, testConfig.AutoDetect.Supabase, loadedConfig.AutoDetect.Supabase)
}

func TestLoadProjectConfig_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent.yaml")

	_, err := LoadProjectConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestLoadGlobalConfig_FirstRun(t *testing.T) {
	// Save original HOME and restore after test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create a temporary directory to act as HOME
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)

	// Verify ~/.lanup doesn't exist yet
	lanupDir := filepath.Join(tmpHome, ".lanup")
	_, err := os.Stat(lanupDir)
	assert.True(t, os.IsNotExist(err), "~/.lanup should not exist before first run")

	// Load global config (should trigger first-run setup)
	config, err := LoadGlobalConfig()
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Verify ~/.lanup directory was created with correct permissions
	info, err := os.Stat(lanupDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
	assert.Equal(t, os.FileMode(0755), info.Mode().Perm())

	// Verify ~/.lanup/logs directory was created with correct permissions
	logsDir := filepath.Join(lanupDir, "logs")
	info, err = os.Stat(logsDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
	assert.Equal(t, os.FileMode(0755), info.Mode().Perm())

	// Verify config.yaml was created with correct permissions
	configPath := filepath.Join(lanupDir, "config.yaml")
	info, err = os.Stat(configPath)
	require.NoError(t, err)
	assert.False(t, info.IsDir())
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	// Verify config has default values
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, 8080, config.DefaultPort)
	assert.Equal(t, 5, config.CheckInterval)
	assert.Contains(t, config.LogPath, ".lanup/logs/lanup.log")

	// Verify config is valid
	err = config.Validate()
	assert.NoError(t, err)
}

func TestLoadGlobalConfig_ExistingConfig(t *testing.T) {
	// Save original HOME and restore after test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create a temporary directory to act as HOME
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)

	// Create ~/.lanup directory structure manually
	lanupDir := filepath.Join(tmpHome, ".lanup")
	err := os.MkdirAll(lanupDir, 0755)
	require.NoError(t, err)

	// Create a custom config file
	configPath := filepath.Join(lanupDir, "config.yaml")
	customConfig := `log_path: /custom/path/lanup.log
log_level: debug
default_port: 9000
check_interval: 10
`
	err = os.WriteFile(configPath, []byte(customConfig), 0600)
	require.NoError(t, err)

	// Load global config (should read existing config)
	config, err := LoadGlobalConfig()
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Verify it loaded the custom values
	assert.Equal(t, "/custom/path/lanup.log", config.LogPath)
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, 9000, config.DefaultPort)
	assert.Equal(t, 10, config.CheckInterval)
}

func TestLoadProjectConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".lanup.yaml")

	// Write invalid YAML (unclosed bracket)
	invalidYAML := `vars:
  API_URL: http://localhost:8000
  [unclosed bracket
output: .env.local
auto_detect:
  docker: true
`
	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	_, err = LoadProjectConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestLoadGlobalConfig_InvalidYAML(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)

	lanupDir := filepath.Join(tmpHome, ".lanup")
	err := os.MkdirAll(lanupDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(lanupDir, "config.yaml")
	// Write invalid YAML (invalid structure with unclosed quotes)
	invalidYAML := `log_path: "/tmp/lanup.log
log_level: "info
default_port: 8080
`
	err = os.WriteFile(configPath, []byte(invalidYAML), 0600)
	require.NoError(t, err)

	_, err = LoadGlobalConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestSaveProjectConfig_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".lanup.yaml")

	// Create invalid config (empty output)
	invalidConfig := &ProjectConfig{
		Vars: map[string]string{
			"API_URL": "http://localhost:8000",
		},
		Output: "", // Invalid: empty output
	}

	err := SaveProjectConfig(configPath, invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid configuration")
}

func TestProjectConfig_Validate_EmptyKey(t *testing.T) {
	config := &ProjectConfig{
		Vars: map[string]string{
			"": "some-value", // Invalid: empty key
		},
		Output: ".env.local",
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key cannot be empty")
}

func TestGlobalConfig_Validate_TildeExpansion(t *testing.T) {
	config := &GlobalConfig{
		LogPath:       "~/.lanup/logs/lanup.log",
		LogLevel:      "info",
		DefaultPort:   8080,
		CheckInterval: 5,
	}

	err := config.Validate()
	assert.NoError(t, err)

	// Verify tilde was expanded
	assert.NotContains(t, config.LogPath, "~")
	assert.Contains(t, config.LogPath, ".lanup/logs/lanup.log")
}

func TestLoadProjectConfig_EmptyPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Create .lanup.yaml in current directory
	testConfig := &ProjectConfig{
		Vars: map[string]string{
			"API_URL": "http://localhost:8000",
		},
		Output: ".env.local",
		AutoDetect: AutoDetectConfig{
			Docker:   true,
			Supabase: false,
		},
	}

	err := SaveProjectConfig("", testConfig)
	require.NoError(t, err)

	// Load with empty path (should default to .lanup.yaml)
	loadedConfig, err := LoadProjectConfig("")
	require.NoError(t, err)
	assert.Equal(t, testConfig.Vars, loadedConfig.Vars)
	assert.Equal(t, testConfig.Output, loadedConfig.Output)
}
