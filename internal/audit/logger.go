package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Level represents the audit log level
type Level int

const (
	LevelInfo Level = iota
	LevelWarning
	LevelError
	LevelCritical
)

// Event represents an audit event
type Event struct {
	Timestamp  time.Time              `json:"timestamp"`
	Level      string                 `json:"level"`
	Category   string                 `json:"category"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	User       string                 `json:"user"`
	SourceIP   string                 `json:"source_ip"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Success    bool                   `json:"success"`
	Error      string                 `json:"error,omitempty"`
}

// Logger handles audit logging
type Logger struct {
	mu        sync.RWMutex
	file      *os.File
	events    []*Event
	maxEvents int
	level     Level
}

// NewLogger creates a new audit logger
func NewLogger(path string, level Level) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open audit log: %w", err)
	}

	return &Logger{
		file:      file,
		events:    make([]*Event, 0, 1000),
		maxEvents: 10000,
		level:     level,
	}, nil
}

// Log logs an audit event
func (l *Logger) Log(level Level, category, action, resource, user, sourceIP string, details map[string]interface{}, success bool, errMsg string) {
	if level < l.level {
		return
	}

	levelStr := "INFO"
	switch level {
	case LevelWarning:
		levelStr = "WARNING"
	case LevelError:
		levelStr = "ERROR"
	case LevelCritical:
		levelStr = "CRITICAL"
	}

	event := &Event{
		Timestamp: time.Now(),
		Level:     levelStr,
		Category:  category,
		Action:    action,
		Resource:  resource,
		User:      user,
		SourceIP:  sourceIP,
		Details:   details,
		Success:   success,
		Error:     errMsg,
	}

	l.mu.Lock()
	l.events = append(l.events, event)
	if len(l.events) > l.maxEvents {
		l.events = l.events[len(l.events)-l.maxEvents:]
	}

	data, _ := json.Marshal(event)
	l.file.Write(append(data, '\n'))
	l.mu.Unlock()
}

// LogInfo logs an info event
func (l *Logger) LogInfo(category, action, resource, user, sourceIP string, details map[string]interface{}) {
	l.Log(LevelInfo, category, action, resource, user, sourceIP, details, true, "")
}

// LogWarning logs a warning event
func (l *Logger) LogWarning(category, action, resource, user, sourceIP string, details map[string]interface{}) {
	l.Log(LevelWarning, category, action, resource, user, sourceIP, details, true, "")
}

// LogError logs an error event
func (l *Logger) LogError(category, action, resource, user, sourceIP string, details map[string]interface{}, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	l.Log(LevelError, category, action, resource, user, sourceIP, details, false, errMsg)
}

// LogCritical logs a critical event
func (l *Logger) LogCritical(category, action, resource, user, sourceIP string, details map[string]interface{}, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	l.Log(LevelCritical, category, action, resource, user, sourceIP, details, false, errMsg)
}

// GetEvents returns recent events
func (l *Logger) GetEvents(count int) []*Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if count > len(l.events) {
		count = len(l.events)
	}
	result := make([]*Event, count)
	copy(result, l.events[len(l.events)-count:])
	return result
}

// Close closes the audit log
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		l.file.Close()
	}
}
