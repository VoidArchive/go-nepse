package nepse

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/voidarchive/go-nepse/internal/auth"
)

// newTestServer creates a mock NEPSE API server
func newTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

// tokenResponse returns a valid token response JSON
func tokenResponse() auth.TokenResponse {
	return auth.TokenResponse{
		Salt1:        1234,
		Salt2:        5678,
		Salt3:        9012,
		Salt4:        3456,
		Salt5:        7890,
		AccessToken:  "testXtokenYwithZjunkAcharsB",
		RefreshToken: "refreshXtokenY",
		ServerTime:   time.Now().UnixMilli(),
	}
}

func TestClient_TokenFetch(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/authenticate/prove" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())
			return
		}
		http.NotFound(w, r)
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	// Token fetch happens internally when making authenticated requests
	ctx := context.Background()
	tokenResp, err := client.Token(ctx)
	if err != nil {
		t.Fatalf("Token() failed: %v", err)
	}
	if tokenResp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
}

func TestClient_TokenRefreshOn401(t *testing.T) {
	var tokenCallCount atomic.Int32
	var apiCallCount atomic.Int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/authenticate/prove":
			tokenCallCount.Add(1)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())

		case "/api/nots/nepse-data/market-open":
			count := apiCallCount.Add(1)
			// First call returns 401, second succeeds
			if count == 1 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"isOpen": "OPEN",
			})

		default:
			http.NotFound(w, r)
		}
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL:   server.URL,
			Endpoints: DefaultEndpoints(),
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	status, err := client.MarketStatus(ctx)
	if err != nil {
		t.Fatalf("MarketStatus() failed: %v", err)
	}

	if status.IsOpen != "OPEN" {
		t.Errorf("expected IsOpen=OPEN, got %q", status.IsOpen)
	}

	// Should have fetched token twice (initial + refresh after 401)
	if tokenCallCount.Load() != 2 {
		t.Errorf("expected 2 token calls, got %d", tokenCallCount.Load())
	}

	// Should have made 2 API calls (first 401, then success)
	if apiCallCount.Load() != 2 {
		t.Errorf("expected 2 API calls, got %d", apiCallCount.Load())
	}
}

func TestClient_RetryOn5xx(t *testing.T) {
	var callCount atomic.Int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/authenticate/prove":
			count := callCount.Add(1)
			// First 2 calls fail with 503, third succeeds
			if count <= 2 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())

		default:
			http.NotFound(w, r)
		}
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  3,
		RetryDelay:  10 * time.Millisecond, // Fast retries for testing
		Config: &Config{
			BaseURL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	tokenResp, err := client.Token(ctx)
	if err != nil {
		t.Fatalf("Token() failed: %v", err)
	}
	if tokenResp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	// Should have made 3 calls (2 failures + 1 success)
	if callCount.Load() != 3 {
		t.Errorf("expected 3 calls with retry, got %d", callCount.Load())
	}
}

func TestClient_RateLimitRetry(t *testing.T) {
	var callCount atomic.Int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/authenticate/prove" {
			count := callCount.Add(1)
			// First call returns 429, second succeeds
			if count == 1 {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())
			return
		}
		http.NotFound(w, r)
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  2,
		RetryDelay:  10 * time.Millisecond,
		Config: &Config{
			BaseURL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	tokenResp, err := client.Token(ctx)
	if err != nil {
		t.Fatalf("Token() failed: %v", err)
	}
	if tokenResp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	if callCount.Load() != 2 {
		t.Errorf("expected 2 calls, got %d", callCount.Load())
	}
}

func TestClient_MaxRetriesExceeded(t *testing.T) {
	var callCount atomic.Int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/authenticate/prove" {
			callCount.Add(1)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		http.NotFound(w, r)
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  2,
		RetryDelay:  10 * time.Millisecond,
		Config: &Config{
			BaseURL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Token(ctx)
	if err == nil {
		t.Fatal("expected error after max retries exceeded")
	}

	// Should have made MaxRetries + 1 calls (initial + retries)
	if callCount.Load() != 3 {
		t.Errorf("expected 3 calls (1 initial + 2 retries), got %d", callCount.Load())
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResponse())
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = client.Token(ctx)
	if err == nil {
		t.Error("expected context cancellation error")
	}
}

func TestClient_InvalidJSONResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/authenticate/prove" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"invalid json`)) // Malformed JSON
			return
		}
		http.NotFound(w, r)
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Token(ctx)
	if err == nil {
		t.Error("expected JSON decode error")
	}
}

func TestClient_CommonHeaders(t *testing.T) {
	var capturedHeaders http.Header

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResponse())
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Token(ctx)
	if err != nil {
		t.Fatalf("Token() failed: %v", err)
	}

	// Verify required headers are set
	requiredHeaders := []string{
		"User-Agent",
		"Accept",
		"Accept-Language",
		"Cache-Control",
	}

	for _, h := range requiredHeaders {
		if capturedHeaders.Get(h) == "" {
			t.Errorf("missing required header: %s", h)
		}
	}
}

func TestClient_AuthenticatedRequestHeader(t *testing.T) {
	var authHeader string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/authenticate/prove":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())

		case "/api/nots/nepse-data/market-open":
			authHeader = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"isOpen": "OPEN",
			})

		default:
			http.NotFound(w, r)
		}
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL:   server.URL,
			Endpoints: DefaultEndpoints(),
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.MarketStatus(ctx)
	if err != nil {
		t.Fatalf("MarketStatus() failed: %v", err)
	}

	// Auth header should start with "Salter "
	if len(authHeader) < 8 || authHeader[:7] != "Salter " {
		t.Errorf("expected auth header to start with 'Salter ', got %q", authHeader)
	}
}

func TestClient_HTTPErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    error
	}{
		{"BadRequest", http.StatusBadRequest, ErrInvalidClientRequest},
		{"Unauthorized", http.StatusUnauthorized, ErrTokenExpired},
		{"Forbidden", http.StatusForbidden, ErrUnauthorized},
		{"NotFound", http.StatusNotFound, ErrNotFound},
		{"TooManyRequests", http.StatusTooManyRequests, ErrRateLimit},
		{"InternalServerError", http.StatusInternalServerError, ErrInvalidServerResponse},
		{"BadGateway", http.StatusBadGateway, ErrInvalidServerResponse},
		{"ServiceUnavailable", http.StatusServiceUnavailable, ErrInvalidServerResponse},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})
			server := newTestServer(handler)
			defer server.Close()

			client, err := NewClient(&Options{
				BaseURL:     server.URL,
				HTTPTimeout: 5 * time.Second,
				MaxRetries:  0, // No retries for error mapping test
				Config: &Config{
					BaseURL: server.URL,
				},
			})
			if err != nil {
				t.Fatalf("NewClient failed: %v", err)
			}
			defer client.Close()

			ctx := context.Background()
			_, err = client.Token(ctx)

			if err == nil {
				t.Fatal("expected error")
			}

			nepseErr, ok := err.(*NepseError)
			if !ok {
				t.Fatalf("expected *NepseError, got %T", err)
			}

			wantNepseErr := tt.wantErr.(*NepseError)
			if nepseErr.Type != wantNepseErr.Type {
				t.Errorf("error type = %v, want %v", nepseErr.Type, wantNepseErr.Type)
			}
		})
	}
}

func TestClient_DebugRawRequest(t *testing.T) {
	expectedBody := `{"test":"data","value":123}`

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/authenticate/prove":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())

		case "/api/test/endpoint":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(expectedBody))

		default:
			http.NotFound(w, r)
		}
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	body, err := client.DebugRawRequest(ctx, "/api/test/endpoint")
	if err != nil {
		t.Fatalf("DebugRawRequest() failed: %v", err)
	}

	if string(body) != expectedBody {
		t.Errorf("body = %q, want %q", string(body), expectedBody)
	}
}

// Benchmark for transport layer
func BenchmarkClient_TokenFetch(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResponse())
	})
	server := newTestServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL: server.URL,
		},
	})
	if err != nil {
		b.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Token(ctx)
		if err != nil {
			b.Fatalf("Token() failed: %v", err)
		}
	}
}
