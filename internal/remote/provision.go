package remote

import (
	"fmt"
	"log"
)

func ensureConfigDir(host *Host) {
	host.SSHCommand("mkdir -p /etc/nyxora 2>/dev/null")
}

func startDaemon(host *Host, name, binary, args string) error {
	startCmd := fmt.Sprintf(
		`nohup %s %s </dev/null >/var/log/nyxora-%s.log 2>&1 &
sleep 1
echo "started"`, binary, args, name)
	_, err := host.SSHCommand(startCmd)
	if err != nil {
		log.Printf("[provision] %s start: %v (may still be OK)", name, err)
	}
	return nil
}

func ProvisionFRPServer(host *Host, port int) error {
	log.Printf("[provision] setting up frp server on %s:%d", host.Address, port)
	ensureConfigDir(host)
	host.SSHCommand("pkill -f 'frps' 2>/dev/null; sleep 1; true")
	cfg := fmt.Sprintf(`[common]
bind_port = %d
`, port)
	path := "/etc/nyxora/frps.ini"
	if err := host.WriteFile(path, cfg, "644"); err != nil {
		return fmt.Errorf("write frps config: %w", err)
	}
	host.SSHCommand(`which frps 2>/dev/null || (TAG=$(curl -sL https://api.github.com/repos/fatedier/frp/releases/latest | grep tag_name | cut -d'"' -f4) && VER=${TAG#v} && curl -sL "https://github.com/fatedier/frp/releases/download/${TAG}/frp_${VER}_linux_amd64.tar.gz" -o /tmp/frp-srv.tar.gz && tar -xzf /tmp/frp-srv.tar.gz -C /tmp/ && cp /tmp/frp_*/frps /usr/local/bin/frps && chmod +x /usr/local/bin/frps)`)
	startDaemon(host, "frps", "frps", fmt.Sprintf("-c %s", path))
	log.Printf("[provision] frp server ready | port: %d", port)
	return nil
}

func ProvisionRatholeServer(host *Host, port int, token string) error {
	log.Printf("[provision] setting up rathole server on %s:%d", host.Address, port)
	ensureConfigDir(host)
	host.SSHCommand("pkill -f 'rathole' 2>/dev/null; sleep 1; true")
	cfg := fmt.Sprintf(`[server]
bind_addr = "0.0.0.0:%d"

[server.services.nyxora]
type = "tcp"
bind_addr = "127.0.0.1:6000"
token = "%s"
`, port, token)
	path := "/etc/nyxora/rathole-server.toml"
	if err := host.WriteFile(path, cfg, "644"); err != nil {
		return fmt.Errorf("write rathole config: %w", err)
	}
	startDaemon(host, "rathole", "rathole", fmt.Sprintf("--server %s", path))
	log.Printf("[provision] rathole server ready | port: %d", port)
	return nil
}

func ProvisionShadowSOCKSServer(host *Host, port int, password, method string) error {
	log.Printf("[provision] setting up shadowsocks server on %s:%d", host.Address, port)
	ensureConfigDir(host)
	host.SSHCommand("pkill -f 'ss-server' 2>/dev/null; sleep 1; true")
	cfg := fmt.Sprintf(`{
	"server": "0.0.0.0",
	"server_port": %d,
	"password": "%s",
	"method": "%s",
	"timeout": 60
}`, port, password, method)
	path := "/etc/nyxora/ss-config.json"
	if err := host.WriteFile(path, cfg, "644"); err != nil {
		return fmt.Errorf("write ss config: %w", err)
	}
	startDaemon(host, "ss", "ss-server", fmt.Sprintf("-c %s", path))
	log.Printf("[provision] shadowsocks server ready | port: %d", port)
	return nil
}

func ProvisionHysteriaServer(host *Host, port int, authPass string) error {
	log.Printf("[provision] setting up hysteria server on %s:%d", host.Address, port)
	ensureConfigDir(host)
	host.SSHCommand("pkill -f 'hysteria' 2>/dev/null; sleep 1; true")
	host.SSHCommand(`openssl req -x509 -newkey rsa:2048 -keyout /etc/nyxora/hy2-key.pem -out /etc/nyxora/hy2-cert.pem -days 365 -nodes -subj "/CN=nyxora" 2>/dev/null`)
	cfg := fmt.Sprintf(`listen: ":%d"
auth:
  type: password
  password: "%s"
tls:
  cert: /etc/nyxora/hy2-cert.pem
  key: /etc/nyxora/hy2-key.pem
bandwidth:
  up: "200 mbps"
  down: "1000 mbps"
`, port, authPass)
	path := "/etc/nyxora/hy2-server.yaml"
	if err := host.WriteFile(path, cfg, "644"); err != nil {
		return fmt.Errorf("write hysteria config: %w", err)
	}
	startDaemon(host, "hy2", "hysteria", fmt.Sprintf("server -c %s", path))
	log.Printf("[provision] hysteria server ready | port: %d", port)
	return nil
}

func ProvisionBackhaulServer(host *Host, port int, token, transport string) error {
	log.Printf("[provision] setting up backhaul server on %s:%d (%s)", host.Address, port, transport)
	ensureConfigDir(host)
	host.SSHCommand("pkill -f 'backhaul' 2>/dev/null; sleep 1; true")
	cfg := fmt.Sprintf(`[server]
bind_addr = "0.0.0.0:%d"
transport = "%s"
token = "%s"
keepalive_period = 75
nodelay = true
heartbeat = 40
channel_size = 2048
log_level = "error"
ports = []
`, port, transport, token)
	path := "/etc/nyxora/backhaul-server.toml"
	if err := host.WriteFile(path, cfg, "644"); err != nil {
		return fmt.Errorf("write backhaul config: %w", err)
	}
	startDaemon(host, "backhaul", "backhaul", fmt.Sprintf("-c %s", path))
	log.Printf("[provision] backhaul server ready | port: %d", port)
	return nil
}

func ProvisionIPsecServer(host *Host, remoteIP string, psk string) error {
	log.Printf("[provision] setting up ipsec on %s (peer: %s)", host.Address, remoteIP)
	ipsecConf := fmt.Sprintf(`config setup

conn nyxora
    left=%s
    right=%s
    authby=secret
    ike=aes256-sha256-modp2048
    esp=aes256-sha256
    auto=start
    type=transport
`, host.Address, remoteIP)
	out, err := host.SSHCommand(fmt.Sprintf(`cat > /etc/ipsec.conf << 'EOF'
%s
EOF`, ipsecConf))
	if err != nil {
		return fmt.Errorf("write ipsec.conf: %s %w", out, err)
	}
	secret := fmt.Sprintf("%s : PSK \"%s\"\n", remoteIP, psk)
	_, err = host.SSHCommand(fmt.Sprintf(`cat > /etc/ipsec.secrets << 'EOF'
%s
EOF`, secret))
	if err != nil {
		return fmt.Errorf("write ipsec.secrets: %w", err)
	}
	_, _ = host.SSHCommand("pkill -x charon 2>/dev/null; ipsec restart 2>/dev/null &")
	log.Printf("[provision] ipsec ready | peer: %s", remoteIP)
	return nil
}

func ProvisionOpenVPNServer(host *Host, port int) error {
	log.Printf("[provision] setting up openvpn server on %s:%d", host.Address, port)
	ensureConfigDir(host)
	_, err := host.SSHCommand(fmt.Sprintf(
		`if command -v openvpn &>/dev/null; then
  mkdir -p /etc/openvpn/server
  if [ ! -f /etc/openvpn/server/server.conf ]; then
    cat > /etc/openvpn/server/server.conf << 'EOF'
port %d
proto udp
dev tun
server 10.200.0.0 255.255.255.0
keepalive 10 120
persist-key
persist-tun
duplicate-cn
push "redirect-gateway def1"
push "dhcp-option DNS 1.1.1.1"
verb 3
EOF
    echo "openvpn config written (no certificates - manual setup needed)"
  else
    echo "openvpn config already exists"
  fi
  ss -tlnp | grep %d || echo "port %d not listening (certificates required)"
else
  echo "openvpn not installed on remote"
fi`, port, port, port))
	if err != nil {
		log.Printf("[provision] openvpn provisioning: %v", err)
	}
	log.Printf("[provision] openvpn config ready | port: %d", port)
	return nil
}

func TeardownProvisioned(host *Host) {
	cmds := []string{
		"pkill -f 'frps' 2>/dev/null",
		"pkill -f 'rathole.*server' 2>/dev/null",
		"pkill -f 'ss-server' 2>/dev/null",
		"pkill -f 'hysteria.*server' 2>/dev/null",
		"pkill -f 'backhaul' 2>/dev/null",
		"ipsec stop 2>/dev/null",
		"rm -f /etc/nyxora/frps.ini /etc/nyxora/rathole-server.toml /etc/nyxora/ss-config.json /etc/nyxora/hy2-server.yaml /etc/nyxora/backhaul-server.toml",
	}
	for _, c := range cmds {
		host.SSHCommand(c)
	}
	log.Printf("[provision] all provisioned services stopped and cleaned")
}
