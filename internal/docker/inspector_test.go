package docker

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDockerPS_EmptyOutput(t *testing.T) {
	output := ""
	services, err := ParseDockerPS(output)

	require.NoError(t, err)
	assert.Empty(t, services)
}

func TestParseDockerPS_SingleContainer(t *testing.T) {
	output := "abc123|my-container|0.0.0.0:8080->80/tcp"

	services, err := ParseDockerPS(output)

	require.NoError(t, err)
	assert.Len(t, services, 1)

	service := services[0]
	assert.Equal(t, "abc123", service.ContainerID)
	assert.Equal(t, "my-container", service.Name)
	assert.Len(t, service.Ports, 1)

	port := service.Ports[0]
	assert.Equal(t, 8080, port.HostPort)
	assert.Equal(t, 80, port.ContainerPort)
	assert.Equal(t, "tcp", port.Protocol)
}

func TestParseDockerPS_MultipleContainers(t *testing.T) {
	output := `abc123|web-server|0.0.0.0:8080->80/tcp
def456|database|0.0.0.0:5432->5432/tcp
ghi789|redis|0.0.0.0:6379->6379/tcp`

	services, err := ParseDockerPS(output)

	require.NoError(t, err)
	assert.Len(t, services, 3)

	// Verify first container
	assert.Equal(t, "abc123", services[0].ContainerID)
	assert.Equal(t, "web-server", services[0].Name)
	assert.Len(t, services[0].Ports, 1)
	assert.Equal(t, 8080, services[0].Ports[0].HostPort)

	// Verify second container
	assert.Equal(t, "def456", services[1].ContainerID)
	assert.Equal(t, "database", services[1].Name)
	assert.Len(t, services[1].Ports, 1)
	assert.Equal(t, 5432, services[1].Ports[0].HostPort)

	// Verify third container
	assert.Equal(t, "ghi789", services[2].ContainerID)
	assert.Equal(t, "redis", services[2].Name)
	assert.Len(t, services[2].Ports, 1)
	assert.Equal(t, 6379, services[2].Ports[0].HostPort)
}

func TestParseDockerPS_MultiplePorts(t *testing.T) {
	output := "abc123|web-server|0.0.0.0:8080->80/tcp, 0.0.0.0:8443->443/tcp"

	services, err := ParseDockerPS(output)

	require.NoError(t, err)
	assert.Len(t, services, 1)

	service := services[0]
	assert.Equal(t, "abc123", service.ContainerID)
	assert.Equal(t, "web-server", service.Name)
	assert.Len(t, service.Ports, 2)

	// Verify first port
	assert.Equal(t, 8080, service.Ports[0].HostPort)
	assert.Equal(t, 80, service.Ports[0].ContainerPort)
	assert.Equal(t, "tcp", service.Ports[0].Protocol)

	// Verify second port
	assert.Equal(t, 8443, service.Ports[1].HostPort)
	assert.Equal(t, 443, service.Ports[1].ContainerPort)
	assert.Equal(t, "tcp", service.Ports[1].Protocol)
}

func TestParseDockerPS_IPv6Format(t *testing.T) {
	output := "abc123|my-container|:::8080->80/tcp"

	services, err := ParseDockerPS(output)

	require.NoError(t, err)
	assert.Len(t, services, 1)

	service := services[0]
	assert.Equal(t, "abc123", service.ContainerID)
	assert.Equal(t, "my-container", service.Name)
	assert.Len(t, service.Ports, 1)

	port := service.Ports[0]
	assert.Equal(t, 8080, port.HostPort)
	assert.Equal(t, 80, port.ContainerPort)
	assert.Equal(t, "tcp", port.Protocol)
}

func TestParseDockerPS_UDPProtocol(t *testing.T) {
	output := "abc123|dns-server|0.0.0.0:53->53/udp"

	services, err := ParseDockerPS(output)

	require.NoError(t, err)
	assert.Len(t, services, 1)

	service := services[0]
	assert.Len(t, service.Ports, 1)

	port := service.Ports[0]
	assert.Equal(t, 53, port.HostPort)
	assert.Equal(t, 53, port.ContainerPort)
	assert.Equal(t, "udp", port.Protocol)
}

func TestParseDockerPS_NoPorts(t *testing.T) {
	output := "abc123|my-container|"

	services, err := ParseDockerPS(output)

	require.NoError(t, err)
	assert.Len(t, services, 1)

	service := services[0]
	assert.Equal(t, "abc123", service.ContainerID)
	assert.Equal(t, "my-container", service.Name)
	assert.Empty(t, service.Ports)
}

func TestParseDockerPS_RealWorldExample(t *testing.T) {
	// Real-world example from docker ps output
	output := `a1b2c3d4e5f6|supabase-db|0.0.0.0:54322->5432/tcp
b2c3d4e5f6a1|supabase-studio|0.0.0.0:54323->3000/tcp
c3d4e5f6a1b2|supabase-kong|0.0.0.0:54321->8000/tcp, 0.0.0.0:54320->8443/tcp
d4e5f6a1b2c3|supabase-auth|9999/tcp
e5f6a1b2c3d4|supabase-rest|3000/tcp`

	services, err := ParseDockerPS(output)

	require.NoError(t, err)
	assert.Len(t, services, 5)

	// Verify supabase-db
	assert.Equal(t, "supabase-db", services[0].Name)
	assert.Len(t, services[0].Ports, 1)
	assert.Equal(t, 54322, services[0].Ports[0].HostPort)
	assert.Equal(t, 5432, services[0].Ports[0].ContainerPort)

	// Verify supabase-studio
	assert.Equal(t, "supabase-studio", services[1].Name)
	assert.Len(t, services[1].Ports, 1)
	assert.Equal(t, 54323, services[1].Ports[0].HostPort)

	// Verify supabase-kong (multiple ports)
	assert.Equal(t, "supabase-kong", services[2].Name)
	assert.Len(t, services[2].Ports, 2)
	assert.Equal(t, 54321, services[2].Ports[0].HostPort)
	assert.Equal(t, 54320, services[2].Ports[1].HostPort)

	// Verify containers with no host port mapping
	assert.Equal(t, "supabase-auth", services[3].Name)
	assert.Empty(t, services[3].Ports)

	assert.Equal(t, "supabase-rest", services[4].Name)
	assert.Empty(t, services[4].Ports)
}

func TestParseDockerPS_ComplexContainerNames(t *testing.T) {
	output := `abc123|my-app-web-1|0.0.0.0:8080->80/tcp
def456|project_database_1|0.0.0.0:5432->5432/tcp
ghi789|test-redis-cache|0.0.0.0:6379->6379/tcp`

	services, err := ParseDockerPS(output)

	require.NoError(t, err)
	assert.Len(t, services, 3)

	assert.Equal(t, "my-app-web-1", services[0].Name)
	assert.Equal(t, "project_database_1", services[1].Name)
	assert.Equal(t, "test-redis-cache", services[2].Name)
}

func TestParseSupabaseStatus_Success(t *testing.T) {
	output := `supabase local development setup is running.

         API URL: http://localhost:54321
      GraphQL URL: http://localhost:54321/graphql/v1
         DB URL: postgresql://postgres:postgres@localhost:54322/postgres
     Studio URL: http://localhost:54323
   Inbucket URL: http://localhost:54324
     JWT secret: super-secret-jwt-token-with-at-least-32-characters-long
       anon key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
service_role key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`

	services, err := parseSupabaseStatus(output)

	require.NoError(t, err)
	assert.NotEmpty(t, services)

	// Verify API URL port
	apiPort, exists := services["api_url"]
	assert.True(t, exists)
	assert.Equal(t, 54321, apiPort)

	// Verify GraphQL URL port
	graphqlPort, exists := services["graphql_url"]
	assert.True(t, exists)
	assert.Equal(t, 54321, graphqlPort)

	// Verify DB URL port
	dbPort, exists := services["db_url"]
	assert.True(t, exists)
	assert.Equal(t, 54322, dbPort)

	// Verify Studio URL port
	studioPort, exists := services["studio_url"]
	assert.True(t, exists)
	assert.Equal(t, 54323, studioPort)

	// Verify Inbucket URL port
	inbucketPort, exists := services["inbucket_url"]
	assert.True(t, exists)
	assert.Equal(t, 54324, inbucketPort)
}

func TestParseSupabaseStatus_NoServices(t *testing.T) {
	output := `supabase local development setup is not running.

Run 'supabase start' to start the local development setup.`

	_, err := parseSupabaseStatus(output)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no supabase services found")
}

func TestParseSupabaseStatus_PartialOutput(t *testing.T) {
	output := `supabase local development setup is running.

         API URL: http://localhost:54321
     Studio URL: http://localhost:54323`

	services, err := parseSupabaseStatus(output)

	require.NoError(t, err)
	assert.Len(t, services, 2)

	apiPort, exists := services["api_url"]
	assert.True(t, exists)
	assert.Equal(t, 54321, apiPort)

	studioPort, exists := services["studio_url"]
	assert.True(t, exists)
	assert.Equal(t, 54323, studioPort)
}

func TestParseSupabaseStatus_CustomPorts(t *testing.T) {
	output := `supabase local development setup is running.

         API URL: http://localhost:8000
         DB URL: postgresql://postgres:postgres@localhost:5432/postgres
     Studio URL: http://localhost:3000`

	services, err := parseSupabaseStatus(output)

	require.NoError(t, err)

	apiPort, exists := services["api_url"]
	assert.True(t, exists)
	assert.Equal(t, 8000, apiPort)

	dbPort, exists := services["db_url"]
	assert.True(t, exists)
	assert.Equal(t, 5432, dbPort)

	studioPort, exists := services["studio_url"]
	assert.True(t, exists)
	assert.Equal(t, 3000, studioPort)
}

func TestParsePortMappings_VariousFormats(t *testing.T) {
	tests := []struct {
		name     string
		portsStr string
		expected []PortMapping
	}{
		{
			name:     "empty string",
			portsStr: "",
			expected: []PortMapping{},
		},
		{
			name:     "single port",
			portsStr: "0.0.0.0:8080->80/tcp",
			expected: []PortMapping{
				{HostPort: 8080, ContainerPort: 80, Protocol: "tcp"},
			},
		},
		{
			name:     "multiple ports",
			portsStr: "0.0.0.0:8080->80/tcp, 0.0.0.0:8443->443/tcp",
			expected: []PortMapping{
				{HostPort: 8080, ContainerPort: 80, Protocol: "tcp"},
				{HostPort: 8443, ContainerPort: 443, Protocol: "tcp"},
			},
		},
		{
			name:     "ipv6 format",
			portsStr: ":::8080->80/tcp",
			expected: []PortMapping{
				{HostPort: 8080, ContainerPort: 80, Protocol: "tcp"},
			},
		},
		{
			name:     "udp protocol",
			portsStr: "0.0.0.0:53->53/udp",
			expected: []PortMapping{
				{HostPort: 53, ContainerPort: 53, Protocol: "udp"},
			},
		},
		{
			name:     "no host binding",
			portsStr: "8080/tcp",
			expected: []PortMapping{},
		},
		{
			name:     "mixed formats",
			portsStr: "0.0.0.0:8080->80/tcp, 9000/tcp, :::8443->443/tcp",
			expected: []PortMapping{
				{HostPort: 8080, ContainerPort: 80, Protocol: "tcp"},
				{HostPort: 8443, ContainerPort: 443, Protocol: "tcp"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePortMappings(tt.portsStr)
			assert.Equal(t, len(tt.expected), len(result))

			for i, expected := range tt.expected {
				assert.Equal(t, expected.HostPort, result[i].HostPort)
				assert.Equal(t, expected.ContainerPort, result[i].ContainerPort)
				assert.Equal(t, expected.Protocol, result[i].Protocol)
			}
		})
	}
}

func TestGetRunningContainers_DockerUnavailable(t *testing.T) {
	// This test will only pass if Docker is not available
	// Skip if Docker is available
	if IsDockerAvailable() {
		t.Skip("Skipping test because Docker is available")
	}

	containers, err := GetRunningContainers()
	assert.Error(t, err)
	assert.Nil(t, containers)
	assert.Contains(t, err.Error(), "docker is not available")
}

func TestGetSupabaseStatus_SupabaseUnavailable(t *testing.T) {
	// This test verifies graceful degradation when Supabase CLI is not available or not running
	services, err := GetSupabaseStatus()

	// If Supabase is not installed or not running, should return error
	if err != nil {
		// Error message could be either "not installed" or "failed to execute" (when not running)
		assert.True(t,
			strings.Contains(err.Error(), "supabase CLI is not installed") ||
				strings.Contains(err.Error(), "failed to execute supabase status") ||
				strings.Contains(err.Error(), "no supabase services found"),
			"Expected error about Supabase unavailability, got: %v", err)
		assert.Nil(t, services)
	} else {
		// If Supabase is installed and running, should return services
		assert.NotNil(t, services)
	}
}
