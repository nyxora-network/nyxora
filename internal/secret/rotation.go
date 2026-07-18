package secret

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// RotationConfig holds rotation settings
type RotationConfig struct {
	Interval    time.Duration
	MaxAge      time.Duration
	AutoRotate  bool
	SecretPath  string
}

// Rotator manages secret rotation
type Rotator struct {
	mu          sync.RWMutex
	config      RotationConfig
	secrets     map[string]*SecretInfo
	running     bool
	stopCh      chan struct{}
	onRotate    func(name, oldValue, newValue string)
}

// SecretInfo holds information about a secret
type SecretInfo struct {
	Name      string
	Value     string
	CreatedAt time.Time
	RotatedAt time.Time
	RotateCount int
}

// DefaultRotationConfig returns default rotation config
func DefaultRotationConfig() RotationConfig {
	return RotationConfig{
		Interval:   24 * time.Hour,
		MaxAge:     30 * 24 * time.Hour,
		AutoRotate: true,
		SecretPath: "/etc/nyxora/secrets",
	}
}

// NewRotator creates a new secret rotator
func NewRotator(config RotationConfig) *Rotator {
	return &Rotator{
		config:  config,
		secrets: make(map[string]*SecretInfo),
		stopCh:  make(chan struct{}),
	}
}

// Register registers a secret for rotation
func (r *Rotator) Register(name, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.secrets[name] = &SecretInfo{
		Name:      name,
		Value:     value,
		CreatedAt: time.Now(),
		RotatedAt: time.Now(),
	}
}

// OnRotate registers a callback for rotation events
func (r *Rotator) OnRotate(fn func(name, oldValue, newValue string)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.onRotate = fn
}

// Start starts the rotation loop
func (r *Rotator) Start() {
	if !r.config.AutoRotate {
		return
	}

	r.mu.Lock()
	r.running = true
	r.mu.Unlock()

	log.Printf("[secret-rotation] started (interval: %s)", r.config.Interval)

	go func() {
		for {
			select {
			case <-r.stopCh:
				return
			case <-time.After(r.config.Interval):
				r.rotate()
			}
		}
	}()
}

// Stop stops the rotation loop
func (r *Rotator) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.running {
		close(r.stopCh)
		r.running = false
		log.Printf("[secret-rotation] stopped")
	}
}

// RotateNow forces immediate rotation
func (r *Rotator) RotateNow() {
	r.rotate()
}

// Get returns the current value of a secret
func (r *Rotator) Get(name string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if secret, ok := r.secrets[name]; ok {
		return secret.Value, true
	}
	return "", false
}

// Stats returns rotation statistics
func (r *Rotator) Stats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return map[string]interface{}{
		"secrets":   len(r.secrets),
		"running":   r.running,
		"interval":  r.config.Interval.String(),
		"max_age":   r.config.MaxAge.String(),
	}
}

func (r *Rotator) rotate() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for name, secret := range r.secrets {
		if time.Since(secret.CreatedAt) > r.config.MaxAge {
			oldValue := secret.Value
			newValue := generateSecret(32)
			secret.Value = newValue
			secret.RotatedAt = time.Now()
			secret.RotateCount++

			log.Printf("[secret-rotation] rotated: %s (count: %d)", name, secret.RotateCount)

			if r.onRotate != nil {
				go r.onRotate(name, oldValue, newValue)
			}
		}
	}
}

func generateSecret(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// SaveSecrets saves secrets to file
func SaveSecrets(path string, secrets map[string]string) error {
	data := ""
	for name, value := range secrets {
		data += fmt.Sprintf("%s=%s\n", name, value)
	}
	return os.WriteFile(path, []byte(data), 0600)
}

// LoadSecrets loads secrets from file
func LoadSecrets(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	secrets := make(map[string]string)
	lines := splitLines(string(data))
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := splitKeyEqualsValue(line)
		if len(parts) == 2 {
			secrets[parts[0]] = parts[1]
		}
	}
	return secrets, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func splitKeyEqualsValue(s string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}
