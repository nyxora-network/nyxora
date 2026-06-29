package failover

import (
	"testing"
	"time"
)

func TestNewFailover(t *testing.T) {
	f := NewFailover(15)
	if f == nil {
		t.Fatal("NewFailover returned nil")
	}
	if f.interval != 15*time.Second {
		t.Errorf("interval should be 15s, got %v", f.interval)
	}
}

func TestUpdateHealthy(t *testing.T) {
	f := NewFailover(15)
	f.Update("wireguard", 50, 2)

	state, ok := f.states["wireguard"]
	if !ok {
		t.Fatal("state should exist after Update")
	}
	if state.Status != StatusHealthy {
		t.Errorf("status should be Healthy, got %v", state.Status)
	}
}

func TestUpdateDegraded(t *testing.T) {
	f := NewFailover(15)
	f.Update("wireguard", 300, 15)

	state := f.states["wireguard"]
	if state.Status != StatusDegraded {
		t.Errorf("status should be Degraded, got %v", state.Status)
	}
	if state.FailCount != 1 {
		t.Errorf("fail count should be 1, got %d", state.FailCount)
	}
}

func TestUpdateDown(t *testing.T) {
	f := NewFailover(15)
	for i := 0; i < 4; i++ {
		f.Update("wireguard", 300, 15)
	}

	state := f.states["wireguard"]
	if state.Status != StatusDown {
		t.Errorf("status should be Down after %d failures, got %v", DefaultThreshold.MaxFailCount+1, state.Status)
	}
}

func TestUpdateRecovery(t *testing.T) {
	f := NewFailover(15)
	f.Update("wireguard", 300, 15)
	f.Update("wireguard", 300, 15)

	f.Update("wireguard", 50, 2)
	state := f.states["wireguard"]
	if state.Status != StatusHealthy {
		t.Errorf("status should recover to Healthy, got %v", state.Status)
	}
	if state.FailCount != 0 {
		t.Errorf("fail count should reset to 0, got %d", state.FailCount)
	}
}

func TestIsHealthy(t *testing.T) {
	f := NewFailover(15)
	if f.IsHealthy("nonexistent") {
		t.Error("nonexistent transport should not be healthy")
	}

	f.Update("wireguard", 50, 2)
	if !f.IsHealthy("wireguard") {
		t.Error("wireguard should be healthy")
	}
}

func TestOnFailoverCallback(t *testing.T) {
	f := NewFailover(15)
	called := false
	var from, to string

	f.OnFailover(func(f, t string) {
		called = true
		from = f
		to = t
	})

	cb := f.GetOnFailover()
	if cb == nil {
		t.Fatal("OnFailover callback should be set")
	}

	cb("wireguard", "ssh")
	if !called {
		t.Error("OnFailover callback should have been called")
	}
	if from != "wireguard" || to != "ssh" {
		t.Errorf("callback args: got from=%s to=%s", from, to)
	}
}

func TestAllStatus(t *testing.T) {
	f := NewFailover(15)
	f.Update("wireguard", 50, 2)
	f.Update("ssh", 300, 15)

	all := f.AllStatus()
	if len(all) != 2 {
		t.Errorf("AllStatus should return 2 entries, got %d", len(all))
	}
	if all["wireguard"] != StatusHealthy {
		t.Error("wireguard should be healthy")
	}
	if all["ssh"] != StatusDegraded {
		t.Error("ssh should be degraded")
	}
}

func TestStatusUnknown(t *testing.T) {
	f := NewFailover(15)
	s := f.Status("nonexistent")
	if s != StatusDown {
		t.Errorf("unknown transport should return StatusDown, got %v", s)
	}
}
