# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

---

## [1.0.0] - 2026-07-18

### Added
- **Cobra CLI framework** — Professional command-line interface with subcommands, flag validation, and help system
- **Concurrent transport connections** — All 12 transports connect simultaneously using `errgroup` (reduces connection time from ~24s to ~2s)
- **30-second connection timeout** — Prevents infinite hangs on unresponsive transports
- **Comprehensive unit tests** — 15+ new tests for transport manager with mock transports
- **CI/CD pipeline improvements** — Updated to Go 1.25, added test timeout, race detection

### Changed
- **Module path fixed** — `go.mod` now correctly declares `github.com/nyxora-network/nyxora`
- **CLI architecture** — Replaced manual `os.Args` parsing with `spf13/cobra` v1.10.2
- **Error handling** — All critical paths now use `fmt.Errorf` wrapping with context
- **Go version** — Minimum required Go version updated from 1.24 to 1.25
- **README.md** — Updated with new CLI commands, architecture highlights, and badges
- **CI/CD workflows** — Updated `ci.yml` and `release.yml` to Go 1.25

### Fixed
- **Import path mismatch** — All internal imports now use correct module path
- **Error context loss** — Bare `return err` patterns replaced with wrapped errors
- **CLI flag parsing** — Fixed fragile manual parsing with proper Cobra flag handling

### Removed
- **Legacy CLI parser** — Removed 400+ lines of manual `os.Args` parsing code

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
