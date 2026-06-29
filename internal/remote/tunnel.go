package remote

import (
	"fmt"
	"log"
	"strings"
)

type TunnelEndpoint struct {
	Host       *Host
	Type       string
	Port       int
	ConfigPath string
	Interface  string
	Status     string
}

func SetupWireGuardRemote(host *Host, localPublicKey string, listenPort int) (string, error) {
	log.Printf("[tunnel] setting up wireguard on remote %s", host.Address)

	privKey, err := host.SSHCommand("wg genkey")
	if err != nil {
		return "", fmt.Errorf("generate wg key: %w", err)
	}
	privKey = strings.TrimSpace(privKey)

	pubKey, err := host.SSHCommand(fmt.Sprintf("echo '%s' | wg pubkey", privKey))
	if err != nil {
		return "", fmt.Errorf("derive pubkey: %w", err)
	}
	pubKey = strings.TrimSpace(pubKey)

	iface := fmt.Sprintf("nyxora%d", listenPort)
	cfg := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = 10.100.%d.1/24
ListenPort = %d
MTU = 1420

[Peer]
PublicKey = %s
AllowedIPs = 10.100.%d.0/24
PersistentKeepalive = 25
`, privKey, listenPort%256, listenPort, localPublicKey, listenPort%256)

	cfgPath := fmt.Sprintf("/etc/wireguard/%s.conf", iface)
	if err := host.WriteFile(cfgPath, cfg, "600"); err != nil {
		return "", fmt.Errorf("write remote wg config: %w", err)
	}

	_, err = host.SSHCommand(fmt.Sprintf("wg-quick up %s 2>&1", iface))
	if err != nil {
		log.Printf("[tunnel] wg-quick up failed on remote (may already be up): %v", err)
	}

	log.Printf("[tunnel] remote wireguard ready | pubkey: %s | iface: %s", pubKey[:16]+"...", iface)
	return pubKey, nil
}

func TeardownRemote(host *Host, iface string) error {
	_, err := host.SSHCommand(fmt.Sprintf("wg-quick down %s 2>/dev/null; ip link delete %s 2>/dev/null", iface, iface))
	return err
}

func CheckTunnelHealth(host *Host, iface string) bool {
	out, err := host.SSHCommand(fmt.Sprintf("wg show %s 2>/dev/null | head -5", iface))
	return err == nil && strings.Contains(out, "interface:")
}

func GetRemotePublicIP(host *Host) (string, error) {
	out, err := host.SSHCommand("curl -s ifconfig.me 2>/dev/null || curl -s ipinfo.io/ip 2>/dev/null || hostname -I 2>/dev/null | awk '{print $1}'")
	if err != nil {
		return "", fmt.Errorf("get remote ip: %w", err)
	}
	return strings.TrimSpace(out), nil
}
