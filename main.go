package main

import (
	"fmt"
	"os"

	"github.com/raucheacho/lanup/cmd"
)

// Version information (set during build with -ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version information in cmd package
	cmd.Version = version

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		// Error is already printed by Cobra, just exit with error code
		os.Exit(getExitCode(err))
	}
}

// getExitCode returns the appropriate exit code based on the error
func getExitCode(err error) int {
	if err == nil {
		return 0
	}

	// Check for specific error types and return appropriate codes
	errMsg := err.Error()

	// Configuration errors
	if contains(errMsg, "config") || contains(errMsg, "configuration") {
		return 2
	}

	// Network errors
	if contains(errMsg, "network") || contains(errMsg, "interface") || contains(errMsg, "IP") {
		return 3
	}

	// Permission errors
	if contains(errMsg, "permission") || contains(errMsg, "denied") {
		return 4
	}

	// Validation errors
	if contains(errMsg, "invalid") || contains(errMsg, "validation") {
		return 5
	}

	// Default error code
	return 1
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func init() {
	// Print version info if requested via environment variable (for debugging)
	if os.Getenv("LANUP_VERSION_INFO") != "" {
		fmt.Printf("lanup version %s (commit: %s, built: %s)\n", version, commit, date)
	}
}
