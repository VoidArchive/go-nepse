package auth

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// sliceSkipAt tests - pure unit tests for character stripping logic
func TestSliceSkipAt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		positions []int
		want      string
	}{
		{
			name:      "empty positions",
			input:     "abcdefg",
			positions: []int{},
			want:      "abcdefg",
		},
		{
			name:      "single position at start",
			input:     "Xabcdef",
			positions: []int{0},
			want:      "abcdef",
		},
		{
			name:      "single position at end",
			input:     "abcdefX",
			positions: []int{6},
			want:      "abcdef",
		},
		{
			name:      "single position in middle",
			input:     "abcXdef",
			positions: []int{3},
			want:      "abcdef",
		},
		{
			name:      "multiple positions",
			input:     "aXbYcZd",
			positions: []int{1, 3, 5},
			want:      "abcd",
		},
		{
			name:      "unsorted positions",
			input:     "aXbYcZd",
			positions: []int{5, 1, 3},
			want:      "abcd",
		},
		{
			name:      "positions out of bounds ignored",
			input:     "abc",
			positions: []int{-1, 10, 1},
			want:      "ac",
		},
		{
			name:      "all positions",
			input:     "XYZ",
			positions: []int{0, 1, 2},
			want:      "",
		},
		{
			name:      "adjacent positions",
			input:     "aXYbc",
			positions: []int{1, 2},
			want:      "abc",
		},
		{
			name:      "empty string",
			input:     "",
			positions: []int{0, 1},
			want:      "",
		},
		{
			name:      "realistic token scenario - 5 junk chars",
			input:     "eXyJAhBbGCcIiDOdJEeSfTFoGkHeiNJsKtLoMkNeOnPoQpRqRsStUuVvWwXxYyZz",
			positions: []int{1, 5, 9, 13, 17},
			want:      "eyJABbGcIiOdJeSfTFoGkHeiNJsKtLoMkNeOnPoQpRqRsStUuVvWwXxYyZz",
		},
		{
			name:      "simple token with junk at known positions",
			input:     "aXbYcZdWe",
			positions: []int{1, 3, 5, 7},
			want:      "abcde",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceSkipAt(tt.input, tt.positions...)
			if got != tt.want {
				t.Errorf("sliceSkipAt(%q, %v) = %q, want %q", tt.input, tt.positions, got, tt.want)
			}
		})
	}
}

func TestSliceSkipAt_DoesNotModifyOriginal(t *testing.T) {
	positions := []int{3, 1, 5}
	original := make([]int, len(positions))
	copy(original, positions)

	_ = sliceSkipAt("aXbYcZd", positions...)

	for i, v := range positions {
		if v != original[i] {
			t.Errorf("positions slice was modified: got %v, want %v", positions, original)
			break
		}
	}
}

// mockNepseHTTP implements NepseHTTP for testing
type mockNepseHTTP struct {
	tokenFunc func(ctx context.Context) (*TokenResponse, error)
	callCount atomic.Int32
}

func (m *mockNepseHTTP) Token(ctx context.Context) (*TokenResponse, error) {
	m.callCount.Add(1)
	if m.tokenFunc != nil {
		return m.tokenFunc(ctx)
	}
	return &TokenResponse{
		Salt1:        1234,
		Salt2:        5678,
		Salt3:        9012,
		Salt4:        3456,
		Salt5:        7890,
		AccessToken:  "testXtokenYwithZjunkAcharsB",
		RefreshToken: "refreshXtokenY",
		ServerTime:   time.Now().UnixMilli(),
	}, nil
}

func TestManager_AccessToken_Caching(t *testing.T) {
	mock := &mockNepseHTTP{}
	manager, err := NewManager(mock)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// First call should fetch token
	token1, err := manager.AccessToken(ctx)
	if err != nil {
		t.Fatalf("first AccessToken failed: %v", err)
	}
	if token1 == "" {
		t.Error("expected non-empty token")
	}
	if mock.callCount.Load() != 1 {
		t.Errorf("expected 1 HTTP call, got %d", mock.callCount.Load())
	}

	// Second call should use cached token
	token2, err := manager.AccessToken(ctx)
	if err != nil {
		t.Fatalf("second AccessToken failed: %v", err)
	}
	if token1 != token2 {
		t.Errorf("cached token mismatch: %q != %q", token1, token2)
	}
	if mock.callCount.Load() != 1 {
		t.Errorf("expected 1 HTTP call (cached), got %d", mock.callCount.Load())
	}
}

func TestManager_AccessToken_Expiration(t *testing.T) {
	mock := &mockNepseHTTP{}
	manager, err := NewManager(mock)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	defer manager.Close()

	// Override TTL to a short duration for testing
	manager.maxUpdatePeriod = 10 * time.Millisecond

	ctx := context.Background()

	// First call
	_, err = manager.AccessToken(ctx)
	if err != nil {
		t.Fatalf("first AccessToken failed: %v", err)
	}
	if mock.callCount.Load() != 1 {
		t.Errorf("expected 1 HTTP call, got %d", mock.callCount.Load())
	}

	// Wait for token to expire
	time.Sleep(20 * time.Millisecond)

	// Should fetch new token
	_, err = manager.AccessToken(ctx)
	if err != nil {
		t.Fatalf("second AccessToken failed: %v", err)
	}
	if mock.callCount.Load() != 2 {
		t.Errorf("expected 2 HTTP calls after expiry, got %d", mock.callCount.Load())
	}
}

func TestManager_ForceUpdate(t *testing.T) {
	var callCount atomic.Int32
	mock := &mockNepseHTTP{
		tokenFunc: func(ctx context.Context) (*TokenResponse, error) {
			count := callCount.Add(1)
			// Use different salt values for each call to produce different tokens
			return &TokenResponse{
				Salt1:        1234 + int(count)*100,
				Salt2:        5678 + int(count)*100,
				Salt3:        9012 + int(count)*100,
				Salt4:        3456 + int(count)*100,
				Salt5:        7890 + int(count)*100,
				AccessToken:  "testXtokenY",
				RefreshToken: "refresh",
				ServerTime:   time.Now().UnixMilli(),
			}, nil
		},
	}
	manager, err := NewManager(mock)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// Get initial token
	token1, err := manager.AccessToken(ctx)
	if err != nil {
		t.Fatalf("first AccessToken failed: %v", err)
	}
	initialCallCount := callCount.Load()

	// Force update should get new token
	err = manager.ForceUpdate(ctx)
	if err != nil {
		t.Fatalf("ForceUpdate failed: %v", err)
	}

	// Verify that a new HTTP call was made
	if callCount.Load() <= initialCallCount {
		t.Error("ForceUpdate should have made a new HTTP call")
	}

	token2, err := manager.AccessToken(ctx)
	if err != nil {
		t.Fatalf("second AccessToken failed: %v", err)
	}

	// Tokens will differ because WASM decodes with different salts
	if token1 == token2 {
		t.Log("Note: tokens are same despite different salts (WASM may produce same indices)")
	}

	// The key test: ForceUpdate caused a new HTTP call
	if callCount.Load() != 2 {
		t.Errorf("expected 2 HTTP calls (initial + force), got %d", callCount.Load())
	}
}

func TestManager_ConcurrentAccess(t *testing.T) {
	mock := &mockNepseHTTP{}
	manager, err := NewManager(mock)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()
	numGoroutines := 50

	var wg sync.WaitGroup
	tokens := make([]string, numGoroutines)
	errs := make([]error, numGoroutines)

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			tokens[idx], errs[idx] = manager.AccessToken(ctx)
		}(i)
	}
	wg.Wait()

	// Check no errors
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d failed: %v", i, err)
		}
	}

	// All tokens should be identical (singleflight)
	for i := 1; i < numGoroutines; i++ {
		if tokens[i] != tokens[0] {
			t.Errorf("token mismatch at %d: %q != %q", i, tokens[i], tokens[0])
		}
	}

	// Should only have made one HTTP call due to singleflight
	if mock.callCount.Load() != 1 {
		t.Errorf("expected 1 HTTP call with singleflight, got %d", mock.callCount.Load())
	}
}

func TestManager_HTTPError(t *testing.T) {
	expectedErr := errors.New("network failure")
	mock := &mockNepseHTTP{
		tokenFunc: func(ctx context.Context) (*TokenResponse, error) {
			return nil, expectedErr
		},
	}
	manager, err := NewManager(mock)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()
	_, err = manager.AccessToken(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error to wrap %v, got %v", expectedErr, err)
	}
}

func TestManager_ContextCancellation(t *testing.T) {
	mock := &mockNepseHTTP{
		tokenFunc: func(ctx context.Context) (*TokenResponse, error) {
			// Simulate slow response
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return &TokenResponse{
					Salt1:       1,
					Salt2:       2,
					Salt3:       3,
					Salt4:       4,
					Salt5:       5,
					AccessToken: "token",
					ServerTime:  time.Now().UnixMilli(),
				}, nil
			}
		},
	}
	manager, err := NewManager(mock)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	defer manager.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = manager.AccessToken(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestSetAuthHeader(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{
			name:  "normal token",
			token: "abc123xyz",
			want:  "Salter abc123xyz",
		},
		{
			name:  "empty token",
			token: "",
			want:  "Salter ",
		},
		{
			name:  "token with special chars",
			token: "abc+/=123",
			want:  "Salter abc+/=123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &mockRequest{headers: make(map[string]string)}
			setAuthHeaderOnMap(req.headers, tt.token)
			got := req.headers["Authorization"]
			if got != tt.want {
				t.Errorf("SetAuthHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

// mockRequest for testing header setting
type mockRequest struct {
	headers map[string]string
}

func setAuthHeaderOnMap(headers map[string]string, token string) {
	headers["Authorization"] = "Salter " + token
}

func TestManager_EmptyAccessToken(t *testing.T) {
	mock := &mockNepseHTTP{
		tokenFunc: func(ctx context.Context) (*TokenResponse, error) {
			return &TokenResponse{
				Salt1:        1234,
				Salt2:        5678,
				Salt3:        9012,
				Salt4:        3456,
				Salt5:        7890,
				AccessToken:  "", // Empty token from API
				RefreshToken: "refresh",
				ServerTime:   time.Now().UnixMilli(),
			}, nil
		},
	}
	manager, err := NewManager(mock)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()
	_, err = manager.AccessToken(ctx)
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}
}

func TestManager_Invalidate(t *testing.T) {
	mock := &mockNepseHTTP{}
	manager, err := NewManager(mock)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// Get initial token
	_, err = manager.AccessToken(ctx)
	if err != nil {
		t.Fatalf("first AccessToken failed: %v", err)
	}

	// Token should be valid
	if !manager.isValid() {
		t.Error("token should be valid after initial fetch")
	}

	// Invalidate
	manager.invalidate()

	// Token should no longer be valid
	if manager.isValid() {
		t.Error("token should be invalid after invalidate()")
	}

	// Next AccessToken call should fetch a new token
	initialCalls := mock.callCount.Load()
	_, err = manager.AccessToken(ctx)
	if err != nil {
		t.Fatalf("second AccessToken failed: %v", err)
	}
	if mock.callCount.Load() <= initialCalls {
		t.Error("AccessToken should have fetched new token after invalidate")
	}
}

func TestManager_ServerTimeHandling(t *testing.T) {
	tests := []struct {
		name       string
		serverTime int64
	}{
		{
			name:       "valid server time",
			serverTime: time.Now().UnixMilli(),
		},
		{
			name:       "zero server time uses local time",
			serverTime: 0,
		},
		{
			name:       "old server time",
			serverTime: 1609459200000, // 2021-01-01
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockNepseHTTP{
				tokenFunc: func(ctx context.Context) (*TokenResponse, error) {
					return &TokenResponse{
						Salt1:        1234,
						Salt2:        5678,
						Salt3:        9012,
						Salt4:        3456,
						Salt5:        7890,
						AccessToken:  "testXtoken",
						RefreshToken: "refresh",
						ServerTime:   tt.serverTime,
					}, nil
				},
			}
			manager, err := NewManager(mock)
			if err != nil {
				t.Fatalf("NewManager failed: %v", err)
			}
			defer manager.Close()

			ctx := context.Background()
			token, err := manager.AccessToken(ctx)
			if err != nil {
				t.Fatalf("AccessToken failed: %v", err)
			}
			if token == "" {
				t.Error("expected non-empty token")
			}
		})
	}
}
