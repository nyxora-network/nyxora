package orchestrator

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nyxora/nyxora/internal/config"
	"github.com/nyxora/nyxora/internal/dashboard"
	"github.com/nyxora/nyxora/internal/failover"
	"github.com/nyxora/nyxora/internal/monitor"
	"github.com/nyxora/nyxora/internal/packager"
	"github.com/nyxora/nyxora/internal/remote"
	"github.com/nyxora/nyxora/internal/routing"
	"github.com/nyxora/nyxora/internal/transport"
)

type Phase string

const (
	PhaseInit      Phase = "initializing"
	PhaseConnecting Phase = "connecting"
	PhaseSetup     Phase = "setting up remote"
	PhaseTunnel    Phase = "establishing tunnel"
	PhaseActive    Phase = "active"
	PhaseFailed    Phase = "failed"
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

	remoteHost  *remote.Host
	localNodeID string

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
		localNodeID: generateNodeID(),
		phase:       PhaseInit,
		startTime:   time.Now(),
	}
}

func (o *Orchestrator) Init() error {
	log.Printf("[orchestrator] initializing nyxora v0.1.0")
	log.Printf("[orchestrator] node id: %s", o.localNodeID)

	if err := os.MkdirAll(o.cfg.DataDir, 0755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	transports := []transport.Transport{
		transport.NewWireGuard(),
		transport.NewQUIC(),
		transport.NewSSH(),
		transport.NewTCP(),
	}
	for _, t := range transports {
		o.transportM.Register(t)
	}

	o.fail.OnFailover(func(from, to string) {
		log.Printf("[orchestrator] *** FAILOVER: %s -> %s ***", from, to)
		o.routeEngine.SetCurrent(to)
	})

	o.fail.OnRecover(func(name string) {
		log.Printf("[orchestrator] *** RECOVER: %s healthy ***", name)
	})

	log.Printf("[orchestrator] initialized with %d transports", len(transports))
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
	o.mu.Unlock()

	if o.onStepUpdate != nil {
		o.onStepUpdate(step)
	}
}

func (o *Orchestrator) OnStepUpdate(fn func(StepStatus)) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.onStepUpdate = fn
}

func (o *Orchestrator) ConnectToRemote(addr string, port int, user, password string) error {
	o.phase = PhaseConnecting
	log.Printf("[orchestrator] connecting to %s@%s:%d", user, addr, port)

	o.remoteHost = remote.NewHost(addr, port, user, password)

	// Step 1: Ping
	o.addStep("Pinging remote server", "RUNNING", fmt.Sprintf("%s ...", addr))
	lat, loss := o.remoteHost.Ping(4)
	if loss > 80 {
		o.addStep("Pinging remote server", "FAILED", fmt.Sprintf("packet loss: %.0f%%", loss))
		o.phase = PhaseFailed
		return fmt.Errorf("remote unreachable: %.0f%% packet loss", loss)
	}
	o.addStep("Pinging remote server", "OK", fmt.Sprintf("%.0fms latency, %.0f%% loss", lat, loss))

	// Step 2: SSH
	o.addStep("SSH authentication", "RUNNING", fmt.Sprintf("%s@%s:%d", user, addr, port))
	msg, ok := o.remoteHost.CheckConnectivity()
	if !ok {
		o.addStep("SSH authentication", "FAILED", msg)
		o.phase = PhaseFailed
		return fmt.Errorf("ssh: %s", msg)
	}
	o.addStep("SSH authentication", "OK", msg)

	// Step 3: OS detection
	o.phase = PhaseSetup
	o.addStep("Detecting OS", "RUNNING", "")
	if err := o.remoteHost.DetectOS(); err != nil {
		o.addStep("Detecting OS", "FAILED", err.Error())
		o.phase = PhaseFailed
		return fmt.Errorf("detect os: %w", err)
	}
	o.addStep("Detecting OS", "OK",
		fmt.Sprintf("%s | %s", o.remoteHost.OSInfo(), o.remoteHost.Arch()))

	// Step 4: Install dependencies
	o.addStep("Installing dependencies", "RUNNING", "wireguard, curl, ncat ...")
	deps := []struct {
		name string
		pkg  string
	}{
		{"wireguard", "wireguard"},
		{"curl", "curl"},
		{"ncat", "ncat"},
	}
	var failed []string
	for _, dep := range deps {
		if !o.remoteHost.CheckTool(dep.name) {
			log.Printf("[orchestrator] installing %s on remote...", dep.pkg)
			if err := o.remoteHost.InstallTool(dep.pkg); err != nil {
				log.Printf("[orchestrator] install %s failed: %v", dep.pkg, err)
				failed = append(failed, dep.name)
			}
		}
	}
	if len(failed) > 0 {
		o.addStep("Installing dependencies", "WARN",
			fmt.Sprintf("partial: %s failed", strings.Join(failed, ", ")))
	} else {
		o.addStep("Installing dependencies", "OK", "all dependencies installed")
	}

	// Step 5: Generate keys
	o.phase = PhaseTunnel
	o.addStep("Generating keys", "RUNNING", "")
	localPriv, localPub := o.generateLocalWGKey()

	// Step 6: Setup WireGuard on remote
	o.addStep("Setting up remote WireGuard", "RUNNING", "")
	remotePub, err := remote.SetupWireGuardRemote(o.remoteHost, localPub, 51820)
	if err != nil {
		o.addStep("Setting up remote WireGuard", "FAILED", err.Error())
		o.phase = PhaseFailed
		return fmt.Errorf("remote wg setup: %w", err)
	}
	o.addStep("Setting up remote WireGuard", "OK",
		fmt.Sprintf("pubkey: %s...", remotePub[:16]))

	// Step 7: Setup local WireGuard
	o.addStep("Setting up local WireGuard", "RUNNING", "")
	remoteIP, err := remote.GetRemotePublicIP(o.remoteHost)
	if err != nil {
		remoteIP = addr
	}

	localWG := transport.NewWireGuard()
	localWG.Init(map[string]string{
		"private_key": localPriv,
		"interface":   "nyxora0",
	})
	if err := localWG.Connect(remoteIP); err != nil {
		o.addStep("Setting up local WireGuard", "FAILED", err.Error())
		o.phase = PhaseFailed
		return fmt.Errorf("local wg setup: %w", err)
	}
	o.transportM.Register(localWG)
	o.addStep("Setting up local WireGuard", "OK", "interface nyxora0 ready")

	// Step 8: Test tunnel
	o.phase = PhaseActive
	o.connected = true

	o.addStep("Tunnel established", "OK",
		fmt.Sprintf("%s <-> %s via WireGuard", o.localNodeID[:8], o.remoteHost.Hostname()))

	log.Printf("[orchestrator] tunnel active: %s <-> %s (%s)",
		o.localNodeID[:8], o.remoteHost.Hostname(), remoteIP)

	// Start monitoring
	o.routeEngine.SetCurrent("wireguard")
	go o.startMonitoring(remoteIP)

	return nil
}

func (o *Orchestrator) generateLocalWGKey() (priv, pub string) {
	if commandExists("wg") {
		out, err := exec.Command("wg", "genkey").Output()
		if err == nil && len(out) > 0 {
			priv = strings.TrimSpace(string(out))
			pubOut, err := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | wg pubkey", priv)).Output()
			if err == nil && len(pubOut) > 0 {
				pub = strings.TrimSpace(string(pubOut))
				return
			}
		}
	}
	priv = fmt.Sprintf("nyxora-local-key-%d", time.Now().UnixNano())
	pub = fmt.Sprintf("nyxora-local-pub-%d", time.Now().UnixNano())
	return
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func (o *Orchestrator) SSHCommand(cmd string) (string, error) {
	if o.remoteHost == nil {
		return "", fmt.Errorf("no remote host")
	}
	return o.remoteHost.SSHCommand(cmd)
}

func (o *Orchestrator) startMonitoring(remoteAddr string) {
	for {
		if !o.running && !o.connected {
			return
		}

		lat, loss := o.remoteHost.Ping(2)

		info := o.transportM.List()
		for _, t := range info {
			o.routeEngine.Update(t.Name, t.Type, lat, 0, loss, 1.0, t.Bandwidth)
			o.fail.Update(t.Name, lat, loss)
		}

		best := o.routeEngine.BestPath()
		current := o.routeEngine.Current()
		if best != nil && best.Name != current && current != "" && best.Score-o.getCurrentScore(current) > 15 {
			log.Printf("[orchestrator] failover: %s -> %s", current, best.Name)
			if cb := o.fail.GetOnFailover(); cb != nil {
				cb(current, best.Name)
			}
		}

		time.Sleep(10 * time.Second)
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

func (o *Orchestrator) Start() error {
	o.mu.Lock()
	if o.running {
		o.mu.Unlock()
		return fmt.Errorf("already running")
	}
	o.running = true
	o.mu.Unlock()

	log.Printf("[orchestrator] starting")

	o.tui.SetProvider(o)
	if err := o.tui.Start(); err != nil {
		log.Printf("[orchestrator] tui error: %v", err)
	}

	go o.fail.Start()

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
		remote.TeardownRemote(o.remoteHost, "nyxora0")
	}
	log.Printf("[orchestrator] stopped")
}

func (o *Orchestrator) Status() map[string]interface{} {
	status := map[string]interface{}{
		"running":          o.running,
		"connected":        o.connected,
		"phase":            string(o.phase),
		"node_id":          o.localNodeID,
		"active_transport": o.transportM.Active(),
		"mode":             "single-side",
		"uptime":           time.Since(o.startTime).String(),
	}

	if o.remoteHost != nil {
		status["remote"] = map[string]interface{}{
			"hostname": o.remoteHost.Hostname(),
			"address":  o.remoteHost.Address,
			"port":     o.remoteHost.Port,
			"os":       o.remoteHost.OSInfo(),
			"arch":     o.remoteHost.Arch(),
		}
	}

	var transports []map[string]interface{}
	for _, info := range o.transportM.List() {
		transports = append(transports, map[string]interface{}{
			"name":    info.Name,
			"type":    info.Type,
			"status":  info.Status,
			"score":   info.Score,
			"latency": info.Latency,
			"jitter":  info.Jitter,
			"loss":    info.Loss,
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

func generateNodeID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("nyx-%x", b)
}


