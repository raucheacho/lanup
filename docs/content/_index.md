---
title: "lanup"
type: docs
---

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

## Quick Start

```bash
# Initialize lanup in your project
cd your-project
lanup init

# Start exposing your services
lanup start
```

Your services are now accessible on your local network! 🎉

## Documentation

- [Installation Guide]({{< ref "/docs/installation" >}})
- [Getting Started]({{< ref "/docs/getting-started" >}})
- [Commands Reference]({{< ref "/docs/commands" >}})
- [Configuration]({{< ref "/docs/configuration" >}})
- [Use Cases]({{< ref "/docs/use-cases" >}})
- [Troubleshooting]({{< ref "/docs/troubleshooting" >}})

## License

MIT License - Copyright (c) 2025 [raucheacho](https://github.com/raucheacho)
