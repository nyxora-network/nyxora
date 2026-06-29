package transport

import (
	"fmt"
	"log"
	"os/exec"
)

type Rathole struct {
	BaseTransport
	token string
}

func NewRathole() *Rathole {
	return &Rathole{
		BaseTransport: NewBase("rathole", "rathole", 2333, ScoringWeights{0.30, 0.35, 0.10, 0.25}, 200),
	}
}

func (r *Rathole) Name() string  { return r.BaseName() }
func (r *Rathole) Type() string { return r.BaseType() }

func (r *Rathole) Init(cfg map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &r.port)
	}
	if t, ok := cfg["token"]; ok {
		r.token = t
	}
	log.Printf("[rathole] initialized (port: %d)", r.port)
	return nil
}

func (r *Rathole) Connect(remoteAddr string) error {
	ctx := r.CancelContext()
	r.KillOldProcess()
	if err := r.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	if !CommandExists("rathole") {
		r.Logf("binary not found, scoring from ping")
		r.SetStatusActive()
		return nil
	}

	token := r.token
	if token == "" {
		token = "nyxora-rathole-fallback"
	}

	cfg := fmt.Sprintf(`[client]
remote_addr = "%s:%d"

[client.services.nyxora]
type = "tcp"
local_addr = "127.0.0.1:22"
token = "%s"
`, remoteAddr, r.port, token)

	tmpPath := fmt.Sprintf("/tmp/nyxora-rathole-%s.toml", remoteAddr)
	if err := WriteConfig(tmpPath, cfg); err != nil {
		r.Logf("config error: %v", err)
		r.SetStatusActive()
		return nil
	}
	r.AddTmpFile(tmpPath)

	cmd := exec.Command("rathole", "--client", tmpPath)
	r.SetCmd(cmd)

	r.RunInBackground(func() {
		if err := cmd.Start(); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			r.Logf("start error: %v", err)
			return
		}
		r.SetStatusActive()
		r.KillOnCancel()
		cmd.Wait()
	})

	r.Logf("connecting to %s:%d", remoteAddr, r.port)
	r.SetStatusActive()
	return nil
}

func (r *Rathole) Disconnect() error { return r.BaseDisconnect() }
func (r *Rathole) Status() Status    { return r.BaseStatus() }
func (r *Rathole) Metrics() *Metrics { return r.BaseMetrics() }
func (r *Rathole) Health() bool      { return r.BaseHealth() }
func (r *Rathole) Score() float64    { return r.BaseScore() }
