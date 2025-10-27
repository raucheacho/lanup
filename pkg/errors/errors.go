package errors

import "fmt"

// ErrorCode represents specific error types in lanup
type ErrorCode int

const (
	// ErrNoNetwork indicates no active network interface was found
	ErrNoNetwork ErrorCode = iota + 1
	// ErrInvalidConfig indicates configuration file is malformed or invalid
	ErrInvalidConfig
	// ErrFileNotFound indicates a required file was not found
	ErrFileNotFound
	// ErrPermissionDenied indicates insufficient permissions for an operation
	ErrPermissionDenied
	// ErrInvalidURL indicates the provided URL is invalid or malformed
	ErrInvalidURL
	// ErrDockerUnavailable indicates Docker is not available or not running
	ErrDockerUnavailable
)

// LanupError represents a structured error with code, message, and cause
type LanupError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

// Error implements the error interface
func (e *LanupError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// NewError creates a new LanupError with the given code, message, and optional cause
func NewError(code ErrorCode, msg string, cause error) *LanupError {
	return &LanupError{
		Code:    code,
		Message: msg,
		Cause:   cause,
	}
}

// ExitCode returns the appropriate exit code for the error
func (e *LanupError) ExitCode() int {
	switch e.Code {
	case ErrNoNetwork:
		return 3
	case ErrInvalidConfig:
		return 2
	case ErrFileNotFound:
		return 1
	case ErrPermissionDenied:
		return 4
	case ErrInvalidURL:
		return 5
	case ErrDockerUnavailable:
		return 1
	default:
		return 1
	}
}
