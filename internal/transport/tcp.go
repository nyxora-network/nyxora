package transport

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type TCP struct {
	BaseTransport
	connMu      sync.Mutex
	localPort   int
	listener    net.Listener
	connections []net.Conn
}

func NewTCP() *TCP {
	return &TCP{
		BaseTransport: NewBase("tcp", "tcp", 9924, ScoringWeights{0.25, 0.35, 0.15, 0.25}, 0),
		localPort:     9925,
	}
}

func (t *TCP) Name() string  { return t.BaseName() }
func (t *TCP) Type() string { return t.BaseType() }

func (t *TCP) Init(cfg map[string]string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &t.port)
	}
	if lp, ok := cfg["local_port"]; ok {
		fmt.Sscanf(lp, "%d", &t.localPort)
	}
	log.Printf("[tcp] initialized (port: %d, local: %d)", t.port, t.localPort)
	return nil
}

func (t *TCP) Connect(remoteAddr string) error {
	ctx := t.CancelContext()
	t.KillOldProcess()
	if err := t.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	port := t.BasePort()
	remoteConn, err := net.DialTimeout("tcp", net.JoinHostPort(remoteAddr, fmt.Sprintf("%d", port)), 5*time.Second)
	if err != nil {
		t.Logf("remote %d unreachable: %v, ping-only mode", port, err)
		t.SetStatusActive()
		return nil
	}

	t.connMu.Lock()
	t.connections = append(t.connections, remoteConn)
	t.connMu.Unlock()

	t.SetStatusActive()
	t.Logf("connected to %s:%d", remoteAddr, port)

	go t.proxy(ctx, remoteConn)
	return nil
}

func (t *TCP) proxy(ctx context.Context, remoteConn net.Conn) {
	defer remoteConn.Close()
	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		remoteConn.SetReadDeadline(time.Now().Add(30 * time.Second))
		n, err := remoteConn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Logf("connection read error: %v", err)
			}
			t.SetStatusFailed()
			return
		}
		if n > 0 {
			remoteConn.Write(buf[:n])
		}
	}
}

func (t *TCP) Disconnect() error {
	t.mu.Lock()
	cancel := t.cancel
	t.status = StatusInactive
	t.connMu.Lock()
	conns := t.connections
	ln := t.listener
	t.connections = nil
	t.listener = nil
	t.connMu.Unlock()
	t.mu.Unlock()
	cancel()
	for _, conn := range conns {
		conn.Close()
	}
	if ln != nil {
		ln.Close()
	}
	t.Logf("disconnected")
	return nil
}

func (t *TCP) Status() Status    { return t.BaseStatus() }
func (t *TCP) Metrics() *Metrics { return t.BaseMetrics() }
func (t *TCP) Health() bool      { return t.BaseHealth() }
func (t *TCP) Score() float64    { return t.BaseScore() }
