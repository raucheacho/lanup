package cmd

import (
	"fmt"

	"github.com/raucheacho/lanup/internal/docker"
	"github.com/raucheacho/lanup/internal/net"
	lanuperrors "github.com/raucheacho/lanup/pkg/errors"
	"github.com/raucheacho/lanup/pkg/utils"
	"github.com/spf13/cobra"
)

// DoctorCmd represents the doctor command
type DoctorCmd struct{}

// HealthCheck represents the result of a health check
type HealthCheck struct {
	Name    string
	Status  bool
	Message string
}

// NewDoctorCmd creates a new doctor command
func NewDoctorCmd() *cobra.Command {
	doctorCmd := &DoctorCmd{}

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose your local environment",
		Long: `Run diagnostic checks to verify that lanup can function properly.

This command checks:
  - Network interfaces and local IP detection
  - Docker availability and running containers
  - Supabase local development setup

Use this command to troubleshoot issues with lanup.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doctorCmd.Run()
		},
	}

	return cmd
}

func init() {
	RootCmd.AddCommand(NewDoctorCmd())
}

// Run executes the doctor command
func (c *DoctorCmd) Run() error {
	utils.PrintSection("Running lanup diagnostics")

	// Run all health checks
	checks := []HealthCheck{
		checkNetworkInterfaces(),
		checkDocker(),
		checkSupabase(),
	}

	// Display results
	allPassed := true
	for _, check := range checks {
		if check.Status {
			utils.Success("%s", check.Name)
		} else {
			utils.Error("%s", check.Name)
			allPassed = false
		}
		if check.Message != "" {
			fmt.Printf("   %s\n", check.Message)
		}
	}

	// Summary
	fmt.Println()
	if allPassed {
		utils.Success("All checks passed! lanup is ready to use.")
		return nil
	} else {
		utils.Warning("Some checks failed. Please review the issues above.")
		return lanuperrors.NewError(lanuperrors.ErrNoNetwork,
			"Health checks failed", nil)
	}
}

// checkNetworkInterfaces verifies that active network interfaces are available
func checkNetworkInterfaces() HealthCheck {
	netInfo, err := net.DetectLocalIP()
	if err != nil {
		return HealthCheck{
			Name:    "Network Interfaces",
			Status:  false,
			Message: fmt.Sprintf("Failed to detect local IP: %v", err),
		}
	}

	return HealthCheck{
		Name:    "Network Interfaces",
		Status:  true,
		Message: fmt.Sprintf("Detected IP: %s on interface %s (%s)", netInfo.IP, netInfo.Interface, netInfo.Type),
	}
}

// checkDocker verifies Docker availability and running containers
func checkDocker() HealthCheck {
	if !docker.IsDockerAvailable() {
		return HealthCheck{
			Name:    "Docker",
			Status:  false,
			Message: "Docker is not installed or not running",
		}
	}

	containers, err := docker.GetRunningContainers()
	if err != nil {
		return HealthCheck{
			Name:    "Docker",
			Status:  false,
			Message: fmt.Sprintf("Docker is available but failed to list containers: %v", err),
		}
	}

	if len(containers) == 0 {
		return HealthCheck{
			Name:    "Docker",
			Status:  true,
			Message: "Docker is running (no containers currently active)",
		}
	}

	return HealthCheck{
		Name:    "Docker",
		Status:  true,
		Message: fmt.Sprintf("Docker is running with %d active container(s)", len(containers)),
	}
}

// checkSupabase verifies Supabase local development status
func checkSupabase() HealthCheck {
	services, err := docker.GetSupabaseStatus()
	if err != nil {
		return HealthCheck{
			Name:    "Supabase",
			Status:  false,
			Message: fmt.Sprintf("Supabase local is not running: %v", err),
		}
	}

	if len(services) == 0 {
		return HealthCheck{
			Name:    "Supabase",
			Status:  false,
			Message: "Supabase CLI is available but no services detected",
		}
	}

	return HealthCheck{
		Name:    "Supabase",
		Status:  true,
		Message: fmt.Sprintf("Supabase local is running with %d service(s)", len(services)),
	}
}
