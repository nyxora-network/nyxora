package loadbalancer

import (
	"crypto/md5"
	"encoding/binary"
	"math"
	"sync"
	"time"
)

// Algorithm defines the load balancing algorithm
type Algorithm int

const (
	AlgorithmRoundRobin Algorithm = iota
	AlgorithmWeightedRoundRobin
	AlgorithmLeastConnections
	AlgorithmLeastLatency
	AlgorithmRandom
	AlgorithmIPHash
)

// Backend represents a backend server
type Backend struct {
	Name           string
	Address        string
	Port           int
	Weight         int
	ActiveConns    int
	TotalConns     int
	Latency        float64
	FailureCount   int
	LastFailure    time.Time
	Healthy        bool
}

// Balancer implements load balancing
type Balancer struct {
	mu          sync.RWMutex
	backends    []*Backend
	algorithm   Algorithm
	current     int
	roundRobin  uint64
}

// NewBalancer creates a new load balancer
func NewBalancer(algorithm Algorithm) *Balancer {
	return &Balancer{
		backends:  make([]*Backend, 0),
		algorithm: algorithm,
	}
}

// AddBackend adds a backend
func (b *Balancer) AddBackend(name, address string, port, weight int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.backends = append(b.backends, &Backend{
		Name:    name,
		Address: address,
		Port:    port,
		Weight:  weight,
		Healthy: true,
	})
}

// RemoveBackend removes a backend
func (b *Balancer) RemoveBackend(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, backend := range b.backends {
		if backend.Name == name {
			b.backends = append(b.backends[:i], b.backends[i+1:]...)
			return
		}
	}
}

// Next selects the next backend
func (b *Balancer) Next() *Backend {
	b.mu.Lock()
	defer b.mu.Unlock()

	healthy := b.getHealthy()
	if len(healthy) == 0 {
		return nil
	}

	switch b.algorithm {
	case AlgorithmRoundRobin:
		return b.roundRobinSelect(healthy)
	case AlgorithmWeightedRoundRobin:
		return b.weightedRoundRobinSelect(healthy)
	case AlgorithmLeastConnections:
		return b.leastConnectionsSelect(healthy)
	case AlgorithmLeastLatency:
		return b.leastLatencySelect(healthy)
	case AlgorithmRandom:
		return b.randomSelect(healthy)
	case AlgorithmIPHash:
		// IPHash requires client IP; fall back to round-robin if not set
		return b.roundRobinSelect(healthy)
	default:
		return b.roundRobinSelect(healthy)
	}
}

// NextWithIP selects the next backend using IP-based hashing for session affinity
func (b *Balancer) NextWithIP(clientIP string) *Backend {
	b.mu.Lock()
	defer b.mu.Unlock()

	healthy := b.getHealthy()
	if len(healthy) == 0 {
		return nil
	}

	if b.algorithm == AlgorithmIPHash {
		return b.ipHashSelect(healthy, clientIP)
	}

	// For non-IPHash algorithms, ignore IP and use normal selection
	switch b.algorithm {
	case AlgorithmRoundRobin:
		return b.roundRobinSelect(healthy)
	case AlgorithmWeightedRoundRobin:
		return b.weightedRoundRobinSelect(healthy)
	case AlgorithmLeastConnections:
		return b.leastConnectionsSelect(healthy)
	case AlgorithmLeastLatency:
		return b.leastLatencySelect(healthy)
	case AlgorithmRandom:
		return b.randomSelect(healthy)
	default:
		return b.roundRobinSelect(healthy)
	}
}

// UpdateMetrics updates backend metrics
func (b *Balancer) UpdateMetrics(name string, latency float64, conns int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, backend := range b.backends {
		if backend.Name == name {
			backend.Latency = latency
			backend.ActiveConns = conns
			return
		}
	}
}

// MarkHealthy marks a backend as healthy
func (b *Balancer) MarkHealthy(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, backend := range b.backends {
		if backend.Name == name {
			backend.Healthy = true
			backend.FailureCount = 0
			return
		}
	}
}

// MarkUnhealthy marks a backend as unhealthy
func (b *Balancer) MarkUnhealthy(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, backend := range b.backends {
		if backend.Name == name {
			backend.Healthy = false
			backend.FailureCount++
			backend.LastFailure = time.Now()
			return
		}
	}
}

// Stats returns balancer statistics
func (b *Balancer) Stats() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()
	healthy := b.getHealthy()
	return map[string]interface{}{
		"algorithm":   b.algorithm,
		"total":       len(b.backends),
		"healthy":     len(healthy),
		"unhealthy":   len(b.backends) - len(healthy),
	}
}

func (b *Balancer) getHealthy() []*Backend {
	var healthy []*Backend
	for _, backend := range b.backends {
		if backend.Healthy {
			healthy = append(healthy, backend)
		}
	}
	return healthy
}

func (b *Balancer) roundRobinSelect(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	backend := backends[b.current%len(backends)]
	b.current++
	return backend
}

func (b *Balancer) weightedRoundRobinSelect(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	var totalWeight int
	for _, backend := range backends {
		totalWeight += backend.Weight
	}
	if totalWeight == 0 {
		return backends[0]
	}

	b.roundRobin = (b.roundRobin + 1) % uint64(totalWeight)
	var cumulative int
	for _, backend := range backends {
		cumulative += backend.Weight
		if int(b.roundRobin) < cumulative {
			return backend
		}
	}
	return backends[len(backends)-1]
}

func (b *Balancer) leastConnectionsSelect(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	var best *Backend
	for _, backend := range backends {
		if best == nil || backend.ActiveConns < best.ActiveConns {
			best = backend
		}
	}
	return best
}

func (b *Balancer) leastLatencySelect(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	var best *Backend
	for _, backend := range backends {
		if best == nil || backend.Latency < best.Latency {
			best = backend
		}
	}
	return best
}

func (b *Balancer) randomSelect(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	idx := int(math.Abs(float64(time.Now().UnixNano()))) % len(backends)
	return backends[idx]
}

// ipHashSelect selects a backend based on consistent hashing of client IP
// This ensures the same client IP always routes to the same backend (session affinity)
func (b *Balancer) ipHashSelect(backends []*Backend, clientIP string) *Backend {
	if len(backends) == 0 {
		return nil
	}

	// Compute MD5 hash of client IP for consistent distribution
	hash := md5.Sum([]byte(clientIP))
	// Use first 8 bytes as uint64 for uniform distribution
	hashVal := binary.BigEndian.Uint64(hash[:8])

	// Weighted IP hash: consider backend weights for better distribution
	var totalWeight int
	for _, backend := range backends {
		totalWeight += backend.Weight
	}

	if totalWeight == 0 {
		// All weights are zero, use simple modulo
		idx := hashVal % uint64(len(backends))
		return backends[idx]
	}

	// Weighted selection: map hash value to weight range
	target := hashVal % uint64(totalWeight)
	var cumulative uint64
	for _, backend := range backends {
		cumulative += uint64(backend.Weight)
		if target < cumulative {
			return backend
		}
	}

	return backends[len(backends)-1]
}
