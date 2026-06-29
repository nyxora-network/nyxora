package monitor

import (
	"testing"
	"time"
)

func TestNewMonitor(t *testing.T) {
	m := NewMonitor(30)
	if m == nil {
		t.Fatal("NewMonitor returned nil")
	}
	if m.interval != 30*time.Second {
		t.Errorf("interval should be 30s, got %v", m.interval)
	}
}

func TestMonitorPing(t *testing.T) {
	m := NewMonitor(30)
	result := m.Ping("127.0.0.1", 2)
	if result.LatencyMs <= 0 {
		t.Errorf("latency to localhost should be positive, got %f", result.LatencyMs)
	}
	if result.PacketLoss != 0 {
		t.Errorf("loss to localhost should be 0, got %f", result.PacketLoss)
	}
	if result.Jitter < 0 {
		t.Errorf("jitter should be non-negative, got %f", result.Jitter)
	}
	if result.Timestamp.IsZero() {
		t.Error("timestamp should be set")
	}
}

func TestMonitorPingEmpty(t *testing.T) {
	m := NewMonitor(30)
	result := m.Ping("", 1)
	if result.LatencyMs != 999 {
		t.Errorf("empty target latency should be 999, got %f", result.LatencyMs)
	}
	if result.PacketLoss != 100 {
		t.Errorf("empty target loss should be 100, got %f", result.PacketLoss)
	}
}

func TestMonitorStopWithoutStart(t *testing.T) {
	m := NewMonitor(30)
	m.Stop()
}

func TestMonitorHistoryEmpty(t *testing.T) {
	m := NewMonitor(30)
	history := m.History("nonexistent")
	if history != nil {
		t.Error("nonexistent target history should be nil")
	}
}

func TestMonitorLastResultEmpty(t *testing.T) {
	m := NewMonitor(30)
	_, ok := m.LastResult("nonexistent")
	if ok {
		t.Error("nonexistent target should return false")
	}
}

func TestMonitorAverageLatencyEmpty(t *testing.T) {
	m := NewMonitor(30)
	_, ok := m.AverageLatency("nonexistent", 1)
	if ok {
		t.Error("nonexistent target should return false")
	}
}
