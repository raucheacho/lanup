package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	lanuperrors "github.com/raucheacho/lanup/pkg/errors"
	"github.com/spf13/cobra"
)

// LogsCmd represents the logs command
type LogsCmd struct {
	Tail   int
	Follow bool
	Clear  bool
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View or manage lanup logs",
	Long: `View or manage lanup logs.

By default, displays all log entries. Use --tail to limit the number of lines,
--follow to stream logs in real-time, or --clear to remove the log file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tail, err := cmd.Flags().GetInt("tail")
		if err != nil {
			return fmt.Errorf("invalid tail value: %w", err)
		}

		follow, err := cmd.Flags().GetBool("follow")
		if err != nil {
			return fmt.Errorf("invalid follow value: %w", err)
		}

		clear, err := cmd.Flags().GetBool("clear")
		if err != nil {
			return fmt.Errorf("invalid clear value: %w", err)
		}

		logsCmd := &LogsCmd{
			Tail:   tail,
			Follow: follow,
			Clear:  clear,
		}

		return logsCmd.Run()
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)

	// Add flags
	logsCmd.Flags().IntP("tail", "n", 0, "show last N lines (0 = show all)")
	logsCmd.Flags().BoolP("follow", "f", false, "follow log output in real-time")
	logsCmd.Flags().Bool("clear", false, "clear the log file (requires confirmation)")
}

// Run executes the logs command
func (c *LogsCmd) Run() error {
	// Get log file path from global config
	config := GetGlobalConfig()
	if config == nil {
		return lanuperrors.NewError(lanuperrors.ErrInvalidConfig,
			"Global configuration not loaded", nil)
	}

	logPath := config.LogPath

	// Expand ~ in path if present
	if strings.HasPrefix(logPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return lanuperrors.NewError(lanuperrors.ErrFileNotFound,
				"Failed to get user home directory", err)
		}
		logPath = filepath.Join(home, logPath[1:])
	}

	// Handle clear flag
	if c.Clear {
		return c.clearLogs(logPath)
	}

	// Handle follow flag
	if c.Follow {
		return c.streamLogs(logPath)
	}

	// Default: display logs with optional tail
	return c.displayLogs(logPath)
}

// displayLogs reads and displays the log file
func (c *LogsCmd) displayLogs(logPath string) error {
	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		fmt.Println("No log file found. Logs will be created when lanup runs.")
		return nil
	}

	// Open log file
	file, err := os.Open(logPath)
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrPermissionDenied,
			"Failed to open log file", err)
	}
	defer file.Close()

	// If tail is specified, read last N lines
	if c.Tail > 0 {
		lines, err := readLastNLines(file, c.Tail)
		if err != nil {
			return lanuperrors.NewError(lanuperrors.ErrFileNotFound,
				"Failed to read log file", err)
		}
		for _, line := range lines {
			fmt.Print(line)
		}
		return nil
	}

	// Otherwise, read entire file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return lanuperrors.NewError(lanuperrors.ErrFileNotFound,
			"Error reading log file", err)
	}

	return nil
}

// streamLogs follows the log file and displays new entries in real-time
func (c *LogsCmd) streamLogs(logPath string) error {
	// Check if log file exists, if not wait for it
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		fmt.Println("Waiting for log file to be created...")
		// Wait for file to be created
		for {
			if _, err := os.Stat(logPath); err == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}

	// Open log file
	file, err := os.Open(logPath)
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrPermissionDenied,
			"Failed to open log file", err)
	}
	defer file.Close()

	// Seek to end of file
	if _, err := file.Seek(0, 2); err != nil {
		return lanuperrors.NewError(lanuperrors.ErrFileNotFound,
			"Failed to seek to end of file", err)
	}

	fmt.Println("Following log file (Ctrl+C to stop)...")

	// Create a scanner and continuously read new lines
	scanner := bufio.NewScanner(file)
	for {
		if scanner.Scan() {
			fmt.Println(scanner.Text())
		} else {
			// No new data, wait a bit
			time.Sleep(500 * time.Millisecond)
			// Check if file still exists (might have been rotated)
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				fmt.Println("Log file was removed or rotated. Exiting.")
				return nil
			}
		}

		if err := scanner.Err(); err != nil {
			return lanuperrors.NewError(lanuperrors.ErrFileNotFound,
				"Error reading log file", err)
		}
	}
}

// clearLogs removes the log file after confirmation
func (c *LogsCmd) clearLogs(logPath string) error {
	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		fmt.Println("No log file found.")
		return nil
	}

	// Ask for confirmation
	fmt.Print("Are you sure you want to clear the log file? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return lanuperrors.NewError(lanuperrors.ErrFileNotFound,
			"Failed to read confirmation", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// Remove the log file
	if err := os.Remove(logPath); err != nil {
		return lanuperrors.NewError(lanuperrors.ErrPermissionDenied,
			"Failed to remove log file", err)
	}

	fmt.Println("Log file cleared successfully.")
	return nil
}

// readLastNLines reads the last N lines from a file
func readLastNLines(file *os.File, n int) ([]string, error) {
	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := stat.Size()

	// Read file in chunks from the end
	const bufferSize = 4096
	var lines []string
	var buffer []byte
	offset := fileSize

	for offset > 0 && len(lines) < n {
		// Calculate how much to read
		readSize := int64(bufferSize)
		if offset < readSize {
			readSize = offset
		}
		offset -= readSize

		// Read chunk
		chunk := make([]byte, readSize)
		_, err := file.ReadAt(chunk, offset)
		if err != nil {
			return nil, err
		}

		// Prepend to buffer
		buffer = append(chunk, buffer...)

		// Split into lines
		text := string(buffer)
		allLines := strings.Split(text, "\n")

		// Keep the incomplete first line in buffer for next iteration
		if offset > 0 {
			buffer = []byte(allLines[0])
			allLines = allLines[1:]
		} else {
			buffer = nil
		}

		// Prepend lines (we're reading backwards)
		lines = append(allLines, lines...)
	}

	// Remove empty lines and trim to requested count
	var result []string
	for i := len(lines) - 1; i >= 0 && len(result) < n; i-- {
		if lines[i] != "" {
			result = append([]string{lines[i] + "\n"}, result...)
		}
	}

	return result, nil
}
