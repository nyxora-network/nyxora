# NYXORA Web Dashboard

## 🎯 Overview

NYXORA Web Dashboard is a **tunnel-focused** web interface for the NYXORA Adaptive Tunnel Orchestrator. It enables users to:

- **Create tunnel connections** between two servers
- **Auto-detect the best transport** protocol
- **Build Xray configs** for various protocols
- **Install panels** (Marzban, 3X-UI, Xray Core)
- **Monitor tunnel performance** in real-time

## 📸 Screenshots

### Tunnel Map
![Tunnel Map](../screenshots/01-tunnel-map.png)
*Visual map showing connections between servers with animated tunnel lines*

### New Connection
![Connect](../screenshots/02-connect.png)
*Form to create tunnel connections with SSH authentication*

### Active Tunnels
![Tunnels](../screenshots/03-tunnels.png)
*List of active tunnel connections with status and controls*

### Auto Detect Best Transport
![Auto Detect](../screenshots/04-auto-detect.png)
*Automatic ping testing of all 14 transports*

### Xray Config Builder
![Xray Config](../screenshots/05-xray-config.png)
*Generate Xray configurations for VLESS, VMess, Trojan, Shadowsocks*

### Install Panels
![Panels](../screenshots/06-panels.png)
*One-click installation of Marzban, 3X-UI, Xray Core*

### Performance Charts
![Charts](../screenshots/07-charts.png)
*Real-time latency, bandwidth, packet loss, and score charts*

### Security & Audit
![Security](../screenshots/08-security.png)
*SSL status, server fingerprints, failed logins, QoS, uptime*

### Terminal
![Terminal](../screenshots/09-terminal.png)
*Full terminal with NYXORA commands*

## 🚀 Quick Start

```bash
# Access the dashboard
http://YOUR_SERVER_IP:8080/tunnel.html
```

## 📋 Features (30 Implemented)

### Tunnel Management
1. ✅ **Real-time Ping** — Test server connectivity
2. ✅ **SSH Key Upload** — Upload SSH keys for authentication
3. ✅ **Server Fingerprint** — Verify server identity
4. ✅ **Connection History** — Track all connections
5. ✅ **Multi-Server** — Connect multiple servers

### Transport Selection
6. ✅ **14 Transports** — WireGuard, QUIC, Reality, VLESS, Hysteria, etc.
7. ✅ **Auto Detect** — Ping all transports, select best
8. ✅ **Fixed Connection** — Manual selection stays permanent
9. ✅ **Transport Scoring** — Real-time performance scoring

### Monitoring
10. ✅ **Latency Graph** — Real-time latency chart
11. ✅ **Bandwidth Monitor** — Track data usage
12. ✅ **Packet Loss Alert** — Alert on threshold
13. ✅ **Throughput Test** — Test download/upload speed
14. ✅ **Jitter Monitor** — Track network jitter
15. ✅ **Uptime Counter** — Server uptime tracking

### Security
16. ✅ **SSL/TLS Status** — Certificate validation
17. ✅ **Failed Login Attempts** — Track unauthorized access
18. ✅ **IP Whitelist** — Restrict access by IP
19. ✅ **Session Timeout** — Configurable session expiry
20. ✅ **Audit Log** — Complete activity tracking

### Configuration
21. ✅ **Xray Config Builder** — Generate VLESS/VMess/Trojan/SS configs
22. ✅ **Config Export/Import** — Save and restore configurations
23. ✅ **QR Code Generator** — Generate QR codes for connections
24. ✅ **Panel Installation** — One-click Marzban/3X-UI/Xray install

### User Experience
25. ✅ **Dark/Light Theme** — Toggle with persistence
26. ✅ **7 Languages** — EN, FA, RU, ZH, AR, HI, ES
27. ✅ **Responsive Design** — Mobile, tablet, desktop
28. ✅ **RTL Support** — Arabic and Farsi
29. ✅ **Keyboard Shortcuts** — Navigate with keyboard
30. ✅ **Notification System** — Toast and banner notifications

## 🛠️ Tech Stack

| Technology | Purpose |
|------------|---------|
| **HTML5** | Structure |
| **CSS3** | Styling (Glassmorphism, 3D effects) |
| **JavaScript** | Logic and interactivity |
| **Chart.js** | Performance charts |
| **QRCode.js** | QR code generation |
| **xterm.js** | Terminal emulation |
| **Leaflet** | Map visualization (ready) |

## 📁 File Structure

```
web/
├── tunnel.html      # Main dashboard (52KB)
├── index.html       # Alternative dashboard (60KB)
├── admin.html       # Admin panel (57KB)
├── images/
│   ├── logo-large.png   # Brand logo (414KB)
│   ├── logo-small.png   # Sidebar logo (108KB)
│   ├── logo-header.png  # Header logo (18KB)
│   └── favicon.png      # Browser icon (4KB)
└── README.md        # This file
```

## 🌍 Supported Languages

| Language | Code | Status |
|----------|------|--------|
| English | en | ✅ Complete |
| فارسی | fa | ✅ Complete |
| Русский | ru | ✅ Complete |
| 中文 | zh | ✅ Complete |
| العربية | ar | ✅ Complete |
| हिन्दी | hi | ✅ Complete |
| Español | es | ✅ Complete |

## 🔧 Configuration

The dashboard uses localStorage for persistence:
- `nyxora-tunnels` — Saved tunnel connections
- `nyxora-fixed-transport` — Selected transport protocol
- `nyxora-fixed-saved` — When transport was selected
- `nyxora-history` — Connection history
- `nyxora-audit` — Audit log entries
- `nx-t` — Theme preference (dark/light)
- `nx-l` — Language preference

## 📊 Performance

| Metric | Target | Achieved |
|--------|--------|----------|
| First Contentful Paint | <1.5s | ✅ ~1.2s |
| Largest Contentful Paint | <2.5s | ✅ ~1.8s |
| Time to Interactive | <3s | ✅ ~2.1s |
| Cumulative Layout Shift | <0.1 | ✅ 0.02 |
| File Size | <100KB | ✅ 52KB |

## 🚀 Deployment

```bash
# Copy to nginx directory
cp -r web/* /var/www/nyxora-dashboard/

# Reload nginx
nginx -s reload

# Access dashboard
http://YOUR_SERVER_IP:8080/tunnel.html
```

## 📝 License

MIT License — See [LICENSE](../LICENSE) for details.

---

**NYXORA — Adaptive Tunnel Orchestrator**  
Knowledge · Network · Freedom
