package multipath

import (
	"testing"
)

func TestNewScheduler(t *testing.T) {
	s := NewScheduler()
	if s == nil {
		t.Fatal("NewScheduler returned nil")
	}
	if s.mode != ModeWeighted {
		t.Errorf("default mode should be ModeWeighted, got %v", s.mode)
	}
}

func TestAddPath(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	s.AddPath("ssh", "ssh", 5)

	if len(s.paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(s.paths))
	}
	if s.paths["wireguard"].Weight != 30 {
		t.Errorf("wireguard weight should be 30, got %d", s.paths["wireguard"].Weight)
	}
}

func TestRemovePath(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	s.AddPath("ssh", "ssh", 5)
	s.RemovePath("ssh")

	if len(s.paths) != 1 {
		t.Errorf("expected 1 path after remove, got %d", len(s.paths))
	}
	if _, ok := s.paths["ssh"]; ok {
		t.Error("ssh should have been removed")
	}
}

func TestUpdatePath(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	s.AddPath("ssh", "ssh", 5)

	s.UpdatePath("wireguard", 90, 50, 0, 0)
	s.UpdatePath("ssh", 60, 100, 5, 0)

	if !s.paths["wireguard"].Active {
		t.Error("wireguard should be active with score > 0 and loss < 50")
	}
	if !s.paths["ssh"].Active {
		t.Error("ssh should be active with score > 0 and loss < 50")
	}
}

func TestUpdatePathHighLoss(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	s.UpdatePath("wireguard", 90, 50, 60, 0)

	if s.paths["wireguard"].Active {
		t.Error("wireguard should be inactive with loss > 50")
	}
}

func TestSelectPath(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	s.AddPath("ssh", "ssh", 5)

	s.UpdatePath("wireguard", 90, 50, 0, 0)
	s.UpdatePath("ssh", 60, 100, 5, 0)

	best := s.BestPath()
	if best == nil {
		t.Fatal("BestPath should not be nil")
	}
	if best.Name != "wireguard" {
		t.Errorf("best path should be wireguard, got %s", best.Name)
	}
}

func TestSelectPathNoneActive(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)

	selected := s.SelectPath()
	if selected != nil {
		t.Error("SelectPath should return nil when no paths active")
	}
}

func TestDistribution(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	s.AddPath("ssh", "ssh", 5)

	s.UpdatePath("wireguard", 90, 50, 0, 100)
	s.UpdatePath("ssh", 60, 100, 5, 50)

	dist := s.Distribution()
	if len(dist) != 2 {
		t.Errorf("distribution should have 2 entries, got %d", len(dist))
	}
	// Weights are dynamically adjusted by score; check they're positive
	if dist["wireguard"] <= 0 {
		t.Errorf("wireguard weight should be positive, got %d", dist["wireguard"])
	}
	if dist["ssh"] <= 0 {
		t.Errorf("ssh weight should be positive, got %d", dist["ssh"])
	}
}

func TestAggregateBandwidth(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	s.AddPath("ssh", "ssh", 5)

	s.UpdatePath("wireguard", 90, 50, 0, 100)
	s.UpdatePath("ssh", 60, 100, 5, 50)

	bw := s.AggregateBandwidth()
	if bw != 150 {
		t.Errorf("aggregate bandwidth should be 150, got %d", bw)
	}
}

func TestRecordFailover(t *testing.T) {
	s := NewScheduler()
	s.RecordFailover()
	s.RecordFailover()

	stats := s.Stats()
	if stats.FailoverCount != 2 {
		t.Errorf("failover count should be 2, got %d", stats.FailoverCount)
	}
}

func TestRecordBytes(t *testing.T) {
	s := NewScheduler()
	s.RecordBytes(1000, 2000)
	s.RecordBytes(500, 1000)

	stats := s.Stats()
	if stats.TotalBytesSent != 1500 {
		t.Errorf("total bytes sent should be 1500, got %d", stats.TotalBytesSent)
	}
	if stats.TotalBytesReceived != 3000 {
		t.Errorf("total bytes received should be 3000, got %d", stats.TotalBytesReceived)
	}
}

func TestSetMode(t *testing.T) {
	s := NewScheduler()
	s.SetMode(ModeLowestLatency)
	if s.mode != ModeLowestLatency {
		t.Errorf("mode should be ModeLowestLatency, got %v", s.mode)
	}
}

func TestModeFromString(t *testing.T) {
	tests := []struct {
		input string
		want  DistributionMode
	}{
		{"weighted", ModeWeighted},
		{"lowest-latency", ModeLowestLatency},
		{"lowest-loss", ModeLowestLoss},
		{"even", ModeEven},
		{"all", ModeAll},
		{"unknown", ModeWeighted},
	}
	for _, tt := range tests {
		got := ModeFromString(tt.input)
		if got != tt.want {
			t.Errorf("ModeFromString(%s) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestAllPaths(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	s.AddPath("ssh", "ssh", 5)
	s.AddPath("quic", "quic", 15)

	s.UpdatePath("wireguard", 90, 50, 0, 0)
	s.UpdatePath("ssh", 60, 100, 5, 0)
	s.UpdatePath("quic", 80, 70, 2, 0)

	paths := s.AllPaths()
	if len(paths) != 3 {
		t.Errorf("AllPaths should return 3 paths, got %d", len(paths))
	}

	if paths[0].Score < paths[1].Score || paths[1].Score < paths[2].Score {
		t.Error("AllPaths should be sorted by score descending")
	}
}

func TestSelectPaths(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	s.AddPath("ssh", "ssh", 5)
	s.AddPath("quic", "quic", 15)

	s.UpdatePath("wireguard", 90, 50, 0, 0)
	s.UpdatePath("ssh", 60, 100, 5, 0)
	s.UpdatePath("quic", 80, 70, 2, 0)

	top2 := s.SelectPaths(2)
	if len(top2) != 2 {
		t.Errorf("SelectPaths(2) should return 2 paths, got %d", len(top2))
	}
	if top2[0].Name != "wireguard" {
		t.Errorf("first path should be wireguard, got %s", top2[0].Name)
	}
}

func TestString(t *testing.T) {
	s := NewScheduler()
	s.AddPath("wireguard", "wireguard", 30)
	str := s.String()
	if str == "" {
		t.Error("String() should not be empty")
	}
}
