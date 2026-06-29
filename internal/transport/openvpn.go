package transport

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

type OpenVPN struct {
	BaseTransport
}

func NewOpenVPN() *OpenVPN {
	return &OpenVPN{
		BaseTransport: NewBase("openvpn", "openvpn", 1194, ScoringWeights{0.30, 0.30, 0.15, 0.25}, 50),
	}
}

func (o *OpenVPN) Name() string  { return o.BaseName() }
func (o *OpenVPN) Type() string { return o.BaseType() }

func (o *OpenVPN) Init(cfg map[string]string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &o.port)
	}
	log.Printf("[openvpn] initialized (port: %d)", o.port)
	return nil
}

func (o *OpenVPN) Connect(remoteAddr string) error {
	ctx := o.CancelContext()
	o.KillOldProcess()
	if err := o.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	if !CommandExists("openvpn") {
		o.Logf("binary not found, scoring from ping only")
		o.SetStatusActive()
		return nil
	}

	config := fmt.Sprintf(`client
dev tun
proto udp
remote %s %d
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
auth SHA256
cipher AES-256-CBC
verb 3
`, remoteAddr, o.port)

	tmpPath := fmt.Sprintf("/tmp/nyxora-openvpn-%d.conf", time.Now().UnixNano())
	if err := WriteConfig(tmpPath, config); err != nil {
		o.Logf("config write error: %v", err)
		o.SetStatusFailed()
		return err
	}
	o.AddTmpFile(tmpPath)

	cmd := exec.Command("openvpn", "--config", tmpPath, "--daemon")
	o.SetCmd(cmd)

	o.RunInBackground(func() {
		if err := cmd.Run(); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			o.Logf("start error: %v", err)
			o.SetStatusFailed()
			return
		}
		o.SetStatusActive()
		o.Logf("connected to %s:%d", remoteAddr, o.port)
	})

	o.SetStatusActive()
	return nil
}

func (o *OpenVPN) Disconnect() error {
	exec.Command("pkill", "-f", "openvpn.*nyxora").Run()
	return o.BaseDisconnect()
}

func (o *OpenVPN) Status() Status  { return o.BaseStatus() }
func (o *OpenVPN) Metrics() *Metrics { return o.BaseMetrics() }
func (o *OpenVPN) Health() bool    { return o.BaseHealth() }
func (o *OpenVPN) Score() float64  { return o.BaseScore() }
