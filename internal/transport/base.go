package transport

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

type BaseTransport struct {
	mu         sync.RWMutex
	name       string
	transport  string
	status     Status
	metrics    *Metrics
	remoteAddr string
	port       int
	weights    ScoringWeights
	bandwidth  int
	ctx        context.Context
	cancel     context.CancelFunc
	cmd        *exec.Cmd
	tmpFiles   []string
	scoringFn  func() float64
}

func NewBase(name, transport string, port int, weights ScoringWeights, bandwidth int) BaseTransport {
	ctx, cancel := context.WithCancel(context.Background())
	return BaseTransport{
		name:      name,
		transport: transport,
		status:    StatusInactive,
		metrics:   &Metrics{},
		port:      port,
		weights:   weights,
		bandwidth: bandwidth,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (b *BaseTransport) BaseName() string         { return b.name }
func (b *BaseTransport) BaseType() string         { return b.transport }
func (b *BaseTransport) BasePort() int            { return b.port }
func (b *BaseTransport) BaseRemoteAddr() string   { return b.remoteAddr }
func (b *BaseTransport) Context() context.Context { return b.ctx }

func (b *BaseTransport) SetScoringFn(fn func() float64) {
	b.scoringFn = fn
}

func (b *BaseTransport) CancelContext() context.Context {
	b.cancel()
	b.ctx, b.cancel = context.WithCancel(context.Background())
	return b.ctx
}

func (b *BaseTransport) BaseStatus() Status {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.status
}

func (b *BaseTransport) BaseMetrics() *Metrics {
	b.mu.RLock()
	defer b.mu.RUnlock()
	m := *b.metrics
	return &m
}

func (b *BaseTransport) BaseHealth() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.status == StatusActive
}

func (b *BaseTransport) BaseScore() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.baseUpdateMetrics()
	b.baseUpdateStability()
	if b.scoringFn != nil {
		return b.scoringFn()
	}
	return ComputeScore(b.metrics, b.weights)
}

func (b *BaseTransport) baseUpdateMetrics() {
	lat, loss, jitter := MeasureLatency(b.remoteAddr, 3)
	b.metrics.LatencyMs = lat
	b.metrics.PacketLoss = loss
	b.metrics.JitterMs = jitter
	b.metrics.Bandwidth = b.bandwidth
}

func (b *BaseTransport) baseUpdateStability() {
	UpdateStability(b.metrics, b.goodLossThreshold(), b.badLatencyThreshold(), b.upRate(), b.downRate())
}

func (b *BaseTransport) goodLossThreshold() float64  { return 10 }
func (b *BaseTransport) badLatencyThreshold() float64 { return 200 }
func (b *BaseTransport) upRate() float64              { return 0.05 }
func (b *BaseTransport) downRate() float64            { return 0.10 }

func (b *BaseTransport) BaseConnectInit(remoteAddr string) error {
	b.KillExisting()
	b.mu.Lock()
	defer b.mu.Unlock()
	b.remoteAddr = remoteAddr
	b.status = StatusTesting
	lat, loss, jitter := MeasureLatency(remoteAddr, 3)
	b.metrics.LatencyMs = lat
	b.metrics.PacketLoss = loss
	b.metrics.JitterMs = jitter
	b.metrics.Bandwidth = b.bandwidth
	if loss > 80 {
		b.status = StatusFailed
		return fmt.Errorf("high packet loss (%.1f%%)", loss)
	}
	return nil
}

func (b *BaseTransport) KillExisting() {
	b.cancel()
}

func (b *BaseTransport) KillOldProcess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cmd != nil && b.cmd.Process != nil {
		b.cmd.Process.Kill()
		b.cmd = nil
	}
}

func (b *BaseTransport) CleanTmpFiles() {
	for _, f := range b.tmpFiles {
		exec.Command("rm", "-f", f).Run()
	}
	b.tmpFiles = nil
}

func (b *BaseTransport) BaseDisconnect() error {
	b.cancel()
	b.KillOldProcess()
	b.CleanTmpFiles()
	b.mu.Lock()
	defer b.mu.Unlock()
	b.status = StatusInactive
	b.Logf("disconnected")
	return nil
}

func (b *BaseTransport) SetStatus(s Status) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.status = s
}

func (b *BaseTransport) SetStatusActive() {
	b.SetStatus(StatusActive)
}

func (b *BaseTransport) SetStatusFailed() {
	b.SetStatus(StatusFailed)
}

func (b *BaseTransport) SetCmd(cmd *exec.Cmd) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.cmd = cmd
}

func (b *BaseTransport) AddTmpFile(path string) {
	b.tmpFiles = append(b.tmpFiles, path)
}

func (b *BaseTransport) Logf(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[%s] ", b.name)+format, args...)
}

func (b *BaseTransport) RunInBackground(fn func()) {
	go func() {
		select {
		case <-b.ctx.Done():
			return
		default:
		}
		fn()
	}()
}

func (b *BaseTransport) KillOnCancel() {
	b.mu.RLock()
	cmd := b.cmd
	b.mu.RUnlock()
	go func() {
		<-b.ctx.Done()
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()
}
