package dns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// DoHClient implements DNS-over-HTTPS
type DoHClient struct {
	mu          sync.RWMutex
	endpoint    string
	client      *http.Client
	cache       map[string]*DoHCacheEntry
	cacheTTL    time.Duration
}

// DoHCacheEntry represents a cached DNS response
type DoHCacheEntry struct {
	Answers   []DoHAnswer
	Expiry    time.Time
}

// DoHAnswer represents a DNS answer from DoH
type DoHAnswer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

// DoHResponse represents a DoH JSON response
type DoHResponse struct {
	Status   int          `json:"Status"`
	TC       bool         `json:"TC"`
	RD       bool         `json:"RD"`
	RA       bool         `json:"RA"`
	AD       bool         `json:"AD"`
	CD       bool         `json:"CD"`
	Question []DoHQuestion `json:"Question"`
	Answer   []DoHAnswer   `json:"Answer"`
}

// DoHQuestion represents a DNS question
type DoHQuestion struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

// DoHProviders lists popular DoH providers
var DoHProviders = map[string]string{
	"cloudflare": "https://cloudflare-dns.com/dns-query",
	"google":     "https://dns.google/resolve",
	"quad9":      "https://dns.quad9.net/dns-query",
	"alidns":     "https://dns.alidns.com/resolve",
}

// NewDoHClient creates a new DoH client
func NewDoHClient(provider string) *DoHClient {
	endpoint, ok := DoHProviders[provider]
	if !ok {
		endpoint = provider // Use as custom endpoint
	}

	return &DoHClient{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache:    make(map[string]*DoHCacheEntry),
		cacheTTL: 5 * time.Minute,
	}
}

// Lookup performs a DNS lookup via DoH
func (d *DoHClient) Lookup(domain string) ([]DoHAnswer, error) {
	d.mu.RLock()
	if entry, ok := d.cache[domain]; ok && time.Now().Before(entry.Expiry) {
		d.mu.RUnlock()
		return entry.Answers, nil
	}
	d.mu.RUnlock()

	// Build request URL
	u, err := url.Parse(d.endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint: %w", err)
	}

	q := u.Query()
	q.Set("name", domain)
	q.Set("type", "1") // A record
	u.RawQuery = q.Encode()

	// Make request
	resp, err := d.client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("doh request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse response
	var dohResp DoHResponse
	if err := json.Unmarshal(body, &dohResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if dohResp.Status != 0 {
		return nil, fmt.Errorf("doh error: status %d", dohResp.Status)
	}

	// Cache result
	d.mu.Lock()
	d.cache[domain] = &DoHCacheEntry{
		Answers: dohResp.Answer,
		Expiry:  time.Now().Add(d.cacheTTL),
	}
	d.mu.Unlock()

	return dohResp.Answer, nil
}

// ClearCache clears the DNS cache
func (d *DoHClient) ClearCache() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cache = make(map[string]*DoHCacheEntry)
}

// Stats returns cache statistics
func (d *DoHClient) Stats() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return map[string]interface{}{
		"endpoint":   d.endpoint,
		"cache_size": len(d.cache),
		"cache_ttl":  d.cacheTTL.String(),
	}
}
