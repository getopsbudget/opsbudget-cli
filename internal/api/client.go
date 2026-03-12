package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const defaultBaseURL = "https://api.opsbudget.com/v1"

// APIError represents an error response from the OpsBudget API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// Monitor represents an uptime monitor.
type Monitor struct {
	ID             string `json:"id,omitempty"`
	URL            string `json:"url"`
	Name           string `json:"name"`
	Method         string `json:"method,omitempty"`
	ExpectedStatus int    `json:"expected_status,omitempty"`
	Interval       int    `json:"interval,omitempty"`
	Enabled        *bool  `json:"enabled,omitempty"`
	Status         string `json:"status,omitempty"`
	LastCheckedAt  string `json:"last_checked_at,omitempty"`
}

// CheckRecord represents a single uptime check result.
type CheckRecord struct {
	Timestamp    string `json:"timestamp"`
	Status       string `json:"status"`
	ResponseTime int    `json:"response_time_ms"`
	StatusCode   int    `json:"status_code"`
}

// Version is set by the CLI at startup for the User-Agent header.
var Version = "dev"

// Client is an HTTP client for the OpsBudget API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new API client with the given auth token.
func NewClient(token string) *Client {
	baseURL := os.Getenv("OPSBUDGET_API_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateMonitor creates a new uptime monitor.
func (c *Client) CreateMonitor(m *Monitor) (*Monitor, error) {
	var result Monitor
	if err := c.do("POST", "/monitors", m, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListMonitors returns all monitors for the authenticated user.
func (c *Client) ListMonitors() ([]Monitor, error) {
	var result []Monitor
	if err := c.do("GET", "/monitors", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteMonitor deletes a monitor by ID.
func (c *Client) DeleteMonitor(id string) error {
	return c.do("DELETE", "/monitors/"+id, nil, nil)
}

// GetMonitorHistory returns recent check history for a monitor.
func (c *Client) GetMonitorHistory(id string, limit int) ([]CheckRecord, error) {
	path := fmt.Sprintf("/monitors/%s/history?limit=%d", id, limit)
	var result []CheckRecord
	if err := c.do("GET", path, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) do(method, path string, body, result any) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", "opsbudget-cli/"+Version)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		msg := string(respBody)
		var errResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if json.Unmarshal(respBody, &errResp) == nil {
			if errResp.Message != "" {
				msg = errResp.Message
			} else if errResp.Error != "" {
				msg = errResp.Error
			}
		}
		return &APIError{StatusCode: resp.StatusCode, Message: msg}
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("parsing response: %w", err)
		}
	}

	return nil
}
