package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/raucheacho/lanup/internal/config"
	"github.com/raucheacho/lanup/internal/docker"
	"github.com/raucheacho/lanup/internal/env"
	"github.com/raucheacho/lanup/internal/net"
	lanuperrors "github.com/raucheacho/lanup/pkg/errors"
	"github.com/raucheacho/lanup/pkg/utils"
	"github.com/spf13/cobra"
)

// StartCmd represents the start command
type StartCmd struct {
	Once   bool
	NoEnv  bool
	DryRun bool
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
	cmd.Flags().BoolVarP(&startCmd.Once, "once", "1", false, "run once and exit without watching for network changes")
	cmd.Flags().BoolVar(&startCmd.NoEnv, "no-env", false, "display variables without writing to file")
	cmd.Flags().BoolVar(&startCmd.DryRun, "dry-run", false, "simulate all operations without writing files")

	return cmd
}

func init() {
	RootCmd.AddCommand(NewStartCmd())
}

// Run executes the start command
func (c *StartCmd) Run() error {
	// Load project configuration
	projectConfig, err := config.LoadProjectConfig("")
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrInvalidConfig,
			"Failed to load project configuration", err)
	}

	if IsVerbose() {
		fmt.Fprintf(os.Stderr, "[verbose] Starting lanup (once=%v)\n", c.Once)
	}

	// Execute the core start logic
	if err := c.executeStart(projectConfig); err != nil {
		return err
	}

	// If not in once mode, start watching for network changes
	if !c.Once {
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

	if IsVerbose() {
		fmt.Fprintf(os.Stderr, "[verbose] Detected IP: %s on %s (%s)\n", netInfo.IP, netInfo.Interface, netInfo.Type)
	}

	// Auto-detected services go to every output
	detected := make(map[string]string)

	// Handle Docker auto-detection if enabled
	if projectConfig.AutoDetect.Docker {
		if docker.IsDockerAvailable() {
			containers, err := docker.GetRunningContainers()
			if err != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Warning: Failed to detect Docker containers: %v\n", err)
			} else {
				if IsVerbose() {
					fmt.Fprintf(os.Stderr, "[verbose] Detected %d Docker container(s)\n", len(containers))
				}
				// Add Docker container ports to variables
				for _, container := range containers {
					for _, port := range container.Ports {
						varName := fmt.Sprintf("DOCKER_%s_PORT", strings.ToUpper(strings.ReplaceAll(container.Name, "-", "_")))
						detected[varName] = fmt.Sprintf("http://localhost:%d", port.HostPort)
					}
				}
			}
		}
	}

	// Handle Supabase auto-detection if enabled
	if projectConfig.AutoDetect.Supabase {
		services, err := docker.GetSupabaseStatus()
		if err != nil {
			if IsVerbose() {
				fmt.Fprintf(os.Stderr, "[verbose] Supabase not available: %v\n", err)
			}
		} else {
			if IsVerbose() {
				fmt.Fprintf(os.Stderr, "[verbose] Detected %d Supabase service(s)\n", len(services))
			}
			// Add Supabase service ports to variables
			for serviceName, port := range services {
				varName := fmt.Sprintf("SUPABASE_%s_PORT", strings.ToUpper(strings.ReplaceAll(serviceName, "_", "_")))
				detected[varName] = fmt.Sprintf("http://localhost:%d", port)
			}
		}
	}

	targets := projectConfig.Targets()
	if c.NoEnv || c.DryRun {
		if c.DryRun {
			utils.Info("Dry run mode - no files will be modified")
			fmt.Println()
		}
		utils.Success("Detected local IP: %s", netInfo.IP)
	}

	var written []writtenTarget
	for _, target := range targets {
		vars := make(map[string]string, len(target.Vars)+len(detected))
		for key, value := range detected {
			vars[key] = value
		}
		for key, value := range target.Vars {
			vars[key] = value
		}

		// Transform URLs from localhost to detected IP, sorted for a stable file
		keys := make([]string, 0, len(vars))
		for key := range vars {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		transformedVars := make([]env.EnvVar, 0, len(vars))
		for _, key := range keys {
			transformedVars = append(transformedVars, env.EnvVar{
				Key:     key,
				Value:   TransformValue(vars[key], netInfo.IP),
				Managed: true,
			})
		}

		// If no-env or dry-run, just display the variables
		if c.NoEnv || c.DryRun {
			c.displayVariables(transformedVars, target.Path)
			continue
		}

		envWriter := env.NewEnvWriter(target.Path)
		if target.InPlace {
			if err := envWriter.Update(transformedVars); err != nil {
				return lanuperrors.NewError(lanuperrors.ErrPermissionDenied,
					"Failed to update env file "+target.Path, err)
			}
		} else {
			// Read existing .env file
			existingVars, err := envWriter.Read()
			if err != nil {
				return lanuperrors.NewError(lanuperrors.ErrFileNotFound,
					"Failed to read existing env file", err)
			}

			// Merge new and existing variables, then write the new .env file
			mergedVars := envWriter.Merge(transformedVars, existingVars)
			if err := envWriter.Write(mergedVars); err != nil {
				return lanuperrors.NewError(lanuperrors.ErrPermissionDenied,
					"Failed to write env file", err)
			}
		}

		if IsVerbose() {
			fmt.Fprintf(os.Stderr, "[verbose] Updated %s with %d variable(s)\n", target.Path, len(transformedVars))
		}
		written = append(written, writtenTarget{path: target.Path, vars: transformedVars})
	}

	if c.NoEnv || c.DryRun {
		return nil
	}

	// Display success message and URLs
	c.displaySuccess(written, netInfo.IP)

	return nil
}

// writtenTarget is an env file lanup just wrote, for the summary
type writtenTarget struct {
	path string
	vars []env.EnvVar
}

// ipPlaceholder marks where the detected IP goes in a value
const ipPlaceholder = "{{ip}}"

// TransformValue puts the detected IP in a value: in place of {{ip}} when the
// value has it (so localhost can stay next to the IP), else in place of localhost
func TransformValue(value string, ip string) string {
	if strings.Contains(value, ipPlaceholder) {
		return strings.ReplaceAll(value, ipPlaceholder, ip)
	}
	return TransformURL(value, ip)
}

// TransformURL replaces localhost or 127.0.0.1 with the detected IP address
func TransformURL(url string, newIP string) string {
	// Replace localhost
	url = strings.ReplaceAll(url, "localhost", newIP)

	// Replace 127.0.0.1
	url = strings.ReplaceAll(url, "127.0.0.1", newIP)

	return url
}

// displayVariables shows the environment variables of one output in the console
func (c *StartCmd) displayVariables(vars []env.EnvVar, path string) {
	fmt.Println()
	utils.PrintSection(path)
	for _, v := range vars {
		fmt.Printf("  %s=%s\n", color.CyanString(v.Key), v.Value)
	}
}

// displaySuccess shows a success message with the exposed URLs of each file
func (c *StartCmd) displaySuccess(targets []writtenTarget, ip string) {
	utils.Success("Successfully exposed services on your LAN!")
	for _, t := range targets {
		utils.Success("Environment file updated: %s", t.path)
	}
	utils.Success("Local IP: %s", ip)
	fmt.Println()

	for _, t := range targets {
		var urls []env.EnvVar
		for _, v := range t.vars {
			// Only display URLs
			if strings.Contains(v.Value, "://") {
				urls = append(urls, v)
			}
		}
		if len(urls) == 0 {
			continue
		}
		utils.PrintSection(t.path)
		for _, v := range urls {
			utils.PrintURL(v.Key, v.Value)
		}
		fmt.Println()
	}

	utils.Info("Tip: Use 'lanup start --once' to run once and exit")
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

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start file watcher for .lanup.yaml
	go func() {
		configPath := ".lanup.yaml"
		stat, err := os.Stat(configPath)
		if err != nil {
			return
		}
		lastModTime := stat.ModTime()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stat, err := os.Stat(configPath)
				if err != nil {
					continue
				}
				if stat.ModTime().After(lastModTime) {
					lastModTime = stat.ModTime()

					fmt.Println()
					utils.Warning("Configuration file changed: %s", configPath)
					utils.Info("Reloading configuration and regenerating environment file...")

					// Reload project config
					newConfig, err := config.LoadProjectConfig(configPath)
					if err != nil {
						utils.Error("Failed to reload configuration: %v", err)
						continue
					}

					// Update config and regenerate
					projectConfig = newConfig
					if err := c.executeStart(projectConfig); err != nil {
						utils.Error("Failed to regenerate env file: %v", err)
					} else {
						utils.Success("Environment file updated successfully!")
						fmt.Println()
					}
				}
			}
		}
	}()

	// Create IP watcher
	watcher := net.NewIPWatcher(interval)

	// Set up the OnChange callback
	watcher.OnChange = func(oldIP, newIP string) {
		fmt.Println()
		utils.Warning("Network change detected!")
		fmt.Printf("  Old IP: %s\n", color.CyanString(oldIP))
		fmt.Printf("  New IP: %s\n", color.CyanString(newIP))
		fmt.Println()
		utils.Info("Regenerating environment file...")

		// Regenerate the .env file with the new IP
		if err := c.executeStart(projectConfig); err != nil {
			utils.Error("Failed to regenerate env file: %v", err)
		} else {
			utils.Success("Environment file updated successfully!")
			fmt.Println()
		}
	}

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
		if IsVerbose() {
			fmt.Fprintf(os.Stderr, "[verbose] Watch mode stopped by user\n")
		}
		return nil
	case err := <-errCh:
		cancel()
		watcher.Stop()
		return fmt.Errorf("watcher error: %w", err)
	}
}
