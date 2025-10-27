package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/raucheacho/lanup/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitCmd_Run_Success(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create and run init command
	initCmd := &InitCmd{
		Format: "yaml",
		Force:  false,
	}

	err = initCmd.Run()
	require.NoError(t, err)

	// Verify file was created
	configPath := filepath.Join(tmpDir, ".lanup.yaml")
	_, err = os.Stat(configPath)
	require.NoError(t, err, "Config file should exist")

	// Load and verify content
	loadedConfig, err := config.LoadProjectConfig(configPath)
	require.NoError(t, err)

	// Verify default values
	assert.NotEmpty(t, loadedConfig.Vars)
	assert.Equal(t, ".env.local", loadedConfig.Output)
	assert.True(t, loadedConfig.AutoDetect.Docker)
	assert.True(t, loadedConfig.AutoDetect.Supabase)

	// Verify config is valid
	err = loadedConfig.Validate()
	assert.NoError(t, err)
}

func TestInitCmd_Run_FileExists_NoForce(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create existing config file
	configPath := filepath.Join(tmpDir, ".lanup.yaml")
	existingContent := "vars:\n  EXISTING: value\noutput: .env.test\n"
	err = os.WriteFile(configPath, []byte(existingContent), 0644)
	require.NoError(t, err)

	// Create and run init command without force
	initCmd := &InitCmd{
		Format: "yaml",
		Force:  false,
	}

	err = initCmd.Run()
	assert.Error(t, err, "Should error when file exists without --force")
	assert.Contains(t, err.Error(), "already exists")

	// Verify original file was not modified
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Equal(t, existingContent, string(content))
}

func TestInitCmd_Run_FileExists_WithForce(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create existing config file
	configPath := filepath.Join(tmpDir, ".lanup.yaml")
	existingContent := "vars:\n  EXISTING: value\noutput: .env.test\n"
	err = os.WriteFile(configPath, []byte(existingContent), 0644)
	require.NoError(t, err)

	// Create and run init command with force
	initCmd := &InitCmd{
		Format: "yaml",
		Force:  true,
	}

	err = initCmd.Run()
	require.NoError(t, err)

	// Verify file was overwritten with default config
	loadedConfig, err := config.LoadProjectConfig(configPath)
	require.NoError(t, err)

	// Should have default values, not the existing ones
	assert.Equal(t, ".env.local", loadedConfig.Output)
	assert.True(t, loadedConfig.AutoDetect.Docker)
	assert.True(t, loadedConfig.AutoDetect.Supabase)
}

func TestInitCmd_Run_InvalidFormat(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create init command with invalid format
	initCmd := &InitCmd{
		Format: "json",
		Force:  false,
	}

	err = initCmd.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported format")
}

func TestInitCmd_Run_TOMLNotSupported(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create init command with TOML format
	initCmd := &InitCmd{
		Format: "toml",
		Force:  false,
	}

	err = initCmd.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TOML format is not yet supported")
}
