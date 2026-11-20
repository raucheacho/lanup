package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0600))
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}

func TestUpdate_KeepsEverythingElseInEnvrc(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".envrc")
	writeFile(t, path, `# Backend
export DB="postgres://${USER:-me}@localhost/db"
export MAIL_FROM="Livro <no-reply@example.com>"

# Realtime
export WS_URL="ws://192.168.1.10:4000/socket"
export ORIGINS=http://localhost:3000,http://192.168.1.10:3000
dotenv_if_exists .env.extra
`)

	w := NewEnvWriter(path)
	w.BackupEnabled = false
	require.NoError(t, w.Update([]EnvVar{
		{Key: "WS_URL", Value: "ws://192.168.11.101:4000/socket", Managed: true},
		{Key: "ORIGINS", Value: "http://localhost:3000,http://192.168.11.101:3000", Managed: true},
	}))

	assert.Equal(t, `# Backend
export DB="postgres://${USER:-me}@localhost/db"
export MAIL_FROM="Livro <no-reply@example.com>"

# Realtime
export WS_URL="ws://192.168.11.101:4000/socket"
export ORIGINS=http://localhost:3000,http://192.168.11.101:3000
dotenv_if_exists .env.extra
`, readFile(t, path))

	// The file mode is kept.
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestUpdate_AppendsMissingAndNeverRemoves(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env.local")
	writeFile(t, path, `# Maps
MAPS_KEY=abc

# lanup:managed
OLD_URL=http://192.168.1.10:1
`)

	w := NewEnvWriter(path)
	w.BackupEnabled = false
	require.NoError(t, w.Update([]EnvVar{
		{Key: "API_URL", Value: "http://192.168.11.101:8084", Managed: true},
	}))

	// OLD_URL is no longer in the config: it stays, lanup never deletes a line.
	assert.Equal(t, `# Maps
MAPS_KEY=abc

# lanup:managed
OLD_URL=http://192.168.1.10:1

# lanup:managed
API_URL=http://192.168.11.101:8084
`, readFile(t, path))
}

func TestUpdate_CreatesFileAndQuotesShellValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", ".envrc")

	w := NewEnvWriter(path)
	w.BackupEnabled = false
	require.NoError(t, w.Update([]EnvVar{
		{Key: "NAME", Value: "Livro app", Managed: true},
		{Key: "URL", Value: "http://192.168.11.101:8084", Managed: true},
	}))

	assert.Equal(t, `# lanup:managed
export NAME="Livro app"
# lanup:managed
export URL=http://192.168.11.101:8084
`, readFile(t, path))
}

func TestUpdate_NoChangeNoBackup(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env.local")
	content := "API_URL=http://192.168.11.101:8084\n"
	writeFile(t, path, content)

	w := NewEnvWriter(path)
	require.NoError(t, w.Update([]EnvVar{{Key: "API_URL", Value: "http://192.168.11.101:8084", Managed: true}}))

	assert.Equal(t, content, readFile(t, path))
	_, err := os.Stat(path + ".bak")
	assert.True(t, os.IsNotExist(err), "an unchanged file needs no backup")
}

func TestUpdate_KeepsQuotesAndEscapes(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	writeFile(t, path, "A='old'\nB=\"old\"\n")

	w := NewEnvWriter(path)
	w.BackupEnabled = false
	require.NoError(t, w.Update([]EnvVar{
		{Key: "A", Value: "it's", Managed: true},
		{Key: "B", Value: `say "hi" $HOME`, Managed: true},
	}))

	assert.Equal(t, "A=\"it's\"\nB=\"say \\\"hi\\\" \\$HOME\"\n", readFile(t, path))
}
