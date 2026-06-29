package transport

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusFailed   Status = "failed"
	StatusTesting  Status = "testing"
)

type Metrics struct {
	LatencyMs   float64 `json:"latency_ms"`
	JitterMs    float64 `json:"jitter_ms"`
	PacketLoss  float64 `json:"packet_loss"`
	Bandwidth   int     `json:"bandwidth"`
	Stability   float64 `json:"stability"`
}

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

const (
	CatVPN    = "vpn"
	CatTunnel = "tunnel"
	CatRelay  = "relay"
	CatProxy  = "proxy"
)
