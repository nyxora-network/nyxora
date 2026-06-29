package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// ServerMode defines the resource profile.
type ServerMode string

const (
	ModeFull    ServerMode = "full"    // all 11 tunnels, 2GB+ RAM
	ModeLite    ServerMode = "lite"    // lightweight tunnels only, 512MB-2GB RAM
	ModeMinimal ServerMode = "minimal" // SSH + Shadowsocks only, <512MB RAM
)

// ValidModes lists all valid ServerMode values.
var ValidModes = []ServerMode{ModeFull, ModeLite, ModeMinimal}

// ValidTransportNames is the set of all valid transport names.
var ValidTransportNames = map[string]bool{
	"wireguard": true, "openvpn": true, "ssh": true, "quic": true,
	"frp": true, "rathole": true, "ipsec": true, "shadowsocks": true,
	"hysteria": true, "backhaul": true, "tcp": true,
}

// ModeThresholds defines RAM thresholds for auto-detection (configurable).
type ModeThresholds struct {
	MinimalMaxMB uint64 // below this → minimal
	LiteMaxMB    uint64 // below this → lite, above → full
}

// DefaultThresholds are the default RAM thresholds.
var DefaultThresholds = ModeThresholds{
	MinimalMaxMB: 512,
	LiteMaxMB:    2048,
}

// Secrets holds all sensitive values. Loaded from env vars with fallback defaults.
type Secrets struct {
	SSPassword    string `json:"ss_password"`
	SSMethod      string `json:"ss_method"`
	RatholeToken  string `json:"rathole_token"`
	HysteriaAuth  string `json:"hysteria_auth"`
	BackhaulToken string `json:"backhaul_token"`
	IPsecPSK      string `json:"ipsec_psk"`
}

type Config struct {
	AgentID    string `json:"agent_id"`
	ListenAddr string `json:"listen_addr"`
	ServerMode bool   `json:"server_mode"`
	RemoteAddr string `json:"remote_addr"`

	MonitorInterval  int `json:"monitor_interval"`
	FailoverInterval int `json:"failover_interval"`
	StabilityWindow  int `json:"stability_window"`

	AllTunnelsActive bool `json:"all_tunnels_active"`
	MaxBandwidth     int  `json:"max_bandwidth"`
	DataDir          string `json:"data_dir"`
	LogLevel         string `json:"log_level"`

	Mode ServerMode `json:"mode"`

	EnabledTransports []string `json:"enabled_transports,omitempty"`

	PortOverrides map[string]int `json:"port_overrides,omitempty"`

	Thresholds *ModeThresholds `json:"thresholds,omitempty"`

	Secrets Secrets `json:"-"`
}

var DefaultConfig = Config{
	ListenAddr:       "0.0.0.0:9922",
	MonitorInterval:  30,
	FailoverInterval: 15,
	StabilityWindow:  5,
	AllTunnelsActive: false,
	MaxBandwidth:     1000,
	DataDir:          "/etc/nyxora",
	LogLevel:         "info",
	Mode:             ModeFull,
}

var LiteTransports = []string{
	"ssh", "shadowsocks", "quic", "tcp", "frp",
}

var MinimalTransports = []string{
	"ssh", "shadowsocks",
}

var AllTransports = []string{
	"wireguard", "openvpn", "ssh", "quic", "frp", "rathole",
	"ipsec", "shadowsocks", "hysteria", "backhaul", "tcp",
}

func GetTransportsForMode(mode ServerMode) []string {
	switch mode {
	case ModeLite:
		return LiteTransports
	case ModeMinimal:
		return MinimalTransports
	default:
		return AllTransports
	}
}

func GetCPULoad() float64 {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(data))
	if len(fields) >= 1 {
		load, err := strconv.ParseFloat(fields[0], 64)
		if err == nil {
			return load
		}
	}
	return 0
}

func GetAvailableRAMMB() uint64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				kb, err := strconv.ParseUint(fields[1], 10, 64)
				if err == nil {
					return kb / 1024
				}
			}
		}
	}
	return 0
}

func DetectMode() ServerMode {
	return DetectModeWithThresholds(DefaultThresholds)
}

func (c *Config) DetectMode() ServerMode {
	if c.Thresholds != nil {
		return DetectModeWithThresholds(*c.Thresholds)
	}
	return DetectModeWithThresholds(DefaultThresholds)
}

func DetectModeWithThresholds(t ModeThresholds) ServerMode {
	availMB := GetAvailableRAMMB()
	cpuLoad := GetCPULoad()
	cpuCount := float64(runtime.NumCPU())

	effectiveRAM := availMB
	if cpuLoad > cpuCount*0.8 {
		effectiveRAM = availMB * 70 / 100
	}

	if effectiveRAM < t.MinimalMaxMB {
		return ModeMinimal
	} else if effectiveRAM < t.LiteMaxMB {
		return ModeLite
	}
	return ModeFull
}

func IsValidMode(mode ServerMode) bool {
	for _, m := range ValidModes {
		if m == mode {
			return true
		}
	}
	return false
}

func IsValidTransport(name string) bool {
	return ValidTransportNames[name]
}

func ValidateTransports(transports []string) error {
	for _, t := range transports {
		if !IsValidTransport(t) {
			return fmt.Errorf("invalid transport: %q (valid: %s)", t, strings.Join(AllTransports, ", "))
		}
	}
	return nil
}

func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", port)
	}
	return nil
}

func ValidatePortOverrides(overrides map[string]int) error {
	usedPorts := make(map[int]string)
	for name, port := range overrides {
		if !IsValidTransport(name) {
			return fmt.Errorf("invalid transport %q in port overrides", name)
		}
		if err := ValidatePort(port); err != nil {
			return fmt.Errorf("port override for %s: %w", name, err)
		}
		if existing, ok := usedPorts[port]; ok {
			return fmt.Errorf("port conflict: %d used by both %q and %q", port, existing, name)
		}
		usedPorts[port] = name
	}
	return nil
}

func (c *Config) GetEffectiveTransports() []string {
	if len(c.EnabledTransports) > 0 {
		return c.EnabledTransports
	}
	return GetTransportsForMode(c.Mode)
}

func (c *Config) GetPort(transportName string, defaultPort int) int {
	if c.PortOverrides != nil {
		if p, ok := c.PortOverrides[transportName]; ok {
			return p
		}
	}
	return defaultPort
}

func (c *Config) Validate() error {
	if c.Mode != "" && !IsValidMode(c.Mode) {
		return fmt.Errorf("invalid mode: %q (valid: full, lite, minimal)", c.Mode)
	}
	if len(c.EnabledTransports) > 0 {
		if err := ValidateTransports(c.EnabledTransports); err != nil {
			return err
		}
	}
	if c.PortOverrides != nil {
		if err := ValidatePortOverrides(c.PortOverrides); err != nil {
			return err
		}
	}
	return nil
}

func generateRandomToken(prefix string) string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return prefix + "-fallback"
	}
	return prefix + "-" + hex.EncodeToString(b)
}

func LoadSecrets() Secrets {
	s := Secrets{
		SSPassword:    os.Getenv("NYXORA_SS_PASSWORD"),
		SSMethod:      os.Getenv("NYXORA_SS_METHOD"),
		RatholeToken:  os.Getenv("NYXORA_RATHOLE_TOKEN"),
		HysteriaAuth:  os.Getenv("NYXORA_HYSTERIA_AUTH"),
		BackhaulToken: os.Getenv("NYXORA_BACKHAUL_TOKEN"),
		IPsecPSK:      os.Getenv("NYXORA_IPSEC_PSK"),
	}
	if s.SSPassword == "" {
		s.SSPassword = generateRandomToken("nyxora-ss")
	}
	if s.SSMethod == "" {
		s.SSMethod = "aes-256-gcm"
	}
	if s.RatholeToken == "" {
		s.RatholeToken = generateRandomToken("nyxora-rathole")
	}
	if s.HysteriaAuth == "" {
		s.HysteriaAuth = generateRandomToken("nyxora-hy2")
	}
	if s.BackhaulToken == "" {
		s.BackhaulToken = generateRandomToken("nyxora-bh")
	}
	if s.IPsecPSK == "" {
		s.IPsecPSK = generateRandomToken("nyxora-ipsec")
	}
	return s
}

func Load(path string) (*Config, error) {
	cfg := DefaultConfig
	if path == "" {
		path = "/etc/nyxora/config.json"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg.Secrets = LoadSecrets()
			return &cfg, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	cfg.Secrets = LoadSecrets()
	return &cfg, nil
}

func (c *Config) Save(path string) error {
	if path == "" {
		path = filepath.Join(c.DataDir, "config.json")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func ServerInfo() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info := map[string]interface{}{
		"cpu_count":     runtime.NumCPU(),
		"cpu_load":      GetCPULoad(),
		"heap_mb":       int(m.HeapAlloc / 1024 / 1024),
		"sys_mb":        int(m.Sys / 1024 / 1024),
		"goroutines":    runtime.NumGoroutine(),
		"suggested_mode": DetectMode(),
	}

	data, err := os.ReadFile("/proc/meminfo")
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					kb, _ := strconv.ParseUint(fields[1], 10, 64)
					info["total_ram_mb"] = kb / 1024
				}
			}
			if strings.HasPrefix(line, "MemAvailable:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					kb, _ := strconv.ParseUint(fields[1], 10, 64)
					info["available_ram_mb"] = kb / 1024
				}
			}
		}
	}

	if cpuLoad, ok := info["cpu_load"].(float64); ok {
		cpuCount := float64(runtime.NumCPU())
		if cpuLoad > cpuCount*0.8 {
			info["cpu_warning"] = fmt.Sprintf("high load (%.1f > %.0f cores)", cpuLoad, cpuCount)
			log.Printf("[config] high CPU load detected: %.1f (may reduce effective RAM)", cpuLoad)
		}
	}
	return info
}
