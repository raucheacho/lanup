package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/raucheacho/lanup/internal/config"
	"github.com/raucheacho/lanup/internal/docker"
	"github.com/raucheacho/lanup/internal/env"
	"github.com/raucheacho/lanup/internal/logger"
	"github.com/raucheacho/lanup/internal/net"
	lanuperrors "github.com/raucheacho/lanup/pkg/errors"
	"github.com/raucheacho/lanup/pkg/utils"
	"github.com/spf13/cobra"
)

// StartCmd represents the start command
type StartCmd struct {
	Watch  bool
	NoEnv  bool
	DryRun bool
	Log    bool
	logger *logger.Logger
}

// NewStartCmd creates a new start command
func NewStartCmd() *cobra.Command {
	startCmd := &StartCmd{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start exposing local services on your LAN",
		Long: `Detect your local IP address and generate environment variables for your services.

This command reads the .lanup.yaml configuration file, detects your local IP address,
and generates a .env file with URLs that can be accessed from any device on your network.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return startCmd.Run()
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&startCmd.Watch, "watch", "w", false, "watch for network changes and update automatically")
	cmd.Flags().BoolVar(&startCmd.NoEnv, "no-env", false, "display variables without writing to file")
	cmd.Flags().BoolVar(&startCmd.DryRun, "dry-run", false, "simulate all operations without writing files")
	cmd.Flags().BoolVar(&startCmd.Log, "log", true, "enable logging to file")

	return cmd
}

func init() {
	RootCmd.AddCommand(NewStartCmd())
}

// Run executes the start command
func (c *StartCmd) Run() error {
	// Initialize logger if enabled
	if c.Log {
		globalCfg := GetGlobalConfig()
		if globalCfg != nil {
			logLevel := logger.INFO
			switch strings.ToLower(globalCfg.LogLevel) {
			case "debug":
				logLevel = logger.DEBUG
			case "warn":
				logLevel = logger.WARN
			case "error":
				logLevel = logger.ERROR
			}

			var err error
			c.logger, err = logger.NewLogger(logger.LoggerConfig{
				Level:      logLevel,
				FilePath:   globalCfg.LogPath,
				MaxSize:    5 * 1024 * 1024, // 5MB
				MaxBackups: 5,
				Console:    false,
				Colors:     false,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to initialize logger: %v\n", err)
			} else {
				defer c.logger.Close()
			}
		}
	}

	// Load project configuration
	projectConfig, err := config.LoadProjectConfig("")
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrInvalidConfig,
			"Failed to load project configuration", err)
	}

	if c.logger != nil {
		c.logger.Info("Starting lanup", logger.Field{Key: "watch", Value: c.Watch})
	}

	// Execute the core start logic
	if err := c.executeStart(projectConfig); err != nil {
		if c.logger != nil {
			c.logger.Error("Start failed", logger.Field{Key: "error", Value: err.Error()})
		}
		return err
	}

	// If watch mode is enabled, start watching for network changes
	if c.Watch {
		return c.watchMode(projectConfig)
	}

	return nil
}

// executeStart performs the core start logic
func (c *StartCmd) executeStart(projectConfig *config.ProjectConfig) error {
	// Detect local IP
	netInfo, err := net.DetectLocalIP()
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrNoNetwork,
			"Failed to detect local IP address", err)
	}

	if c.logger != nil {
		c.logger.Info("Detected IP",
			logger.Field{Key: "ip", Value: netInfo.IP},
			logger.Field{Key: "interface", Value: netInfo.Interface},
			logger.Field{Key: "type", Value: netInfo.Type})
	}

	// Collect variables from configuration
	vars := make(map[string]string)
	for key, value := range projectConfig.Vars {
		vars[key] = value
	}

	// Handle Docker auto-detection if enabled
	if projectConfig.AutoDetect.Docker {
		if docker.IsDockerAvailable() {
			containers, err := docker.GetRunningContainers()
			if err != nil {
				if c.logger != nil {
					c.logger.Warn("Failed to get Docker containers", logger.Field{Key: "error", Value: err.Error()})
				}
				fmt.Fprintf(os.Stderr, "⚠️  Warning: Failed to detect Docker containers: %v\n", err)
			} else {
				if c.logger != nil {
					c.logger.Info("Detected Docker containers", logger.Field{Key: "count", Value: len(containers)})
				}
				// Add Docker container ports to variables
				for _, container := range containers {
					for _, port := range container.Ports {
						varName := fmt.Sprintf("DOCKER_%s_PORT", strings.ToUpper(strings.ReplaceAll(container.Name, "-", "_")))
						vars[varName] = fmt.Sprintf("http://localhost:%d", port.HostPort)
					}
				}
			}
		}
	}

	// Handle Supabase auto-detection if enabled
	if projectConfig.AutoDetect.Supabase {
		services, err := docker.GetSupabaseStatus()
		if err != nil {
			if c.logger != nil {
				c.logger.Warn("Failed to get Supabase status", logger.Field{Key: "error", Value: err.Error()})
			}
			// Don't show warning for Supabase as it's optional
		} else {
			if c.logger != nil {
				c.logger.Info("Detected Supabase services", logger.Field{Key: "count", Value: len(services)})
			}
			// Add Supabase service ports to variables
			for serviceName, port := range services {
				varName := fmt.Sprintf("SUPABASE_%s_PORT", strings.ToUpper(strings.ReplaceAll(serviceName, "_", "_")))
				vars[varName] = fmt.Sprintf("http://localhost:%d", port)
			}
		}
	}

	// Transform URLs from localhost to detected IP
	transformedVars := make([]env.EnvVar, 0, len(vars))
	for key, value := range vars {
		transformedValue := transformURL(value, netInfo.IP)
		transformedVars = append(transformedVars, env.EnvVar{
			Key:     key,
			Value:   transformedValue,
			Managed: true,
		})
	}

	// If no-env or dry-run, just display the variables
	if c.NoEnv || c.DryRun {
		c.displayVariables(transformedVars, netInfo.IP, c.DryRun)
		return nil
	}

	// Read existing .env file
	envWriter := env.NewEnvWriter(projectConfig.Output)
	existingVars, err := envWriter.Read()
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrFileNotFound,
			"Failed to read existing env file", err)
	}

	// Merge new and existing variables
	mergedVars := envWriter.Merge(transformedVars, existingVars)

	// Write the new .env file
	if err := envWriter.Write(mergedVars); err != nil {
		return lanuperrors.NewError(lanuperrors.ErrPermissionDenied,
			"Failed to write env file", err)
	}

	if c.logger != nil {
		c.logger.Info("Updated env file",
			logger.Field{Key: "path", Value: projectConfig.Output},
			logger.Field{Key: "vars", Value: len(transformedVars)})
	}

	// Display success message and URLs
	c.displaySuccess(transformedVars, netInfo.IP, projectConfig.Output)

	return nil
}

// transformURL replaces localhost or 127.0.0.1 with the detected IP address
func transformURL(url string, newIP string) string {
	// Replace localhost
	url = strings.ReplaceAll(url, "localhost", newIP)

	// Replace 127.0.0.1
	url = strings.ReplaceAll(url, "127.0.0.1", newIP)

	return url
}

// displayVariables shows the environment variables in the console
func (c *StartCmd) displayVariables(vars []env.EnvVar, ip string, isDryRun bool) {
	if isDryRun {
		utils.Info("Dry run mode - no files will be modified")
		fmt.Println()
	}

	utils.Success("Detected local IP: %s", ip)
	fmt.Println()

	if len(vars) > 0 {
		utils.PrintSection("Environment Variables")
		for _, v := range vars {
			fmt.Printf("  %s=%s\n", color.CyanString(v.Key), v.Value)
		}
	}
}

// displaySuccess shows a success message with the exposed URLs
func (c *StartCmd) displaySuccess(vars []env.EnvVar, ip string, outputPath string) {
	utils.Success("Successfully exposed services on your LAN!")
	utils.Success("Environment file updated: %s", outputPath)
	utils.Success("Local IP: %s", ip)
	fmt.Println()

	if len(vars) > 0 {
		utils.PrintSection("Your services are now accessible at")
		for _, v := range vars {
			// Only display URLs (values that start with http)
			if strings.HasPrefix(v.Value, "http") {
				utils.PrintURL(v.Key, v.Value)
			}
		}
		fmt.Println()
	}

	utils.Info("Tip: Use 'lanup start --watch' to automatically update when your network changes")
}

// watchMode starts watching for network changes and regenerates the .env file
func (c *StartCmd) watchMode(projectConfig *config.ProjectConfig) error {
	fmt.Println()
	utils.Info("Watch mode enabled - monitoring network changes...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Get check interval from global config
	globalCfg := GetGlobalConfig()
	interval := 5 * time.Second
	if globalCfg != nil && globalCfg.CheckInterval > 0 {
		interval = time.Duration(globalCfg.CheckInterval) * time.Second
	}

	// Create IP watcher
	watcher := net.NewIPWatcher(interval)

	// Set up the OnChange callback
	watcher.OnChange = func(oldIP, newIP string) {
		if c.logger != nil {
			c.logger.Warn("Network interface changed",
				logger.Field{Key: "old_ip", Value: oldIP},
				logger.Field{Key: "new_ip", Value: newIP})
		}

		fmt.Println()
		utils.Warning("Network change detected!")
		fmt.Printf("  Old IP: %s\n", color.CyanString(oldIP))
		fmt.Printf("  New IP: %s\n", color.CyanString(newIP))
		fmt.Println()
		utils.Info("Regenerating environment file...")

		// Regenerate the .env file with the new IP
		if err := c.executeStart(projectConfig); err != nil {
			utils.Error("Failed to regenerate env file: %v", err)
			if c.logger != nil {
				c.logger.Error("Failed to regenerate env file", logger.Field{Key: "error", Value: err.Error()})
			}
		} else {
			utils.Success("Environment file updated successfully!")
			fmt.Println()
		}
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Start the watcher in a goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := watcher.Start(ctx); err != nil && err != context.Canceled {
			errCh <- err
		}
	}()

	// Wait for signal or error
	select {
	case <-sigCh:
		fmt.Println()
		fmt.Println("Shutting down gracefully...")
		cancel()
		watcher.Stop()
		if c.logger != nil {
			c.logger.Info("Watch mode stopped by user")
		}
		return nil
	case err := <-errCh:
		cancel()
		watcher.Stop()
		return fmt.Errorf("watcher error: %w", err)
	}
}
