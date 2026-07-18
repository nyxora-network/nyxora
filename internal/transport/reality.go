package transport

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Reality implements the Transport interface for Reality protocol
type Reality struct {
	BaseTransport
	mu            sync.RWMutex
	privateKey    *ecdh.PrivateKey
	publicKey     *ecdh.PublicKey
	serverName    string
	dest          string
	privateKeyStr string
	publicKeyStr  string
	shortID       string
	conn          net.Conn
	tlsConn       *tls.Conn
	// Encryption state
	encKey        []byte // AES-256-GCM encryption key
	decKey        []byte // AES-256-GCM decryption key
	sendNonce     uint64 // Counter for send nonce
	recvNonce     uint64 // Counter for receive nonce
	sendMu        sync.Mutex
	recvMu        sync.Mutex
}

// NewReality creates a new Reality transport
func NewReality() *Reality {
	return &Reality{
		BaseTransport: NewBase("reality", "reality", 443, ScoringWeights{0.35, 0.25, 0.15, 0.25}, 120),
		serverName:    "",
		dest:          "",
		shortID:       "",
	}
}

func (r *Reality) Name() string  { return r.BaseName() }
func (r *Reality) Type() string { return r.BaseType() }

// Init initializes the Reality transport with config
func (r *Reality) Init(cfg map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if sn, ok := cfg["server_name"]; ok {
		r.serverName = sn
	}
	if dest, ok := cfg["dest"]; ok {
		r.dest = dest
	}
	if pk, ok := cfg["private_key"]; ok {
		r.privateKeyStr = pk
	}
	if pub, ok := cfg["public_key"]; ok {
		r.publicKeyStr = pub
	}
	if sid, ok := cfg["short_id"]; ok {
		r.shortID = sid
	}

	// Generate key pair if not provided
	if r.privateKeyStr == "" {
		privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("generate key: %w", err)
		}
		r.privateKey = privateKey
		r.publicKey = privateKey.PublicKey()
		r.privateKeyStr = hex.EncodeToString(privateKey.Bytes())
		r.publicKeyStr = hex.EncodeToString(r.publicKey.Bytes())
	}

	log.Printf("[reality] initialized (server: %s, dest: %s, short_id: %s)", r.serverName, r.dest, r.shortID)
	return nil
}

// Connect establishes a Reality connection
func (r *Reality) Connect(remoteAddr string) error {
	r.CancelContext()
	if err := r.BaseConnectInit(remoteAddr); err != nil {
		return err
	}

	r.mu.Lock()
	serverName := r.serverName
	dest := r.dest
	shortID := r.shortID
	privateKey := r.privateKey
	publicKey := r.publicKey
	r.mu.Unlock()

	if serverName == "" {
		serverName = remoteAddr
	}

	// Establish TCP connection
	addr := fmt.Sprintf("%s:%d", remoteAddr, r.port)
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		r.Logf("tcp connection failed: %v", err)
		r.SetStatusFailed()
		return err
	}

	// Perform Reality handshake
	if err := r.performHandshake(conn, serverName, dest, shortID, privateKey, publicKey); err != nil {
		conn.Close()
		r.Logf("reality handshake failed: %v", err)
		r.SetStatusFailed()
		return err
	}

	r.mu.Lock()
	r.conn = conn
	r.mu.Unlock()

	r.SetStatusActive()
	r.Logf("connected to %s via Reality", remoteAddr)
	return nil
}

// performHandshake performs the Reality protocol handshake with X25519 ECDH and AES-GCM encryption
func (r *Reality) performHandshake(conn net.Conn, serverName, dest, shortID string, privateKey *ecdh.PrivateKey, publicKey *ecdh.PublicKey) error {
	// Generate ephemeral key pair for this session
	ephemeralPriv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("generate ephemeral key: %w", err)
	}
	ephemeralPub := ephemeralPriv.PublicKey()

	// Generate random nonce for handshake
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	// Create handshake packet (Reality protocol format)
	// Version (1 byte)
	packet := []byte{0x01}

	// Ephemeral public key (32 bytes for X25519)
	pubKeyBytes := ephemeralPub.Bytes()
	packet = append(packet, pubKeyBytes...)

	// Short ID (8 bytes)
	shortIDBytes := make([]byte, 8)
	if shortID != "" {
		sidBytes, _ := hex.DecodeString(shortID)
		if len(sidBytes) >= 8 {
			copy(shortIDBytes, sidBytes[:8])
		}
	} else {
		rand.Read(shortIDBytes)
	}
	packet = append(packet, shortIDBytes...)

	// Nonce (32 bytes)
	packet = append(packet, nonce...)

	// Server name length + server name
	nameBytes := []byte(serverName)
	packet = append(packet, byte(len(nameBytes)))
	packet = append(packet, nameBytes...)

	// HMAC for authentication (using static private key)
	h := hmac.New(sha256.New, privateKey.Bytes())
	h.Write(packet)
	mac := h.Sum(nil)[:16] // Truncate to 16 bytes
	packet = append(packet, mac...)

	// Send handshake
	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("send handshake: %w", err)
	}

	// Read response
	response := make([]byte, 1024)
	n, err := io.ReadFull(conn, response[:512])
	if err != nil && err != io.ErrUnexpectedEOF {
		return fmt.Errorf("read response: %w", err)
	}

	// Parse response: version(1) + server ephemeral pub(32) + encrypted tag(16)
	if n < 49 {
		return fmt.Errorf("invalid response length: %d", n)
	}

	// Extract server's ephemeral public key
	serverEphemeralPubBytes := response[1:33]
	serverEphemeralPub, err := ecdh.X25519().NewPublicKey(serverEphemeralPubBytes)
	if err != nil {
		return fmt.Errorf("parse server public key: %w", err)
	}

	// Derive shared secret using X25519 ECDH
	sharedSecret, err := ephemeralPriv.ECDH(serverEphemeralPub)
	if err != nil {
		return fmt.Errorf("derive shared secret: %w", err)
	}

	// Derive encryption keys using HKDF-like construction
	// enc_key = HMAC-SHA256(shared_secret, "reality-enc" || nonce)
	// dec_key = HMAC-SHA256(shared_secret, "reality-dec" || nonce)
	r.encKey = r.deriveKey(sharedSecret, nonce, "reality-enc")
	r.decKey = r.deriveKey(sharedSecret, nonce, "reality-dec")

	// Reset nonces
	r.sendNonce = 0
	r.recvNonce = 0

	// Verify server authentication (check MAC in response)
	serverMAC := response[33:49]
	h = hmac.New(sha256.New, r.decKey)
	h.Write(response[:33])
	expectedMAC := h.Sum(nil)[:16]

	if !hmac.Equal(serverMAC, expectedMAC) {
		return fmt.Errorf("server authentication failed")
	}

	return nil
}

// deriveKey derives an encryption key using HMAC-SHA256
func (r *Reality) deriveKey(sharedSecret, nonce []byte, info string) []byte {
	h := hmac.New(sha256.New, sharedSecret)
	h.Write([]byte(info))
	h.Write(nonce)
	return h.Sum(nil)[:32] // 256-bit key for AES-256-GCM
}

// Disconnect closes the Reality connection
func (r *Reality) Disconnect() error {
	r.mu.Lock()
	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
	if r.tlsConn != nil {
		r.tlsConn.Close()
		r.tlsConn = nil
	}
	r.mu.Unlock()
	return r.BaseDisconnect()
}

func (r *Reality) Status() Status  { return r.BaseStatus() }
func (r *Reality) Metrics() *Metrics { return r.BaseMetrics() }
func (r *Reality) Health() bool    { return r.BaseHealth() }
func (r *Reality) Score() float64  { return r.BaseScore() }

// GetPublicKey returns the public key as hex string
func (r *Reality) GetPublicKey() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.publicKeyStr
}

// GetPrivateKey returns the private key as hex string
func (r *Reality) GetPrivateKey() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.privateKeyStr
}

// Read reads and decrypts data from the Reality connection
func (r *Reality) Read(b []byte) (int, error) {
	r.mu.RLock()
	conn := r.conn
	r.mu.RUnlock()

	if conn == nil {
		return 0, fmt.Errorf("not connected")
	}

	// Read frame: length(4) + encrypted_data
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lengthBuf); err != nil {
		return 0, fmt.Errorf("read frame length: %w", err)
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	if length > uint32(len(b)+aes.BlockSize+16) {
		return 0, fmt.Errorf("frame too large: %d", length)
	}

	encrypted := make([]byte, length)
	if _, err := io.ReadFull(conn, encrypted); err != nil {
		return 0, fmt.Errorf("read encrypted data: %w", err)
	}

	// Decrypt
	plaintext, err := r.decrypt(encrypted)
	if err != nil {
		return 0, fmt.Errorf("decrypt: %w", err)
	}

	n := copy(b, plaintext)
	return n, nil
}

// Write encrypts and writes data to the Reality connection
func (r *Reality) Write(b []byte) (int, error) {
	r.mu.RLock()
	conn := r.conn
	r.mu.RUnlock()

	if conn == nil {
		return 0, fmt.Errorf("not connected")
	}

	// Encrypt
	encrypted, err := r.encrypt(b)
	if err != nil {
		return 0, fmt.Errorf("encrypt: %w", err)
	}

	// Write frame: length(4) + encrypted_data
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, uint32(len(encrypted)))

	if _, err := conn.Write(lengthBuf); err != nil {
		return 0, fmt.Errorf("write frame length: %w", err)
	}

	if _, err := conn.Write(encrypted); err != nil {
		return 0, fmt.Errorf("write encrypted data: %w", err)
	}

	return len(b), nil
}

// encrypt encrypts data using AES-256-GCM
func (r *Reality) encrypt(plaintext []byte) ([]byte, error) {
	r.sendMu.Lock()
	defer r.sendMu.Unlock()

	block, err := aes.NewCipher(r.encKey)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	// Nonce: 8-byte counter + 4-byte random = 12 bytes (GCM standard)
	nonce := make([]byte, gcm.NonceSize())
	binary.LittleEndian.PutUint64(nonce[:8], r.sendNonce)
	r.sendNonce++

	// Seal: nonce || ciphertext || tag
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// decrypt decrypts data using AES-256-GCM
func (r *Reality) decrypt(ciphertext []byte) ([]byte, error) {
	r.recvMu.Lock()
	defer r.recvMu.Unlock()

	block, err := aes.NewCipher(r.decKey)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and encrypted data
	nonce := ciphertext[:gcm.NonceSize()]
	encrypted := ciphertext[gcm.NonceSize():]

	// Open: verify tag and decrypt
	plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %w", err)
	}

	return plaintext, nil
}

// GenerateKeys generates a new X25519 key pair for Reality
func GenerateRealityKeys() (privateKey, publicKey string, err error) {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	pub := priv.PublicKey()
	return hex.EncodeToString(priv.Bytes()), hex.EncodeToString(pub.Bytes()), nil
}
