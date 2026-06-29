package remote

import (
	"testing"
)

func TestNewHost(t *testing.T) {
	h := NewHost("1.2.3.4", 22, "root", "password")
	if h == nil {
		t.Fatal("NewHost returned nil")
	}
	if h.Address != "1.2.3.4" {
		t.Errorf("address should be 1.2.3.4, got %s", h.Address)
	}
	if h.Port != 22 {
		t.Errorf("port should be 22, got %d", h.Port)
	}
	if h.User != "root" {
		t.Errorf("user should be root, got %s", h.User)
	}
	if h.Password != "password" {
		t.Errorf("password should be password, got %s", h.Password)
	}
}

func TestHostAccessors(t *testing.T) {
	h := NewHost("1.2.3.4", 22, "root", "pass")
	if h.Hostname() != "" {
		t.Error("hostname should be empty initially")
	}
	if h.OSInfo() != "" {
		t.Error("OSInfo should be empty initially")
	}
	if h.Arch() != "" {
		t.Error("Arch should be empty initially")
	}
	if h.Latency() != 0 {
		t.Error("Latency should be 0 initially")
	}
	if h.Loss() != 0 {
		t.Error("Loss should be 0 initially")
	}
}

func TestSSHKeyExists(t *testing.T) {
	exists := SSHKeyExists()
	_ = exists // may or may not exist in test env
}
