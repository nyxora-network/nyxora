package transport

import (
	"fmt"
	"log"
	"os/exec"
)

type Backhaul struct {
	BaseTransport
	token     string
	transport string
}

func NewBackhaul() *Backhaul {
	return &Backhaul{
		BaseTransport: NewBase("backhaul", "backhaul", 3080, ScoringWeights{0.30, 0.30, 0.15, 0.25}, 150),
		transport:    "tcp",
	}
}

func (b *Backhaul) Name() string  { return b.BaseName() }
func (b *Backhaul) Type() string { return b.BaseType() }

func (b *Backhaul) Init(cfg map[string]string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &b.port)
	}
	if t, ok := cfg["token"]; ok {
		b.token = t
	}
	if tr, ok := cfg["backhaul_transport"]; ok {
		b.transport = tr
	}
	log.Printf("[backhaul] initialized (port: %d, transport: %s)", b.port, b.transport)
	return nil
}

func (b *Backhaul) Connect(remoteAddr string) error {
	ctx := b.CancelContext()
	b.KillOldProcess()
	if err := b.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	if !CommandExists("backhaul") {
		b.Logf("binary not found locally, ping score only")
		b.SetStatusActive()
		return nil
	}

	token := b.token
	if token == "" {
		token = "nyxora-bh-fallback"
	}

	clientCfg := fmt.Sprintf(`[client]
remote_addr = "%s:%d"
transport = "%s"
token = "%s"
connection_pool = 4
keepalive_period = 75
dial_timeout = 10
nodelay = true
retry_interval = 3
log_level = "error"
`, remoteAddr, b.port, b.transport, token)

	clientCfgPath := "/etc/nyxora/backhaul-client.toml"
	if err := WriteConfig(clientCfgPath, clientCfg); err != nil {
		b.Logf("client config error: %v", err)
		b.SetStatusFailed()
		return err
	}
	b.AddTmpFile(clientCfgPath)

	cmd := exec.Command("backhaul", "-c", clientCfgPath)
	b.SetCmd(cmd)

	b.RunInBackground(func() {
		if err := cmd.Start(); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			b.Logf("client start error: %v", err)
			return
		}
		b.SetStatusActive()
		b.KillOnCancel()
		cmd.Wait()
	})

	b.Logf("connecting %s to %s:%d", b.transport, remoteAddr, b.port)
	b.SetStatusActive()
	return nil
}

func (b *Backhaul) Disconnect() error {
	exec.Command("pkill", "-f", "backhaul.*nyxora").Run()
	return b.BaseDisconnect()
}

func (b *Backhaul) Status() Status  { return b.BaseStatus() }
func (b *Backhaul) Metrics() *Metrics { return b.BaseMetrics() }
func (b *Backhaul) Health() bool    { return b.BaseHealth() }
func (b *Backhaul) Score() float64  { return b.BaseScore() }
