# NYXORA

**Adaptive Tunnel Orchestrator** — Self-healing multi-transport VPN/tunnel manager.

Install on ONE server. It SSHs into a remote server, auto-provisions tunnels, monitors them, and fails over if needed. No agent required on the remote side.

## Features

- **11 transport types**: WireGuard, OpenVPN, SSH, QUIC, FRP, Rathole, IPsec, Shadowsocks, Hysteria, Backhaul, TCP
- **Automatic failover**: Detects degraded tunnels and switches to healthier ones
- **Multipath scheduling**: 5 distribution modes (weighted, lowest-latency, lowest-loss, even, all-active)
- **Real-time TUI dashboard**: Live terminal monitoring
- **Zero remote setup**: Just SSH access needed (root password or key)

## Quick Start

```bash
# Build
go build -ldflags="-s -w" -o nyxora ./cmd/nyxora

# Install
./nyxora install

# Connect to remote server
./nyxora connect 91.107.243.237 --user root --password <pass>

# Live dashboard
./nyxora dashboard
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NYXORA_SS_PASSWORD` | Shadowsocks password | random |
| `NYXORA_SS_METHOD` | Shadowsocks cipher | aes-256-gcm |
| `NYXORA_RATHOLE_TOKEN` | Rathole auth token | random |
| `NYXORA_HYSTERIA_AUTH` | Hysteria auth password | random |
| `NYXORA_BACKHAUL_TOKEN` | Backhaul auth token | random |
| `NYXORA_IPSEC_PSK` | IPsec pre-shared key | random |
| `NYXORA_ALL_ACTIVE` | Enable all tunnels simultaneously | false |

## Project Structure

```
nyxora/
├── cmd/
│   ├── nyxora/main.go          CLI entrypoint
│   └── quic-server/main.go     QUIC echo server
├── internal/
│   ├── config/                 Config + secrets management
│   ├── orchestrator/           Brain: connect, monitor, failover
│   ├── transport/              11 transport implementations
│   ├── remote/                 SSH client + remote provisioning
│   ├── multipath/              Multipath scheduler (5 modes)
│   ├── failover/               Automatic failover engine
│   ├── routing/                Scorer + routing engine
│   ├── monitor/                Ping-based monitoring
│   ├── dashboard/              TUI terminal dashboard
│   └── packager/               Tar.gz pack/unpack
├── tunnels/                    Install scripts per tunnel
├── Makefile                    Build, test, install
└── PROJECT.txt                 Full documentation (EN/FA)
```

## Transports

| # | Name | Port | Category | Score | Weight |
|---|------|------|----------|-------|--------|
| 1 | wireguard | 51820 | VPN | 95 | 30 |
| 2 | openvpn | 1194 | VPN | 75 | 10 |
| 3 | ssh | 22 | tunnel | 60 | 5 |
| 4 | quic | 9923 | tunnel | 80 | 15 |
| 5 | frp | 7000 | relay | 70 | 10 |
| 6 | rathole | 2333 | relay | 85 | 12 |
| 7 | ipsec | 500 | VPN | 70 | 5 |
| 8 | shadowsocks | 8388 | proxy | 55 | 3 |
| 9 | hysteria | 8444 | tunnel | 90 | 12 |
| 10 | backhaul | 3080 | relay | 82 | 10 |
| 11 | tcp | 9924 | tunnel | 50 | 3 |

## Building

```bash
go build -ldflags="-s -w" -o nyxora ./cmd/nyxora
go test ./...
go vet ./...
```

## Docker

```bash
docker build -t nyxora .
docker run --rm -it nyxora connect <ip> --user root
```

## License

MIT
