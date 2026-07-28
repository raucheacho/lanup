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

// --- .envrc (direnv) support tests ---

func TestEnvWriter_IsEnvRC(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{".envrc", true},
		{"/home/user/project/.envrc", true},
		{".env", false},
		{".env.local", false},
		{".env.production", false},
	}

	for _, tt := range tests {
		w := NewEnvWriter(tt.path)
		assert.Equal(t, tt.expected, w.isEnvRC(), "isEnvRC() for %q", tt.path)
	}
}

func TestEnvWriter_Write_EnvRC_Format(t *testing.T) {
	tmpDir := t.TempDir()
	envrcPath := filepath.Join(tmpDir, ".envrc")

	vars := []EnvVar{
		{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
		{Key: "DB_URL", Value: "postgresql://192.168.1.100:5432/db", Managed: true},
	}

	writer := NewEnvWriter(envrcPath)
	err := writer.Write(vars)
	require.NoError(t, err)

	content, err := os.ReadFile(envrcPath)
	require.NoError(t, err)

	contentStr := string(content)

	// Should use export syntax
	assert.Contains(t, contentStr, "export API_URL=http://192.168.1.100:8000")
	assert.Contains(t, contentStr, "export DB_URL=postgresql://192.168.1.100:5432/db")

	// Should NOT have plain KEY=VALUE (without export)
	lines := strings.Split(contentStr, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		assert.True(t, strings.HasPrefix(trimmed, "export "),
			"Expected 'export ' prefix in .envrc, got: %s", trimmed)
	}
}

func TestEnvWriter_Write_DotEnv_NoExport(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env.local")

	vars := []EnvVar{
		{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
	}

	writer := NewEnvWriter(envPath)
	err := writer.Write(vars)
	require.NoError(t, err)

	content, err := os.ReadFile(envPath)
	require.NoError(t, err)

	contentStr := string(content)

	// Should NOT use export syntax for .env files
	assert.NotContains(t, contentStr, "export ")
	assert.Contains(t, contentStr, "API_URL=http://192.168.1.100:8000")
}

func TestEnvWriter_Read_EnvRC_WithExport(t *testing.T) {
	tmpDir := t.TempDir()
	envrcPath := filepath.Join(tmpDir, ".envrc")

	// Write a .envrc with export syntax
	content := `# Generated by lanup
# lanup:managed
export API_URL=http://192.168.1.100:8000
# lanup:managed
export DB_URL=postgresql://192.168.1.100:5432/db

# User variables
export MY_SECRET=supersecret
`
	err := os.WriteFile(envrcPath, []byte(content), 0644)
	require.NoError(t, err)

	writer := NewEnvWriter(envrcPath)
	vars, err := writer.Read()
	require.NoError(t, err)

	assert.Len(t, vars, 3)

	// Verify managed variables
	assert.Equal(t, "API_URL", vars[0].Key)
	assert.Equal(t, "http://192.168.1.100:8000", vars[0].Value)
	assert.True(t, vars[0].Managed)

	assert.Equal(t, "DB_URL", vars[1].Key)
	assert.Equal(t, "postgresql://192.168.1.100:5432/db", vars[1].Value)
	assert.True(t, vars[1].Managed)

	// Verify user variable
	assert.Equal(t, "MY_SECRET", vars[2].Key)
	assert.Equal(t, "supersecret", vars[2].Value)
	assert.False(t, vars[2].Managed)
}

func TestEnvWriter_RoundTrip_EnvRC(t *testing.T) {
	tmpDir := t.TempDir()
	envrcPath := filepath.Join(tmpDir, ".envrc")

	// Write variables
	original := []EnvVar{
		{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
		{Key: "CUSTOM_VAR", Value: "my_value", Managed: false},
	}

	writer := NewEnvWriter(envrcPath)
	writer.BackupEnabled = false
	err := writer.Write(original)
	require.NoError(t, err)

	// Read them back
	readBack, err := writer.Read()
	require.NoError(t, err)

	assert.Len(t, readBack, 2)
	assert.Equal(t, "API_URL", readBack[0].Key)
	assert.Equal(t, "http://192.168.1.100:8000", readBack[0].Value)
	assert.True(t, readBack[0].Managed)

	assert.Equal(t, "CUSTOM_VAR", readBack[1].Key)
	assert.Equal(t, "my_value", readBack[1].Value)
	assert.False(t, readBack[1].Managed)
}

func TestEnvWriter_Write_EnvRC_UserVars(t *testing.T) {
	tmpDir := t.TempDir()
	envrcPath := filepath.Join(tmpDir, ".envrc")

	vars := []EnvVar{
		{Key: "API_URL", Value: "http://192.168.1.100:8000", Managed: true},
		{Key: "MY_SECRET", Value: "supersecret", Managed: false},
	}

	writer := NewEnvWriter(envrcPath)
	err := writer.Write(vars)
	require.NoError(t, err)

	content, err := os.ReadFile(envrcPath)
	require.NoError(t, err)

	contentStr := string(content)

	// Both managed and user vars should use export syntax
	assert.Contains(t, contentStr, "export API_URL=http://192.168.1.100:8000")
	assert.Contains(t, contentStr, "export MY_SECRET=supersecret")
	assert.Contains(t, contentStr, "# User variables (preserved)")
}
