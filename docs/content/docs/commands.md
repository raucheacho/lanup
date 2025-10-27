---
title: "Commands Reference"
weight: 3
---

# Commands Reference

Complete reference for all lanup commands.

## lanup init

Initialize lanup configuration in your project.

```bash
lanup init [flags]
```

### Flags

- `--format string` - Configuration file format (yaml or toml) (default "yaml")
- `--force` - Overwrite existing configuration file

### Examples

```bash
# Create default configuration
lanup init

# Force overwrite existing config
lanup init --force
```

---

## lanup start

Start exposing local services on your LAN.

```bash
lanup start [flags]
```

### Flags

- `-w, --watch` - Watch for network changes and update automatically
- `--no-env` - Display variables without writing to file
- `--dry-run` - Simulate all operations without writing files
- `--log` - Enable logging to file (default true)

### Examples

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

## lanup expose

Quickly expose a single service without configuration.

```bash
lanup expose [URL] [flags]
```

### Flags

- `--name string` - Assign an alias to the exposed service
- `--port int` - Use a custom port instead of the original
- `--https` - Use HTTPS protocol instead of HTTP

### Examples

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

## lanup logs

View or manage lanup logs.

```bash
lanup logs [flags]
```

### Flags

- `-n, --tail int` - Show last N lines (0 = show all)
- `-f, --follow` - Follow log output in real-time
- `--clear` - Clear the log file (requires confirmation)

### Examples

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

## lanup doctor

Diagnose your local environment.

```bash
lanup doctor
```

Checks:

- Network interfaces and local IP detection
- Docker availability and running containers
- Supabase local development setup

### Example Output

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

## Global Flags

These flags are available for all commands:

- `--config string` - Config file (default is $HOME/.lanup/config.yaml)
- `-v, --verbose` - Enable verbose output
- `-h, --help` - Help for any command

## Exit Codes

- `0` - Success
- `1` - General error
- `2` - Configuration error
- `3` - Network error
- `4` - Permission error
