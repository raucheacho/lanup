package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnvWriter(t *testing.T) {
	writer := NewEnvWriter(".env.test")

	assert.NotNil(t, writer)
	assert.Equal(t, ".env.test", writer.FilePath)
	assert.True(t, writer.BackupEnabled)
}

func TestEnvWriter_Read(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []EnvVar
	}{
		{
			name:     "empty file",
			content:  "",
			expected: []EnvVar{},
		},
		{
			name: "simple variables",
			content: `API_URL=http://localhost:8000
DATABASE_URL=postgresql://localhost:5432/db`,
			expected: []EnvVar{
				{Key: "API_URL", Value: "http://localhost:8000", Managed: false},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
			},
		},
		{
			name: "managed variables",
			content: `# lanup:managed
API_URL=http://192.168.1.100:8000
# lanup:managed
SUPABASE_URL=http://192.168.1.100:54321`,
			expected: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
				{Key: "SUPABASE_URL", Value: "http://192.168.1.100:54321", Managed: true},
			},
		},
		{
			name: "mixed managed and user variables",
			content: `# lanup:managed
API_URL=http://192.168.1.100:8000

# User variables
DATABASE_URL=postgresql://localhost:5432/db
SECRET_KEY=my-secret`,
			expected: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
				{Key: "SECRET_KEY", Value: "my-secret", Managed: false},
			},
		},
		{
			name: "variables with quotes",
			content: `API_URL="http://localhost:8000"
SECRET_KEY='my-secret'`,
			expected: []EnvVar{
				{Key: "API_URL", Value: "http://localhost:8000", Managed: false},
				{Key: "SECRET_KEY", Value: "my-secret", Managed: false},
			},
		},
		{
			name: "with comments and empty lines",
			content: `# This is a comment
API_URL=http://localhost:8000

# Another comment
DATABASE_URL=postgresql://localhost:5432/db

`,
			expected: []EnvVar{
				{Key: "API_URL", Value: "http://localhost:8000", Managed: false},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			envPath := filepath.Join(tmpDir, ".env")

			// Write test content
			err := os.WriteFile(envPath, []byte(tt.content), 0644)
			require.NoError(t, err)

			writer := NewEnvWriter(envPath)
			vars, err := writer.Read()

			require.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(vars))

			for i, expected := range tt.expected {
				assert.Equal(t, expected.Key, vars[i].Key)
				assert.Equal(t, expected.Value, vars[i].Value)
				assert.Equal(t, expected.Managed, vars[i].Managed)
			}
		})
	}
}

func TestEnvWriter_Read_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	writer := NewEnvWriter(envPath)
	vars, err := writer.Read()

	require.NoError(t, err)
	assert.Empty(t, vars)
}

func TestEnvWriter_Backup(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	backupPath := envPath + ".bak"

	// Create original file
	originalContent := "API_URL=http://localhost:8000\n"
	err := os.WriteFile(envPath, []byte(originalContent), 0644)
	require.NoError(t, err)

	writer := NewEnvWriter(envPath)
	err = writer.Backup()
	require.NoError(t, err)

	// Verify backup was created
	_, err = os.Stat(backupPath)
	require.NoError(t, err)

	// Verify backup content matches original
	backupContent, err := os.ReadFile(backupPath)
	require.NoError(t, err)
	assert.Equal(t, originalContent, string(backupContent))
}

func TestEnvWriter_Backup_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	writer := NewEnvWriter(envPath)
	err := writer.Backup()

	// Should not error when file doesn't exist
	require.NoError(t, err)
}

func TestEnvWriter_Merge(t *testing.T) {
	tests := []struct {
		name     string
		newVars  []EnvVar
		existing []EnvVar
		expected []EnvVar
	}{
		{
			name: "merge with empty existing",
			newVars: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
			},
			existing: []EnvVar{},
			expected: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
			},
		},
		{
			name: "preserve non-managed variables",
			newVars: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
			},
			existing: []EnvVar{
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
				{Key: "SECRET_KEY", Value: "my-secret", Managed: false},
			},
			expected: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
				{Key: "SECRET_KEY", Value: "my-secret", Managed: false},
			},
		},
		{
			name: "replace managed variables",
			newVars: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
			},
			existing: []EnvVar{
				{Key: "API_URL", Value: "http://localhost:8000", Managed: true},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
			},
			expected: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
			},
		},
		{
			name: "complex merge scenario",
			newVars: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
				{Key: "SUPABASE_URL", Value: "http://192.168.1.100:54321", Managed: true},
			},
			existing: []EnvVar{
				{Key: "API_URL", Value: "http://localhost:8000", Managed: true},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
				{Key: "SECRET_KEY", Value: "my-secret", Managed: false},
			},
			expected: []EnvVar{
				{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
				{Key: "SUPABASE_URL", Value: "http://192.168.1.100:54321", Managed: true},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
				{Key: "SECRET_KEY", Value: "my-secret", Managed: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewEnvWriter(".env")
			result := writer.Merge(tt.newVars, tt.existing)

			assert.Equal(t, len(tt.expected), len(result))

			// Create maps for easier comparison
			resultMap := make(map[string]EnvVar)
			for _, v := range result {
				resultMap[v.Key] = v
			}

			for _, expected := range tt.expected {
				actual, exists := resultMap[expected.Key]
				assert.True(t, exists, "Expected key %s not found", expected.Key)
				assert.Equal(t, expected.Value, actual.Value)
				assert.Equal(t, expected.Managed, actual.Managed)
			}
		})
	}
}

func TestEnvWriter_Write(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	vars := []EnvVar{
		{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
		{Key: "SUPABASE_URL", Value: "http://192.168.1.100:54321", Managed: true},
		{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
		{Key: "SECRET_KEY", Value: "my-secret", Managed: false},
	}

	writer := NewEnvWriter(envPath)
	err := writer.Write(vars)
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(envPath)
	require.NoError(t, err)

	// Read and verify content
	content, err := os.ReadFile(envPath)
	require.NoError(t, err)

	contentStr := string(content)

	// Check header
	assert.Contains(t, contentStr, "# Generated by lanup on")
	assert.Contains(t, contentStr, "# Do not edit the managed variables manually")

	// Check managed variables have markers
	assert.Contains(t, contentStr, "# lanup:managed\nAPI_URL=http://192.168.1.100:8000")
	assert.Contains(t, contentStr, "# lanup:managed\nSUPABASE_URL=http://192.168.1.100:54321")

	// Check user variables section
	assert.Contains(t, contentStr, "# User variables (preserved)")
	assert.Contains(t, contentStr, "DATABASE_URL=postgresql://localhost:5432/db")
	assert.Contains(t, contentStr, "SECRET_KEY=my-secret")
}

func TestEnvWriter_Write_WithBackup(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	backupPath := envPath + ".bak"

	// Create original file
	originalContent := "OLD_VAR=old_value\n"
	err := os.WriteFile(envPath, []byte(originalContent), 0644)
	require.NoError(t, err)

	vars := []EnvVar{
		{Key: "NEW_VAR", Value: "new_value", Managed: true},
	}

	writer := NewEnvWriter(envPath)
	writer.BackupEnabled = true
	err = writer.Write(vars)
	require.NoError(t, err)

	// Verify backup was created
	backupContent, err := os.ReadFile(backupPath)
	require.NoError(t, err)
	assert.Equal(t, originalContent, string(backupContent))

	// Verify new content
	newContent, err := os.ReadFile(envPath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), "NEW_VAR=new_value")
	assert.NotContains(t, string(newContent), "OLD_VAR=old_value")
}

func TestTransformURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		newIP    string
		expected string
	}{
		{
			name:     "replace localhost",
			url:      "http://localhost:8000",
			newIP:    "192.168.1.100",
			expected: "http://192.168.1.100:8000",
		},
		{
			name:     "replace 127.0.0.1",
			url:      "http://127.0.0.1:8000",
			newIP:    "192.168.1.100",
			expected: "http://192.168.1.100:8000",
		},
		{
			name:     "replace localhost with https",
			url:      "https://localhost:8443",
			newIP:    "192.168.1.100",
			expected: "https://192.168.1.100:8443",
		},
		{
			name:     "replace localhost without port",
			url:      "http://localhost",
			newIP:    "192.168.1.100",
			expected: "http://192.168.1.100",
		},
		{
			name:     "replace localhost with path",
			url:      "http://localhost:8000/api/v1",
			newIP:    "192.168.1.100",
			expected: "http://192.168.1.100:8000/api/v1",
		},
		{
			name:     "replace multiple occurrences",
			url:      "http://localhost:8000?redirect=http://localhost:3000",
			newIP:    "192.168.1.100",
			expected: "http://192.168.1.100:8000?redirect=http://192.168.1.100:3000",
		},
		{
			name:     "no replacement needed",
			url:      "http://192.168.1.50:8000",
			newIP:    "192.168.1.100",
			expected: "http://192.168.1.50:8000",
		},
		{
			name:     "replace with different private IP",
			url:      "http://localhost:54321",
			newIP:    "10.0.0.5",
			expected: "http://10.0.0.5:54321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformURL(tt.url, tt.newIP)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvWriter_Write_OnlyManagedVars(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	vars := []EnvVar{
		{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
		{Key: "SUPABASE_URL", Value: "http://192.168.1.100:54321", Managed: true},
	}

	writer := NewEnvWriter(envPath)
	err := writer.Write(vars)
	require.NoError(t, err)

	content, err := os.ReadFile(envPath)
	require.NoError(t, err)

	contentStr := string(content)

	// Should not have user variables section
	assert.NotContains(t, contentStr, "# User variables (preserved)")

	// Should have managed variables
	assert.Contains(t, contentStr, "# lanup:managed")
	assert.Contains(t, contentStr, "API_URL=http://192.168.1.100:8000")
}

func TestEnvWriter_Write_OnlyUserVars(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	vars := []EnvVar{
		{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/db", Managed: false},
		{Key: "SECRET_KEY", Value: "my-secret", Managed: false},
	}

	writer := NewEnvWriter(envPath)
	err := writer.Write(vars)
	require.NoError(t, err)

	content, err := os.ReadFile(envPath)
	require.NoError(t, err)

	contentStr := string(content)

	// Should have user variables section
	assert.Contains(t, contentStr, "# User variables (preserved)")

	// Should not have managed markers before user vars
	lines := strings.Split(contentStr, "\n")
	for i, line := range lines {
		if strings.Contains(line, "DATABASE_URL") || strings.Contains(line, "SECRET_KEY") {
			// Check previous line is not managed marker
			if i > 0 {
				assert.NotContains(t, lines[i-1], "# lanup:managed")
			}
		}
	}
}
