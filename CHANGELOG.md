# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CONTRIBUTING.md in Persian, Russian, and Chinese
- Enhanced API Reference documentation
- Comprehensive test coverage improvements

### Changed
- Updated transport count from 11 to 12 across all READMEs
- Added WebSocket transport to all language documentation
- Fixed Go version badge from 1.25 to 1.24 in all READMEs

---

## [0.2.0] - 2026-06-23

### Added
- **Interactive Bubble Tea TUI** with full keyboard navigation
  - Main menu with 9 options
  - Connect wizard with step-by-step guidance
  - Transport status viewer with animated score bars
  - Tunnel topology view
  - Help screen with keyboard shortcuts
- **3 professional color themes**
  - Catppuccin Mocha (dark)
  - Tokyo Night (dark)
  - Catppuccin Latte (light)
- **ASCII art boot splash** with animated gradient progress bar
- **Live system monitoring**
  - CPU load with color-coded status
  - RAM usage with percentage bar
  - Goroutine count
- **TrueColor gradient rendering engine** for smooth color transitions
- **Backhaul transport** implementation (12th transport)
- **Single-side install flow** - no agent required on remote
- **TUI wizard** for first-time setup

### Changed
- Major refactor of internal architecture
- Dashboard now uses Catppuccin TrueColor palette
- Improved ANSI escape handling in dashboard
- Better error handling and user feedback
- Optimized transport scoring algorithm

### Fixed
- WireGuard IPv6 endpoint formatting
- WireGuard remote key passing and subnet alignment
- WireGuard iptables rules to prevent SSH loss
- FRP/Rathole install scripts with API-based URL resolution
- Race condition in failover engine
- Memory leak in transport metrics collection

### Security
- Secret files now use 0600 permissions
- Config file permissions hardened
- Added input validation for port overrides

---

## [0.1.0] - 2026-06-01

### Added
- **Core orchestrator** with connect/disconnect flow
- **11 initial tunnel transports**
  - WireGuard (full remote provisioning)
  - OpenVPN
  - SSH tunnel
  - QUIC
  - FRP (Fast Reverse Proxy)
  - Rathole
  - IPsec/strongSwan
  - Shadowsocks
  - Hysteria 2
  - TCP tunnel
  - (Backhaul added in v0.2.0)
- **SSH-based remote setup** and management
  - Auto-detect OS (Ubuntu, Debian, CentOS)
  - Auto-install tunnel binaries
  - Password and key authentication
- **Ping-based monitoring system**
  - Latency measurement
  - Packet loss detection
  - Jitter calculation
  - Stability scoring
- **Automatic failover engine**
  - Health status tracking (healthy/degraded/down)
  - Configurable thresholds
  - Callback-based notifications
- **Multipath scheduler** with 5 distribution modes
  - Weighted (based on transport weights)
  - Lowest-latency (route through fastest path)
  - Lowest-loss (route through most reliable path)
  - Even (equal distribution)
  - All-active (all tunnels simultaneously)
- **Scoring engine**
  - Latency scoring
  - Packet loss scoring
  - Jitter scoring
  - Stability scoring
  - Configurable weights
- **ANSI terminal dashboard**
  - Real-time transport status
  - Score visualization
  - Connection info
- **Configuration management**
  - JSON load/save
  - Environment variable support
  - Mode detection (full/lite/minimal)
  - Port overrides
- **Secret/token auto-generation**
  - Shadowsocks password
  - Rathole token
  - Hysteria auth
  - Backhaul token
  - IPsec PSK
- **Tar.gz packaging** for tunnel assets
- **Docker multi-stage build**
- **Makefile** with build, test, install, clean targets
- **GitHub Actions CI/CD**
  - Build and test on push/PR
  - CodeQL security analysis
  - golangci-lint
  - Release automation

---

## Roadmap

### [0.3.0] - Planned
- [ ] VLESS transport support
- [ ] Reality protocol
- [ ] DNS-over-HTTPS resolver
- [ ] Web UI for browser-based management
- [ ] Homebrew formula for macOS
- [ ] AUR package for Arch Linux
- [ ] Enhanced logging with structured output
- [ ] Performance benchmarks

### [0.4.0] - Planned
- [ ] Load balancing algorithms
- [ ] Certificate rotation
- [ ] Audit logging
- [ ] Multi-node clustering
- [ ] Prometheus/Grafana integration
- [ ] REST API for external tools

### [1.0.0] - Future
- [ ] Stable API
- [ ] Full test coverage (>80%)
- [ ] Production-ready status
- [ ] Enterprise features
