package transport

import (
	"testing"
)

func TestNewWireGuard(t *testing.T) {
	wg := NewWireGuard()
	if wg == nil {
		t.Fatal("NewWireGuard returned nil")
	}
	if wg.Name() != "wireguard" {
		t.Errorf("name should be wireguard, got %s", wg.Name())
	}
	if wg.Type() != "wireguard" {
		t.Errorf("type should be wireguard, got %s", wg.Type())
	}
}

func TestWireGuardInit(t *testing.T) {
	wg := NewWireGuard()
	err := wg.Init(map[string]string{
		"private_key": "test-key",
		"remote_pub":  "test-pub",
		"local_addr":  "10.100.5.2/24",
		"interface":   "nyxora5",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if wg.iface != "nyxora5" {
		t.Errorf("interface should be nyxora5, got %s", wg.iface)
	}
}

func TestWireGuardStatus(t *testing.T) {
	wg := NewWireGuard()
	if wg.Status() != StatusInactive {
		t.Errorf("initial status should be inactive, got %s", wg.Status())
	}
}

func TestWireGuardHealth(t *testing.T) {
	wg := NewWireGuard()
	if wg.Health() {
		t.Error("initial health should be false")
	}
}

func TestWireGuardMetrics(t *testing.T) {
	wg := NewWireGuard()
	m := wg.Metrics()
	if m == nil {
		t.Fatal("Metrics should not be nil")
	}
}

func TestWireGuardScore(t *testing.T) {
	wg := NewWireGuard()
	s := wg.Score()
	if s < 0 || s > 100 {
		t.Errorf("score should be 0-100, got %f", s)
	}
}

func TestNewSSH(t *testing.T) {
	s := NewSSH()
	if s == nil {
		t.Fatal("NewSSH returned nil")
	}
	if s.Name() != "ssh" {
		t.Errorf("name should be ssh, got %s", s.Name())
	}
}

func TestSSHInit(t *testing.T) {
	s := NewSSH()
	err := s.Init(map[string]string{
		"port":     "2222",
		"user":     "admin",
		"password": "secret",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
}

func TestNewOpenVPN(t *testing.T) {
	o := NewOpenVPN()
	if o == nil {
		t.Fatal("NewOpenVPN returned nil")
	}
	if o.Name() != "openvpn" {
		t.Errorf("name should be openvpn, got %s", o.Name())
	}
}

func TestNewQUIC(t *testing.T) {
	q := NewQUIC()
	if q == nil {
		t.Fatal("NewQUIC returned nil")
	}
	if q.Name() != "quic" {
		t.Errorf("name should be quic, got %s", q.Name())
	}
}

func TestNewFRP(t *testing.T) {
	f := NewFRP()
	if f == nil {
		t.Fatal("NewFRP returned nil")
	}
	if f.Name() != "frp" {
		t.Errorf("name should be frp, got %s", f.Name())
	}
}

func TestNewRathole(t *testing.T) {
	r := NewRathole()
	if r == nil {
		t.Fatal("NewRathole returned nil")
	}
	if r.Name() != "rathole" {
		t.Errorf("name should be rathole, got %s", r.Name())
	}
}

func TestRatholeInit(t *testing.T) {
	r := NewRathole()
	err := r.Init(map[string]string{
		"port":  "2333",
		"token": "my-secret-token",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if r.token != "my-secret-token" {
		t.Errorf("token should be my-secret-token, got %s", r.token)
	}
}

func TestNewIPsec(t *testing.T) {
	i := NewIPsec()
	if i == nil {
		t.Fatal("NewIPsec returned nil")
	}
	if i.Name() != "ipsec" {
		t.Errorf("name should be ipsec, got %s", i.Name())
	}
}

func TestNewShadowSOCKS(t *testing.T) {
	s := NewShadowSOCKS()
	if s == nil {
		t.Fatal("NewShadowSOCKS returned nil")
	}
	if s.Name() != "shadowsocks" {
		t.Errorf("name should be shadowsocks, got %s", s.Name())
	}
}

func TestShadowSOCKSInit(t *testing.T) {
	s := NewShadowSOCKS()
	err := s.Init(map[string]string{
		"port":     "8388",
		"password": "mypass",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if s.password != "mypass" {
		t.Errorf("password should be mypass, got %s", s.password)
	}
}

func TestNewHysteria(t *testing.T) {
	h := NewHysteria()
	if h == nil {
		t.Fatal("NewHysteria returned nil")
	}
	if h.Name() != "hysteria" {
		t.Errorf("name should be hysteria, got %s", h.Name())
	}
}

func TestHysteriaInit(t *testing.T) {
	h := NewHysteria()
	err := h.Init(map[string]string{
		"port": "8444",
		"auth": "my-auth",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if h.authPass != "my-auth" {
		t.Errorf("auth should be my-auth, got %s", h.authPass)
	}
}

func TestNewBackhaul(t *testing.T) {
	b := NewBackhaul()
	if b == nil {
		t.Fatal("NewBackhaul returned nil")
	}
	if b.Name() != "backhaul" {
		t.Errorf("name should be backhaul, got %s", b.Name())
	}
}

func TestBackhaulInit(t *testing.T) {
	b := NewBackhaul()
	err := b.Init(map[string]string{
		"port":              "3080",
		"token":             "my-token",
		"backhaul_transport": "ws",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if b.token != "my-token" {
		t.Errorf("token should be my-token, got %s", b.token)
	}
	if b.transport != "ws" {
		t.Errorf("transport should be ws, got %s", b.transport)
	}
}

func TestNewTCP(t *testing.T) {
	tcp := NewTCP()
	if tcp == nil {
		t.Fatal("NewTCP returned nil")
	}
	if tcp.Name() != "tcp" {
		t.Errorf("name should be tcp, got %s", tcp.Name())
	}
}

func TestTCPInit(t *testing.T) {
	tcp := NewTCP()
	err := tcp.Init(map[string]string{
		"port":       "9924",
		"local_port": "9925",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
}

func TestBaseTransportDisconnect(t *testing.T) {
	wg := NewWireGuard()
	wg.SetStatusActive()
	if wg.Status() != StatusActive {
		t.Fatal("should be active")
	}
	err := wg.Disconnect()
	if err != nil {
		t.Fatalf("Disconnect failed: %v", err)
	}
	if wg.Status() != StatusInactive {
		t.Errorf("should be inactive after disconnect, got %s", wg.Status())
	}
}

func TestBaseTransportCancelContext(t *testing.T) {
	wg := NewWireGuard()
	ctx1 := wg.CancelContext()
	if ctx1 == nil {
		t.Fatal("CancelContext returned nil")
	}
	ctx2 := wg.CancelContext()
	if ctx2 == ctx1 {
		t.Error("CancelContext should return new context")
	}
}

func TestBaseTransportSetStatus(t *testing.T) {
	wg := NewWireGuard()
	wg.SetStatus(StatusTesting)
	if wg.Status() != StatusTesting {
		t.Errorf("status should be testing, got %s", wg.Status())
	}
	wg.SetStatusActive()
	if wg.Status() != StatusActive {
		t.Errorf("status should be active, got %s", wg.Status())
	}
	wg.SetStatusFailed()
	if wg.Status() != StatusFailed {
		t.Errorf("status should be failed, got %s", wg.Status())
	}
}

func TestBaseTransportKillOldProcess(t *testing.T) {
	wg := NewWireGuard()
	wg.KillOldProcess() // should not panic
}

func TestBaseTransportCleanTmpFiles(t *testing.T) {
	wg := NewWireGuard()
	wg.CleanTmpFiles() // should not panic
}

func TestBaseTransportLogf(t *testing.T) {
	wg := NewWireGuard()
	wg.Logf("test message %s", "hello") // should not panic
}

func TestAllTransportNames(t *testing.T) {
	transports := []Transport{
		NewWireGuard(),
		NewOpenVPN(),
		NewSSH(),
		NewQUIC(),
		NewFRP(),
		NewRathole(),
		NewIPsec(),
		NewShadowSOCKS(),
		NewHysteria(),
		NewBackhaul(),
		NewTCP(),
	}

	seen := make(map[string]bool)
	for _, tr := range transports {
		name := tr.Name()
		if seen[name] {
			t.Errorf("duplicate transport name: %s", name)
		}
		seen[name] = true

		if tr.Type() == "" {
			t.Errorf("transport %s has empty type", name)
		}
	}

	if len(transports) != 11 {
		t.Errorf("expected 11 transports, got %d", len(transports))
	}
}

func TestAllTransportsImplementInterface(t *testing.T) {
	transports := []Transport{
		NewWireGuard(),
		NewOpenVPN(),
		NewSSH(),
		NewQUIC(),
		NewFRP(),
		NewRathole(),
		NewIPsec(),
		NewShadowSOCKS(),
		NewHysteria(),
		NewBackhaul(),
		NewTCP(),
	}

	for _, tr := range transports {
		_ = tr.Name()
		_ = tr.Type()
		_ = tr.Status()
		_ = tr.Metrics()
		_ = tr.Health()
		_ = tr.Score()
	}
}
