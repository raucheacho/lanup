package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestProjectConfig_Outputs(t *testing.T) {
	raw := `
vars:
  API_URL: http://localhost:8084
outputs:
  - path: backend/.envrc
    vars:
      ORIGINS: http://localhost:3000,http://{{ip}}:3000
  - path: apps/web/.env.local
    vars:
      API_URL: http://localhost:9000
`
	var cfg ProjectConfig
	require.NoError(t, yaml.Unmarshal([]byte(raw), &cfg))
	require.NoError(t, cfg.Validate())

	targets := cfg.Targets()
	require.Len(t, targets, 2)

	assert.Equal(t, "backend/.envrc", targets[0].Path)
	assert.True(t, targets[0].InPlace)
	assert.Equal(t, map[string]string{
		"API_URL": "http://localhost:8084",
		"ORIGINS": "http://localhost:3000,http://{{ip}}:3000",
	}, targets[0].Vars)

	// An output's own vars win over the shared ones.
	assert.Equal(t, "http://localhost:9000", targets[1].Vars["API_URL"])
	// The shared map is not changed by an output.
	assert.Equal(t, "http://localhost:8084", cfg.Vars["API_URL"])
}

func TestProjectConfig_LegacyOutputIsRegenerated(t *testing.T) {
	cfg := ProjectConfig{Vars: map[string]string{"A": "http://localhost:1"}, Output: ".env.local"}
	require.NoError(t, cfg.Validate())

	targets := cfg.Targets()
	require.Len(t, targets, 1)
	assert.Equal(t, ".env.local", targets[0].Path)
	assert.False(t, targets[0].InPlace)
}

func TestProjectConfig_OutputsValidation(t *testing.T) {
	tests := []struct {
		name string
		cfg  ProjectConfig
	}{
		{"no output at all", ProjectConfig{Vars: map[string]string{"A": "b"}}},
		{"empty path", ProjectConfig{Outputs: []OutputConfig{{Vars: map[string]string{"A": "b"}}}}},
		{"same file twice", ProjectConfig{Outputs: []OutputConfig{
			{Path: "a/.env", Vars: map[string]string{"A": "b"}},
			{Path: "./a/.env", Vars: map[string]string{"B": "c"}},
		}}},
		{"nothing to write", ProjectConfig{Outputs: []OutputConfig{{Path: ".env"}}}},
		{"empty value", ProjectConfig{Outputs: []OutputConfig{{Path: ".env", Vars: map[string]string{"A": ""}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.cfg.Validate())
		})
	}
}
