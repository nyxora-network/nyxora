package metrics

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Counter struct {
	mu    sync.Mutex
	value float64
	tags  map[string]string
}

func NewCounter(tags map[string]string) *Counter {
	return &Counter{tags: tags}
}

func (c *Counter) Inc()           { c.Add(1) }
func (c *Counter) Add(v float64)  { c.mu.Lock(); c.value += v; c.mu.Unlock() }
func (c *Counter) Value() float64 { c.mu.Lock(); defer c.mu.Unlock(); return c.value }
func (c *Counter) Reset()         { c.mu.Lock(); c.value = 0; c.mu.Unlock() }

type Gauge struct {
	mu    sync.Mutex
	value float64
	tags  map[string]string
}

func NewGauge(tags map[string]string) *Gauge {
	return &Gauge{tags: tags}
}

func (g *Gauge) Set(v float64)  { g.mu.Lock(); g.value = v; g.mu.Unlock() }
func (g *Gauge) Value() float64 { g.mu.Lock(); defer g.mu.Unlock(); return g.value }
func (g *Gauge) Inc()           { g.Add(1) }
func (g *Gauge) Dec()           { g.Add(-1) }
func (g *Gauge) Add(v float64)  { g.mu.Lock(); g.value += v; g.mu.Unlock() }

type Histogram struct {
	mu      sync.Mutex
	buckets []float64
	counts  []int
	total   float64
	count   int
}

func NewHistogram(buckets []float64) *Histogram {
	if len(buckets) == 0 {
		buckets = []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000}
	}
	return &Histogram{
		buckets: buckets,
		counts:  make([]int, len(buckets)+1),
	}
}

func (h *Histogram) Observe(v float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.total += v
	h.count++
	for i, b := range h.buckets {
		if v <= b {
			h.counts[i]++
			return
		}
	}
	h.counts[len(h.buckets)]++
}

func (h *Histogram) Count() int     { h.mu.Lock(); defer h.mu.Unlock(); return h.count }
func (h *Histogram) Sum() float64   { h.mu.Lock(); defer h.mu.Unlock(); return h.total }
func (h *Histogram) Mean() float64  { h.mu.Lock(); defer h.mu.Unlock(); if h.count == 0 { return 0 }; return h.total / float64(h.count) }

type Collector struct {
	mu          sync.RWMutex
	counters    map[string]*Counter
	gauges      map[string]*Gauge
	histograms  map[string]*Histogram
	startTime   time.Time
}

func NewCollector() *Collector {
	return &Collector{
		counters:   make(map[string]*Counter),
		gauges:     make(map[string]*Gauge),
		histograms: make(map[string]*Histogram),
		startTime:  time.Now(),
	}
}

func (c *Collector) Counter(name string, tags map[string]string) *Counter {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.counters[name]; !ok {
		c.counters[name] = NewCounter(tags)
	}
	return c.counters[name]
}

func (c *Collector) Gauge(name string, tags map[string]string) *Gauge {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.gauges[name]; !ok {
		c.gauges[name] = NewGauge(tags)
	}
	return c.gauges[name]
}

func (c *Collector) Histogram(name string, buckets []float64) *Histogram {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.histograms[name]; !ok {
		c.histograms[name] = NewHistogram(buckets)
	}
	return c.histograms[name]
}

func (c *Collector) UptimeSeconds() float64 {
	return time.Since(c.startTime).Seconds()
}

func (c *Collector) Snapshot() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snap := map[string]interface{}{
		"uptime_seconds": c.UptimeSeconds(),
		"timestamp":      time.Now().Format(time.RFC3339),
	}

	counters := make(map[string]float64)
	for name, ctr := range c.counters {
		counters[name] = ctr.Value()
	}
	snap["counters"] = counters

	gauges := make(map[string]float64)
	for name, g := range c.gauges {
		gauges[name] = g.Value()
	}
	snap["gauges"] = gauges

	histograms := make(map[string]map[string]interface{})
	for name, h := range c.histograms {
		histograms[name] = map[string]interface{}{
			"count": h.Count(),
			"sum":   h.Sum(),
			"mean":  h.Mean(),
		}
	}
	snap["histograms"] = histograms

	return snap
}

func (c *Collector) Promethize() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var out string
	for name, ctr := range c.counters {
		out += fmt.Sprintf("# TYPE nyxora_%s counter\n", name)
		out += fmt.Sprintf("nyxora_%s %f\n", name, ctr.Value())
	}
	for name, g := range c.gauges {
		out += fmt.Sprintf("# TYPE nyxora_%s gauge\n", name)
		out += fmt.Sprintf("nyxora_%s %f\n", name, g.Value())
	}
	for name, h := range c.histograms {
		out += fmt.Sprintf("# TYPE nyxora_%s histogram\n", name)
		out += fmt.Sprintf("nyxora_%s_count %d\n", name, h.Count())
		out += fmt.Sprintf("nyxora_%s_sum %f\n", name, h.Sum())
	}
	out += "# TYPE nyxora_uptime_seconds gauge\n"
	out += fmt.Sprintf("nyxora_uptime_seconds %f\n", c.UptimeSeconds())
	return out
}

type Server struct {
	collector *Collector
	addr      string
}

func NewServer(addr string, collector *Collector) *Server {
	return &Server{
		collector: collector,
		addr:      addr,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", s.handleMetrics)
	mux.HandleFunc("/health", s.handleHealth)

	log.Printf("[metrics] server starting on %s", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprint(w, s.collector.Promethize())
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
