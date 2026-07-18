# NYXORA Security Hardening Guide

This guide provides recommendations for hardening NYXORA deployments.

## Table of Contents

1. [SSH Hardening](#ssh-hardening)
2. [Firewall Configuration](#firewall-configuration)
3. [Secret Management](#secret-management)
4. [Network Security](#network-security)
5. [System Hardening](#system-hardening)
6. [Monitoring and Logging](#monitoring-and-logging)

---

## SSH Hardening

### Disable Password Authentication

Edit `/etc/ssh/sshd_config`:

```bash
PasswordAuthentication no
ChallengeResponseAuthentication no
UsePAM no
```

### Use Key-Based Authentication Only

```bash
# Generate strong key
ssh-keygen -t ed25519 -a 100

# Copy to server
ssh-copy-id -i ~/.ssh/id_ed25519.pub user@server
```

### Restrict SSH Access

```bash
# /etc/ssh/sshd_config
AllowUsers nyxora
MaxAuthTries 3
LoginGraceTime 30
ClientAliveInterval 300
ClientAliveCountMax 2
```

---

## Firewall Configuration

### UFW (Ubuntu/Debian)

```bash
# Allow SSH
ufw allow 22/tcp

# Allow NYXORA ports
ufw allow 51820/udp  # WireGuard
ufw allow 1194/udp   # OpenVPN
ufw allow 9923/udp   # QUIC
ufw allow 7000/tcp   # FRP
ufw allow 2333/tcp   # Rathole
ufw allow 8388/tcp   # Shadowsocks
ufw allow 8444/udp   # Hysteria
ufw allow 3080/tcp   # Backhaul
ufw allow 9924/tcp   # TCP tunnel
ufw allow 9925/tcp   # WebSocket

# Enable firewall
ufw enable
```

### iptables

```bash
# Basic rules
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -j ACCEPT
iptables -A INPUT -p udp --dport 51820 -j ACCEPT
iptables -A INPUT -p tcp --dport 1194 -j ACCEPT
iptables -A INPUT -j DROP

# Save rules
iptables-save > /etc/iptables/rules.v4
```

---

## Secret Management

### Environment Variables

Set secrets via environment variables:

```bash
export NYXORA_SS_PASSWORD="your-strong-password"
export NYXORA_RATHOLE_TOKEN="your-secure-token"
export NYXORA_HYSTERIA_AUTH="your-hysteria-password"
export NYXORA_BACKHAUL_TOKEN="your-backhaul-token"
export NYXORA_IPSEC_PSK="your-ipsec-psk"
```

### File Permissions

```bash
# Set proper permissions
chmod 600 /etc/nyxora/config.json
chmod 600 /etc/nyxora/secrets/*
chown -R root:root /etc/nyxora
```

### Secret Rotation

Enable automatic secret rotation in config:

```json
{
  "secret_rotation": {
    "enabled": true,
    "interval_hours": 24,
    "max_age_days": 30
  }
}
```

---

## Network Security

### Disable IPv6 (if not needed)

```bash
# /etc/sysctl.conf
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
```

### Enable SYN Cookies

```bash
# /etc/sysctl.conf
net.ipv4.tcp_syncookies = 1
```

### Limit Connection Tracking

```bash
# /etc/sysctl.conf
net.netfilter.nf_conntrack_max = 262144
```

---

## System Hardening

### Update System

```bash
# Ubuntu/Debian
apt update && apt upgrade -y

# CentOS/RHEL
yum update -y
```

### Install Security Updates Only

```bash
# Ubuntu/Debian
unattended-upgrades
```

### Disable Unnecessary Services

```bash
# List running services
systemctl list-units --type=service --state=running

# Disable unnecessary services
systemctl disable <service-name>
```

### Set Kernel Parameters

```bash
# /etc/sysctl.conf
# Disable IP forwarding (unless required)
net.ipv4.ip_forward = 0

# Enable reverse path filtering
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1

# Ignore ICMP redirects
net.ipv4.conf.all.accept_redirects = 0
net.ipv4.conf.default.accept_redirects = 0

# Log Martian packets
net.ipv4.conf.all.log_martians = 1
```

---

## Monitoring and Logging

### Enable Audit Logging

Configure NYXORA to log all actions:

```json
{
  "audit": {
    "enabled": true,
    "log_path": "/var/log/nyxora/audit.log",
    "level": "info"
  }
}
```

### Monitor Logs

```bash
# Watch NYXORA logs
tail -f /var/log/nyxora/nyxora.log

# Watch audit logs
tail -f /var/log/nyxora/audit.log
```

### Set Up Log Rotation

```bash
# /etc/logrotate.d/nyxora
/var/log/nyxora/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 root root
}
```

---

## Checklist

- [ ] SSH password authentication disabled
- [ ] Key-based authentication enabled
- [ ] Firewall configured and enabled
- [ ] Secrets rotated regularly
- [ ] File permissions set correctly
- [ ] System updated with security patches
- [ ] Unnecessary services disabled
- [ ] Kernel parameters hardened
- [ ] Audit logging enabled
- [ ] Log rotation configured
- [ ] Monitoring alerts set up

---

## Additional Resources

- [CIS Benchmarks](https://www.cisecurity.org/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
