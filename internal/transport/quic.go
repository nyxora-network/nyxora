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
	connections []*quic.Conn
	mu          sync.Mutex
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
	log.Printf("[quic] initialized (port: %d)", q.port)
	return nil
}

func (q *QUIC) Connect(remoteAddr string) error {
	ctx := q.CancelContext()
	q.KillOldProcess()
	if err := q.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"nyxora-quic"},
	}

	endpoint := fmt.Sprintf("%s:%d", remoteAddr, q.port)
	conn, err := quic.DialAddr(ctx, endpoint, tlsConf, &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	})
	if err != nil {
		q.Logf("QUIC dial failed to %s: %v, falling back to ping-only", endpoint, err)
		q.SetStatusActive()
		return nil
	}

	q.mu.Lock()
	q.connections = append(q.connections, conn)
	q.mu.Unlock()

	go q.acceptStreams(ctx, conn)

	q.SetStatusActive()
	q.Logf("connected to %s:%d via QUIC", remoteAddr, q.port)
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
	q.cancel()
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, conn := range q.connections {
		conn.CloseWithError(0, "disconnect")
	}
	q.connections = nil
	q.status = StatusInactive
	q.Logf("disconnected")
	return nil
}

func (q *QUIC) Status() Status    { return q.BaseStatus() }
func (q *QUIC) Metrics() *Metrics { return q.BaseMetrics() }
func (q *QUIC) Health() bool      { return q.BaseHealth() }
func (q *QUIC) Score() float64    { return q.BaseScore() }
