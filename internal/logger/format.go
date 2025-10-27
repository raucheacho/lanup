package logger

import (
	"fmt"
	"os"
	"time"
)

// ColorScheme defines ANSI color codes for different log levels
type ColorScheme struct {
	Info    string
	Warn    string
	Error   string
	Success string
	Debug   string
	Reset   string
}

// GetColorScheme returns the color scheme with ANSI color codes
func GetColorScheme() ColorScheme {
	return ColorScheme{
		Info:    "\033[34m", // Blue
		Warn:    "\033[33m", // Yellow
		Error:   "\033[31m", // Red
		Success: "\033[32m", // Green
		Debug:   "\033[36m", // Cyan
		Reset:   "\033[0m",  // Reset
	}
}

// FormatLogEntry formats a log entry with timestamp, level, and optional colorization
func FormatLogEntry(level LogLevel, module string, msg string, fields ...Field) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var entry string

	// Add color if terminal supports it
	if IsTerminal() {
		colors := GetColorScheme()
		var color string

		switch level {
		case DEBUG:
			color = colors.Debug
		case INFO:
			color = colors.Info
		case WARN:
			color = colors.Warn
		case ERROR:
			color = colors.Error
		default:
			color = colors.Reset
		}

		// Format with color
		if module != "" {
			entry = fmt.Sprintf("[%s] %s%-5s%s %s: %s",
				timestamp, color, level.String(), colors.Reset, module, msg)
		} else {
			entry = fmt.Sprintf("[%s] %s%-5s%s %s",
				timestamp, color, level.String(), colors.Reset, msg)
		}
	} else {
		// Format without color
		if module != "" {
			entry = fmt.Sprintf("[%s] %-5s %s: %s",
				timestamp, level.String(), module, msg)
		} else {
			entry = fmt.Sprintf("[%s] %-5s %s",
				timestamp, level.String(), msg)
		}
	}

	// Add fields if present
	if len(fields) > 0 {
		for _, field := range fields {
			entry += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	entry += "\n"
	return entry
}

// IsTerminal checks if the output is a terminal (TTY)
// This is used to determine whether to use colored output
func IsTerminal() bool {
	// Check if stdout is a terminal
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	// Check if it's a character device (terminal)
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
