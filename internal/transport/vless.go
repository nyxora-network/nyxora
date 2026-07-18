package transport

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// VLESS implements the Transport interface for VLESS protocol
type VLESS struct {
	BaseTransport
	mu          sync.RWMutex
	tlsConfig   *tls.Config
	listenAddr  string
	destination string
	uuid        string
	serverName  string
	flow        string
	conn        net.Conn
}

// NewVLESS creates a new VLESS transport
func NewVLESS() *VLESS {
	return &VLESS{
		BaseTransport: NewBase("vless", "vless", 443, ScoringWeights{0.30, 0.30, 0.15, 0.25}, 100),
		uuid:          "",
		serverName:    "",
		flow:          "xtls-rprx-vision",
	}
}

func (v *VLESS) Name() string  { return v.BaseName() }
func (v *VLESS) Type() string { return v.BaseType() }

// Init initializes the VLESS transport with config
func (v *VLESS) Init(cfg map[string]string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if uuid, ok := cfg["uuid"]; ok {
		v.uuid = uuid
	}
	if sn, ok := cfg["server_name"]; ok {
		v.serverName = sn
	}
	if flow, ok := cfg["flow"]; ok {
		v.flow = flow
	}
	if addr, ok := cfg["listen_addr"]; ok {
		v.listenAddr = addr
	}
	if dest, ok := cfg["destination"]; ok {
		v.destination = dest
	}

	if v.uuid == "" {
		return fmt.Errorf("vless: uuid is required")
	}

	log.Printf("[vless] initialized (server: %s, flow: %s)", v.serverName, v.flow)
	return nil
}

// Connect establishes a VLESS connection
func (v *VLESS) Connect(remoteAddr string) error {
	v.CancelContext()
	if err := v.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	v.mu.Lock()
	serverName := v.serverName
	flow := v.flow
	uuid := v.uuid
	v.mu.Unlock()

	if serverName == "" {
		serverName = remoteAddr
	}

	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		NextProtos: []string{"h2", "http/1.1"},
	}

	v.mu.Lock()
	v.tlsConfig = tlsConfig
	v.mu.Unlock()

	// Establish TCP connection
	addr := fmt.Sprintf("%s:%d", remoteAddr, v.port)
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		v.Logf("tcp connection failed: %v", err)
		v.SetStatusFailed()
		return err
	}

	// Upgrade to TLS
	tlsConn := tls.Client(conn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		v.Logf("tls handshake failed: %v", err)
		v.SetStatusFailed()
		return err
	}

	// Send VLESS header
	if err := v.sendVLESSHeader(tlsConn, uuid, remoteAddr, flow); err != nil {
		tlsConn.Close()
		v.Logf("vless header failed: %v", err)
		v.SetStatusFailed()
		return err
	}

	v.mu.Lock()
	v.conn = tlsConn
	v.mu.Unlock()

	v.SetStatusActive()
	v.Logf("connected to %s via VLESS (flow: %s)", remoteAddr, flow)
	return nil
}

// sendVLESSHeader sends the VLESS protocol header
func (v *VLESS) sendVLESSHeader(conn net.Conn, uuid, address, flow string) error {
	// VLESS version
	conn.Write([]byte{0x00})

	// UUID (16 bytes)
	uuidBytes, err := parseUUID(uuid)
	if err != nil {
		return fmt.Errorf("parse uuid: %w", err)
	}
	conn.Write(uuidBytes)

	// Payload
	payload := make([]byte, 0)

	// Addon length
	payload = append(payload, 0x00)

	// Command: TCP
	payload = append(payload, 0x01)

	// Port (big endian)
	port := uint16(v.port)
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, port)
	payload = append(payload, portBytes...)

	// Address type and address
	addrType, addrBytes := encodeAddress(address)
	payload = append(payload, addrType)
	payload = append(payload, addrBytes...)

	// Flow
	if flow != "" {
		flowBytes := []byte(flow)
		payload = append(payload, byte(len(flowBytes)))
		payload = append(payload, flowBytes...)
	} else {
		payload = append(payload, 0x00)
	}

	// Write payload length and payload
	lengthByte := byte(len(payload))
	conn.Write([]byte{lengthByte})
	conn.Write(payload)

	return nil
}

// Disconnect closes the VLESS connection
func (v *VLESS) Disconnect() error {
	v.mu.Lock()
	if v.conn != nil {
		v.conn.Close()
		v.conn = nil
	}
	v.tlsConfig = nil
	v.mu.Unlock()
	return v.BaseDisconnect()
}

func (v *VLESS) Status() Status  { return v.BaseStatus() }
func (v *VLESS) Metrics() *Metrics { return v.BaseMetrics() }
func (v *VLESS) Health() bool    { return v.BaseHealth() }
func (v *VLESS) Score() float64  { return v.BaseScore() }

// parseUUID parses a UUID string to 16 bytes
func parseUUID(uuid string) ([]byte, error) {
	// Simple UUID parsing (remove hyphens, convert hex)
	clean := ""
	for _, c := range uuid {
		if c != '-' {
			clean += string(c)
		}
	}

	if len(clean) != 32 {
		return nil, fmt.Errorf("invalid UUID length: %d", len(clean))
	}

	b := make([]byte, 16)
	for i := 0; i < 16; i++ {
		var hi, lo byte
		if clean[i*2] >= '0' && clean[i*2] <= '9' {
			hi = clean[i*2] - '0'
		} else if clean[i*2] >= 'a' && clean[i*2] <= 'f' {
			hi = clean[i*2] - 'a' + 10
		} else if clean[i*2] >= 'A' && clean[i*2] <= 'F' {
			hi = clean[i*2] - 'A' + 10
		} else {
			return nil, fmt.Errorf("invalid hex character: %c", clean[i*2])
		}

		if clean[i*2+1] >= '0' && clean[i*2+1] <= '9' {
			lo = clean[i*2+1] - '0'
		} else if clean[i*2+1] >= 'a' && clean[i*2+1] <= 'f' {
			lo = clean[i*2+1] - 'a' + 10
		} else if clean[i*2+1] >= 'A' && clean[i*2+1] <= 'F' {
			lo = clean[i*2+1] - 'A' + 10
		} else {
			return nil, fmt.Errorf("invalid hex character: %c", clean[i*2+1])
		}

		b[i] = hi<<4 | lo
	}
	return b, nil
}

// encodeAddress encodes an address for VLESS protocol
func encodeAddress(addr string) (byte, []byte) {
	// Check if IPv4
	ip := net.ParseIP(addr)
	if ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			return 0x01, ip4
		}
		if ip6 := ip.To16(); ip6 != nil {
			return 0x03, ip6
		}
	}

	// Domain
	return 0x02, []byte(addr)
}

// GenerateRandomUUID generates a random UUID
func GenerateRandomUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// Read reads data from the VLESS connection
func (v *VLESS) Read(b []byte) (int, error) {
	v.mu.RLock()
	conn := v.conn
	v.mu.RUnlock()

	if conn == nil {
		return 0, fmt.Errorf("not connected")
	}
	return conn.Read(b)
}

// Write writes data to the VLESS connection
func (v *VLESS) Write(b []byte) (int, error) {
	v.mu.RLock()
	conn := v.conn
	v.mu.RUnlock()

	if conn == nil {
		return 0, fmt.Errorf("not connected")
	}
	return conn.Write(b)
}
