package garmin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const baseURL = "https://connectapi.garmin.com"

// Client is the Garmin Connect API client using DI Bearer token auth.
type Client struct {
	httpClient *http.Client
	token      *Token
	email      string
	password   string
	mu         sync.Mutex
}

// NewClient creates a new Garmin client. Reads GARMIN_EMAIL and GARMIN_PASSWORD from env.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		email:      os.Getenv("GARMIN_EMAIL"),
		password:   os.Getenv("GARMIN_PASSWORD"),
	}
}

// EnsureAuthenticated loads a cached token, refreshes if expired, or performs full login.
func (c *Client) EnsureAuthenticated() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Already have a valid token in memory.
	if c.token != nil && !c.token.IsExpired() {
		return nil
	}

	// Try loading from file.
	if c.token == nil {
		if t, err := loadToken(); err == nil {
			c.token = t
		}
	}

	// If we have a token but it's expired, try refreshing.
	if c.token != nil && c.token.IsExpired() {
		t, err := refreshToken(c.token)
		if err == nil {
			c.token = t
			_ = saveToken(t)
			return nil
		}
		// Refresh failed, fall through to full login.
		c.token = nil
	}

	// Still valid after loading from file.
	if c.token != nil && !c.token.IsExpired() {
		return nil
	}

	// Full login flow.
	return c.fullLogin()
}

func (c *Client) fullLogin() error {
	if c.email == "" || c.password == "" {
		return fmt.Errorf("GARMIN_EMAIL and GARMIN_PASSWORD environment variables must be set")
	}

	ticket, err := ssoLogin(c.email, c.password)
	if err != nil {
		return fmt.Errorf("garmin SSO login failed: %w", err)
	}

	token, err := exchangeTicket(ticket)
	if err != nil {
		return fmt.Errorf("garmin DI token exchange failed: %w", err)
	}

	c.token = token
	_ = saveToken(token)
	return nil
}

// forceRefresh forces a token refresh (used on 401 retry).
func (c *Client) forceRefresh() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != nil {
		t, err := refreshToken(c.token)
		if err == nil {
			c.token = t
			_ = saveToken(t)
			return nil
		}
	}

	return c.fullLogin()
}

// doRequest makes an authenticated API request to the Garmin Connect API.
func (c *Client) doRequest(method, path string, body interface{}) (json.RawMessage, error) {
	if err := c.EnsureAuthenticated(); err != nil {
		return nil, err
	}

	resp, err := c.executeRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	// On 401, refresh and retry once.
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		if err := c.forceRefresh(); err != nil {
			return nil, fmt.Errorf("re-authentication failed: %w", err)
		}
		resp, err = c.executeRequest(method, path, body)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request %s %s failed (HTTP %d): %s", method, path, resp.StatusCode, string(respBody))
	}

	return json.RawMessage(respBody), nil
}

func (c *Client) executeRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.mu.Lock()
	accessToken := c.token.AccessToken
	c.mu.Unlock()

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GCM-Android-5.23")
	req.Header.Set("X-Garmin-User-Agent", "com.garmin.android.apps.connectmobile/5.23; ; Google/sdk_gphone64_arm64/google; Android/33; Dalvik/2.1.0")
	req.Header.Set("X-Garmin-Paired-App-Version", "10861")
	req.Header.Set("X-Garmin-Client-Platform", "Android")
	req.Header.Set("X-App-Ver", "10861")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-GCExperience", "GC5")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}
