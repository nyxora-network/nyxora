package failover

import (
	"log"
	"sync"
	"time"
)

type TransportStatus int

const (
	StatusHealthy TransportStatus = iota
	StatusDegraded
	StatusDown
)

type TransportState struct {
	Name      string
	Status    TransportStatus
	Latency   float64
	PacketLoss float64
	LastGood  time.Time
	FailCount int
}

type Failover struct {
	mu          sync.RWMutex
	states      map[string]*TransportState
	threshold   Threshold
	interval    time.Duration
	onFailoverCb  func(from, to string)
	onRecoverCb   func(name string)
	running     bool
	stopCh      chan struct{}
}

func (f *Failover) GetOnFailover() func(from, to string) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.onFailoverCb
}

type Threshold struct {
	MaxLatency    float64
	MaxPacketLoss float64
	MaxJitter     float64
	MaxFailCount  int
	ScoreDiff     float64
}

var DefaultThreshold = Threshold{
	MaxLatency:    200,
	MaxPacketLoss: 10,
	MaxJitter:     50,
	MaxFailCount:  3,
	ScoreDiff:     20,
}

func NewFailover(intervalSec int) *Failover {
	return &Failover{
		states:   make(map[string]*TransportState),
		threshold: DefaultThreshold,
		interval: time.Duration(intervalSec) * time.Second,
		stopCh:   make(chan struct{}),
	}
}

func (f *Failover) OnFailover(fn func(from, to string)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.onFailoverCb = fn
}

func (f *Failover) OnRecover(fn func(name string)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.onRecoverCb = fn
}

func (f *Failover) Start() {
	f.mu.Lock()
	f.running = true
	f.mu.Unlock()

	log.Printf("[failover] started (interval: %s)", f.interval)

	go func() {
		for {
			select {
			case <-f.stopCh:
				return
			case <-time.After(f.interval):
				f.evaluate()
			}
		}
	}()
}

func (f *Failover) Stop() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.running {
		close(f.stopCh)
		f.running = false
		log.Printf("[failover] stopped")
	}
}

func (f *Failover) Update(name string, latency, packetLoss float64) {
	f.mu.Lock()
	defer f.mu.Unlock()

	state, ok := f.states[name]
	if !ok {
		state = &TransportState{Name: name}
		f.states[name] = state
	}

		state.Latency = latency
	state.PacketLoss = packetLoss

	if latency < f.threshold.MaxLatency && packetLoss < f.threshold.MaxPacketLoss {
		if state.Status != StatusHealthy {
			state.Status = StatusHealthy
			state.FailCount = 0
			state.LastGood = time.Now()
			log.Printf("[failover] %s recovered", name)
			if f.onRecoverCb != nil {
				go f.onRecoverCb(name)
			}
		}
		state.LastGood = time.Now()
	} else if latency >= f.threshold.MaxLatency || packetLoss >= f.threshold.MaxPacketLoss {
		state.FailCount++
		if state.FailCount >= f.threshold.MaxFailCount {
			state.Status = StatusDown
		} else {
			state.Status = StatusDegraded
		}
	}
}

func (f *Failover) evaluate() {
	f.mu.RLock()
	states := make(map[string]TransportState)
	for k, v := range f.states {
		states[k] = *v
	}
	f.mu.RUnlock()

	var healthy, degraded, down []string
	for name, state := range states {
		switch state.Status {
		case StatusHealthy:
			healthy = append(healthy, name)
		case StatusDegraded:
			degraded = append(degraded, name)
		case StatusDown:
			down = append(down, name)
		}
	}

	if len(down) > 0 && len(healthy) > 0 {
		log.Printf("[failover] transports down: %v, switching to healthy", down)
		if f.onFailoverCb != nil {
			for _, d := range down {
				go f.onFailoverCb(d, healthy[0])
			}
		}
	}

	if len(degraded) > 0 && len(healthy) > 0 {
		log.Printf("[failover] transports degraded: %v", degraded)
	} else if len(down) > 0 && len(healthy) == 0 {
		log.Printf("[failover] all transports down!")
	}
}

func (f *Failover) IsHealthy(name string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	state, ok := f.states[name]
	if !ok {
		return false
	}
	return state.Status == StatusHealthy
}

func (f *Failover) Status(name string) TransportStatus {
	f.mu.RLock()
	defer f.mu.RUnlock()
	state, ok := f.states[name]
	if !ok {
		return StatusDown
	}
	return state.Status
}

func (f *Failover) AllStatus() map[string]TransportStatus {
	f.mu.RLock()
	defer f.mu.RUnlock()
	result := make(map[string]TransportStatus)
	for name, state := range f.states {
		result[name] = state.Status
	}
	return result
}
