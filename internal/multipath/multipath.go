package multipath

import (
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

type PathState struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Score     float64 `json:"score"`
	Latency   float64 `json:"latency"`
	Loss      float64 `json:"loss"`
	Weight    int     `json:"weight"`
	Active    bool    `json:"active"`
	Bandwidth int     `json:"bandwidth"`
}

type Scheduler struct {
	mu          sync.RWMutex
	paths       map[string]*PathState
	totalWeight int
	mode        DistributionMode
	stats       Stats
}

type DistributionMode int

const (
	ModeWeighted DistributionMode = iota
	ModeLowestLatency
	ModeLowestLoss
	ModeEven
	ModeAll
)

type Stats struct {
	TotalBytesSent     int64     `json:"total_bytes_sent"`
	TotalBytesReceived int64     `json:"total_bytes_received"`
	ActivePaths        int       `json:"active_paths"`
	BestPath           string    `json:"best_path"`
	FailoverCount      int       `json:"failover_count"`
	LastSwitch         time.Time `json:"last_switch"`
	Uptime             string    `json:"uptime"`
	startTime          time.Time
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		paths:  make(map[string]*PathState),
		mode:   ModeWeighted,
		stats:  Stats{startTime: time.Now()},
	}
}

func (s *Scheduler) SetMode(mode DistributionMode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mode = mode
	s.recalculate()
}

func (s *Scheduler) AddPath(name, transportType string, weight int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paths[name] = &PathState{
		Name:   name,
		Type:   transportType,
		Weight: weight,
		Active: false,
	}
	s.recalculate()
	log.Printf("[multipath] added path: %s (weight: %d%%)", name, weight)
}

func (s *Scheduler) RemovePath(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.paths, name)
	s.recalculate()
}

func (s *Scheduler) UpdatePath(name string, score, latency, loss float64, bandwidth int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.paths[name]
	if !ok {
		return
	}

	p.Score = score
	p.Latency = latency
	p.Loss = loss
	p.Bandwidth = bandwidth
	p.Active = score > 0 && loss < 50

	if s.mode == ModeWeighted {
		s.adjustWeights()
	}
}

func (s *Scheduler) SetActive(name string, active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, ok := s.paths[name]; ok {
		p.Active = active
		s.recalculate()
	}
}

func (s *Scheduler) SelectPath() *PathState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var active []*PathState
	for _, p := range s.paths {
		if p.Active {
			active = append(active, p)
		}
	}

	if len(active) == 0 {
		return nil
	}

	switch s.mode {
	case ModeLowestLatency:
		sort.Slice(active, func(i, j int) bool {
			return active[i].Latency < active[j].Latency
		})
		return active[0]

	case ModeLowestLoss:
		sort.Slice(active, func(i, j int) bool {
			return active[i].Loss < active[j].Loss
		})
		return active[0]

	case ModeEven:
		idx := int(time.Now().UnixMilli()) % len(active)
		return active[idx]

	case ModeAll:
		return active[0]

	default:
		r := time.Now().UnixNano() % int64(s.totalWeight)
		var cumulative int64
		for _, p := range active {
			cumulative += int64(p.Weight)
			if r < cumulative {
				return p
			}
		}
		return active[len(active)-1]
	}
}

func (s *Scheduler) SelectPaths(count int) []*PathState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var active []*PathState
	for _, p := range s.paths {
		if p.Active {
			active = append(active, p)
		}
	}

	if count > len(active) {
		count = len(active)
	}

	sort.Slice(active, func(i, j int) bool {
		return active[i].Score > active[j].Score
	})

	return active[:count]
}

func (s *Scheduler) Distribution() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]int)
	for _, p := range s.paths {
		if p.Active {
			result[p.Name] = p.Weight
		}
	}
	return result
}

func (s *Scheduler) BestPath() *PathState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var best *PathState
	for _, p := range s.paths {
		if !p.Active {
			continue
		}
		if best == nil || p.Score > best.Score {
			best = p
		}
	}
	return best
}

func (s *Scheduler) AllPaths() []PathState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []PathState
	for _, p := range s.paths {
		list = append(list, *p)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Score > list[j].Score
	})
	return list
}

func (s *Scheduler) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	active := 0
	bestName := ""
	var bestScore float64
	for _, p := range s.paths {
		if p.Active {
			active++
		}
		if p.Score > bestScore {
			bestScore = p.Score
			bestName = p.Name
		}
	}

	s.stats.ActivePaths = active
	s.stats.BestPath = bestName
	s.stats.Uptime = time.Since(s.stats.startTime).Round(time.Second).String()

	return s.stats
}

func (s *Scheduler) AggregateBandwidth() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := 0
	for _, p := range s.paths {
		if p.Active {
			total += p.Bandwidth
		}
	}
	return total
}

func (s *Scheduler) recalculate() {
	total := 0
	for _, p := range s.paths {
		if p.Active {
			total += p.Weight
		}
	}
	s.totalWeight = total
}

func (s *Scheduler) adjustWeights() {
	var activeCount int
	var totalScore float64
	for _, p := range s.paths {
		if p.Active && p.Score > 0 {
			activeCount++
			totalScore += p.Score
		}
	}

	if activeCount == 0 || totalScore == 0 {
		return
	}

	avgScore := totalScore / float64(activeCount)

	for _, p := range s.paths {
		if !p.Active {
			p.Weight = 0
			continue
		}

		ratio := p.Score / avgScore
		p.Weight = int(math.Round(ratio * 20))
		if p.Weight < 5 {
			p.Weight = 5
		}
		if p.Weight > 60 {
			p.Weight = 60
		}
	}

	s.recalculate()

	targetSum := 100
	if s.totalWeight > 0 {
		for _, p := range s.paths {
			p.Weight = (p.Weight * targetSum) / s.totalWeight
		}
	}

	var finalTotal int
	for _, p := range s.paths {
		if p.Active {
			finalTotal += p.Weight
		}
	}
	s.totalWeight = finalTotal
}

func (s *Scheduler) RecordFailover() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.FailoverCount++
	s.stats.LastSwitch = time.Now()
}

func (s *Scheduler) RecordBytes(sent, received int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.TotalBytesSent += sent
	s.stats.TotalBytesReceived += received
}

func ModeFromString(mode string) DistributionMode {
	switch mode {
	case "weighted":
		return ModeWeighted
	case "lowest-latency":
		return ModeLowestLatency
	case "lowest-loss":
		return ModeLowestLoss
	case "even":
		return ModeEven
	case "all":
		return ModeAll
	default:
		return ModeWeighted
	}
}

func (m DistributionMode) String() string {
	switch m {
	case ModeWeighted:
		return "weighted"
	case ModeLowestLatency:
		return "lowest-latency"
	case ModeLowestLoss:
		return "lowest-loss"
	case ModeEven:
		return "even"
	case ModeAll:
		return "all"
	default:
		return "unknown"
	}
}

func (s *Scheduler) String() string {
	s.mu.RLock()
	mode := s.mode
	s.mu.RUnlock()

	paths := s.AllPaths()
	active := 0
	for _, p := range paths {
		if p.Active {
			active++
		}
	}
	return fmt.Sprintf("multipath: %d/%d active, best=%s, mode=%s, bw=%d",
		active, len(paths),
		s.Stats().BestPath,
		mode.String(),
		s.AggregateBandwidth(),
	)
}
