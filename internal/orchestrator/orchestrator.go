package orchestrator

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nyxora-network/nyxora/internal/config"
	"github.com/nyxora-network/nyxora/internal/dashboard"
	"github.com/nyxora-network/nyxora/internal/failover"
	"github.com/nyxora-network/nyxora/internal/monitor"
	"github.com/nyxora-network/nyxora/internal/multipath"
	"github.com/nyxora-network/nyxora/internal/packager"
	"github.com/nyxora-network/nyxora/internal/remote"
	"github.com/nyxora-network/nyxora/internal/routing"
	"github.com/nyxora-network/nyxora/internal/transport"
)

type Phase string

const (
	PhaseInit       Phase = "initializing"
	PhaseConnecting Phase = "connecting"
	PhaseSetup      Phase = "setting up remote"
	PhaseTunnel     Phase = "establishing tunnel"
	PhaseMultipath  Phase = "multipath active"
	PhaseActive     Phase = "active"
	PhaseFailed     Phase = "failed"
)

type StepStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
	Done   bool   `json:"done"`
	TimeMs int64  `json:"time_ms"`
}

type Orchestrator struct {
	cfg         *config.Config
	transportM  *transport.Manager
	mon         *monitor.Monitor
	routeEngine *routing.Engine
	fail        *failover.Failover
	pkg         *packager.Packager
	tui         *dashboard.TUI
	scheduler   *multipath.Scheduler

	remoteHost   *remote.Host
	remoteIface  string
	localNodeID  string

	mu        sync.Mutex
	running   bool
	connected bool
	phase     Phase
	steps     []StepStatus
	startTime time.Time

	onStepUpdate func(StepStatus)
}

func New(cfg *config.Config) *Orchestrator {
	return &Orchestrator{
		cfg:         cfg,
		transportM:  transport.NewManager(cfg.AllTunnelsActive),
		mon:         monitor.NewMonitor(cfg.MonitorInterval),
		routeEngine: routing.NewEngine(),
		fail:        failover.NewFailover(cfg.FailoverInterval),
		pkg:         packager.NewPackager(cfg.DataDir),
		tui:         dashboard.NewTUI(2),
		scheduler:   multipath.NewScheduler(),
		localNodeID: generateNodeID(),
		phase:       PhaseInit,
		startTime:   time.Now(),
	}
}

var version = "dev"

func (o *Orchestrator) Init() error {
	log.Printf("[orchestrator] initializing nyxora v%s", version)
	log.Printf("[orchestrator] node id: %s", o.localNodeID)

	if err := os.MkdirAll(o.cfg.DataDir, 0755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	allTransports := []transport.Transport{
		transport.NewWireGuard(),
		transport.NewOpenVPN(),
		transport.NewSSH(),
		transport.NewQUIC(),
		transport.NewFRP(),
		transport.NewRathole(),
		transport.NewIPsec(),
		transport.NewShadowSOCKS(),
		transport.NewHysteria(),
		transport.NewBackhaul(),
		transport.NewTCP(),
		transport.NewWebSocket(),
	}

	// Filter transports based on mode
	effectiveTransports := o.cfg.GetEffectiveTransports()
	enabledSet := make(map[string]bool)
	for _, name := range effectiveTransports {
		enabledSet[name] = true
	}

	for _, t := range allTransports {
		if enabledSet[t.Name()] {
			o.transportM.Register(t)
		}
	}

	for _, meta := range transport.ListTunnels() {
		o.scheduler.AddPath(meta.Name, meta.Type, meta.Weight)
	}

	o.fail.OnFailover(func(from, to string) {
		log.Printf("[orchestrator] *** FAILOVER: %s -> %s ***", from, to)
		o.routeEngine.SetCurrent(to)
		o.scheduler.RecordFailover()
	})

	o.fail.OnRecover(func(name string) {
		log.Printf("[orchestrator] *** RECOVER: %s healthy ***", name)
	})

	log.Printf("[orchestrator] initialized with %d transports, %d paths",
		len(allTransports), len(transport.ListTunnels()))
	return nil
}

func (o *Orchestrator) addStep(name, status, detail string) {
	step := StepStatus{
		Name:   name,
		Status: status,
		Detail: detail,
		Done:   status == "OK",
		TimeMs: time.Since(o.startTime).Milliseconds(),
	}
	o.mu.Lock()
	o.steps = append(o.steps, step)
	cb := o.onStepUpdate
	o.mu.Unlock()
	if cb != nil {
		cb(step)
	}
}

func (o *Orchestrator) OnStepUpdate(fn func(StepStatus)) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.onStepUpdate = fn
}

func generateNodeID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("nyx-%x", b)
}

func generateLocalWGKey() (priv, pub string) {
	if transport.CommandExists("wg") {
		out, err := execCmd("wg", "genkey")
		if err == nil && len(out) > 0 {
			priv = trimNewline(out)
			pubOut, err := execCmdShell(fmt.Sprintf("echo '%s' | wg pubkey", priv))
			if err == nil && len(pubOut) > 0 {
				pub = trimNewline(pubOut)
				return
			}
		}
	}
	return "", ""
}

func getLocalIP() string {
	out, err := execCmdShell("ip route get 1 | awk '{print $7;exit}'")
	if err == nil && len(out) > 0 {
		return trimNewline(out)
	}
	return ""
}
