package transport

import (
	"fmt"
	"log"
	"os/exec"
)

type SSH struct {
	BaseTransport
	user      string
	password  string
	localPort int
}

func NewSSH() *SSH {
	return &SSH{
		BaseTransport: NewBase("ssh", "ssh", 22, ScoringWeights{0.30, 0.30, 0.15, 0.25}, 0),
		user:          "root",
		localPort:     1080,
	}
}

func (s *SSH) Name() string  { return s.BaseName() }
func (s *SSH) Type() string { return s.BaseType() }

func (s *SSH) Init(cfg map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &s.port)
	}
	if u, ok := cfg["user"]; ok {
		s.user = u
	}
	if pw, ok := cfg["password"]; ok {
		s.password = pw
	}
	log.Printf("[ssh] initialized (port: %d, user: %s)", s.port, s.user)
	return nil
}

func (s *SSH) Connect(remoteAddr string) error {
	ctx := s.CancelContext()
	s.KillOldProcess()
	if err := s.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	usePassword := s.password != "" && CommandExists("sshpass")

	if usePassword {
		tunnelCmd := exec.Command("sshpass",
			"-p", s.password,
			"ssh",
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"-o", "ServerAliveInterval=10",
			"-o", "ServerAliveCountMax=3",
			"-o", "ExitOnForwardFailure=yes",
			"-N",
			"-D", fmt.Sprintf("127.0.0.1:%d", s.localPort),
			"-p", fmt.Sprintf("%d", s.port),
			fmt.Sprintf("%s@%s", s.user, remoteAddr),
		)

		s.SetCmd(tunnelCmd)
		s.Logf("starting password tunnel %s@%s:%d -> 127.0.0.1:%d",
			s.user, remoteAddr, s.port, s.localPort)

		if err := tunnelCmd.Start(); err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			s.Logf("tunnel start error: %v", err)
			s.SetStatusFailed()
			return nil
		}

		s.SetStatusActive()
		s.KillOnCancel()

		go func() {
			if err := tunnelCmd.Wait(); err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				s.Logf("tunnel exited: %v", err)
				s.SetStatusFailed()
			}
		}()

		s.Logf("connecting via sshpass tunnel to %s:%d", remoteAddr, s.port)
		return nil
	}

	s.Logf("sshpass not available, using ping-only mode")
	s.SetStatusActive()
	return nil
}

func (s *SSH) Disconnect() error {
	return s.BaseDisconnect()
}

func (s *SSH) Status() Status  { return s.BaseStatus() }
func (s *SSH) Metrics() *Metrics { return s.BaseMetrics() }
func (s *SSH) Health() bool    { return s.BaseHealth() }
func (s *SSH) Score() float64  { return s.BaseScore() }
