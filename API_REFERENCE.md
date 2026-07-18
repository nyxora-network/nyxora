# NYXORA API Reference

## Table of Contents

- [Package: config](#package-config)
- [Package: transport](#package-transport)
- [Package: orchestrator](#package-orchestrator)
- [Package: multipath](#package-multipath)
- [Package: failover](#package-failover)
- [Package: monitor](#package-monitor)
- [Package: routing](#package-routing)
- [Package: remote](#package-remote)
- [Package: dashboard](#package-dashboard)
- [Package: interactive](#package-interactive)
- [Package: metrics](#package-metrics)
- [Package: dns](#package-dns)
- [Package: ratelimit](#package-ratelimit)
- [Package: tls](#package-tls)
- [Package: packager](#package-packager)

---

## Package: `internal/config`

### Types

#### `Config`
Main configuration struct.

```go
type Config struct {
    AgentID           string            `json:"agent_id"`
    ListenAddr        string            `json:"listen_addr"`
    ServerMode        bool              `json:"server_mode"`
    RemoteAddr        string            `json:"remote_addr"`
    MonitorInterval   int               `json:"monitor_interval"`
    FailoverInterval  int               `json:"failover_interval"`
    StabilityWindow   int               `json:"stability_window"`
    AllTunnelsActive  bool              `json:"all_tunnels_active"`
    MaxBandwidth      int               `json:"max_bandwidth"`
    DataDir           string            `json:"data_dir"`
    LogLevel          string            `json:"log_level"`
    Mode              ServerMode        `json:"mode"`
    EnabledTransports []string          `json:"enabled_transports,omitempty"`
    PortOverrides     map[string]int    `json:"port_overrides,omitempty"`
    Thresholds        *ModeThresholds   `json:"thresholds,omitempty"`
    Secrets           Secrets           `json:"-"`
}
```

#### `ServerMode`
```go
type ServerMode string
const (
    ModeFull    ServerMode = "full"    // all 12 tunnels, 2GB+ RAM
    ModeLite    ServerMode = "lite"    // lightweight tunnels only, 512MB-2GB RAM
    ModeMinimal ServerMode = "minimal" // SSH + Shadowsocks only, <512MB RAM
)
```

#### `Secrets`
```go
type Secrets struct {
    SSPassword    string `json:"ss_password"`
    SSMethod      string `json:"ss_method"`
    RatholeToken  string `json:"rathole_token"`
    HysteriaAuth  string `json:"hysteria_auth"`
    BackhaulToken string `json:"backhaul_token"`
    IPsecPSK      string `json:"ipsec_psk"`
}
```

#### `ModeThresholds`
```go
type ModeThresholds struct {
    MinimalMaxMB uint64 // below this → minimal
    LiteMaxMB    uint64 // below this → lite, above → full
}
```

### Functions

| Function | Description |
|----------|-------------|
| `Load(path string) (*Config, error)` | Load config from JSON file |
| `(c *Config) Save(path string) error` | Save config to JSON file |
| `(c *Config) Validate() error` | Validate configuration |
| `(c *Config) GetEffectiveTransports() []string` | Get enabled transport list |
| `(c *Config) GetPort(transportName string, defaultPort int) int` | Get port with override |
| `ServerInfo() map[string]interface{}` | Get local server information |
| `GetTransportsForMode(mode ServerMode) []string` | Get transports for a mode |
| `DetectMode() ServerMode` | Auto-detect mode based on RAM |
| `LoadSecrets() Secrets` | Load secrets from env vars |
| `ValidateTransports(transports []string) error` | Validate transport names |
| `ValidatePortOverrides(overrides map[string]int) error` | Validate port conflicts |

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NYXORA_SS_PASSWORD` | Shadowsocks password | auto-generated |
| `NYXORA_SS_METHOD` | Shadowsocks cipher | `aes-256-gcm` |
| `NYXORA_RATHOLE_TOKEN` | Rathole auth token | auto-generated |
| `NYXORA_HYSTERIA_AUTH` | Hysteria auth password | auto-generated |
| `NYXORA_BACKHAUL_TOKEN` | Backhaul auth token | auto-generated |
| `NYXORA_IPSEC_PSK` | IPsec pre-shared key | auto-generated |
| `NYXORA_ALL_ACTIVE` | Enable all tunnels simultaneously | `false` |
| `NYXORA_CONFIG` | Custom config path | `/etc/nyxora/config.json` |

---

## Package: `internal/transport`

### Interface: `Transport`

```go
type Transport interface {
    Name() string
    Type() string
    Init(cfg map[string]string) error
    Connect(remoteAddr string) error
    Disconnect() error
    Status() Status
    Metrics() *Metrics
    Health() bool
    Score() float64
}
```

### Types

#### `Status`
```go
type Status string
const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
    StatusFailed   Status = "failed"
    StatusTesting  Status = "testing"
)
```

#### `Metrics`
```go
type Metrics struct {
    LatencyMs   float64 `json:"latency_ms"`
    JitterMs    float64 `json:"jitter_ms"`
    PacketLoss  float64 `json:"packet_loss"`
    Bandwidth   int     `json:"bandwidth"`
    Stability   float64 `json:"stability"`
}
```

#### `Info`
```go
type Info struct {
    Name      string  `json:"name"`
    Type      string  `json:"type"`
    Status    Status  `json:"status"`
    Score     float64 `json:"score"`
    Latency   float64 `json:"latency"`
    Jitter    float64 `json:"jitter"`
    Loss      float64 `json:"packet_loss"`
    Stability float64 `json:"stability"`
    Bandwidth int     `json:"bandwidth"`
    Weight    int     `json:"weight"`
}
```

#### `ScoringWeights`
```go
type ScoringWeights struct {
    Latency   float64
    Loss      float64
    Jitter    float64
    Stability float64
}
```

### Transport Categories

| Category | Constant |
|----------|----------|
| VPN | `CatVPN = "vpn"` |
| Tunnel | `CatTunnel = "tunnel"` |
| Relay | `CatRelay = "relay"` |
| Proxy | `CatProxy = "proxy"` |

### Available Transports

| Name | Port | Protocol | Category | Score | Weight |
|------|------|----------|----------|-------|--------|
| wireguard | 51820 | UDP | VPN | 95 | 30 |
| openvpn | 1194 | UDP | VPN | 75 | 10 |
| ssh | 22 | TCP | Tunnel | 60 | 5 |
| quic | 9923 | UDP | Tunnel | 80 | 15 |
| frp | 7000 | TCP | Relay | 70 | 10 |
| rathole | 2333 | TCP | Relay | 85 | 12 |
| ipsec | 500 | UDP | VPN | 70 | 5 |
| shadowsocks | 8388 | TCP | Proxy | 55 | 3 |
| hysteria | 8444 | UDP | Tunnel | 90 | 12 |
| backhaul | 3080 | TCP | Relay | 82 | 10 |
| tcp | 9924 | TCP | Tunnel | 50 | 3 |
| websocket | 9925 | TCP | Tunnel | 70 | 8 |

### Manager Methods

| Method | Description |
|--------|-------------|
| `NewManager(allActive bool) *Manager` | Create new manager |
| `(m *Manager) Register(t Transport)` | Register a transport |
| `(m *Manager) Get(name string) (Transport, bool)` | Get transport by name |
| `(m *Manager) List() []Info` | List all transports with metrics |
| `(m *Manager) ActiveList() []Info` | List active transports |
| `(m *Manager) Best() (Transport, float64)` | Get best transport by score |
| `(m *Manager) ConnectAll(remoteAddr string) error` | Connect all transports |
| `(m *Manager) DisconnectAll()` | Disconnect all transports |
| `(m *Manager) IsActive(name string) bool` | Check if transport is active |
| `(m *Manager) ActiveCount() int` | Get active transport count |
| `(m *Manager) SetWeight(name string, weight int)` | Set transport weight |
| `(m *Manager) NormalizeWeights()` | Normalize all weights to 100% |

### Utility Functions

| Function | Description |
|----------|-------------|
| `CommandExists(name string) bool` | Check if command exists in PATH |
| `MeasureLatency(addr string, count int) (latency, packetLoss, jitter float64)` | Measure latency via ping |
| `WriteConfig(path, content string) error` | Write config file (0644) |
| `WriteSecret(path, content string) error` | Write secret file (0600) |
| `FormatEndpoint(addr string, port int) string` | Format address:port |
| `ExtractSubnet(addr string) int` | Extract subnet from IP |
| `ComputeScore(m *Metrics, w ScoringWeights) float64` | Compute transport score |
| `UpdateStability(m *Metrics, ...)` | Update stability metric |

---

## Package: `internal/orchestrator`

### Types

#### `Orchestrator`
```go
type Orchestrator struct {
    cfg         *config.Config
    transportM  *transport.Manager
    mon         *monitor.Monitor
    routeEngine *routing.Engine
    fail        *failover.Failover
    pkg         *packager.Packager
    tui         *dashboard.TUI
    scheduler   *multipath.Scheduler
    remoteHost  *remote.Host
    localNodeID string
    phase       Phase
    steps       []StepStatus
    startTime   time.Time
    onStepUpdate func(StepStatus)
}
```

#### `Phase`
```go
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
```

#### `StepStatus`
```go
type StepStatus struct {
    Name   string `json:"name"`
    Status string `json:"status"`
    Detail string `json:"detail"`
    Done   bool   `json:"done"`
    TimeMs int64  `json:"time_ms"`
}
```

### Methods

| Method | Description |
|--------|-------------|
| `New(cfg *config.Config) *Orchestrator` | Create new orchestrator |
| `(o *Orchestrator) Init() error` | Initialize all components |
| `(o *Orchestrator) ConnectToRemote(addr string, port int, user, password string) error` | Full connection flow |
| `(o *Orchestrator) Start() error` | Start monitoring loop |
| `(o *Orchestrator) Stop()` | Stop all components |
| `(o *Orchestrator) Status() map[string]interface{}` | Get current status |
| `(o *Orchestrator) OnStepUpdate(fn func(StepStatus))` | Register step callback |

---

## Package: `internal/multipath`

### Types

#### `Scheduler`
```go
type Scheduler struct {
    mu          sync.RWMutex
    paths       map[string]*PathState
    totalWeight int
    mode        DistributionMode
    stats       Stats
}
```

#### `DistributionMode`
```go
type DistributionMode int
const (
    ModeWeighted      DistributionMode = iota // Weighted distribution
    ModeLowestLatency                         // Lowest latency path
    ModeLowestLoss                            // Lowest loss path
    ModeEven                                  // Equal distribution
    ModeAll                                   // All active simultaneously
)
```

#### `PathState`
```go
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
```

#### `Stats`
```go
type Stats struct {
    TotalBytesSent     int64     `json:"total_bytes_sent"`
    TotalBytesReceived int64     `json:"total_bytes_received"`
    ActivePaths        int       `json:"active_paths"`
    BestPath           string    `json:"best_path"`
    FailoverCount      int       `json:"failover_count"`
    LastSwitch         time.Time `json:"last_switch"`
    Uptime             string    `json:"uptime"`
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewScheduler() *Scheduler` | Create new scheduler |
| `(s *Scheduler) SetMode(mode DistributionMode)` | Set scheduling mode |
| `(s *Scheduler) AddPath(name, transportType string, weight int)` | Add a path |
| `(s *Scheduler) RemovePath(name string)` | Remove a path |
| `(s *Scheduler) UpdatePath(name string, score, latency, loss float64, bandwidth int)` | Update path metrics |
| `(s *Scheduler) SetActive(name string, active bool)` | Set path active state |
| `(s *Scheduler) SelectPath() *PathState` | Select best path |
| `(s *Scheduler) SelectPaths(count int) []*PathState` | Select top N paths |
| `(s *Scheduler) BestPath() *PathState` | Get best path |
| `(s *Scheduler) AllPaths() []PathState` | Get all paths |
| `(s *Scheduler) Stats() Stats` | Get scheduler stats |
| `(s *Scheduler) AggregateBandwidth() int` | Get total bandwidth |
| `(s *Scheduler) RecordFailover()` | Record a failover event |
| `(s *Scheduler) RecordBytes(sent, received int64)` | Record bytes transferred |
| `ModeFromString(mode string) DistributionMode` | Parse mode from string |

---

## Package: `internal/failover`

### Types

#### `Failover`
```go
type Failover struct {
    mu            sync.RWMutex
    states        map[string]*TransportState
    threshold     Threshold
    interval      time.Duration
    onFailoverCb  func(from, to string)
    onRecoverCb   func(name string)
    running       bool
    stopCh        chan struct{}
}
```

#### `TransportStatus`
```go
type TransportStatus int
const (
    StatusHealthy TransportStatus = iota
    StatusDegraded
    StatusDown
)
```

#### `Threshold`
```go
type Threshold struct {
    MaxLatency    float64
    MaxPacketLoss float64
    MaxJitter     float64
    MaxFailCount  int
    ScoreDiff     float64
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewFailover(intervalSec int) *Failover` | Create new failover engine |
| `(f *Failover) Start()` | Start evaluation loop |
| `(f *Failover) Stop()` | Stop evaluation loop |
| `(f *Failover) Update(name string, latency, packetLoss float64)` | Update transport metrics |
| `(f *Failover) OnFailover(fn func(from, to string))` | Register failover callback |
| `(f *Failover) OnRecover(fn func(name string))` | Register recovery callback |
| `(f *Failover) IsHealthy(name string) bool` | Check if transport is healthy |
| `(f *Failover) Status(name string) TransportStatus` | Get transport status |
| `(f *Failover) AllStatus() map[string]TransportStatus` | Get all statuses |

---

## Package: `internal/monitor`

### Types

#### `Monitor`
```go
type Monitor struct {
    mu          sync.RWMutex
    targets     map[string]*Target
    interval    time.Duration
    history     map[string][]Result
    maxHistory  int
    running     bool
    stopCh      chan struct{}
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewMonitor(intervalSec int) *Monitor` | Create new monitor |
| `(m *Monitor) AddTarget(name, addr string)` | Add monitoring target |
| `(m *Monitor) RemoveTarget(name string)` | Remove target |
| `(m *Monitor) Start()` | Start monitoring loop |
| `(m *Monitor) Stop()` | Stop monitoring |
| `(m *Monitor) LastResult(name string) *Result` | Get last result |
| `(m *Monitor) History(name string) []Result` | Get result history |
| `(m *Monitor) AverageLatency(name string) float64` | Get average latency |

---

## Package: `internal/routing`

### Types

#### `Engine`
```go
type Engine struct {
    mu      sync.RWMutex
    scorer  *Scorer
    paths   map[string]*PathScore
    current string
}
```

#### `PathScore`
```go
type PathScore struct {
    Name      string
    Type      string
    Score     float64
    Latency   float64
    Jitter    float64
    Loss      float64
    Stability float64
    Bandwidth int
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewEngine() *Engine` | Create new routing engine |
| `(e *Engine) Update(name, transportType string, latency, jitter, loss, stability float64, bandwidth int)` | Update path score |
| `(e *Engine) Current() string` | Get current path |
| `(e *Engine) SetCurrent(name string)` | Set current path |
| `(e *Engine) BestPath() *PathScore` | Get best path |
| `(e *Engine) AllPaths() []PathScore` | Get all paths |
| `(e *Engine) NeedsFailover(threshold float64) bool` | Check if failover needed |

---

## Package: `internal/remote`

### Types

#### `Host`
```go
type Host struct {
    Address  string `json:"address"`
    Port     int    `json:"port"`
    User     string `json:"user"`
    Password string `json:"password,omitempty"`
    KeyPath  string `json:"key_path,omitempty"`
    hostname string
    osInfo   string
    arch     string
    latency  float64
    loss     float64
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewHost(addr string, port int, user, password string) *Host` | Create new host |
| `(h *Host) SSHCommand(cmd string) (string, error)` | Execute SSH command |
| `(h *Host) SCP(localPath, remotePath string) error` | Copy file via SCP |
| `(h *Host) Ping(count int) (latency, loss float64)` | Ping the host |
| `(h *Host) DetectOS() error` | Detect remote OS |
| `(h *Host) CheckTool(tool string) bool` | Check if tool exists |
| `(h *Host) InstallTool(tool string) error` | Install tool on remote |
| `(h *Host) WriteFile(path, content string, mode string) error` | Write file on remote |
| `(h *Host) ReadFile(path string) (string, error)` | Read file from remote |
| `(h *Host) RunDaemon(command string) (string, error)` | Run daemon on remote |
| `(h *Host) CheckPort(port int, proto string) bool` | Check if port is open |
| `(h *Host) CheckConnectivity() (string, bool)` | Check SSH connectivity |
| `(h *Host) CheckResources() (*RemoteResources, error)` | Check remote resources |

---

## Package: `internal/dashboard`

### Types

#### `TUI`
```go
type TUI struct {
    mu       sync.Mutex
    width    int
    height   int
    interval time.Duration
    running  bool
    stopCh   chan struct{}
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewTUI(intervalSec int) *TUI` | Create new TUI |
| `(t *TUI) Start() error` | Start dashboard |
| `(t *TUI) Stop()` | Stop dashboard |

---

## Package: `internal/interactive`

### Functions

| Function | Description |
|----------|-------------|
| `RunMenu() (int, error)` | Launch interactive menu |
| `RunTransportStatus(transports []TransportStatus) error` | Show transport status |
| `RunUpdateChecker() error` | Check for updates |

### Theme Colors (TrueColor hex)

| Theme | Primary | Surface | Success | Warning | Error |
|-------|---------|---------|---------|---------|-------|
| Catppuccin Mocha | `#CBA6F7` | `#313244` | `#A6E3A1` | `#F9E2AF` | `#F38BA8` |
| Tokyo Night | `#7AA2F7` | `#24283B` | `#9ECE6A` | `#E0AF68` | `#F7768E` |
| Catppuccin Latte | `#8839EF` | `#E6E9EF` | `#40A02B` | `#DF8E1D` | `#D20F39` |

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Navigate up |
| `↓` / `j` | Navigate down |
| `Enter` | Select item |
| `q` / `Esc` | Go back / Quit |
| `1` | Catppuccin Mocha theme |
| `2` | Tokyo Night theme |
| `3` | Catppuccin Latte theme |
| `s` | Toggle status bar |
| `?` | Open help |
| `t` | Tunnel topology |

---

## Package: `internal/metrics`

### Types

#### `Collector`
```go
type Collector struct {
    mu       sync.RWMutex
    counters map[string]*Counter
    gauges   map[string]*Gauge
    histograms map[string]*Histogram
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewCollector() *Collector` | Create new collector |
| `(c *Collector) Counter(name string, labels map[string]string) *Counter` | Get/create counter |
| `(c *Collector) Gauge(name string, labels map[string]string) *Gauge` | Get/create gauge |
| `(c *Collector) Histogram(name string, labels map[string]string) *Histogram` | Get/create histogram |
| `(c *Collector) Promethize() string` | Export as Prometheus format |

---

## Package: `internal/dns`

### Types

#### `Resolver`
```go
type Resolver struct {
    mu        sync.RWMutex
    cache     map[string]*CacheEntry
    nameserver string
    timeout   time.Duration
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewResolver(nameserver string) *Resolver` | Create new resolver |
| `(r *Resolver) Lookup(domain string) ([]string, error)` | DNS lookup |
| `(r *Resolver) ClearCache()` | Clear DNS cache |
| `(r *Resolver) Stats() map[string]interface{}` | Get resolver stats |

---

## Package: `internal/ratelimit`

### Types

#### `Limiter`
```go
type Limiter struct {
    mu       sync.Mutex
    rate     float64
    burst    int
    tokens   float64
    lastTime time.Time
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewLimiter(rate float64, burst int) *Limiter` | Create new limiter |
| `(l *Limiter) Allow() bool` | Check if request is allowed |
| `(l *Limiter) Wait()` | Block until allowed |

---

## Package: `internal/tls`

### Types

#### `CertManager`
```go
type CertManager struct {
    mu      sync.RWMutex
    certDir string
    cert    *tls.Certificate
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewCertManager(certDir string) *CertManager` | Create new cert manager |
| `(cm *CertManager) GetCert() (*tls.Certificate, error)` | Get or generate certificate |
| `(cm *CertManager) TLSConfig() *tls.Config` | Get TLS config |

---

## Package: `internal/packager`

### Types

#### `Packager`
```go
type Packager struct {
    dataDir string
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewPackager(dataDir string) *Packager` | Create new packager |
| `(p *Packager) Pack(srcDir, destPath string) error` | Pack directory to tar.gz |
| `(p *Packager) Unpack(archivePath, destDir string) error` | Unpack tar.gz |
| `(p *Packager) List archives() ([]string, error)` | List available archives |
