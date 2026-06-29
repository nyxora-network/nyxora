package transport

import (
	"fmt"
	"log"
	"os/exec"
)

type FRP struct {
	BaseTransport
}

func NewFRP() *FRP {
	return &FRP{
		BaseTransport: NewBase("frp", "frp", 7000, ScoringWeights{0.35, 0.30, 0.15, 0.20}, 100),
	}
}

func (f *FRP) Name() string  { return f.BaseName() }
func (f *FRP) Type() string { return f.BaseType() }

func (f *FRP) Init(cfg map[string]string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &f.port)
	}
	log.Printf("[frp] initialized (port: %d)", f.port)
	return nil
}

func (f *FRP) Connect(remoteAddr string) error {
	ctx := f.CancelContext()
	f.KillOldProcess()
	if err := f.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	if !CommandExists("frpc") {
		f.Logf("frpc not found, using ping score")
		f.SetStatusActive()
		return nil
	}

	cfg := fmt.Sprintf(`[common]
server_addr = %s
server_port = %d

[nyxora-tunnel]
type = tcp
local_ip = 127.0.0.1
local_port = 22
remote_port = 6000
`, remoteAddr, f.port)

	tmpPath := fmt.Sprintf("/tmp/nyxora-frpc-%s.ini", remoteAddr)
	if err := WriteConfig(tmpPath, cfg); err != nil {
		f.Logf("config error: %v", err)
		f.SetStatusActive()
		return nil
	}
	f.AddTmpFile(tmpPath)

	cmd := exec.Command("frpc", "-c", tmpPath)
	f.SetCmd(cmd)

	f.RunInBackground(func() {
		if err := cmd.Start(); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			f.Logf("start error: %v", err)
			return
		}
		f.SetStatusActive()
		f.KillOnCancel()
		cmd.Wait()
	})

	f.Logf("connecting to %s:%d", remoteAddr, f.port)
	f.SetStatusActive()
	return nil
}

func (f *FRP) Disconnect() error { return f.BaseDisconnect() }
func (f *FRP) Status() Status    { return f.BaseStatus() }
func (f *FRP) Metrics() *Metrics { return f.BaseMetrics() }
func (f *FRP) Health() bool      { return f.BaseHealth() }
func (f *FRP) Score() float64    { return f.BaseScore() }
