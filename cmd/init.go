package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/raucheacho/lanup/internal/config"
	lanuperrors "github.com/raucheacho/lanup/pkg/errors"
	"github.com/raucheacho/lanup/pkg/utils"
	"github.com/spf13/cobra"
)

// InitCmd represents the init command
type InitCmd struct {
	Format string
	Force  bool
}

// NewInitCmd creates a new init command
func NewInitCmd() *cobra.Command {
	initCmd := &InitCmd{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize lanup configuration in the current project",
		Long: `Initialize lanup configuration by creating a .lanup.yaml file in the current directory.

This file defines which services should be exposed on your local network.
You can customize the variables, output file path, and auto-detection settings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initCmd.Run()
		},
	}

	// Add flags
	cmd.Flags().StringVar(&initCmd.Format, "format", "yaml", "configuration file format (yaml or toml)")
	cmd.Flags().BoolVar(&initCmd.Force, "force", false, "overwrite existing configuration file")

	return cmd
}

func init() {
	RootCmd.AddCommand(NewInitCmd())
}

// Run executes the init command
func (c *InitCmd) Run() error {
	// Validate format
	if c.Format != "yaml" && c.Format != "toml" {
		return lanuperrors.NewError(lanuperrors.ErrInvalidConfig,
			fmt.Sprintf("Unsupported format: %s (supported: yaml, toml)", c.Format), nil)
	}

	// Note: Currently only YAML is implemented
	if c.Format == "toml" {
		return lanuperrors.NewError(lanuperrors.ErrInvalidConfig,
			"TOML format is not yet supported, please use yaml", nil)
	}

	// Determine config file path
	configPath := ".lanup.yaml"

	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		if !c.Force {
			return lanuperrors.NewError(lanuperrors.ErrFileNotFound,
				fmt.Sprintf("Configuration file already exists at %s\nUse --force to overwrite", configPath), nil)
		}
		utils.Warning("Overwriting existing configuration file at %s", configPath)
	}

	// Generate default configuration
	defaultConfig := config.GetDefaultProjectConfig()

	// Save configuration to file
	if err := config.SaveProjectConfig(configPath, defaultConfig); err != nil {
		return lanuperrors.NewError(lanuperrors.ErrPermissionDenied,
			"Failed to create configuration file", err)
	}

	// Get absolute path for display
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		absPath = configPath
	}

	// Display success message
	utils.Success("Configuration file created successfully!")
	utils.Info("Location: %s", absPath)
	fmt.Println()
	utils.PrintSection("Next steps")
	fmt.Printf("  1. Edit %s to configure your services\n", configPath)
	fmt.Printf("  2. Run 'lanup start' to expose your services on the LAN\n")

	return nil
}
