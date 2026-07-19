---
title: NYXORA
description: Adaptive Tunnel Orchestrator — Documentation
---

# NYXORA Documentation

Welcome to the NYXORA documentation. NYXORA is an adaptive tunnel orchestrator that helps you manage and route traffic through multiple tunnel transports.

## Quick Links

- [Quickstart Guide](quickstart.md) — Get started in 5 minutes
- [Architecture Overview](architecture.md) — Understand how NYXORA works
- [Transports](transports.md) — Supported tunnel protocols

## What is NYXORA?

NYXORA is a Go-based tunnel orchestrator designed for:

- **Multi-transport routing** — Automatically select the best tunnel
- **High latency tolerance** — Optimized for high-RTT links
- **Packet loss resilience** — Smart retransmission and fallback
- **Cross-platform** — Linux and macOS support

## Features

- 11 transport protocols
- Adaptive routing
- CLI interface
- Docker support
- REST API

## Get Started

```bash
curl -fsSL https://raw.githubusercontent.com/nyxora-network/nyxora/main/install.sh | bash
nyxora --help
```
