package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Field represents a structured log field
type Field struct {
	Key   string
	Value interface{}
}

// Logger provides structured logging with rotation support
type Logger struct {
	Level      LogLevel
	FilePath   string
	MaxSize    int64 // bytes
	MaxBackups int
	Console    bool
	Colors     bool
	mu         sync.Mutex
	file       *os.File
	size       int64
}

// LoggerConfig holds configuration for creating a new logger
type LoggerConfig struct {
	Level      LogLevel
	FilePath   string
	MaxSize    int64
	MaxBackups int
	Console    bool
	Colors     bool
}

// NewLogger creates a new logger instance with the given configuration
func NewLogger(config LoggerConfig) (*Logger, error) {
	// Set defaults
	if config.MaxSize == 0 {
		config.MaxSize = 5 * 1024 * 1024 // 5MB
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 5
	}

	logger := &Logger{
		Level:      config.Level,
		FilePath:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		Console:    config.Console,
		Colors:     config.Colors,
	}

	// Create log directory if it doesn't exist
	if config.FilePath != "" {
		dir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open or create log file
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logger.file = file

		// Get current file size
		info, err := file.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to stat log file: %w", err)
		}
		logger.size = info.Size()
	}

	return logger, nil
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(DEBUG, msg, fields...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...Field) {
	l.log(INFO, msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(WARN, msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...Field) {
	l.log(ERROR, msg, fields...)
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, msg string, fields ...Field) {
	// Check if we should log this level
	if level < l.Level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Format the log entry
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("[%s] %-5s %s", timestamp, level.String(), msg)

	// Add fields if present
	if len(fields) > 0 {
		for _, field := range fields {
			entry += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}
	entry += "\n"

	// Write to file if configured
	if l.file != nil {
		n, err := l.file.WriteString(entry)
		if err != nil {
			// If we can't write to the log file, write to stderr
			fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
		} else {
			l.size += int64(n)

			// Check if rotation is needed
			if l.size >= l.MaxSize {
				if err := l.rotate(); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to rotate log file: %v\n", err)
				}
			}
		}
	}

	// Write to console if configured
	if l.Console {
		var output io.Writer = os.Stdout
		if level == ERROR {
			output = os.Stderr
		}

		// Use colored output if enabled
		if l.Colors && IsTerminal() {
			entry = FormatLogEntry(level, "", msg, fields...)
		}

		fmt.Fprint(output, entry)
	}
}

// rotate performs log rotation
func (l *Logger) rotate() error {
	// Close current file
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("failed to close log file: %w", err)
		}
	}

	// Rotate existing backup files
	for i := l.MaxBackups - 1; i >= 1; i-- {
		oldPath := fmt.Sprintf("%s.%d", l.FilePath, i)
		newPath := fmt.Sprintf("%s.%d", l.FilePath, i+1)

		// Check if old backup exists
		if _, err := os.Stat(oldPath); err == nil {
			// Remove the oldest backup if it exists
			if i == l.MaxBackups-1 {
				os.Remove(newPath)
			}
			// Rename the backup
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("failed to rotate backup %d: %w", i, err)
			}
		}
	}

	// Rename current log file to .1
	backupPath := fmt.Sprintf("%s.1", l.FilePath)
	if err := os.Rename(l.FilePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	// Clean up old backups beyond MaxBackups
	l.cleanupOldBackups()

	// Create new log file
	file, err := os.OpenFile(l.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	l.file = file
	l.size = 0

	return nil
}

// cleanupOldBackups removes backup files beyond MaxBackups
func (l *Logger) cleanupOldBackups() {
	dir := filepath.Dir(l.FilePath)
	base := filepath.Base(l.FilePath)

	// Find all backup files
	pattern := filepath.Join(dir, base+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	// Sort by modification time (oldest first)
	type fileInfo struct {
		path    string
		modTime time.Time
	}
	var files []fileInfo
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		files = append(files, fileInfo{path: match, modTime: info.ModTime()})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.Before(files[j].modTime)
	})

	// Remove oldest files if we exceed MaxBackups
	if len(files) > l.MaxBackups {
		for i := 0; i < len(files)-l.MaxBackups; i++ {
			os.Remove(files[i].path)
		}
	}
}
