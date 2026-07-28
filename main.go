package main

import (
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
	// Exit codes are handled by cmd.Execute() via pkg/errors
	cmd.Execute()
}
