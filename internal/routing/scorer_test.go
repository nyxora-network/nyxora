package routing

import (
	"testing"

	"github.com/nyxora-network/nyxora/internal/transport"
)

func TestNewScorer(t *testing.T) {
	s := NewScorer()
	if s == nil {
		t.Fatal("NewScorer returned nil")
	}
}

func TestScorerScore(t *testing.T) {
	s := NewScorer()

	score := s.Score(10, 1, 0, 1.0, 1000)
	if score < 80 {
		t.Errorf("perfect conditions score should be > 80, got %f", score)
	}

	score = s.Score(50, 5, 60, 0.8, 500)
	if score != 0 {
		t.Errorf("high loss score should be 0, got %f", score)
	}

	score = s.Score(0, 0, 0, 1.0, 100)
	if score != 5 {
		t.Errorf("zero latency score should be 5, got %f", score)
	}
}

func TestScorerRank(t *testing.T) {
	s := NewScorer()
	scores := []PathScore{
		{Name: "ssh", Score: 60},
		{Name: "wireguard", Score: 95},
		{Name: "quic", Score: 80},
	}

	ranked := s.Rank(scores)
	if ranked[0].Name != "wireguard" {
		t.Errorf("first should be wireguard, got %s", ranked[0].Name)
	}
	if ranked[1].Name != "quic" {
		t.Errorf("second should be quic, got %s", ranked[1].Name)
	}
	if ranked[2].Name != "ssh" {
		t.Errorf("third should be ssh, got %s", ranked[2].Name)
	}
}

func TestScorerBest(t *testing.T) {
	s := NewScorer()

	best := s.Best(nil)
	if best != nil {
		t.Error("Best of nil should be nil")
	}

	best = s.Best([]PathScore{})
	if best != nil {
		t.Error("Best of empty should be nil")
	}

	scores := []PathScore{
		{Name: "ssh", Score: 60},
		{Name: "wireguard", Score: 95},
	}
	best = s.Best(scores)
	if best == nil || best.Name != "wireguard" {
		t.Errorf("Best should be wireguard, got %v", best)
	}
}

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	if e == nil {
		t.Fatal("NewEngine returned nil")
	}
}

func TestEngineUpdate(t *testing.T) {
	e := NewEngine()
	e.Update("wireguard", "wireguard", 50, 5, 0, 0.9, 100)

	paths := e.AllPaths()
	if len(paths) != 1 {
		t.Errorf("expected 1 path, got %d", len(paths))
	}
	if paths[0].Name != "wireguard" {
		t.Errorf("expected wireguard, got %s", paths[0].Name)
	}
}

func TestEngineBestPath(t *testing.T) {
	e := NewEngine()
	e.Update("wireguard", "wireguard", 50, 5, 0, 0.9, 100)
	e.Update("ssh", "ssh", 150, 20, 5, 0.6, 50)

	best := e.BestPath()
	if best == nil {
		t.Fatal("BestPath should not be nil")
	}
	if best.Name != "wireguard" {
		t.Errorf("best path should be wireguard, got %s", best.Name)
	}
}

func TestEngineCurrent(t *testing.T) {
	e := NewEngine()
	if e.Current() != "" {
		t.Error("initial current should be empty")
	}

	e.SetCurrent("wireguard")
	if e.Current() != "wireguard" {
		t.Errorf("current should be wireguard, got %s", e.Current())
	}
}

func TestEngineNeedsFailover(t *testing.T) {
	e := NewEngine()
	e.Update("wireguard", "wireguard", 50, 5, 0, 0.9, 100)
	e.Update("ssh", "ssh", 200, 30, 20, 0.3, 30)
	e.SetCurrent("ssh")

	if !e.NeedsFailover(15) {
		t.Error("should need failover when score diff > threshold")
	}

	e.SetCurrent("wireguard")
	if e.NeedsFailover(15) {
		t.Error("should not need failover when best is current")
	}
}

func TestNewScorerWithWeights(t *testing.T) {
	w := transport.ScoringWeights{
		Latency:   0.50,
		Loss:      0.20,
		Jitter:    0.10,
		Stability: 0.20,
	}
	s := NewScorerWithWeights(w)

	scoreDefault := NewScorer().Score(100, 10, 5, 0.8, 100)
	scoreCustom := s.Score(100, 10, 5, 0.8, 100)

	if scoreDefault == scoreCustom {
		t.Error("custom weights should produce different score than default")
	}
}
