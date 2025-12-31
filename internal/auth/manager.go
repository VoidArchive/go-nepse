// Package auth handles NEPSE API authentication using embedded WASM
// to decode obfuscated tokens returned by the NEPSE API.
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// DefaultTokenTTL defines when to proactively refresh tokens.
// NEPSE tokens expire after ~60 seconds; we refresh at 45s to avoid
// mid-request expiration.
const DefaultTokenTTL = 45 * time.Second

// NepseHTTP abstracts HTTP calls needed for token acquisition.
type NepseHTTP interface {
	Token(ctx context.Context) (*TokenResponse, error)
}

// TokenResponse is the JSON structure from /api/authenticate/prove.
// Salt values are used to compute which characters to strip from tokens.
type TokenResponse struct {
	Salt1        int    `json:"salt1"`
	Salt2        int    `json:"salt2"`
	Salt3        int    `json:"salt3"`
	Salt4        int    `json:"salt4"`
	Salt5        int    `json:"salt5"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ServerTime   int64  `json:"serverTime"`
}

// Salts holds the 5 salt values from the NEPSE token response.
// These are used to compute POST payload IDs for certain endpoints.
type Salts struct {
	Salt1 int
	Salt2 int
	Salt3 int
	Salt4 int
	Salt5 int
}

// Manager provides thread-safe access to NEPSE authentication tokens.
// It uses singleflight to prevent thundering herd when multiple goroutines
// request tokens simultaneously during refresh.
type Manager struct {
	http   NepseHTTP
	parser *tokenParser

	maxUpdatePeriod time.Duration

	mu          sync.RWMutex
	accessToken string
	salts       Salts
	tokenTS     time.Time

	sf singleflight.Group
}

// NewManager creates a Manager with the embedded WASM token parser.
func NewManager(httpClient NepseHTTP) (*Manager, error) {
	parser, err := newTokenParser()
	if err != nil {
		return nil, fmt.Errorf("init wasm parser: %w", err)
	}
	return &Manager{
		http:            httpClient,
		parser:          parser,
		maxUpdatePeriod: DefaultTokenTTL,
	}, nil
}

// Close must be called to release WASM runtime memory.
func (m *Manager) Close() error {
	if m.parser != nil {
		return m.parser.close()
	}
	return nil
}

// AccessToken returns a valid access token, refreshing if expired.
func (m *Manager) AccessToken(ctx context.Context) (string, error) {
	if m.isValid() {
		m.mu.RLock()
		t := m.accessToken
		m.mu.RUnlock()
		return t, nil
	}
	if err := m.update(ctx); err != nil {
		return "", err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.accessToken == "" {
		return "", errors.New("empty access token after update")
	}
	return m.accessToken, nil
}

// GetSalts returns the current salt values, refreshing the token if expired.
// Salts are used to compute POST payload IDs for graph endpoints.
func (m *Manager) GetSalts(ctx context.Context) (Salts, error) {
	if m.isValid() {
		m.mu.RLock()
		s := m.salts
		m.mu.RUnlock()
		return s, nil
	}
	if err := m.update(ctx); err != nil {
		return Salts{}, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.salts, nil
}

// ForceUpdate invalidates the cache and fetches fresh tokens.
// Used after receiving 401 to force re-authentication.
func (m *Manager) ForceUpdate(ctx context.Context) error {
	m.invalidate()
	return m.update(ctx)
}

// invalidate clears the cached token, forcing the next update to fetch fresh tokens.
func (m *Manager) invalidate() {
	m.mu.Lock()
	m.accessToken = ""
	m.tokenTS = time.Time{}
	m.mu.Unlock()
}

func (m *Manager) isValid() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.accessToken == "" || m.tokenTS.IsZero() {
		return false
	}
	return time.Since(m.tokenTS) < m.maxUpdatePeriod
}

func (m *Manager) update(ctx context.Context) error {
	_, err, _ := m.sf.Do("token_update", func() (any, error) {
		if m.isValid() {
			return nil, nil
		}

		resp, err := m.http.Token(ctx)
		if err != nil {
			return nil, fmt.Errorf("token update: %w", err)
		}

		access, ts, err := m.parseResponse(*resp)
		if err != nil {
			return nil, err
		}

		m.mu.Lock()
		m.accessToken = access
		m.salts = Salts{
			Salt1: resp.Salt1,
			Salt2: resp.Salt2,
			Salt3: resp.Salt3,
			Salt4: resp.Salt4,
			Salt5: resp.Salt5,
		}
		if ts > 0 {
			m.tokenTS = time.Unix(ts, 0)
		} else {
			m.tokenTS = time.Now()
		}
		m.mu.Unlock()

		return nil, nil
	})
	return err
}

func (m *Manager) parseResponse(tr TokenResponse) (string, int64, error) {
	salts := [5]int{tr.Salt1, tr.Salt2, tr.Salt3, tr.Salt4, tr.Salt5}

	idx, err := m.parser.indicesFromSalts(salts)
	if err != nil {
		return "", 0, fmt.Errorf("wasm parse: %w", err)
	}

	parsedAccess := sliceSkipAt(tr.AccessToken, idx.access...)
	return parsedAccess, tr.ServerTime / 1000, nil
}

// sliceSkipAt strips junk characters inserted by NEPSE's token obfuscation.
func sliceSkipAt(s string, positions ...int) string {
	if len(positions) == 0 {
		return s
	}

	ps := make([]int, len(positions))
	copy(ps, positions)
	sort.Ints(ps)

	b := []byte(s)
	var out []byte
	prev := 0
	for _, p := range ps {
		if p < 0 || p >= len(b) {
			continue
		}
		out = append(out, b[prev:p]...)
		prev = p + 1
	}
	out = append(out, b[prev:]...)
	return string(out)
}

// SetAuthHeader adds the NEPSE-specific "Salter" authorization header.
func SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Salter "+token)
}
