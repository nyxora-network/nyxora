package transport

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	latencyCmdCache = make(map[string]bool)
	latencyCmdMu    sync.Mutex
)

func CommandExists(name string) bool {
	latencyCmdMu.Lock()
	defer latencyCmdMu.Unlock()
	if cached, ok := latencyCmdCache[name]; ok {
		return cached
	}
	_, err := exec.LookPath(name)
	latencyCmdCache[name] = err == nil
	return err == nil
}

func MeasureLatency(addr string, count int) (latency, packetLoss, jitter float64) {
	if addr == "" {
		return 999, 100, 999
	}
	var rtts []float64
	var lossCount float64
	for i := 0; i < count; i++ {
		start := time.Now()
		cmd := exec.Command("ping", "-c", "1", "-W", "2", addr)
		if err := cmd.Run(); err == nil {
			rtt := time.Since(start).Seconds() * 1000
			rtts = append(rtts, rtt)
		} else {
			lossCount++
		}
	}
	packetLoss = (lossCount / float64(count)) * 100
	if len(rtts) == 0 {
		return 999, 100, 999
	}
	var sum float64
	for _, r := range rtts {
		sum += r
	}
	latency = sum / float64(len(rtts))
	if len(rtts) > 1 {
		var jSum float64
		for i := 1; i < len(rtts); i++ {
			diff := rtts[i] - rtts[i-1]
			if diff < 0 {
				diff = -diff
			}
			jSum += diff
		}
		jitter = jSum / float64(len(rtts)-1)
	}
	return
}

func WriteConfig(path, content string) error {
	dir := path[:len(path)-len(path[len(path)-1:])]
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			dir = path[:i]
			break
		}
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func FormatEndpoint(addr string, port int) string {
	if strings.Contains(addr, ":") {
		return fmt.Sprintf("[%s]:%d", addr, port)
	}
	return fmt.Sprintf("%s:%d", addr, port)
}

func ExtractSubnet(addr string) int {
	parts := strings.Split(addr, ".")
	if len(parts) >= 3 {
		var n int
		if _, err := fmt.Sscanf(parts[2], "%d", &n); err == nil && n >= 0 && n <= 255 {
			return n
		}
	}
	return 0
}

type ScoringWeights struct {
	Latency   float64
	Loss      float64
	Jitter    float64
	Stability float64
}

var DefaultScoringWeights = ScoringWeights{
	Latency:   0.30,
	Loss:      0.30,
	Jitter:    0.15,
	Stability: 0.25,
}

func ComputeScore(m *Metrics, w ScoringWeights) float64 {
	if m.PacketLoss > 50 {
		return 0
	}
	if m.LatencyMs <= 0 {
		return 5
	}
	latencyScore := math.Max(0, 100-m.LatencyMs/2)
	lossScore := math.Max(0, 100-m.PacketLoss*2)
	jitterScore := math.Max(0, 100-m.JitterMs*3)
	stabilityScore := m.Stability * 100
	return latencyScore*w.Latency + lossScore*w.Loss + jitterScore*w.Jitter + stabilityScore*w.Stability
}

func UpdateStability(m *Metrics, goodThreshold float64, badLatency float64, upRate, downRate float64) {
	if m.PacketLoss < goodThreshold && m.LatencyMs < badLatency {
		m.Stability = math.Min(1, m.Stability+upRate)
	} else {
		m.Stability = math.Max(0, m.Stability-downRate)
	}
}
