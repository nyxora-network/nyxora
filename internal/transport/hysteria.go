package transport

import (
	"fmt"
	"log"
	"os/exec"
)

type Hysteria struct {
	BaseTransport
	authPass string
}

func NewHysteria() *Hysteria {
	return &Hysteria{
		BaseTransport: NewBase("hysteria", "hysteria", 8444, ScoringWeights{0.20, 0.35, 0.15, 0.30}, 300),
	}
}

func (h *Hysteria) Name() string  { return h.BaseName() }
func (h *Hysteria) Type() string { return h.BaseType() }

func (h *Hysteria) Init(cfg map[string]string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &h.port)
	}
	if a, ok := cfg["auth"]; ok {
		h.authPass = a
	}
	log.Printf("[hysteria] initialized (port: %d)", h.port)
	return nil
}

func (h *Hysteria) Connect(remoteAddr string) error {
	ctx := h.CancelContext()
	h.KillOldProcess()
	if err := h.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	if !CommandExists("hysteria") {
		h.Logf("binary not found, ping score only")
		h.SetStatusActive()
		return nil
	}

	auth := h.authPass
	if auth == "" {
		auth = "nyxora-hy2-fallback"
	}

	config := fmt.Sprintf(`server: %s:%d
auth: %s
tls:
  insecure: true
bandwidth:
  up: "100 mbps"
  down: "500 mbps"
socks5:
  listen: "127.0.0.1:1082"
`, remoteAddr, h.port, auth)

	tmpPath := fmt.Sprintf("/tmp/nyxora-hy2-%s.yaml", remoteAddr)
	if err := WriteConfig(tmpPath, config); err != nil {
		h.Logf("config error: %v", err)
		h.SetStatusActive()
		return nil
	}
	h.AddTmpFile(tmpPath)

	cmd := exec.Command("hysteria", "client", "-c", tmpPath)
	h.SetCmd(cmd)

	h.RunInBackground(func() {
		if err := cmd.Start(); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			h.Logf("start error: %v", err)
			return
		}
		h.SetStatusActive()
		h.KillOnCancel()
		cmd.Wait()
	})

	h.Logf("connecting to %s:%d", remoteAddr, h.port)
	h.SetStatusActive()
	return nil
}

func (h *Hysteria) Disconnect() error { return h.BaseDisconnect() }
func (h *Hysteria) Status() Status    { return h.BaseStatus() }
func (h *Hysteria) Metrics() *Metrics { return h.BaseMetrics() }
func (h *Hysteria) Health() bool      { return h.BaseHealth() }
func (h *Hysteria) Score() float64    { return h.BaseScore() }
