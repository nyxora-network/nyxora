package transport

import (
	"fmt"
	"log"
	"os/exec"
)

type IPsec struct {
	BaseTransport
	localIP  string
	targetIP string
	psk      string
}

func NewIPsec() *IPsec {
	return &IPsec{
		BaseTransport: NewBase("ipsec", "ipsec", 500, ScoringWeights{0.25, 0.25, 0.20, 0.30}, 500),
	}
}

func (i *IPsec) Name() string  { return i.BaseName() }
func (i *IPsec) Type() string { return i.BaseType() }

func (i *IPsec) Init(cfg map[string]string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &i.port)
	}
	if lip, ok := cfg["local_ip"]; ok {
		i.localIP = lip
	}
	if rip, ok := cfg["remote_ip"]; ok {
		i.targetIP = rip
	}
	if p, ok := cfg["psk"]; ok {
		i.psk = p
	}
	log.Printf("[ipsec] initialized (port: %d)", i.port)
	return nil
}

func (i *IPsec) Connect(remoteAddr string) error {
	ctx := i.CancelContext()
	if err := i.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	if !CommandExists("ipsec") {
		i.Logf("strongswan not installed, ping score only")
		i.SetStatusActive()
		return nil
	}

	connectTo := i.targetIP
	if connectTo == "" {
		connectTo = remoteAddr
	}

	psk := i.psk
	if psk == "" {
		psk = "nyxora-ipsec-fallback"
	}

	secret := fmt.Sprintf("%s : PSK \"%s\"\n", connectTo, psk)
	WriteConfig("/etc/ipsec.secrets", secret)

	leftIP := i.localIP
	if leftIP == "" {
		leftIP = "%any"
	}

	config := fmt.Sprintf(`config setup

conn nyxora
    left=%s
    right=%s
    authby=secret
    ike=aes256-sha256-modp2048
    esp=aes256-sha256
    auto=start
    type=transport
`, leftIP, connectTo)

	WriteConfig("/etc/ipsec.conf", config)

	i.RunInBackground(func() {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err := exec.Command("ipsec", "restart").Run(); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			i.Logf("restart error: %v", err)
			i.SetStatusFailed()
			return
		}
		exec.Command("ipsec", "up", "nyxora").Run()
	})

	i.Logf("connecting to %s", connectTo)
	i.SetStatusActive()
	return nil
}

func (i *IPsec) Disconnect() error {
	exec.Command("ipsec", "down", "nyxora").Run()
	exec.Command("ipsec", "stop").Run()
	return i.BaseDisconnect()
}

func (i *IPsec) Status() Status  { return i.BaseStatus() }
func (i *IPsec) Metrics() *Metrics { return i.BaseMetrics() }
func (i *IPsec) Health() bool    { return i.BaseHealth() }
func (i *IPsec) Score() float64  { return i.BaseScore() }
