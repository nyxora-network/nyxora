package transport

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

type ShadowSOCKS struct {
	BaseTransport
	password string
	method   string
}

func NewShadowSOCKS() *ShadowSOCKS {
	return &ShadowSOCKS{
		BaseTransport: NewBase("shadowsocks", "shadowsocks", 8388, ScoringWeights{0.20, 0.30, 0.20, 0.30}, 30),
		method:        "aes-256-gcm",
		password:      fmt.Sprintf("nyxora-ss-%d", time.Now().Unix()),
	}
}

func (s *ShadowSOCKS) Name() string  { return s.BaseName() }
func (s *ShadowSOCKS) Type() string { return s.BaseType() }

func (s *ShadowSOCKS) Init(cfg map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &s.port)
	}
	if pw, ok := cfg["password"]; ok {
		s.password = pw
	}
	log.Printf("[shadowsocks] initialized (port: %d)", s.port)
	return nil
}

func (s *ShadowSOCKS) Connect(remoteAddr string) error {
	ctx := s.CancelContext()
	s.KillOldProcess()
	if err := s.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	if !CommandExists("ss-local") && !CommandExists("ss-redir") {
		s.Logf("binary not found, ping score only")
		s.SetStatusActive()
		return nil
	}

	config := fmt.Sprintf(`{
    "server": "%s",
    "server_port": %d,
    "local_port": 1081,
    "password": "%s",
    "method": "%s",
    "timeout": 60
}`, remoteAddr, s.port, s.password, s.method)

	tmpPath := fmt.Sprintf("/tmp/nyxora-ss-%s.json", remoteAddr)
	if err := WriteConfig(tmpPath, config); err != nil {
		s.Logf("config error: %v", err)
		s.SetStatusActive()
		return nil
	}
	s.AddTmpFile(tmpPath)

	binary := "ss-local"
	if !CommandExists(binary) {
		binary = "ss-redir"
	}

	cmd := exec.Command(binary, "-c", tmpPath)
	s.SetCmd(cmd)

	s.RunInBackground(func() {
		if err := cmd.Start(); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			s.Logf("start error: %v", err)
			return
		}
		s.SetStatusActive()
		s.KillOnCancel()
		cmd.Wait()
	})

	s.Logf("connecting to %s:%d", remoteAddr, s.port)
	s.SetStatusActive()
	return nil
}

func (s *ShadowSOCKS) Disconnect() error { return s.BaseDisconnect() }
func (s *ShadowSOCKS) Status() Status    { return s.BaseStatus() }
func (s *ShadowSOCKS) Metrics() *Metrics { return s.BaseMetrics() }
func (s *ShadowSOCKS) Health() bool      { return s.BaseHealth() }
func (s *ShadowSOCKS) Score() float64    { return s.BaseScore() }
