package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Output utilities for console formatting with colors and emojis

var (
	// Color functions
	successColor   = color.New(color.FgGreen, color.Bold)
	infoColor      = color.New(color.FgBlue)
	warningColor   = color.New(color.FgYellow, color.Bold)
	errorColor     = color.New(color.FgRed, color.Bold)
	highlightColor = color.New(color.FgCyan, color.Bold)
)

// Success prints a success message with green color and checkmark emoji
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if isTerminal() {
		successColor.Printf("✅ %s\n", msg)
	} else {
		fmt.Printf("[SUCCESS] %s\n", msg)
	}
}

// Info prints an informational message with blue color and info emoji
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if isTerminal() {
		infoColor.Printf("ℹ️  %s\n", msg)
	} else {
		fmt.Printf("[INFO] %s\n", msg)
	}
}

// Warning prints a warning message with yellow color and warning emoji
func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if isTerminal() {
		warningColor.Printf("⚠️  %s\n", msg)
	} else {
		fmt.Printf("[WARNING] %s\n", msg)
	}
}

// Error prints an error message with red color and error emoji
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if isTerminal() {
		errorColor.Fprintf(os.Stderr, "❌ %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", msg)
	}
}

// Highlight prints a highlighted message with cyan color
func Highlight(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if isTerminal() {
		highlightColor.Printf("🔗 %s\n", msg)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

// PrintURL prints a URL with special formatting
func PrintURL(name, url string) {
	if isTerminal() {
		fmt.Printf("  %s %s\n",
			color.New(color.FgCyan, color.Bold).Sprint(name+":"),
			color.New(color.FgWhite, color.Underline).Sprint(url))
	} else {
		fmt.Printf("  %s %s\n", name+":", url)
	}
}

// PrintSection prints a section header
func PrintSection(title string) {
	if isTerminal() {
		fmt.Println()
		color.New(color.FgMagenta, color.Bold).Printf("═══ %s ═══\n", title)
		fmt.Println()
	} else {
		fmt.Printf("\n=== %s ===\n\n", title)
	}
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
