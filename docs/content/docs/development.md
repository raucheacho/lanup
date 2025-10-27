---
title: "Development Guide"
weight: 7
---

# Development Guide

Guide for contributing to lanup development.

## Prerequisites

- Go 1.22 or higher
- Git
- Make (optional, but recommended)

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/raucheacho/lanup.git
cd lanup
```

### Install Dependencies

```bash
go mod download
```

### Build

```bash
# Build for current platform
make build

# Or use go directly
go build -o lanup .
```

### Run Tests

```bash
# Run all tests
make test

# Or use go directly
go test ./...

# Run with coverage
make test-coverage
```

## Project Structure

```
lanup/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command
│   ├── init.go            # Init command
│   ├── start.go           # Start command
│   ├── expose.go          # Expose command
│   ├── logs.go            # Logs command
│   └── doctor.go          # Doctor command
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   ├── net/               # Network detection
│   ├── env/               # Environment file management
│   ├── logger/            # Logging system
│   └── docker/            # Docker integration
├── pkg/                   # Public packages
│   ├── errors/            # Error handling
│   └── utils/             # Utility functions
├── docs/                  # Hugo documentation
├── main.go                # Entry point
├── Makefile               # Build scripts
└── README.md              # Project README
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/my-feature
```

### 2. Make Changes

Follow Go best practices:
- Use `gofmt` for formatting
- Add comments for exported functions
- Write tests for new features
- Keep functions small and focused

### 3. Run Tests

```bash
make test
```

### 4. Format Code

```bash
make fmt
```

### 5. Commit Changes

```bash
git add .
git commit -m "Add my feature"
```

### 6. Push and Create PR

```bash
git push origin feature/my-feature
```

Then create a Pull Request on GitHub.

## Building

### Build for Current Platform

```bash
make build
```

### Build for All Platforms

```bash
make build-all
```

This creates binaries for:
- macOS (Intel and Apple Silicon)
- Linux (amd64)
- Windows (amd64)

### Install Locally

```bash
make install
```

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run specific package
go test ./internal/net/...

# Run with verbose output
go test -v ./...
```

### Integration Tests

```bash
# Run integration tests
go test ./cmd/... -v
```

### Coverage

```bash
# Generate coverage report
make test-coverage

# View coverage in browser
go tool cover -html=coverage.out
```

## Code Style

### Formatting

Use `gofmt` for consistent formatting:

```bash
gofmt -w .
```

### Linting

Use `golangci-lint` for linting:

```bash
golangci-lint run
```

### Comments

Add comments for exported functions:

```go
// DetectLocalIP detects the local IP address of the machine.
// It returns the IP address, interface name, and interface type.
func DetectLocalIP() (*NetworkInfo, error) {
    // Implementation
}
```

## Adding New Commands

### 1. Create Command File

```go
// cmd/mycommand.go
package cmd

import (
    "github.com/spf13/cobra"
)

func NewMyCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "mycommand",
        Short: "Short description",
        Long:  "Long description",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation
            return nil
        },
    }
    
    return cmd
}

func init() {
    RootCmd.AddCommand(NewMyCommand())
}
```

### 2. Add Tests

```go
// cmd/mycommand_test.go
package cmd

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestMyCommand(t *testing.T) {
    // Test implementation
}
```

### 3. Update Documentation

Add documentation in `docs/content/docs/commands.md`.

## Adding New Features

### 1. Plan the Feature

- Create an issue on GitHub
- Discuss the approach
- Get feedback from maintainers

### 2. Implement

- Write code following project conventions
- Add tests
- Update documentation

### 3. Submit PR

- Create a pull request
- Describe the changes
- Link to related issues

## Release Process

### 1. Update Version

```bash
# Update version in main.go
version = "1.2.0"
```

### 2. Create Tag

```bash
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

### 3. Build Releases

```bash
make build-all
```

### 4. Create GitHub Release

- Go to GitHub Releases
- Create new release
- Upload binaries
- Write release notes

## Documentation

### Building Documentation

```bash
cd docs
hugo server
```

Visit `http://localhost:1313` to view the documentation.

### Adding Documentation

1. Create new markdown file in `docs/content/docs/`
2. Add frontmatter with title and weight
3. Write content in markdown
4. Test locally with `hugo server`

## Contributing Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback

### Pull Request Guidelines

- One feature per PR
- Include tests
- Update documentation
- Follow code style
- Write clear commit messages

### Commit Message Format

```
type: subject

body

footer
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `test`: Tests
- `refactor`: Code refactoring
- `chore`: Maintenance

Example:
```
feat: add watch mode for network changes

Implement watch mode that monitors network interfaces
and automatically updates environment files when the
IP address changes.

Closes #123
```

## Getting Help

- [GitHub Issues](https://github.com/raucheacho/lanup/issues)
- [GitHub Discussions](https://github.com/raucheacho/lanup/discussions)

## License

MIT License - see [LICENSE](https://github.com/raucheacho/lanup/blob/main/LICENSE)
