package cmd

import (
	"os"

	"github.com/raucheacho/lanup/internal/config"
	lanuperrors "github.com/raucheacho/lanup/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	verbose bool

	// Global configuration loaded at startup
	globalConfig *config.GlobalConfig

	// Version information (set during build)
	Version = "dev"
)

// RootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:   "lanup",
	Short: "lanup - Expose local services on your LAN",
	Long: `lanup is a CLI tool that automatically exposes your local backend services on your local network.

It detects your local IP address, updates environment variables, and allows you to test
your applications from any device on the same network without manual configuration.`,
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
}

// Execute runs the root command
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		// If it's a LanupError, exit with the appropriate code
		if lanupErr, ok := err.(*lanuperrors.LanupError); ok {
			os.Exit(lanupErr.ExitCode())
		}
		// Otherwise, exit with generic error code
		os.Exit(1)
	}
}

func init() {
	// Add persistent flags available to all commands
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}

// initConfig reads in config file and ENV variables if set
func initConfig() error {
	var err error

	// Load global configuration
	globalConfig, err = config.LoadGlobalConfig()
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrInvalidConfig, "Failed to load global configuration", err)
	}

	// If verbose flag is set, override log level
	if verbose {
		globalConfig.LogLevel = "debug"
	}

	return nil
}

// GetGlobalConfig returns the loaded global configuration
func GetGlobalConfig() *config.GlobalConfig {
	return globalConfig
}

// IsVerbose returns whether verbose mode is enabled
func IsVerbose() bool {
	return verbose
}
