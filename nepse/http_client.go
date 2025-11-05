package nepse

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/voidarchive/nepseauth/auth"
)

// HTTPClient implements the NEPSE HTTP client with authentication
type HTTPClient struct {
	client      *http.Client
	config      *Config
	authManager *auth.Manager
	options     *Options
}

// NewHTTPClient creates a new HTTP client for NEPSE API
func NewHTTPClient(options *Options) (*HTTPClient, error) {
	if options == nil {
		options = DefaultOptions()
	}

	if options.Config == nil {
		options.Config = DefaultConfig()
	}

	// Create or use provided HTTP client
	httpClient := options.HTTPClient
	if httpClient == nil {
		// Only construct transport if no client supplied
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{ //nolint:gosec // user controls via TLSVerification
				InsecureSkipVerify: !options.TLSVerification,
			},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			// Rely on Go's transparent gzip decompression (DisableCompression=false)
		}
		httpClient = &http.Client{
			Timeout:   options.HTTPTimeout,
			Transport: transport,
		}
	} else if httpClient.Timeout == 0 {
		// Ensure a reasonable default timeout if caller didn't set one
		httpClient.Timeout = options.HTTPTimeout
	}

	nepseClient := &HTTPClient{
		client:  httpClient,
		config:  options.Config,
		options: options,
	}

	// Create auth manager
	authManager, err := auth.NewManager(nepseClient)
	if err != nil {
		return nil, NewInternalError("failed to create auth manager", err)
	}
	nepseClient.authManager = authManager

	return nepseClient, nil
}

// GetTokens implements the auth.NepseHTTP interface for the auth package
func (h *HTTPClient) GetTokens(ctx context.Context) (*auth.TokenResponse, error) {
	url := h.config.BaseURL + "/api/authenticate/prove"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, NewInternalError("failed to create request", err)
	}

	h.setCommonHeaders(req, false)

	resp, err := h.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, MapHTTPStatusToError(resp.StatusCode, resp.Status)
	}

	body, err := h.getResponseBody(resp)
	if err != nil {
		return nil, NewInternalError("failed to read response body", err)
	}
	defer body.Close()

	var tokenResp auth.TokenResponse
	if err := json.NewDecoder(body).Decode(&tokenResp); err != nil {
		return nil, NewInternalError("failed to decode token response", err)
	}

	return &tokenResp, nil
}

// RefreshTokens implements the auth.NepseHTTP interface for the auth package
func (h *HTTPClient) RefreshTokens(ctx context.Context, refreshToken string) (*auth.TokenResponse, error) {
	url := h.config.BaseURL + "/api/authenticate/refresh-token"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, NewInternalError("failed to create request", err)
	}

	req.Header.Set("Authorization", "Salter "+refreshToken)
	h.setCommonHeaders(req, true)

	resp, err := h.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, MapHTTPStatusToError(resp.StatusCode, resp.Status)
	}

	body, err := h.getResponseBody(resp)
	if err != nil {
		return nil, NewInternalError("failed to read response body", err)
	}
	defer body.Close()

	var tokenResp auth.TokenResponse
	if err := json.NewDecoder(body).Decode(&tokenResp); err != nil {
		return nil, NewInternalError("failed to decode token response", err)
	}

	return &tokenResp, nil
}

// doRequest performs HTTP request with retry logic
func (h *HTTPClient) doRequest(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= h.options.MaxRetries; attempt++ {
		if attempt > 0 {
			// Calculate backoff delay
			delay := minDuration(h.options.RetryDelay*time.Duration(1<<uint(attempt-1)), 30*time.Second)
			time.Sleep(delay)
		}

		resp, err := h.client.Do(req)
		if err != nil {
			lastErr = NewNetworkError(err)
			continue
		}

		// Check if we should retry based on status code
		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = MapHTTPStatusToError(resp.StatusCode, resp.Status)
			if !lastErr.(*NepseError).IsRetryable() {
				return nil, lastErr
			}
			continue
		}

		return resp, nil
	}

	return nil, lastErr
}

// getResponseBody handles gzip decompression
func (h *HTTPClient) getResponseBody(resp *http.Response) (io.ReadCloser, error) {
	// Let net/http handle decompression transparently.
	return resp.Body, nil
}

// setCommonHeaders sets common HTTP headers for requests
func (h *HTTPClient) setCommonHeaders(req *http.Request, _ bool) {
	// Set headers from config
	for key, value := range h.config.Headers {
		if key == "Host" {
			req.Header.Set(key, strings.Replace(h.config.BaseURL, "https://", "", 1))
		} else if key == "Referer" {
			req.Header.Set(key, h.config.BaseURL+"/")
		} else if value != "" {
			req.Header.Set(key, value)
		}
	}

	// Set additional headers for better browser mimicking
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Origin", h.config.BaseURL)
}

// apiRequest performs an authenticated API request
func (h *HTTPClient) apiRequest(ctx context.Context, endpoint string, result any) error {
	return h.apiRequestWithRetry(ctx, endpoint, result, 0)
}

// apiRequestWithRetry performs an authenticated API request with token refresh retry
func (h *HTTPClient) apiRequestWithRetry(ctx context.Context, endpoint string, result any, retryCount int) error {
	token, err := h.authManager.AccessToken(ctx)
	if err != nil {
		return NewInternalError("failed to get access token", err)
	}

	url := h.config.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NewInternalError("failed to create request", err)
	}

	// Set authenticated headers
	auth.AuthHeader(req, token)
	req.Header.Set("Content-Type", "application/json")
	h.setCommonHeaders(req, true)

	resp, err := h.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle token expiration
	if resp.StatusCode == http.StatusUnauthorized && retryCount == 0 {
		if err := h.authManager.ForceUpdate(ctx); err != nil {
			return NewInternalError("failed to refresh token", err)
		}
		return h.apiRequestWithRetry(ctx, endpoint, result, retryCount+1)
	}

	if resp.StatusCode != http.StatusOK {
		return MapHTTPStatusToError(resp.StatusCode, resp.Status)
	}

	body, err := h.getResponseBody(resp)
	if err != nil {
		return NewInternalError("failed to read response body", err)
	}
	defer body.Close()

	if err := json.NewDecoder(body).Decode(result); err != nil {
		return NewInternalError("failed to decode response", err)
	}

	return nil
}

// SetTLSVerification sets TLS verification on/off
func (h *HTTPClient) SetTLSVerification(enabled bool) {
	if transport, ok := h.client.Transport.(*http.Transport); ok {
		transport.TLSClientConfig.InsecureSkipVerify = !enabled
	}
	h.options.TLSVerification = enabled
}

// GetConfig returns the current configuration
func (h *HTTPClient) GetConfig() *Config {
	return h.config
}

// TestGetRequest performs a test GET request to any endpoint (for debugging)
func (h *HTTPClient) TestGetRequest(ctx context.Context, endpoint string, result any) error {
	return h.apiRequest(ctx, endpoint, result)
}

// Close closes the HTTP client and auth manager
func (h *HTTPClient) Close(ctx context.Context) error {
	if h.authManager != nil {
		return h.authManager.Close(ctx)
	}
	return nil
}
