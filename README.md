# lanup

**lanup** is a CLI tool that automatically exposes your local backend services on your local network (LAN). It detects your local IP address, updates environment variables, and allows you to test your applications from any device on the same network without manual configuration.

Perfect for testing mobile apps, sharing work with teammates, or accessing your development environment from multiple devices.

## Features

- 🌐 **Automatic IP Detection** - Detects your local network IP address
- 🔄 **Watch Mode** - Automatically updates when your network changes
- 🐳 **Docker Integration** - Auto-detects running Docker containers
- 🗄️ **Supabase Support** - Detects Supabase local development services
- 📝 **Environment File Management** - Safely updates .env files with backups
- 🎨 **Colorized Output** - Beautiful, easy-to-read console output
- 🔍 **Diagnostics** - Built-in health checks for troubleshooting
- 🔒 **Secure by Default** - Only uses private IP addresses (RFC 1918)

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap raucheacho/tap
brew install lanup
```

### npm (All platforms)

```bash
npm install -g lanup
```

### Scoop (Windows)

```powershell
scoop bucket add raucheacho https://github.com/raucheacho/scoop-bucket
scoop install lanup
```

### Go Install

If you have Go 1.22+ installed:

```bash
go install github.com/raucheacho/lanup@latest
```

### Pre-built Binaries

Download from the [releases page](https://github.com/raucheacho/lanup/releases).

#### macOS

```bash
# Intel Mac
curl -L https://github.com/raucheacho/lanup/releases/latest/download/lanup_Darwin_x86_64.tar.gz | tar xz
sudo mv lanup /usr/local/bin/

# Apple Silicon (M1/M2)
curl -L https://github.com/raucheacho/lanup/releases/latest/download/lanup_Darwin_arm64.tar.gz | tar xz
sudo mv lanup /usr/local/bin/
```

#### Linux

```bash
curl -L https://github.com/raucheacho/lanup/releases/latest/download/lanup_Linux_x86_64.tar.gz | tar xz
sudo mv lanup /usr/local/bin/
```

#### Windows

Download `lanup_Windows_x86_64.zip` from the releases page and add it to your PATH.

## Quick Start

1. **Initialize lanup in your project**

```bash
cd your-project
lanup init
```

This creates a `.lanup.yaml` configuration file.

2. **Edit the configuration**

```yaml
vars:
  SUPABASE_URL: "http://localhost:54321"
  API_URL: "http://localhost:8000"
  DASHBOARD_URL: "http://localhost:3000"

output: ".env.local"

auto_detect:
  docker: true
  supabase: true
```

3. **Start exposing your services**

```bash
lanup start
```

Your services are now accessible on your local network! 🎉

## Usage

### Commands

#### `lanup init`

Initialize lanup configuration in your project.

```bash
lanup init [flags]
```

**Flags:**

- `--format string` - Configuration file format (yaml or toml) (default "yaml")
- `--force` - Overwrite existing configuration file

**Examples:**

```bash
# Create default configuration
lanup init

# Force overwrite existing config
lanup init --force
```

---

#### `lanup start`

Start exposing local services on your LAN.

```bash
lanup start [flags]
```

**Flags:**

- `-w, --watch` - Watch for network changes and update automatically
- `--no-env` - Display variables without writing to file
- `--dry-run` - Simulate all operations without writing files
- `--log` - Enable logging to file (default true)

**Examples:**

```bash
# Basic usage
lanup start

# Watch mode - auto-update on network changes
lanup start --watch

# Preview without modifying files
lanup start --dry-run

# Display variables without writing .env file
lanup start --no-env
```

---

#### `lanup expose`

Quickly expose a single service without configuration.

```bash
lanup expose [URL] [flags]
```

**Flags:**

- `--name string` - Assign an alias to the exposed service
- `--port int` - Use a custom port instead of the original
- `--https` - Use HTTPS protocol instead of HTTP

**Examples:**

```bash
# Expose a single service
lanup expose http://localhost:3000

# Expose with a custom name
lanup expose http://localhost:8080 --name api

# Expose with a different port
lanup expose http://localhost:5000 --port 8000

# Expose with HTTPS
lanup expose http://localhost:3000 --https
```

---

#### `lanup logs`

View or manage lanup logs.

```bash
lanup logs [flags]
```

**Flags:**

- `-n, --tail int` - Show last N lines (0 = show all)
- `-f, --follow` - Follow log output in real-time
- `--clear` - Clear the log file (requires confirmation)

**Examples:**

```bash
# View all logs
lanup logs

# View last 50 lines
lanup logs --tail 50

# Follow logs in real-time
lanup logs --follow

# Clear log file
lanup logs --clear
```

---

#### `lanup doctor`

Diagnose your local environment.

```bash
lanup doctor
```

Checks:

- Network interfaces and local IP detection
- Docker availability and running containers
- Supabase local development setup

**Example output:**

```
Running lanup diagnostics
✓ Network Interfaces
   Detected IP: 192.168.1.100 on interface en0 (wifi)
✓ Docker
   Docker is running with 3 active container(s)
✗ Supabase
   Supabase local is not running

⚠️  Some checks failed. Please review the issues above.
```

---

### Global Flags

These flags are available for all commands:

- `--config string` - Config file (default is $HOME/.lanup/config.yaml)
- `-v, --verbose` - Enable verbose output

## Configuration

### Project Configuration (`.lanup.yaml`)

Created in your project directory with `lanup init`.

```yaml
# Environment variables to expose
vars:
  SUPABASE_URL: "http://localhost:54321"
  SUPABASE_ANON_KEY: "your-anon-key"
  API_URL: "http://localhost:8000"
  DASHBOARD_URL: "http://localhost:3000"

# Output file path
output: ".env.local"

# Auto-detection settings
auto_detect:
  docker: true # Auto-detect Docker containers
  supabase: true # Auto-detect Supabase services
```

### Global Configuration (`~/.lanup/config.yaml`)

Created automatically on first run.

```yaml
# Log file path
log_path: "~/.lanup/logs/lanup.log"

# Log level (debug, info, warn, error)
log_level: "info"

# Default port for services
default_port: 8080

# Check interval for watch mode (seconds)
check_interval: 5
```

### Generated Environment File

lanup generates a `.env.local` file (or your configured output path) with transformed URLs:

```bash
# Generated by lanup on 2025-10-27 23:50:12
# Do not edit the managed variables manually

# lanup:managed
SUPABASE_URL=http://192.168.1.100:54321
# lanup:managed
SUPABASE_ANON_KEY=your-anon-key
# lanup:managed
API_URL=http://192.168.1.100:8000
# lanup:managed
DASHBOARD_URL=http://192.168.1.100:3000

# User variables (preserved)
DATABASE_URL=postgresql://localhost:5432/mydb
SECRET_KEY=my-secret
```

Variables marked with `# lanup:managed` are updated by lanup. Other variables are preserved.

## Use Cases

### Mobile App Development

Test your backend API from your phone or tablet:

```bash
# Start your backend
npm run dev

# Expose it on your network
lanup start --watch

# Access from your mobile device
# http://192.168.1.100:3000
```

### Supabase Local Development

Automatically expose Supabase services:

```bash
# Start Supabase
supabase start

# Expose with auto-detection
lanup start

# Your Supabase services are now accessible on your network
```

### Docker Development

Auto-detect and expose Docker containers:

```yaml
# .lanup.yaml
auto_detect:
  docker: true
```

```bash
# Start your Docker containers
docker-compose up

# Expose them
lanup start
```

### Team Collaboration

Share your local development environment with teammates on the same network:

```bash
# Start with watch mode
lanup start --watch

# Share the URLs with your team
# They can access your services from their devices
```

## Troubleshooting

### No Network Interface Detected

**Problem:** lanup can't detect your local IP address.

**Solutions:**

- Ensure you're connected to a network (Wi-Fi or Ethernet)
- Check that your network interface is active
- Run `lanup doctor` to diagnose the issue
- Try manually checking your IP with `ifconfig` (macOS/Linux) or `ipconfig` (Windows)

### Docker Auto-Detection Not Working

**Problem:** Docker containers are not being detected.

**Solutions:**

- Ensure Docker is running: `docker ps`
- Check Docker permissions (you may need to add your user to the docker group)
- Disable auto-detection and manually configure ports in `.lanup.yaml`
- Run `lanup doctor` to verify Docker status

### Environment File Not Updating

**Problem:** The `.env.local` file is not being updated.

**Solutions:**

- Check file permissions: `ls -la .env.local`
- Ensure the output path in `.lanup.yaml` is correct
- Try `lanup start --dry-run` to see what would be written
- Check logs: `lanup logs --tail 50`

### Network Changes Not Detected in Watch Mode

**Problem:** Watch mode doesn't update when switching networks.

**Solutions:**

- Check the `check_interval` in `~/.lanup/config.yaml` (default is 5 seconds)
- Ensure you have a stable network connection
- Try restarting watch mode
- Check logs for errors: `lanup logs --follow`

### Permission Denied Errors

**Problem:** lanup can't write files or access logs.

**Solutions:**

- Check file permissions in your project directory
- Ensure `~/.lanup/` directory has correct permissions
- Try running with appropriate permissions
- Check disk space: `df -h`

### Services Not Accessible from Other Devices

**Problem:** Other devices can't access the exposed URLs.

**Solutions:**

- Ensure all devices are on the same network
- Check firewall settings on your development machine
- Verify the services are actually running: `curl http://localhost:PORT`
- Try accessing from your development machine first using the network IP
- Some networks (corporate, public Wi-Fi) may block device-to-device communication

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Project Structure

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
├── main.go                # Entry point
├── Makefile               # Build scripts
└── README.md              # This file
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/net/...
```

### Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes and add tests
4. Run tests: `make test`
5. Format code: `make fmt`
6. Commit your changes: `git commit -am 'Add my feature'`
7. Push to the branch: `git push origin feature/my-feature`
8. Create a Pull Request

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Add comments for exported functions
- Write tests for new features
- Keep functions small and focused

## Security

### Private IP Addresses Only

lanup only uses private IP addresses (RFC 1918):

- `192.168.0.0/16`
- `10.0.0.0/8`
- `172.16.0.0/12`

Public IP addresses are automatically rejected to prevent accidental exposure to the internet.

### Environment File Safety

- Automatic backups before modifications (`.env.local.bak`)
- Preserves user-defined variables
- Clear markers for managed variables
- No logging of secrets (detects `*_KEY`, `*_SECRET`, `*_TOKEN` patterns)

### File Permissions

- Global config: `0600` (read/write for owner only)
- Environment files: `0644` (read for all, write for owner)
- Log directory: `0755` (accessible by owner)

## License

MIT License - see [LICENSE](LICENSE) file for details.

Copyright (c) 2025 [Rauche Acho](https://github.com/raucheacho)

## Author

Created and maintained by **Rauche Acho** - [GitHub](https://github.com/raucheacho)

## Documentation

Full documentation is available at: **[https://raucheacho.github.io/lanup/](https://raucheacho.github.io/lanup/)**

Or build locally:
```bash
cd docs
./setup.sh
hugo server
```

## Support

- **Documentation:** [https://raucheacho.github.io/lanup/](https://raucheacho.github.io/lanup/)
- **Issues:** [GitHub Issues](https://github.com/raucheacho/lanup/issues)
- **Discussions:** [GitHub Discussions](https://github.com/raucheacho/lanup/discussions)

## Acknowledgments

Built with:

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Color](https://github.com/fatih/color) - Terminal colors
- [YAML](https://github.com/go-yaml/yaml) - YAML parsing

---

Made with ❤️ for developers who want to test on real devices without the hassle.
