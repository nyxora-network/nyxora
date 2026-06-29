package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig
	if cfg.MonitorInterval != 30 {
		t.Errorf("default MonitorInterval should be 30, got %d", cfg.MonitorInterval)
	}
	if cfg.FailoverInterval != 15 {
		t.Errorf("default FailoverInterval should be 15, got %d", cfg.FailoverInterval)
	}
	if cfg.DataDir != "/etc/nyxora" {
		t.Errorf("default DataDir should be /etc/nyxora, got %s", cfg.DataDir)
	}
	if cfg.AllTunnelsActive {
		t.Error("default AllTunnelsActive should be false")
	}
}

func TestLoadNonExistent(t *testing.T) {
	cfg, err := Load("/nonexistent/config.json")
	if err != nil {
		t.Errorf("Load of nonexistent file should return default config, got error: %v", err)
	}
	if cfg.MonitorInterval != DefaultConfig.MonitorInterval {
		t.Error("should return default config for nonexistent file")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg := DefaultConfig
	cfg.MonitorInterval = 60
	cfg.AllTunnelsActive = true

	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.MonitorInterval != 60 {
		t.Errorf("MonitorInterval should be 60, got %d", loaded.MonitorInterval)
	}
	if !loaded.AllTunnelsActive {
		t.Error("AllTunnelsActive should be true")
	}
}

func TestSaveEmptyPath(t *testing.T) {
	dir := t.TempDir()
	cfg := DefaultConfig
	cfg.DataDir = dir

	if err := cfg.Save(""); err != nil {
		t.Errorf("Save with empty path should use DataDir, got error: %v", err)
	}

	path := filepath.Join(dir, "config.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("config file should exist at DataDir/config.json")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte("not json"), 0600)

	_, err := Load(path)
	if err == nil {
		t.Error("Load of invalid JSON should return error")
	}
}

func TestLoadSecretsFromEnv(t *testing.T) {
	os.Setenv("NYXORA_SS_PASSWORD", "test-ss-pass")
	os.Setenv("NYXORA_SS_METHOD", "chacha20-ietf-poly1305")
	os.Setenv("NYXORA_RATHOLE_TOKEN", "test-rathole-token")
	os.Setenv("NYXORA_HYSTERIA_AUTH", "test-hy-auth")
	os.Setenv("NYXORA_BACKHAUL_TOKEN", "test-bh-token")
	os.Setenv("NYXORA_IPSEC_PSK", "test-ipsec-psk")
	defer func() {
		os.Unsetenv("NYXORA_SS_PASSWORD")
		os.Unsetenv("NYXORA_SS_METHOD")
		os.Unsetenv("NYXORA_RATHOLE_TOKEN")
		os.Unsetenv("NYXORA_HYSTERIA_AUTH")
		os.Unsetenv("NYXORA_BACKHAUL_TOKEN")
		os.Unsetenv("NYXORA_IPSEC_PSK")
	}()

	s := LoadSecrets()
	if s.SSPassword != "test-ss-pass" {
		t.Errorf("SSPassword should be test-ss-pass, got %s", s.SSPassword)
	}
	if s.SSMethod != "chacha20-ietf-poly1305" {
		t.Errorf("SSMethod should be chacha20-ietf-poly1305, got %s", s.SSMethod)
	}
	if s.RatholeToken != "test-rathole-token" {
		t.Errorf("RatholeToken should be test-rathole-token, got %s", s.RatholeToken)
	}
	if s.HysteriaAuth != "test-hy-auth" {
		t.Errorf("HysteriaAuth should be test-hy-auth, got %s", s.HysteriaAuth)
	}
	if s.BackhaulToken != "test-bh-token" {
		t.Errorf("BackhaulToken should be test-bh-token, got %s", s.BackhaulToken)
	}
	if s.IPsecPSK != "test-ipsec-psk" {
		t.Errorf("IPsecPSK should be test-ipsec-psk, got %s", s.IPsecPSK)
	}
}

func TestLoadSecretsFallback(t *testing.T) {
	os.Unsetenv("NYXORA_SS_PASSWORD")
	os.Unsetenv("NYXORA_SS_METHOD")
	os.Unsetenv("NYXORA_RATHOLE_TOKEN")
	os.Unsetenv("NYXORA_HYSTERIA_AUTH")
	os.Unsetenv("NYXORA_BACKHAUL_TOKEN")
	os.Unsetenv("NYXORA_IPSEC_PSK")

	s := LoadSecrets()
	if s.SSPassword == "" {
		t.Error("SSPassword should have fallback value")
	}
	if s.SSMethod != "aes-256-gcm" {
		t.Errorf("SSMethod default should be aes-256-gcm, got %s", s.SSMethod)
	}
	if s.RatholeToken == "" {
		t.Error("RatholeToken should have fallback value")
	}
	if s.HysteriaAuth == "" {
		t.Error("HysteriaAuth should have fallback value")
	}
	if s.BackhaulToken == "" {
		t.Error("BackhaulToken should have fallback value")
	}
	if s.IPsecPSK == "" {
		t.Error("IPsecPSK should have fallback value")
	}
}

func TestConfigWithSecrets(t *testing.T) {
	os.Setenv("NYXORA_SS_PASSWORD", "env-pass")
	defer os.Unsetenv("NYXORA_SS_PASSWORD")

	cfg, err := Load("/nonexistent/path.json")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Secrets.SSPassword != "env-pass" {
		t.Errorf("Secrets.SSPassword should be env-pass, got %s", cfg.Secrets.SSPassword)
	}
}

func TestGetTransportsForMode(t *testing.T) {
	full := GetTransportsForMode(ModeFull)
	if len(full) != 11 {
		t.Errorf("full mode should have 11 transports, got %d", len(full))
	}

	lite := GetTransportsForMode(ModeLite)
	if len(lite) != 5 {
		t.Errorf("lite mode should have 5 transports, got %d", len(lite))
	}

	minimal := GetTransportsForMode(ModeMinimal)
	if len(minimal) != 2 {
		t.Errorf("minimal mode should have 2 transports, got %d", len(minimal))
	}
}

func TestGetEffectiveTransports(t *testing.T) {
	cfg := DefaultConfig
	cfg.Mode = ModeLite
	transports := cfg.GetEffectiveTransports()
	if len(transports) != 5 {
		t.Errorf("effective transports for lite should be 5, got %d", len(transports))
	}

	cfg.EnabledTransports = []string{"ssh", "quic"}
	transports = cfg.GetEffectiveTransports()
	if len(transports) != 2 {
		t.Errorf("explicit transports should override mode, got %d", len(transports))
	}
}

func TestGetPort(t *testing.T) {
	cfg := DefaultConfig
	if cfg.GetPort("ssh", 22) != 22 {
		t.Error("should return default port")
	}

	cfg.PortOverrides = map[string]int{"ssh": 2222}
	if cfg.GetPort("ssh", 22) != 2222 {
		t.Error("should return overridden port")
	}
}

func TestServerInfo(t *testing.T) {
	info := ServerInfo()
	if info["cpu_count"].(int) == 0 {
		t.Error("cpu_count should be > 0")
	}
	if _, ok := info["suggested_mode"].(ServerMode); !ok {
		t.Error("suggested_mode should be ServerMode")
	}
}

// === Validation Tests ===

func TestIsValidMode(t *testing.T) {
	tests := []struct {
		mode ServerMode
		want bool
	}{
		{ModeFull, true},
		{ModeLite, true},
		{ModeMinimal, true},
		{"invalid", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsValidMode(tt.mode); got != tt.want {
			t.Errorf("IsValidMode(%q) = %v, want %v", tt.mode, got, tt.want)
		}
	}
}

func TestIsValidTransport(t *testing.T) {
	valid := []string{"wireguard", "openvpn", "ssh", "quic", "frp", "rathole", "ipsec", "shadowsocks", "hysteria", "backhaul", "tcp"}
	for _, name := range valid {
		if !IsValidTransport(name) {
			t.Errorf("IsValidTransport(%q) should be true", name)
		}
	}

	invalid := []string{"foo", "bar", "", "CloudFlare", "WireGuard"}
	for _, name := range invalid {
		if IsValidTransport(name) {
			t.Errorf("IsValidTransport(%q) should be false", name)
		}
	}
}

func TestValidateTransports(t *testing.T) {
	err := ValidateTransports([]string{"ssh", "wireguard", "quic"})
	if err != nil {
		t.Errorf("valid transports should pass, got: %v", err)
	}

	err = ValidateTransports([]string{"ssh", "invalid_one"})
	if err == nil {
		t.Error("invalid transport should fail validation")
	}

	err = ValidateTransports([]string{})
	if err != nil {
		t.Errorf("empty list should pass, got: %v", err)
	}
}

func TestValidatePort(t *testing.T) {
	validPorts := []int{1, 80, 443, 8080, 51820, 65535}
	for _, p := range validPorts {
		if err := ValidatePort(p); err != nil {
			t.Errorf("port %d should be valid, got: %v", p, err)
		}
	}

	invalidPorts := []int{-1, 0, 65536, 99999, 100000}
	for _, p := range invalidPorts {
		if err := ValidatePort(p); err == nil {
			t.Errorf("port %d should be invalid", p)
		}
	}
}

func TestValidatePortOverrides(t *testing.T) {
	valid := map[string]int{"ssh": 2222, "wireguard": 51820}
	if err := ValidatePortOverrides(valid); err != nil {
		t.Errorf("valid overrides should pass, got: %v", err)
	}

	invalidPort := map[string]int{"ssh": 99999}
	if err := ValidatePortOverrides(invalidPort); err == nil {
		t.Error("invalid port should fail")
	}

	invalidTransport := map[string]int{"invalid_transport": 8080}
	if err := ValidatePortOverrides(invalidTransport); err == nil {
		t.Error("invalid transport should fail")
	}

	conflict := map[string]int{"ssh": 8080, "wireguard": 8080}
	if err := ValidatePortOverrides(conflict); err == nil {
		t.Error("port conflict should fail")
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := DefaultConfig
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid, got: %v", err)
	}

	cfg.Mode = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("invalid mode should fail validation")
	}

	cfg.Mode = ModeFull
	cfg.EnabledTransports = []string{"invalid"}
	if err := cfg.Validate(); err == nil {
		t.Error("invalid transport should fail validation")
	}

	cfg.EnabledTransports = nil
	cfg.PortOverrides = map[string]int{"ssh": 99999}
	if err := cfg.Validate(); err == nil {
		t.Error("invalid port should fail validation")
	}

	cfg.PortOverrides = map[string]int{"ssh": 8080, "wireguard": 8080}
	if err := cfg.Validate(); err == nil {
		t.Error("port conflict should fail validation")
	}
}

func TestLoadWithInvalidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(`{"mode": "invalid"}`), 0600)

	_, err := Load(path)
	if err == nil {
		t.Error("Load with invalid mode should fail")
	}
}

func TestDetectModeWithThresholds(t *testing.T) {
	tests := []struct {
		name       string
		thresholds ModeThresholds
	}{
		{"default", DefaultThresholds},
		{"custom", ModeThresholds{MinimalMaxMB: 256, LiteMaxMB: 1024}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := DetectModeWithThresholds(tt.thresholds)
			if !IsValidMode(mode) {
				t.Errorf("DetectModeWithThresholds returned invalid mode: %s", mode)
			}
		})
	}
}

func TestConfigThresholds(t *testing.T) {
	cfg := DefaultConfig
	cfg.Thresholds = &ModeThresholds{MinimalMaxMB: 256, LiteMaxMB: 1024}
	mode := cfg.DetectMode()
	if !IsValidMode(mode) {
		t.Errorf("config DetectMode should return valid mode, got %s", mode)
	}
}

func TestConfigThresholdsNil(t *testing.T) {
	cfg := DefaultConfig
	cfg.Thresholds = nil
	mode := cfg.DetectMode()
	if !IsValidMode(mode) {
		t.Errorf("config DetectMode with nil thresholds should work, got %s", mode)
	}
}

func TestSaveAndLoadWithThresholds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg := DefaultConfig
	cfg.Thresholds = &ModeThresholds{MinimalMaxMB: 256, LiteMaxMB: 1024}
	cfg.PortOverrides = map[string]int{"ssh": 2222}
	cfg.EnabledTransports = []string{"ssh", "quic"}

	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Thresholds == nil {
		t.Fatal("thresholds should be loaded")
	}
	if loaded.Thresholds.MinimalMaxMB != 256 {
		t.Errorf("MinimalMaxMB should be 256, got %d", loaded.Thresholds.MinimalMaxMB)
	}
	if loaded.PortOverrides["ssh"] != 2222 {
		t.Errorf("port override ssh should be 2222, got %d", loaded.PortOverrides["ssh"])
	}
	if len(loaded.EnabledTransports) != 2 {
		t.Errorf("enabled transports should be 2, got %d", len(loaded.EnabledTransports))
	}
}
