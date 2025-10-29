---
title: "Installation"
weight: 1
---

# Installation

## Prerequisites

- Go 1.22+ (for building from source)
- macOS, Linux, or Windows

## Installation Methods

### Homebrew (macOS/Linux)

The recommended way to install lanup on macOS and Linux:

```bash
brew tap raucheacho/tap
brew install lanup
```

### Scoop (Windows)

The recommended way to install lanup on Windows:

```powershell
scoop bucket add raucheacho https://github.com/raucheacho/scoop-bucket
scoop install lanup
```

### Build from Source

If you have Go 1.22+ installed:

```bash
# Clone the repository
git clone https://github.com/raucheacho/lanup.git
cd lanup

# Build for your platform
make build

# Or install to GOPATH/bin
make install
```

Or install directly with Go:

```bash
go install github.com/raucheacho/lanup@latest
```

## Verify Installation

```bash
lanup --version
```

## Next Steps

- [Getting Started Guide]({{< ref "getting-started" >}})
- [Configuration]({{< ref "configuration" >}})
