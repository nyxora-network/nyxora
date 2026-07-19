package orchestrator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nyxora-network/nyxora/internal/failover"
	"github.com/nyxora-network/nyxora/internal/remote"
	"github.com/nyxora-network/nyxora/internal/transport"
)

func (o *Orchestrator) Start() error {
	o.mu.Lock()
	if o.running {
		o.mu.Unlock()
		return fmt.Errorf("orchestrator start: already running")
	}
	o.running = true
	o.mu.Unlock()

	log.Printf("[orchestrator] starting")

	o.tui.SetProvider(o)
	if err := o.tui.Start(); err != nil {
		log.Printf("[orchestrator] tui error: %v", err)
	}

	go o.fail.Start()

	log.Printf("[orchestrator] multipath scheduler: %s", o.scheduler.String())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Printf("[orchestrator] shutting down...")
	o.Stop()
	return nil
}

func (o *Orchestrator) Stop() {
	o.mu.Lock()
	defer o.mu.Unlock()
	if !o.running {
		return
	}
	o.running = false
	o.connected = false
	o.fail.Stop()
	o.mon.Stop()
	o.tui.Stop()
	o.transportM.DisconnectAll()
	if o.remoteHost != nil {
		remote.TeardownRemote(o.remoteHost, o.remoteIface)
		remote.TeardownProvisioned(o.remoteHost)
	}
	log.Printf("[orchestrator] stopped (uptime: %s)", time.Since(o.startTime).Round(time.Second))
}

func (o *Orchestrator) Status() map[string]interface{} {
	o.mu.Lock()
	running := o.running
	connected := o.connected
	phase := o.phase
	remoteHost := o.remoteHost
	o.mu.Unlock()

	status := map[string]interface{}{
		"running":          running,
		"connected":        connected,
		"phase":            string(phase),
		"node_id":          o.localNodeID,
		"active_transport": o.transportM.ActiveNames(),
		"all_active":       o.cfg.AllTunnelsActive,
		"mode":             "single-side",
		"uptime":           time.Since(o.startTime).Round(time.Second).String(),
	}

	if remoteHost != nil {
		status["remote"] = map[string]interface{}{
			"hostname": remoteHost.Hostname(),
			"address":  remoteHost.Address,
			"port":     remoteHost.Port,
			"os":       remoteHost.OSInfo(),
			"arch":     remoteHost.Arch(),
		}
	}

	var transports []map[string]interface{}
	for _, info := range o.transportM.List() {
		transports = append(transports, map[string]interface{}{
			"name":      info.Name,
			"type":      info.Type,
			"status":    info.Status,
			"score":     info.Score,
			"latency":   info.Latency,
			"jitter":    info.Jitter,
			"loss":      info.Loss,
			"stable":    info.Stability,
			"bandwidth": info.Bandwidth,
			"weight":    info.Weight,
		})
	}
	status["transports"] = transports

	best := o.routeEngine.BestPath()
	if best != nil {
		status["best_path"] = best.Name
		status["best_score"] = best.Score
	}

	failoverStatus := make(map[string]string)
	for name, s := range o.fail.AllStatus() {
		switch s {
		case failover.StatusHealthy:
			failoverStatus[name] = "healthy"
		case failover.StatusDegraded:
			failoverStatus[name] = "degraded"
		case failover.StatusDown:
			failoverStatus[name] = "down"
		}
	}
	status["failover"] = failoverStatus

	status["multipath"] = map[string]interface{}{
		"active":    o.scheduler.Stats().ActivePaths,
		"total":     len(transport.ListTunnels()),
		"best":      o.scheduler.Stats().BestPath,
		"failovers": o.scheduler.Stats().FailoverCount,
		"bandwidth": o.scheduler.AggregateBandwidth(),
		"mode":      "weighted",
		"paths":     o.scheduler.AllPaths(),
	}

	o.mu.Lock()
	steps := make([]StepStatus, len(o.steps))
	copy(steps, o.steps)
	o.mu.Unlock()
	status["steps"] = steps

	return status
}

func (o *Orchestrator) Steps() []StepStatus {
	o.mu.Lock()
	defer o.mu.Unlock()
	steps := make([]StepStatus, len(o.steps))
	copy(steps, o.steps)
	return steps
}

func (o *Orchestrator) startMonitoring(remoteAddr string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		o.mu.Lock()
		running := o.running
		connected := o.connected
		host := o.remoteHost
		o.mu.Unlock()

		if !running && !connected {
			return
		}

		lat, loss := host.Ping(2)

		for _, info := range o.transportM.List() {
			o.routeEngine.Update(info.Name, info.Type, lat, info.Jitter, loss, info.Stability, info.Bandwidth)
			o.fail.Update(info.Name, lat, loss)
			o.scheduler.UpdatePath(info.Name, info.Score, lat, loss, info.Bandwidth)
		}

		if o.cfg.AllTunnelsActive {
			weights := o.scheduler.Distribution()
			for name, w := range weights {
				o.transportM.SetWeight(name, w)
			}
		}

		best := o.routeEngine.BestPath()
		current := o.routeEngine.Current()
		if best != nil && best.Name != current && current != "" {
			diff := best.Score - o.getCurrentScore(current)
			if diff > 15 {
				log.Printf("[orchestrator] failover: %s (%.1f) -> %s (%.1f)",
					current, o.getCurrentScore(current), best.Name, best.Score)
				if cb := o.fail.GetOnFailover(); cb != nil {
					cb(current, best.Name)
				}
				o.scheduler.RecordFailover()
			}
		}

		<-ticker.C
	}
}

func (o *Orchestrator) getCurrentScore(name string) float64 {
	for _, p := range o.routeEngine.AllPaths() {
		if p.Name == name {
			return p.Score
		}
	}
	return 0
}

func execCmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	return string(out), err
}

func execCmdShell(cmd string) (string, error) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	return string(out), err
}

func trimNewline(s string) string {
	return strings.TrimSpace(s)
}
