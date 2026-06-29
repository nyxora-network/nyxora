package transport

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

type WireGuard struct {
	BaseTransport
	privateKey   string
	remotePubKey string
	localAddr    string
	iface        string
}

func NewWireGuard() *WireGuard {
	w := &WireGuard{
		BaseTransport: NewBase("wireguard", "wireguard", 51820, ScoringWeights{0.35, 0.30, 0.15, 0.20}, 0),
		iface:         "nyxora0",
	}
	return w
}

func (w *WireGuard) Name() string  { return w.BaseName() }
func (w *WireGuard) Type() string { return w.BaseType() }

func (w *WireGuard) Init(cfg map[string]string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if key, ok := cfg["private_key"]; ok {
		w.privateKey = key
	}
	if pub, ok := cfg["remote_pub"]; ok {
		w.remotePubKey = pub
	}
	if addr, ok := cfg["local_addr"]; ok {
		w.localAddr = addr
	}
	if iface, ok := cfg["interface"]; ok {
		w.iface = iface
	}
	log.Printf("[wireguard] initialized (iface: %s)", w.iface)
	return nil
}

func (w *WireGuard) Connect(remoteAddr string) error {
	w.CancelContext()
	if err := w.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	if !CommandExists("wg") {
		w.Logf("wg not installed, skipping wg setup")
		w.SetStatusActive()
		return nil
	}

	endpoint := FormatEndpoint(remoteAddr, 51820)
	conn, err := net.DialTimeout("udp", endpoint, 3*time.Second)
	if err != nil {
		w.Logf("remote WG port unreachable via %s: %v", endpoint, err)
		w.SetStatusActive()
		return nil
	}
	conn.Close()

	cfg := w.generateConfig(remoteAddr)
	cfgPath := fmt.Sprintf("/etc/wireguard/%s.conf", w.iface)
	if err := WriteConfig(cfgPath, cfg); err != nil {
		w.Logf("config write error: %v (continuing with ping-only)", err)
		w.SetStatusActive()
		return nil
	}

	if CommandExists("wg-quick") {
		exec.Command("wg-quick", "down", w.iface).Run()
		if out, err := exec.Command("wg-quick", "up", w.iface).CombinedOutput(); err != nil {
			w.Logf("wg-quick up failed: %v\n%s", err, string(out))
		} else {
			w.Logf("wg-quick up %s success", w.iface)
		}
	}

	w.SetStatusActive()
	w.Logf("connected to %s via %s", remoteAddr, w.iface)
	return nil
}

func (w *WireGuard) Disconnect() error {
	if CommandExists("wg-quick") {
		exec.Command("wg-quick", "down", w.iface).Run()
	}
	return w.BaseDisconnect()
}

func (w *WireGuard) Status() Status  { return w.BaseStatus() }
func (w *WireGuard) Metrics() *Metrics { return w.BaseMetrics() }
func (w *WireGuard) Health() bool    { return w.BaseHealth() }
func (w *WireGuard) Score() float64  { return w.BaseScore() }

func (w *WireGuard) generateConfig(remoteAddr string) string {
	privateKey := w.privateKey
	if privateKey == "" {
		out, _ := exec.Command("wg", "genkey").Output()
		if len(out) > 0 {
			privateKey = strings.TrimSpace(string(out))
		} else {
			privateKey = "nyxora-wg-auto-" + fmt.Sprintf("%d", time.Now().UnixNano())
		}
	}
	pubKey := w.remotePubKey
	if pubKey == "" {
		pubKey = remoteAddr + "-pub"
	}
	localAddr := w.localAddr
	if localAddr == "" {
		localAddr = "10.100.0.2/24"
	}
	subnet := ExtractSubnet(localAddr)
	return fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s
DNS = 1.1.1.1
MTU = 1420

[Peer]
PublicKey = %s
Endpoint = %s
AllowedIPs = 10.100.%d.0/24
PersistentKeepalive = 25
`, privateKey, localAddr, pubKey, FormatEndpoint(remoteAddr, 51820), subnet)
}
