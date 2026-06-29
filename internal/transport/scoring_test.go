package transport

import (
	"testing"
)

func TestComputeScore(t *testing.T) {
	tests := []struct {
		name     string
		metrics  *Metrics
		weights  ScoringWeights
		expected float64
	}{
		{
			name:     "zero packet loss, low latency",
			metrics:  &Metrics{LatencyMs: 50, PacketLoss: 0, JitterMs: 5, Stability: 0.9},
			weights:  DefaultScoringWeights,
			expected: 87.75,
		},
		{
			name:     "high packet loss returns 0",
			metrics:  &Metrics{LatencyMs: 50, PacketLoss: 60, JitterMs: 5, Stability: 0.9},
			weights:  DefaultScoringWeights,
			expected: 0,
		},
		{
			name:     "zero latency returns 5",
			metrics:  &Metrics{LatencyMs: 0, PacketLoss: 0, JitterMs: 0, Stability: 1.0},
			weights:  DefaultScoringWeights,
			expected: 5,
		},
		{
			name:     "perfect conditions",
			metrics:  &Metrics{LatencyMs: 10, PacketLoss: 0, JitterMs: 1, Stability: 1.0},
			weights:  DefaultScoringWeights,
			expected: 98.05,
		},
		{
			name:     "terrible conditions",
			metrics:  &Metrics{LatencyMs: 500, PacketLoss: 60, JitterMs: 50, Stability: 0.1},
			weights:  DefaultScoringWeights,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputeScore(tt.metrics, tt.weights)
			diff := result - tt.expected
			if diff < -0.1 || diff > 0.1 {
				t.Errorf("ComputeScore() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestComputeScoreWeights(t *testing.T) {
	m := &Metrics{LatencyMs: 100, PacketLoss: 5, JitterMs: 10, Stability: 0.8}

	w1 := ScoringWeights{Latency: 1.0, Loss: 0, Jitter: 0, Stability: 0}
	s1 := ComputeScore(m, w1)

	w2 := ScoringWeights{Latency: 0, Loss: 1.0, Jitter: 0, Stability: 0}
	s2 := ComputeScore(m, w2)

	if s1 == s2 {
		t.Error("Different weights should produce different scores")
	}
}

func TestCommandExists(t *testing.T) {
	if !CommandExists("sh") {
		t.Error("sh should exist on any Unix system")
	}
	if CommandExists("definitely-not-a-real-binary-xyz123") {
		t.Error("non-existent binary should not be found")
	}
}

func TestMeasureLatency(t *testing.T) {
	lat, loss, jitter := MeasureLatency("127.0.0.1", 2)
	if lat <= 0 {
		t.Errorf("latency to localhost should be positive, got %f", lat)
	}
	if loss != 0 {
		t.Errorf("loss to localhost should be 0, got %f", loss)
	}
	if jitter < 0 {
		t.Errorf("jitter should be non-negative, got %f", jitter)
	}
}

func TestMeasureLatencyEmpty(t *testing.T) {
	lat, loss, jitter := MeasureLatency("", 1)
	if lat != 999 {
		t.Errorf("empty addr should return 999 latency, got %f", lat)
	}
	if loss != 100 {
		t.Errorf("empty addr should return 100 loss, got %f", loss)
	}
	if jitter != 999 {
		t.Errorf("empty addr should return 999 jitter, got %f", jitter)
	}
}

func TestFormatEndpoint(t *testing.T) {
	tests := []struct {
		addr string
		port int
		want string
	}{
		{"1.2.3.4", 8080, "1.2.3.4:8080"},
		{"2001:db8::1", 443, "[2001:db8::1]:443"},
		{"::1", 80, "[::1]:80"},
	}
	for _, tt := range tests {
		got := FormatEndpoint(tt.addr, tt.port)
		if got != tt.want {
			t.Errorf("FormatEndpoint(%s, %d) = %s, want %s", tt.addr, tt.port, got, tt.want)
		}
	}
}

func TestExtractSubnet(t *testing.T) {
	tests := []struct {
		addr string
		want int
	}{
		{"10.100.108.2/24", 108},
		{"10.100.0.2/24", 0},
		{"10.100.255.2/24", 255},
		{"10.100.99.2/24", 99},
		{"invalid", 0},
	}
	for _, tt := range tests {
		got := ExtractSubnet(tt.addr)
		if got != tt.want {
			t.Errorf("ExtractSubnet(%s) = %d, want %d", tt.addr, got, tt.want)
		}
	}
}

func TestUpdateStability(t *testing.T) {
	m := &Metrics{PacketLoss: 2, LatencyMs: 50, Stability: 0.5}
	UpdateStability(m, 10, 200, 0.1, 0.1)
	if m.Stability != 0.6 {
		t.Errorf("Stability should increase, got %f", m.Stability)
	}

	m2 := &Metrics{PacketLoss: 30, LatencyMs: 500, Stability: 0.5}
	UpdateStability(m2, 10, 200, 0.1, 0.1)
	if m2.Stability != 0.4 {
		t.Errorf("Stability should decrease, got %f", m2.Stability)
	}

	m3 := &Metrics{PacketLoss: 0, LatencyMs: 10, Stability: 0.99}
	UpdateStability(m3, 10, 200, 0.1, 0.1)
	if m3.Stability != 1.0 {
		t.Errorf("Stability should cap at 1.0, got %f", m3.Stability)
	}

	m4 := &Metrics{PacketLoss: 50, LatencyMs: 500, Stability: 0.01}
	UpdateStability(m4, 10, 200, 0.1, 0.1)
	if m4.Stability != 0.0 {
		t.Errorf("Stability should floor at 0.0, got %f", m4.Stability)
	}
}
