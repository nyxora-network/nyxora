# Hacker News Post Draft

## Title

NYXORA: Self-healing multi-transport VPN/tunnel orchestrator (12 transports, auto-failover, interactive TUI)

## URL

https://github.com/nyxorammd-lgtm/nyxora

## Post Content

Show HN: I built NYXORA, a self-healing multi-transport VPN/tunnel orchestrator

NYXORA is a Go-based tool that manages multiple tunnel transports simultaneously with automatic failover and a beautiful interactive TUI.

**Key Features:**
- 12 tunnel transports: WireGuard, OpenVPN, SSH, QUIC, FRP, Rathole, IPsec, Shadowsocks, Hysteria, Backhaul, TCP, WebSocket
- Zero-config remote: No agent required, just SSH access
- Auto-failover: Detects degraded tunnels and switches instantly
- Interactive TUI with 3 professional themes
- Prometheus metrics and DNS resolver

**Use Cases:**
- Censorship bypass (auto-switches to obfuscated protocols)
- Unstable networks (continuous scoring and failover)
- DevOps automation (JSON config, daemon mode)

**Technical Highlights:**
- Written in pure Go, minimal dependencies
- Cross-platform (Linux, macOS)
- Docker support
- Comprehensive tests

Would love feedback from the community!

## Best Times to Post

- **Tuesday-Thursday**, 9-11 AM EST
- Avoid weekends and holidays
