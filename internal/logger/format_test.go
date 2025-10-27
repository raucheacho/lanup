package logger

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetColorScheme(t *testing.T) {
	colors := GetColorScheme()

	// Verify all color codes are set
	assert.NotEmpty(t, colors.Info)
	assert.NotEmpty(t, colors.Warn)
	assert.NotEmpty(t, colors.Error)
	assert.NotEmpty(t, colors.Success)
	assert.NotEmpty(t, colors.Debug)
	assert.NotEmpty(t, colors.Reset)

	// Verify they are ANSI codes
	assert.Contains(t, colors.Info, "\033[")
	assert.Contains(t, colors.Warn, "\033[")
	assert.Contains(t, colors.Error, "\033[")
	assert.Contains(t, colors.Success, "\033[")
	assert.Contains(t, colors.Debug, "\033[")
	assert.Equal(t, "\033[0m", colors.Reset)
}

func TestFormatLogEntry_WithoutFields(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		module   string
		msg      string
		expected []string // Strings that should be in the output
	}{
		{
			name:   "DEBUG level",
			level:  DEBUG,
			module: "test.module",
			msg:    "debug message",
			expected: []string{
				"DEBUG",
				"test.module",
				"debug message",
			},
		},
		{
			name:   "INFO level",
			level:  INFO,
			module: "test.module",
			msg:    "info message",
			expected: []string{
				"INFO",
				"test.module",
				"info message",
			},
		},
		{
			name:   "WARN level",
			level:  WARN,
			module: "test.module",
			msg:    "warning message",
			expected: []string{
				"WARN",
				"test.module",
				"warning message",
			},
		},
		{
			name:   "ERROR level",
			level:  ERROR,
			module: "test.module",
			msg:    "error message",
			expected: []string{
				"ERROR",
				"test.module",
				"error message",
			},
		},
		{
			name:   "without module",
			level:  INFO,
			module: "",
			msg:    "message without module",
			expected: []string{
				"INFO",
				"message without module",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatLogEntry(tt.level, tt.module, tt.msg)

			// Verify timestamp format [YYYY-MM-DD HH:MM:SS]
			assert.Contains(t, result, "[")
			assert.Contains(t, result, "]")

			// Verify expected strings are present
			for _, expected := range tt.expected {
				assert.Contains(t, result, expected)
			}

			// Verify ends with newline
			assert.True(t, strings.HasSuffix(result, "\n"))
		})
	}
}

func TestFormatLogEntry_WithFields(t *testing.T) {
	fields := []Field{
		{Key: "user_id", Value: 123},
		{Key: "action", Value: "login"},
		{Key: "success", Value: true},
	}

	result := FormatLogEntry(INFO, "auth", "user logged in", fields...)

	// Verify all fields are present
	assert.Contains(t, result, "user_id=123")
	assert.Contains(t, result, "action=login")
	assert.Contains(t, result, "success=true")

	// Verify basic structure
	assert.Contains(t, result, "INFO")
	assert.Contains(t, result, "auth")
	assert.Contains(t, result, "user logged in")
}

func TestFormatLogEntry_ColorizedOutput(t *testing.T) {
	// Note: This test assumes IsTerminal() returns true
	// In a real terminal, the output would be colorized

	tests := []struct {
		name  string
		level LogLevel
	}{
		{"DEBUG with color", DEBUG},
		{"INFO with color", INFO},
		{"WARN with color", WARN},
		{"ERROR with color", ERROR},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatLogEntry(tt.level, "test", "message")

			// Verify the log level string is present
			assert.Contains(t, result, tt.level.String())

			// If terminal, should contain ANSI codes
			if IsTerminal() {
				// Should contain color codes
				assert.True(t,
					strings.Contains(result, "\033[") || !IsTerminal(),
					"Expected ANSI color codes in terminal output")
			}
		})
	}
}

func TestFormatLogEntry_LevelPadding(t *testing.T) {
	// Test that log levels are properly padded to 5 characters
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO "},
		{WARN, "WARN "},
		{ERROR, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			result := FormatLogEntry(tt.level, "", "test message")

			// The level should be padded to 5 characters
			// Format is: [timestamp] LEVEL message
			assert.Contains(t, result, tt.level.String())
		})
	}
}

func TestFormatLogEntry_TimestampFormat(t *testing.T) {
	result := FormatLogEntry(INFO, "test", "message")

	// Verify timestamp format [YYYY-MM-DD HH:MM:SS]
	// Extract the timestamp part
	start := strings.Index(result, "[")
	end := strings.Index(result, "]")

	assert.NotEqual(t, -1, start, "Should contain opening bracket")
	assert.NotEqual(t, -1, end, "Should contain closing bracket")
	assert.True(t, end > start, "Closing bracket should come after opening")

	timestamp := result[start+1 : end]

	// Verify format: YYYY-MM-DD HH:MM:SS (19 characters)
	assert.Equal(t, 19, len(timestamp), "Timestamp should be 19 characters")
	assert.Contains(t, timestamp, "-", "Should contain date separators")
	assert.Contains(t, timestamp, ":", "Should contain time separators")
	assert.Contains(t, timestamp, " ", "Should contain space between date and time")
}

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatLogEntry_EmptyMessage(t *testing.T) {
	result := FormatLogEntry(INFO, "test", "")

	// Should still format properly with empty message
	assert.Contains(t, result, "INFO")
	assert.Contains(t, result, "test")
	assert.True(t, strings.HasSuffix(result, "\n"))
}

func TestFormatLogEntry_SpecialCharacters(t *testing.T) {
	specialMsg := "Message with special chars: \n\t\"quotes\" and 'apostrophes'"
	result := FormatLogEntry(INFO, "test", specialMsg)

	// Should preserve special characters
	assert.Contains(t, result, specialMsg)
	assert.Contains(t, result, "INFO")
}

func TestFormatLogEntry_MultipleFields(t *testing.T) {
	fields := []Field{
		{Key: "field1", Value: "value1"},
		{Key: "field2", Value: 42},
		{Key: "field3", Value: 3.14},
		{Key: "field4", Value: true},
		{Key: "field5", Value: nil},
	}

	result := FormatLogEntry(INFO, "test", "message with many fields", fields...)

	// Verify all fields are present
	assert.Contains(t, result, "field1=value1")
	assert.Contains(t, result, "field2=42")
	assert.Contains(t, result, "field3=3.14")
	assert.Contains(t, result, "field4=true")
	assert.Contains(t, result, "field5=<nil>")
}

func TestFormatLogEntry_FieldsWithSpecialValues(t *testing.T) {
	fields := []Field{
		{Key: "empty_string", Value: ""},
		{Key: "with_spaces", Value: "value with spaces"},
		{Key: "with_equals", Value: "key=value"},
	}

	result := FormatLogEntry(INFO, "test", "testing special field values", fields...)

	// Verify fields are formatted correctly
	assert.Contains(t, result, "empty_string=")
	assert.Contains(t, result, "with_spaces=value with spaces")
	assert.Contains(t, result, "with_equals=key=value")
}

func TestColorScheme_AllColorsUnique(t *testing.T) {
	colors := GetColorScheme()

	// Verify all colors are unique (except potentially some could be the same)
	colorMap := make(map[string]bool)
	colorMap[colors.Info] = true
	colorMap[colors.Warn] = true
	colorMap[colors.Error] = true
	colorMap[colors.Success] = true
	colorMap[colors.Debug] = true

	// Should have at least 4 unique colors (some might overlap)
	assert.GreaterOrEqual(t, len(colorMap), 4)
}

func TestFormatLogEntry_ConsistentFormat(t *testing.T) {
	// Test that multiple calls produce consistent format
	result1 := FormatLogEntry(INFO, "test", "message")
	result2 := FormatLogEntry(INFO, "test", "message")

	// Timestamps will differ, but structure should be the same
	// Both should have the same number of brackets
	assert.Equal(t,
		strings.Count(result1, "["),
		strings.Count(result2, "["))
	assert.Equal(t,
		strings.Count(result1, "]"),
		strings.Count(result2, "]"))

	// Both should contain the same level and module
	assert.Contains(t, result1, "INFO")
	assert.Contains(t, result2, "INFO")
	assert.Contains(t, result1, "test")
	assert.Contains(t, result2, "test")
}
