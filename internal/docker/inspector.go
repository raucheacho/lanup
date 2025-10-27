package docker

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// DockerService represents a running Docker container with its port mappings
type DockerService struct {
	ContainerID string
	Name        string
	Ports       []PortMapping
}

// PortMapping represents a port mapping between host and container
type PortMapping struct {
	HostPort      int
	ContainerPort int
	Protocol      string
}

// IsDockerAvailable checks if Docker is installed and running
func IsDockerAvailable() bool {
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	return err == nil
}

// GetRunningContainers returns a list of running Docker containers with their port mappings
func GetRunningContainers() ([]DockerService, error) {
	if !IsDockerAvailable() {
		return nil, fmt.Errorf("docker is not available")
	}

	cmd := exec.Command("docker", "ps", "--format", "{{.ID}}|{{.Names}}|{{.Ports}}")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute docker ps: %w", err)
	}

	return ParseDockerPS(out.String())
}

// ParseDockerPS parses the output of docker ps command and extracts container information
func ParseDockerPS(output string) ([]DockerService, error) {
	if strings.TrimSpace(output) == "" {
		return []DockerService{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	services := make([]DockerService, 0, len(lines))

	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}

		service := DockerService{
			ContainerID: strings.TrimSpace(parts[0]),
			Name:        strings.TrimSpace(parts[1]),
			Ports:       parsePortMappings(parts[2]),
		}

		services = append(services, service)
	}

	return services, nil
}

// parsePortMappings extracts port mappings from the docker ps ports column
// Format examples:
// - "0.0.0.0:8080->80/tcp"
// - "0.0.0.0:8080->80/tcp, 0.0.0.0:8443->443/tcp"
// - ":::8080->80/tcp"
func parsePortMappings(portsStr string) []PortMapping {
	if strings.TrimSpace(portsStr) == "" {
		return []PortMapping{}
	}

	mappings := []PortMapping{}

	// Split by comma for multiple port mappings
	portParts := strings.Split(portsStr, ",")

	// Regex to match port mappings: 0.0.0.0:8080->80/tcp or :::8080->80/tcp
	portRegex := regexp.MustCompile(`(?:0\.0\.0\.0|:::)?:?(\d+)->(\d+)/(tcp|udp)`)

	for _, part := range portParts {
		matches := portRegex.FindStringSubmatch(strings.TrimSpace(part))
		if len(matches) == 4 {
			hostPort, _ := strconv.Atoi(matches[1])
			containerPort, _ := strconv.Atoi(matches[2])
			protocol := matches[3]

			mappings = append(mappings, PortMapping{
				HostPort:      hostPort,
				ContainerPort: containerPort,
				Protocol:      protocol,
			})
		}
	}

	return mappings
}

// GetSupabaseStatus returns a map of Supabase service names to their ports
func GetSupabaseStatus() (map[string]int, error) {
	// Check if supabase CLI is available
	cmd := exec.Command("supabase", "--version")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("supabase CLI is not installed or not available in PATH")
	}

	// Execute supabase status command
	cmd = exec.Command("supabase", "status")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute supabase status: %w", err)
	}

	return parseSupabaseStatus(out.String())
}

// parseSupabaseStatus parses the output of supabase status command
// Expected format:
//
//	supabase local development setup is running.
//
//	        API URL: http://localhost:54321
//	   GraphQL URL: http://localhost:54321/graphql/v1
//	        DB URL: postgresql://postgres:postgres@localhost:54322/postgres
//	    Studio URL: http://localhost:54323
//	  Inbucket URL: http://localhost:54324
//	    JWT secret: ...
//	      anon key: ...
//
// service_role key: ...
func parseSupabaseStatus(output string) (map[string]int, error) {
	services := make(map[string]int)

	lines := strings.Split(output, "\n")

	// Regex to match service URLs with ports
	// Matches patterns like "API URL: http://localhost:54321" and "DB URL: postgresql://...@localhost:54322/..."
	// This regex looks for the last occurrence of :port before a slash or end of line
	urlRegex := regexp.MustCompile(`^\s*([^:]+):\s*\S+@?[^:@]+:(\d+)`)

	for _, line := range lines {
		matches := urlRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			serviceName := strings.TrimSpace(matches[1])
			port, err := strconv.Atoi(matches[2])
			if err != nil {
				continue
			}

			// Normalize service names
			serviceName = strings.ToLower(serviceName)
			serviceName = strings.ReplaceAll(serviceName, " ", "_")

			services[serviceName] = port
		}
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no supabase services found in status output")
	}

	return services, nil
}
