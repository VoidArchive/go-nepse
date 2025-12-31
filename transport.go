package nepse

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/voidarchive/go-nepse/internal/auth"
)

// initClient initializes the HTTP transport and auth manager.
func initClient(options *Options) (*Client, error) {
	hc := options.HTTPClient
	if hc == nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !options.TLSVerification, //nolint:gosec
			},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		}
		hc = &http.Client{
			Timeout:   options.HTTPTimeout,
			Transport: transport,
		}
	}
	// NOTE: Don't modify user-provided http.Client; users are responsible for setting timeout.

	c := &Client{
		httpClient: hc,
		config:     options.Config,
		options:    options,
	}

	authManager, err := auth.NewManager(c)
	if err != nil {
		return nil, NewInternalError("failed to create auth manager", err)
	}
	c.authManager = authManager

	return c, nil
}

// Token implements auth.NepseHTTP interface.
func (c *Client) Token(ctx context.Context) (*auth.TokenResponse, error) {
	url := c.config.BaseURL + "/api/authenticate/prove"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, NewInternalError("failed to create request", err)
	}

	c.setCommonHeaders(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, MapHTTPStatusToError(resp.StatusCode, resp.Status)
	}

	var tokenResp auth.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, NewInternalError("failed to decode token response", err)
	}

	return &tokenResp, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	var lastErr error
	maxDelay := 30 * time.Second

	for attempt := 0; attempt <= c.options.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := min(c.options.RetryDelay*time.Duration(1<<uint(attempt-1)), maxDelay)

			timer := time.NewTimer(delay)
			select {
			case <-req.Context().Done():
				timer.Stop()
				return nil, req.Context().Err()
			case <-timer.C:
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = NewNetworkError(err)
			continue
		}

		// Retry on server errors and rate limits
		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			_ = resp.Body.Close()
			lastErr = MapHTTPStatusToError(resp.StatusCode, resp.Status)
			continue
		}

		return resp, nil
	}

	return nil, lastErr
}

func (c *Client) setCommonHeaders(req *http.Request) {
	// Standard headers - use pure browser UA without library identifier
	// Some NEPSE endpoints reject requests with non-browser user agents
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	// Browser fingerprint headers (required by NEPSE)
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	// Dynamic headers derived from BaseURL
	req.Header.Set("Host", strings.TrimPrefix(c.config.BaseURL, "https://"))
	req.Header.Set("Origin", c.config.BaseURL)
	req.Header.Set("Referer", c.config.BaseURL+"/")
}

// doAuthenticatedRequest executes an authenticated API request with automatic token refresh on 401.
func (c *Client) doAuthenticatedRequest(ctx context.Context, endpoint string, tokenRetry bool) (*http.Response, error) {
	token, err := c.authManager.AccessToken(ctx)
	if err != nil {
		return nil, NewInternalError("failed to get access token", err)
	}

	url := c.config.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, NewInternalError("failed to create request", err)
	}

	auth.SetAuthHeader(req, token)
	req.Header.Set("Accept", "application/json")
	c.setCommonHeaders(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Retry once on 401 with fresh token
	if resp.StatusCode == http.StatusUnauthorized && !tokenRetry {
		_ = resp.Body.Close()
		if err := c.authManager.ForceUpdate(ctx); err != nil {
			return nil, NewInternalError("failed to refresh token", err)
		}
		return c.doAuthenticatedRequest(ctx, endpoint, true)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, MapHTTPStatusToError(resp.StatusCode, resp.Status)
	}

	return resp, nil
}

func (c *Client) apiRequest(ctx context.Context, endpoint string, result any) error {
	resp, err := c.doAuthenticatedRequest(ctx, endpoint, false)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return NewInternalError("failed to decode response", err)
	}
	return nil
}

func (c *Client) apiRequestRaw(ctx context.Context, endpoint string) ([]byte, error) {
	resp, err := c.doAuthenticatedRequest(ctx, endpoint, false)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	return io.ReadAll(resp.Body)
}

// DebugRawRequest makes an authenticated request and returns the raw response.
// This is for debugging API responses.
func (c *Client) DebugRawRequest(ctx context.Context, endpoint string) ([]byte, error) {
	return c.apiRequestRaw(ctx, endpoint)
}

// doAuthenticatedPostRequest executes an authenticated POST API request.
func (c *Client) doAuthenticatedPostRequest(ctx context.Context, endpoint string, body any, tokenRetry bool) (*http.Response, error) {
	token, err := c.authManager.AccessToken(ctx)
	if err != nil {
		return nil, NewInternalError("failed to get access token", err)
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, NewInternalError("failed to marshal request body", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	url := c.config.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, NewInternalError("failed to create request", err)
	}

	auth.SetAuthHeader(req, token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	c.setCommonHeaders(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Retry once on 401 with fresh token
	if resp.StatusCode == http.StatusUnauthorized && !tokenRetry {
		_ = resp.Body.Close()
		if err := c.authManager.ForceUpdate(ctx); err != nil {
			return nil, NewInternalError("failed to refresh token", err)
		}
		return c.doAuthenticatedPostRequest(ctx, endpoint, body, true)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, MapHTTPStatusToError(resp.StatusCode, resp.Status)
	}

	return resp, nil
}

// apiPostRequest makes an authenticated POST request and decodes the JSON response.
func (c *Client) apiPostRequest(ctx context.Context, endpoint string, body any, result any) error {
	resp, err := c.doAuthenticatedPostRequest(ctx, endpoint, body, false)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return NewInternalError("failed to decode response", err)
	}
	return nil
}

// apiPostRequestRaw makes an authenticated POST request and returns raw bytes.
func (c *Client) apiPostRequestRaw(ctx context.Context, endpoint string, body any) ([]byte, error) {
	resp, err := c.doAuthenticatedPostRequest(ctx, endpoint, body, false)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	return io.ReadAll(resp.Body)
}

// DebugRawPostRequest makes an authenticated POST request and returns the raw response.
// This is for debugging API responses.
func (c *Client) DebugRawPostRequest(ctx context.Context, endpoint string, body any) ([]byte, error) {
	return c.apiPostRequestRaw(ctx, endpoint, body)
}

// DebugDecodedToken returns the WASM-decoded access token for debugging.
// This is the token that would be sent in Authorization headers.
func (c *Client) DebugDecodedToken(ctx context.Context) (string, error) {
	return c.authManager.AccessToken(ctx)
}
