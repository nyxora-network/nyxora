package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

type QUIC struct {
	BaseTransport
	connMu      sync.Mutex
	connections []*quic.Conn
	tlsVerify   bool
}

func NewQUIC() *QUIC {
	return &QUIC{
		BaseTransport: NewBase("quic", "quic", 9923, ScoringWeights{0.35, 0.30, 0.15, 0.20}, 0),
	}
}

func (q *QUIC) Name() string  { return q.BaseName() }
func (q *QUIC) Type() string { return q.BaseType() }

func (q *QUIC) Init(cfg map[string]string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if p, ok := cfg["port"]; ok {
		fmt.Sscanf(p, "%d", &q.port)
	}
	q.connMu.Lock()
	if v, ok := cfg["tls_verify"]; ok {
		q.tlsVerify = v == "true" || v == "1"
	}
	q.connMu.Unlock()
	log.Printf("[quic] initialized (port: %d, tls_verify: %v)", q.port, q.tlsVerify)
	return nil
}

func (q *QUIC) Connect(remoteAddr string) error {
	ctx := q.CancelContext()
	q.KillOldProcess()
	if err := q.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	q.connMu.Lock()
	tlsVerify := q.tlsVerify
	q.connMu.Unlock()

	tlsConf := &tls.Config{
		InsecureSkipVerify: !tlsVerify,
		NextProtos:         []string{"nyxora-quic"},
	}

	port := q.BasePort()
	endpoint := fmt.Sprintf("%s:%d", remoteAddr, port)
	conn, err := quic.DialAddr(ctx, endpoint, tlsConf, &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	})
	if err != nil {
		q.Logf("QUIC dial failed to %s: %v, falling back to ping-only", endpoint, err)
		q.SetStatusActive()
		return nil
	}

	q.connMu.Lock()
	q.connections = append(q.connections, conn)
	q.connMu.Unlock()

	go q.acceptStreams(ctx, conn)

	q.SetStatusActive()
	q.Logf("connected to %s:%d via QUIC", remoteAddr, port)
	return nil
}

func (q *QUIC) acceptStreams(ctx context.Context, conn *quic.Conn) {
	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				q.Logf("accept stream error: %v", err)
				return
			}
		}
		go q.handleStream(stream)
	}
}

func (q *QUIC) handleStream(stream *quic.Stream) {
	defer stream.Close()
	buf := make([]byte, 4096)
	for {
		_, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF {
				q.Logf("stream read error: %v", err)
			}
			return
		}
	}
}

func (q *QUIC) Disconnect() error {
	q.mu.Lock()
	cancel := q.cancel
	q.status = StatusInactive
	q.connMu.Lock()
	conns := q.connections
	q.connections = nil
	q.connMu.Unlock()
	q.mu.Unlock()
	cancel()
	for _, conn := range conns {
		conn.CloseWithError(0, "disconnect")
	}
	q.Logf("disconnected")
	return nil
}

func (q *QUIC) Status() Status    { return q.BaseStatus() }
func (q *QUIC) Metrics() *Metrics { return q.BaseMetrics() }
func (q *QUIC) Health() bool      { return q.BaseHealth() }
func (q *QUIC) Score() float64    { return q.BaseScore() }
