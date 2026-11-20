# lanup

**lanup** is a CLI tool that automatically exposes your local backend services on your local network (LAN). It detects your local IP address, updates environment variables, and allows you to test your applications from any device on the same network without manual configuration.

Perfect for testing mobile apps, sharing work with teammates, or accessing your development environment from multiple devices.

## Features

- 🌐 **Automatic IP Detection** - Detects your local network IP address
- 🔄 **Watch Mode** - Automatically updates when your network changes (default behavior)
- 🐳 **Docker Integration** - Auto-detects running Docker containers
- 🗄️ **Supabase Support** - Detects Supabase local development services
- 📝 **Environment File Management** - Safely updates `.env` and `.envrc` files with backups
- 📱 **QR Code Generation** - Instantly share URLs to mobile devices via terminal QR codes
- 🎨 **Colorized Output** - Beautiful, easy-to-read console output
- 🔍 **Diagnostics** - Built-in health checks for troubleshooting
- 🔒 **Secure by Default** - Only uses private IP addresses (RFC 1918)

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap raucheacho/tap
brew install lanup
```

### Scoop (Windows)

```powershell
scoop bucket add raucheacho https://github.com/raucheacho/scoop-bucket
scoop install lanup
```

### Build from Source

If you have Go 1.22+ installed:

```bash
git clone https://github.com/raucheacho/lanup.git
cd lanup
go build -o lanup .
sudo mv lanup /usr/local/bin/  # or add to your PATH
```

Or install directly with Go:

```bash
go install github.com/raucheacho/lanup@latest
```

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

- `-1, --once` - Run once and exit without watching for network changes
- `--no-env` - Display variables without writing to file
- `--dry-run` - Simulate all operations without writing files

**Examples:**

```bash
# Start watching for network changes (default behavior)
lanup start

# Update once and exit
lanup start --once

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

- `--port int` - Use a custom port instead of the original
- `--qr` - Display a QR code for the exposed URL

**Examples:**

```bash
# Expose a single service
lanup expose http://localhost:3000

# Expose and display a QR code
lanup expose http://localhost:3000 --qr

# Expose with a different port
lanup expose http://localhost:5000 --port 8000

# Expose with HTTPS
lanup expose http://localhost:3000 --https
```


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

### Several Files (monorepos)

One `.lanup.yaml` at the root can update every env file of a monorepo with `outputs`. Paths are relative to the directory lanup runs in; top-level `vars` go to every output, and an output's own `vars` win.

```yaml
outputs:
  - path: backend/.envrc
    vars:
      # {{ip}} is replaced by the detected IP; localhost is then left as is
      ALLOWED_ORIGINS: "http://localhost:3000,http://{{ip}}:3000"
      WS_URL: "ws://localhost:4000/socket"
  - path: apps/mobile/.env.local
    vars:
      API_URL: "http://localhost:8084"
```

Each output is **updated in place**: only the lines of its variables are rewritten (keeping their `export` prefix and quotes); comments, other variables and shell commands stay untouched. A variable not yet in the file is appended under a `# lanup:managed` marker, and no line is ever removed. A `.envrc` output is written with `export` and quotes values that need it.

The single `output` keeps its behavior: the file is regenerated with a header. Both can be used together.

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

# Expose it on your network (watches for changes automatically)
lanup start

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
# Start lanup (watch mode is the default)
lanup start

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
│   └── doctor.go          # Doctor command
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   ├── net/               # Network detection
│   ├── env/               # Environment file management
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

All documentation is contained within this README. Check out the commands and use cases above!

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
