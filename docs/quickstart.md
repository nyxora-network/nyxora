# Quick Start Guide

## Prerequisites

- Linux server (Ubuntu 22.04+ recommended)
- Root SSH access to a remote server
- Go 1.25+ (only if building from source)

## Installation

### One-Line Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/nyxora-network/nyxora/main/install.sh | sudo bash
```

### Manual Install

```bash
git clone https://github.com/nyxora-network/nyxora.git
cd nyxora
make build
sudo make install
```

## First Connection

```bash
# Check dependencies
nyxora install

# Connect to a remote server
nyxora connect 192.168.1.100 --user root --password your_password

# Launch interactive menu
nyxora tui

# Live dashboard
nyxora dashboard
```

## What Happens When You Connect

1. **Ping** — measures latency and packet loss
2. **SSH** — authenticates to the remote server
3. **OS Detection** — detects Ubuntu/Debian/CentOS
4. **Install** — installs tunnel binaries on the remote
5. **WireGuard** — sets up kernel-level VPN
6. **Provision** — starts 5 daemons on the remote
7. **Test** — activates and scores all tunnels
8. **Monitor** — continuous health checking with auto-failover

## Next Steps

- Explore the interactive TUI: `nyxora tui`
- View live stats: `nyxora dashboard`
- Check status: `nyxora status`
- Run as a service: `nyxora daemon`
