package garmin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	ssoLoginURL   = "https://sso.garmin.com/mobile/api/login"
	diOAuthURL    = "https://diauth.garmin.com/di-oauth2-service/oauth/token"
	tokenFileName = "garmin_token.json"
)

var diClientIDs = []string{
	"GARMIN_CONNECT_MOBILE_ANDROID_DI_2025Q2",
	"GARMIN_CONNECT_MOBILE_ANDROID_DI_2024Q4",
	"GARMIN_CONNECT_MOBILE_ANDROID_DI",
	"GARMIN_CONNECT_MOBILE_IOS_DI",
}

// Token holds the DI OAuth tokens and metadata.
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ExpiresAt    int64  `json:"expires_at"`
}

// IsExpired returns true if the token has expired or will expire within 60 seconds.
func (t *Token) IsExpired() bool {
	return time.Now().Unix() >= t.ExpiresAt-60
}

func tokenFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".trainer-app", tokenFileName)
}

func loadToken() (*Token, error) {
	data, err := os.ReadFile(tokenFilePath())
	if err != nil {
		return nil, err
	}
	var t Token
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func saveToken(t *Token) error {
	dir := filepath.Dir(tokenFilePath())
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(tokenFilePath(), data, 0600)
}

// ssoLogin performs the mobile SSO login and returns a service ticket.
func ssoLogin(email, password string) (string, error) {
	params := url.Values{
		"clientId": {"GCM_IOS_DARK"},
		"locale":   {"en-US"},
		"service":  {"https://mobile.integration.garmin.com/gcm/ios"},
	}

	body := map[string]interface{}{
		"username":     email,
		"password":     password,
		"rememberMe":   true,
		"captchaToken": "",
	}
	bodyJSON, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", ssoLoginURL+"?"+params.Encode(), bytes.NewReader(bodyJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://sso.garmin.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("SSO login request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("SSO login failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		ServiceTicketID string `json:"serviceTicketId"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse SSO response: %w", err)
	}
	if result.ServiceTicketID == "" {
		return "", fmt.Errorf("SSO login returned empty service ticket, response: %s", string(respBody))
	}

	return result.ServiceTicketID, nil
}

func setDIHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "GCM-Android-5.23")
	req.Header.Set("X-Garmin-User-Agent", "com.garmin.android.apps.connectmobile/5.23; ; Google/sdk_gphone64_arm64/google; Android/33; Dalvik/2.1.0")
	req.Header.Set("X-Garmin-Paired-App-Version", "10861")
	req.Header.Set("X-Garmin-Client-Platform", "Android")
	req.Header.Set("X-App-Ver", "10861")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-GCExperience", "GC5")
	req.Header.Set("Accept", "application/json")
}

// exchangeTicket exchanges a service ticket for DI OAuth tokens, trying multiple client IDs.
func exchangeTicket(ticket string) (*Token, error) {
	var lastErr error
	for _, clientID := range diClientIDs {
		token, err := tryExchangeTicket(ticket, clientID)
		if err == nil {
			return token, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("all DI client IDs failed, last error: %w", lastErr)
}

func tryExchangeTicket(ticket, clientID string) (*Token, error) {
	form := url.Values{
		"client_id":      {clientID},
		"service_ticket": {ticket},
		"grant_type":     {"https://connectapi.garmin.com/di-oauth2-service/oauth/grant/service_ticket"},
		"service_url":    {"https://mobile.integration.garmin.com/gcm/ios"},
	}

	req, err := http.NewRequest("POST", diOAuthURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":"))
	req.Header.Set("Authorization", "Basic "+authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setDIHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DI token exchange failed for %s (HTTP %d): %s", clientID, resp.StatusCode, string(respBody))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ClientID:     clientID,
		ExpiresAt:    time.Now().Unix() + tokenResp.ExpiresIn,
	}, nil
}

// refreshToken refreshes the DI OAuth token.
func refreshToken(t *Token) (*Token, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {t.ClientID},
		"refresh_token": {t.RefreshToken},
	}

	req, err := http.NewRequest("POST", diOAuthURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(t.ClientID + ":"))
	req.Header.Set("Authorization", "Basic "+authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setDIHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	return &Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ClientID:     t.ClientID,
		ExpiresAt:    time.Now().Unix() + tokenResp.ExpiresIn,
	}, nil
}
