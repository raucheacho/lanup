---
title: "Installation"
weight: 1
---

# Installation

## Prerequisites

- Go 1.22+ (for building from source)
- macOS, Linux, or Windows

## Installation Methods

### From Source (Go Install)

If you have Go 1.22+ installed:

```bash
go install github.com/raucheacho/lanup@latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/raucheacho/lanup.git
cd lanup

# Build for your platform
make build

# Or install to GOPATH/bin
make install
```

### Pre-built Binaries

Download pre-built binaries from the [releases page](https://github.com/raucheacho/lanup/releases).

#### macOS

**Intel Mac:**
```bash
curl -L https://github.com/raucheacho/lanup/releases/latest/download/lanup-darwin-amd64 -o lanup
chmod +x lanup
sudo mv lanup /usr/local/bin/
```

**Apple Silicon (M1/M2):**
```bash
curl -L https://github.com/raucheacho/lanup/releases/latest/download/lanup-darwin-arm64 -o lanup
chmod +x lanup
sudo mv lanup /usr/local/bin/
```

#### Linux

```bash
curl -L https://github.com/raucheacho/lanup/releases/latest/download/lanup-linux-amd64 -o lanup
chmod +x lanup
sudo mv lanup /usr/local/bin/
```

#### Windows

Download `lanup-windows-amd64.exe` from the releases page and add it to your PATH.

## Verify Installation

```bash
lanup --version
```

## Next Steps

- [Getting Started Guide]({{< ref "getting-started" >}})
- [Configuration]({{< ref "configuration" >}})
