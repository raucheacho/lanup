package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/fatih/color"
	"github.com/raucheacho/lanup/internal/net"
	lanuperrors "github.com/raucheacho/lanup/pkg/errors"
	"github.com/spf13/cobra"
)

// ExposeCmd represents the expose command
type ExposeCmd struct {
	URL   string
	Name  string
	Port  int
	HTTPS bool
}

// NewExposeCmd creates a new expose command
func NewExposeCmd() *cobra.Command {
	exposeCmd := &ExposeCmd{}

	cmd := &cobra.Command{
		Use:   "expose [URL]",
		Short: "Quickly expose a single service without configuration",
		Long: `Expose a single localhost URL on your local network without creating a configuration file.

This command detects your local IP address and transforms a localhost URL to be accessible
from any device on your network.

Examples:
  lanup expose http://localhost:3000
  lanup expose http://localhost:8080 --name api
  lanup expose http://localhost:5000 --port 8000
  lanup expose http://localhost:3000 --https`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			exposeCmd.URL = args[0]
			return exposeCmd.Run()
		},
	}

	// Add flags
	cmd.Flags().StringVar(&exposeCmd.Name, "name", "", "assign an alias to the exposed service")
	cmd.Flags().IntVar(&exposeCmd.Port, "port", 0, "use a custom port instead of the original")
	cmd.Flags().BoolVar(&exposeCmd.HTTPS, "https", false, "use HTTPS protocol instead of HTTP")

	return cmd
}

func init() {
	RootCmd.AddCommand(NewExposeCmd())
}

// Run executes the expose command
func (c *ExposeCmd) Run() error {
	// Validate the URL
	if err := c.validateURL(); err != nil {
		return err
	}

	// Detect local IP
	netInfo, err := net.DetectLocalIP()
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrNoNetwork,
			"Failed to detect local IP address", err)
	}

	// Transform the URL
	transformedURL, err := c.transformURL(netInfo.IP)
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrInvalidURL,
			"Failed to transform URL", err)
	}

	// Display the result
	c.displayResult(netInfo.IP, transformedURL)

	return nil
}

// validateURL checks if the URL is valid and uses localhost or 127.0.0.1
func (c *ExposeCmd) validateURL() error {
	// Parse the URL
	parsedURL, err := url.Parse(c.URL)
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrInvalidURL, "Invalid URL format", err)
	}

	// Check if scheme is present
	if parsedURL.Scheme == "" {
		return lanuperrors.NewError(lanuperrors.ErrInvalidURL,
			"Invalid URL: missing protocol (http:// or https://)", nil)
	}

	// Check if scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return lanuperrors.NewError(lanuperrors.ErrInvalidURL,
			fmt.Sprintf("Invalid URL: protocol must be http or https, got %s", parsedURL.Scheme), nil)
	}

	// Extract hostname
	hostname := parsedURL.Hostname()
	if hostname == "" {
		return lanuperrors.NewError(lanuperrors.ErrInvalidURL, "Invalid URL: missing hostname", nil)
	}

	// Check if hostname is localhost or 127.0.0.1
	if hostname != "localhost" && hostname != "127.0.0.1" {
		return lanuperrors.NewError(lanuperrors.ErrInvalidURL,
			fmt.Sprintf("Invalid URL: hostname must be localhost or 127.0.0.1, got %s", hostname), nil)
	}

	return nil
}

// transformURL replaces localhost/127.0.0.1 with the detected IP and applies custom settings
func (c *ExposeCmd) transformURL(localIP string) (string, error) {
	// Parse the original URL
	parsedURL, err := url.Parse(c.URL)
	if err != nil {
		return "", err
	}

	// Replace hostname with local IP
	parsedURL.Host = strings.Replace(parsedURL.Host, "localhost", localIP, 1)
	parsedURL.Host = strings.Replace(parsedURL.Host, "127.0.0.1", localIP, 1)

	// Apply custom port if specified
	if c.Port > 0 {
		parsedURL.Host = fmt.Sprintf("%s:%d", localIP, c.Port)
	}

	// Apply HTTPS if specified
	if c.HTTPS {
		parsedURL.Scheme = "https"
	}

	return parsedURL.String(), nil
}

// displayResult shows the transformed URL in a user-friendly format
func (c *ExposeCmd) displayResult(localIP, transformedURL string) {
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Printf("%s %s\n", green("✓"), "Successfully exposed service on your LAN!")
	fmt.Printf("%s %s\n\n", green("✓"), "Local IP: "+cyan(localIP))

	if c.Name != "" {
		fmt.Printf("%s %s\n", yellow("📌"), "Service name: "+bold(c.Name))
	}

	fmt.Printf("%s %s\n", yellow("🌐"), "Original URL:")
	fmt.Printf("  %s\n\n", c.URL)

	fmt.Printf("%s %s\n", yellow("🌐"), "Network URL:")
	fmt.Printf("  %s\n\n", cyan(transformedURL))

	fmt.Println("💡 Tip: Use 'lanup init' to configure multiple services in your project")
}
