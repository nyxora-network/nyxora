package transport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Manager struct {
	mu         sync.RWMutex
	transports map[string]Transport
	active     map[string]bool
	weights    map[string]int
	allActive  bool
}

func NewManager(allActive bool) *Manager {
	return &Manager{
		transports: make(map[string]Transport),
		active:     make(map[string]bool),
		weights:    make(map[string]int),
		allActive:  allActive,
	}
}

func (m *Manager) Register(t Transport) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transports[t.Name()] = t
	for _, meta := range TunnelRegistry {
		if meta.Name == t.Name() {
			m.weights[t.Name()] = meta.Weight
			break
		}
	}
	log.Printf("[transport] registered: %s (%s)", t.Name(), t.Type())
}

func (m *Manager) Get(name string) (Transport, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.transports[name]
	return t, ok
}

func (m *Manager) List() []Info {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []Info
	for _, t := range m.transports {
		metrics := t.Metrics()
		weight := m.weights[t.Name()]
		list = append(list, Info{
			Name:      t.Name(),
			Type:      t.Type(),
			Status:    t.Status(),
			Score:     t.Score(),
			Latency:   metrics.LatencyMs,
			Jitter:    metrics.JitterMs,
			Loss:      metrics.PacketLoss,
			Stability: metrics.Stability,
			Bandwidth: metrics.Bandwidth,
			Weight:    weight,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Score > list[j].Score
	})
	return list
}

func (m *Manager) ActiveList() []Info {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []Info
	for name, t := range m.transports {
		if !m.active[name] {
			continue
		}
		metrics := t.Metrics()
		list = append(list, Info{
			Name:      t.Name(),
			Type:      t.Type(),
			Status:    t.Status(),
			Score:     t.Score(),
			Latency:   metrics.LatencyMs,
			Jitter:    metrics.JitterMs,
			Loss:      metrics.PacketLoss,
			Stability: metrics.Stability,
			Bandwidth: metrics.Bandwidth,
			Weight:    m.weights[name],
		})
	}
	return list
}

func (m *Manager) Best() (Transport, float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.bestLocked()
}

func (m *Manager) ConnectAll(remoteAddr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("[manager] probing %s concurrently...", remoteAddr)
	latency, loss, _ := MeasureLatency(remoteAddr, 4)
	log.Printf("[manager] remote base latency: %.1fms, loss: %.1f%%", latency, loss)

	type result struct {
		name  string
		t     Transport
		score float64
		err   error
	}

	// Timeout to prevent infinite hangs
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	resultChan := make(chan result, len(m.transports))

	for name, t := range m.transports {
		t := t
		name := name

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			log.Printf("[manager] trying %s...", name)
			err := t.Connect(remoteAddr)
			score := 0.0
			if err == nil {
				score = t.Score()
			}
			resultChan <- result{name: name, t: t, score: score, err: err}
			return nil
		})
	}

	_ = g.Wait()
	close(resultChan)

	var candidates []result
	var allErrors []error

	for res := range resultChan {
		if res.err != nil {
			allErrors = append(allErrors, fmt.Errorf("%s: %w", res.name, res.err))
			log.Printf("[manager] %s failed: %v", res.name, res.err)
			res.t.Disconnect()
		} else {
			log.Printf("[manager] %s success (score: %.1f)", res.name, res.score)
			candidates = append(candidates, res)
		}
	}

	if len(candidates) == 0 {
		return fmt.Errorf("connect all transports to %s: %w", remoteAddr, errors.Join(allErrors...))
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	if m.allActive {
		for _, c := range candidates {
			m.active[c.name] = true
		}
		log.Printf("[manager] all-active: %d transports active", len(candidates))
		if len(candidates) > 0 {
			log.Printf("[manager] primary: %s (score: %.1f)", candidates[0].name, candidates[0].score)
		}
	} else {
		best := candidates[0]
		m.active[best.name] = true

		for _, c := range candidates[1:] {
			c.t.Disconnect()
		}
		log.Printf("[manager] selected primary: %s (score: %.1f)", best.name, best.score)
	}

	return nil
}

func (m *Manager) SetAllActive(active bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.allActive = active
}

func (m *Manager) IsAllActive() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.allActive
}

func (m *Manager) Activate(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.transports[name]; ok {
		m.active[name] = true
		return true
	}
	return false
}

func (m *Manager) Deactivate(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.active, name)
}

func (m *Manager) IsActive(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.active[name]
}

func (m *Manager) ActiveNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var names []string
	for name := range m.active {
		names = append(names, name)
	}
	return names
}

func (m *Manager) ActiveCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.active)
}

func (m *Manager) GetWeights() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	w := make(map[string]int)
	for k, v := range m.weights {
		w[k] = v
	}
	return w
}

func (m *Manager) SetWeight(name string, weight int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.weights[name] = weight
}

func (m *Manager) NormalizeWeights() {
	m.mu.Lock()
	defer m.mu.Unlock()
	var total int
	for _, w := range m.weights {
		total += w
	}
	if total == 0 {
		return
	}
	for name := range m.weights {
		m.weights[name] = (m.weights[name] * 100) / total
	}
}

func (m *Manager) DisconnectAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for name, t := range m.transports {
		if m.active[name] {
			t.Disconnect()
		}
	}
	m.active = make(map[string]bool)
}

func (m *Manager) ActiveCount_() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("%d/%d", len(m.active), len(m.transports))
}

func (m *Manager) bestLocked() (Transport, float64) {
	var best Transport
	var bestScore float64
	for _, t := range m.transports {
		s := t.Score()
		if s > bestScore {
			bestScore = s
			best = t
		}
	}
	return best, bestScore
}

func (m *Manager) RefreshScores() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, t := range m.transports {
		t.Score()
	}
}

// probeRemote is exported via scoring.go's MeasureLatency
