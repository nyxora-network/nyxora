# Reddit Post Draft

## Subreddits

- r/selfhosted
- r/golang
- r/privacy
- r/VPN

## Post Title

I built NYXORA - a self-healing multi-transport VPN/tunnel orchestrator with 12 transports and auto-failover

## Post Content

Hey everyone!

I've been working on NYXORA, a Go-based tool that manages multiple tunnel transports with automatic failover. Here's what it does:

**What is NYXORA?**
- Self-healing tunnel manager that automatically switches between 12 different transports
- Zero-config remote: Just SSH into any server and it auto-provisions everything
- Interactive TUI with beautiful themes

**Supported Transports:**
WireGuard, OpenVPN, SSH, QUIC, FRP, Rathole, IPsec, Shadowsocks, Hysteria, Backhaul, TCP, WebSocket

**Key Features:**
- Automatic failover (detects degraded tunnels and switches)
- 5 multipath scheduling modes
- Real-time scoring engine
- Prometheus metrics
- Docker support

**Use Cases:**
- Censorship bypass (auto-switches to obfuscated protocols)
- Unstable networks (continuous scoring)
- Low-resource VPS (minimal mode uses only SSH + Shadowsocks)

**Installation:**
```bash
curl -L github.com/nyxorammd-lgtm/nyxora/releases/download/v0.2.0/nyxora_linux_amd64 -o /usr/local/bin/nyxora && chmod +x /usr/local/bin/nyxora
```

**GitHub:** https://github.com/nyxorammd-lgtm/nyxora

Would love to hear your feedback! Happy to answer any questions.

## Posting Tips

- Post during peak hours (US mornings)
- Engage with comments quickly
- Cross-post to relevant subreddits
