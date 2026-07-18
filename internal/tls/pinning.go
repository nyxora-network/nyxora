package tls

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"
)

// PinnedCert holds a pinned certificate
type PinnedCert struct {
	Name     string
	Fingerprint string
	Expiry   time.Time
}

// CertPinner implements certificate pinning
type CertPinner struct {
	mu          sync.RWMutex
	pinned      map[string]*PinnedCert
	allowedHashes []string
	verifyHostname bool
}

// NewCertPinner creates a new certificate pinner
func NewCertPinner() *CertPinner {
	return &CertPinner{
		pinned:      make(map[string]*PinnedCert),
		verifyHostname: true,
	}
}

// Pin pins a certificate by name and fingerprint
func (cp *CertPinner) Pin(name, fingerprint string, expiry time.Time) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.pinned[name] = &PinnedCert{
		Name:        name,
		Fingerprint: fingerprint,
		Expiry:      expiry,
	}
}

// PinFromConn pins a certificate from an existing connection
func (cp *CertPinner) PinFromConn(name string, conn net.Conn) error {
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return fmt.Errorf("not a TLS connection")
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return fmt.Errorf("no peer certificates")
	}

	cert := state.PeerCertificates[0]
	fp := certFingerprint(cert)

	cp.Pin(name, fp, cert.NotAfter)
	return nil
}

// PinFromPEM pins a certificate from PEM bytes
func (cp *CertPinner) PinFromPEM(name string, pemBytes []byte) error {
	cert, err := tls.X509KeyPair(pemBytes, pemBytes)
	if err != nil {
		return fmt.Errorf("parse certificate: %w", err)
	}

	if len(cert.Certificate) == 0 {
		return fmt.Errorf("no certificates")
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("parse x509: %w", err)
	}

	fp := certFingerprint(x509Cert)
	cp.Pin(name, fp, x509Cert.NotAfter)
	return nil
}

// Verify verifies a certificate against pinned certificates
func (cp *CertPinner) Verify(name string, cert *x509.Certificate) error {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	pinned, ok := cp.pinned[name]
	if !ok {
		return fmt.Errorf("no pinned certificate for %s", name)
	}

	if time.Now().After(pinned.Expiry) {
		return fmt.Errorf("pinned certificate expired")
	}

	fp := certFingerprint(cert)
	if fp != pinned.Fingerprint {
		return fmt.Errorf("certificate fingerprint mismatch: got %s, want %s", fp, pinned.Fingerprint)
	}

	return nil
}

// TLSConfig returns a TLS config with certificate pinning
func (cp *CertPinner) TLSConfig(serverName string) *tls.Config {
	return &tls.Config{
		ServerName: serverName,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			if len(rawCerts) == 0 {
				return fmt.Errorf("no certificates provided")
			}

			cert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return fmt.Errorf("parse certificate: %w", err)
			}

			return cp.Verify(serverName, cert)
		},
		MinVersion: tls.VersionTLS12,
	}
}

// List returns all pinned certificates
func (cp *CertPinner) List() []*PinnedCert {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	result := make([]*PinnedCert, 0, len(cp.pinned))
	for _, cert := range cp.pinned {
		result = append(result, cert)
	}
	return result
}

func certFingerprint(cert *x509.Certificate) string {
	hash := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(hash[:])
}
