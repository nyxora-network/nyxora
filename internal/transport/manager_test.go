package transport

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// mockTransport implements the Transport interface for testing.
type mockTransport struct {
	name       string
	transport  string
	connectErr error
	score      float64
	connected  atomic.Bool
}

func (m *mockTransport) Name() string             { return m.name }
func (m *mockTransport) Type() string             { return m.transport }
func (m *mockTransport) Status() Status           { return StatusActive }
func (m *mockTransport) Health() bool             { return m.connected.Load() }
func (m *mockTransport) Score() float64           { return m.score }
func (m *mockTransport) Metrics() *Metrics        { return &Metrics{Bandwidth: 1000} }
func (m *mockTransport) Init(cfg map[string]string) error { return nil }

func (m *mockTransport) Connect(addr string) error {
	if m.connectErr != nil {
		return m.connectErr
	}
	m.connected.Store(true)
	return nil
}

func (m *mockTransport) Disconnect() error {
	m.connected.Store(false)
	return nil
}

func newMock(name string, score float64, connectErr error) *mockTransport {
	return &mockTransport{
		name:       name,
		transport:  name,
		score:      score,
		connectErr: connectErr,
	}
}

// --- Manager Tests ---

func TestManagerAllActive(t *testing.T) {
	m := NewManager(true)
	if !m.IsAllActive() {
		t.Error("allActive should be true")
	}
}

func TestManagerGet(t *testing.T) {
	m := NewManager(false)
	tr := newMock("wg", 95, nil)
	m.Register(tr)

	got, ok := m.Get("wg")
	if !ok {
		t.Fatal("Get should find wg")
	}
	if got.Name() != "wg" {
		t.Errorf("expected wg, got %s", got.Name())
	}

	_, ok = m.Get("nonexistent")
	if ok {
		t.Error("Get should not find nonexistent")
	}
}

func TestManagerBestMode_SingleActive(t *testing.T) {
	m := NewManager(false) // best mode

	t1 := newMock("ssh", 60, nil)
	t2 := newMock("wireguard", 95, nil)
	t3 := newMock("quic", 80, nil)

	m.Register(t1)
	m.Register(t2)
	m.Register(t3)

	err := m.ConnectAll("192.168.1.1:443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	active := m.ActiveNames()
	if len(active) != 1 {
		t.Fatalf("expected 1 active transport, got %d: %v", len(active), active)
	}
	if active[0] != "wireguard" {
		t.Errorf("expected wireguard (highest score), got %s", active[0])
	}
}

func TestManagerAllMode_MultipleActive(t *testing.T) {
	m := NewManager(true) // all-active mode

	t1 := newMock("ssh", 60, nil)
	t2 := newMock("wireguard", 95, nil)
	t3 := newMock("quic", 80, nil)

	m.Register(t1)
	m.Register(t2)
	m.Register(t3)

	err := m.ConnectAll("192.168.1.1:443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	active := m.ActiveNames()
	if len(active) != 3 {
		t.Errorf("expected 3 active transports, got %d: %v", len(active), active)
	}
}

func TestManagerConnectAll_AllFail(t *testing.T) {
	m := NewManager(false)

	t1 := newMock("ssh", 0, errors.New("connection refused"))
	t2 := newMock("wireguard", 0, errors.New("timeout"))

	m.Register(t1)
	m.Register(t2)

	err := m.ConnectAll("192.168.1.1:443")
	if err == nil {
		t.Fatal("expected error when all transports fail")
	}

	active := m.ActiveNames()
	if len(active) != 0 {
		t.Errorf("expected 0 active transports, got %d", len(active))
	}
}

func TestManagerConnectAll_PartialFailure(t *testing.T) {
	m := NewManager(false)

	t1 := newMock("ssh", 60, errors.New("refused"))
	t2 := newMock("wireguard", 95, nil)
	t3 := newMock("quic", 80, nil)

	m.Register(t1)
	m.Register(t2)
	m.Register(t3)

	err := m.ConnectAll("192.168.1.1:443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	active := m.ActiveNames()
	if len(active) != 1 {
		t.Errorf("expected 1 active, got %d: %v", len(active), active)
	}
}

func TestManagerActivateDeactivate(t *testing.T) {
	m := NewManager(false)
	tr := newMock("wg", 95, nil)
	m.Register(tr)

	if !m.Activate("wg") {
		t.Error("Activate should return true for registered transport")
	}
	if !m.IsActive("wg") {
		t.Error("IsActive should return true after Activate")
	}

	m.Deactivate("wg")
	if m.IsActive("wg") {
		t.Error("IsActive should return false after Deactivate")
	}

	if m.Activate("nonexistent") {
		t.Error("Activate should return false for unregistered transport")
	}
}

func TestManagerDisconnectAll(t *testing.T) {
	m := NewManager(false)
	t1 := newMock("ssh", 60, nil)
	t2 := newMock("wg", 95, nil)
	m.Register(t1)
	m.Register(t2)

	m.ConnectAll("192.168.1.1:443")
	if m.ActiveCount() == 0 {
		t.Fatal("should have active transports after ConnectAll")
	}

	m.DisconnectAll()
	if m.ActiveCount() != 0 {
		t.Errorf("expected 0 active after DisconnectAll, got %d", m.ActiveCount())
	}
}

func TestManagerWeights(t *testing.T) {
	m := NewManager(false)
	tr := newMock("wg", 95, nil)
	m.Register(tr)

	m.SetWeight("wg", 30)
	w := m.GetWeights()
	if w["wg"] != 30 {
		t.Errorf("expected weight 30, got %d", w["wg"])
	}

	m.NormalizeWeights()
	w = m.GetWeights()
	total := 0
	for _, v := range w {
		total += v
	}
	if total != 100 {
		t.Errorf("expected total weight 100 after normalize, got %d", total)
	}
}

func TestManagerActiveCount(t *testing.T) {
	m := NewManager(false)
	m.Register(newMock("a", 50, nil))
	m.Register(newMock("b", 60, nil))
	m.Register(newMock("c", 70, nil))

	if m.ActiveCount() != 0 {
		t.Errorf("expected 0 active initially, got %d", m.ActiveCount())
	}

	m.ConnectAll("192.168.1.1:443")
	if m.ActiveCount() != 1 {
		t.Errorf("expected 1 active in best mode, got %d", m.ActiveCount())
	}
}

func TestManagerList_SortedByScore(t *testing.T) {
	m := NewManager(false)
	m.Register(newMock("low", 30, nil))
	m.Register(newMock("high", 90, nil))
	m.Register(newMock("mid", 60, nil))

	list := m.List()
	if len(list) != 3 {
		t.Fatalf("expected 3, got %d", len(list))
	}
	if list[0].Name != "high" || list[1].Name != "mid" || list[2].Name != "low" {
		t.Errorf("list not sorted by score: %v", list)
	}
}

// --- Concurrent Connection Test ---

func TestManagerConnectAll_Concurrent(t *testing.T) {
	m := NewManager(false)

	for i := 0; i < 10; i++ {
		tr := newMock(fmt.Sprintf("t%d", i), float64(i*10), nil)
		m.Register(tr)
	}

	start := time.Now()
	err := m.ConnectAll("192.168.1.1:443")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if elapsed > 5*time.Second {
		t.Errorf("connections took too long: %v (expected parallel execution)", elapsed)
	}

	active := m.ActiveNames()
	if len(active) != 1 {
		t.Errorf("expected 1 active in best mode, got %d", len(active))
	}
}

func TestManagerConnectAll_Timeout(t *testing.T) {
	m := NewManager(false)

	slowTr := &slowMockTransport{
		mockTransport: mockTransport{name: "slow", transport: "slow", score: 99},
		delay:         35 * time.Second,
	}
	m.Register(slowTr)

	done := make(chan error, 1)
	go func() {
		done <- m.ConnectAll("192.168.1.1:443")
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Log("slow transport might have connected within timeout")
		}
	case <-time.After(40 * time.Second):
		t.Fatal("ConnectAll did not return within timeout")
	}
}

// slowMockTransport simulates a transport with slow connection.
type slowMockTransport struct {
	mockTransport
	delay time.Duration
}

func (s *slowMockTransport) Connect(addr string) error {
	time.Sleep(s.delay)
	return nil
}

// --- Registry Tests ---

func TestListTunnels(t *testing.T) {
	tunnels := ListTunnels()
	if len(tunnels) != 12 {
		t.Errorf("expected 12 tunnels, got %d", len(tunnels))
	}
}
