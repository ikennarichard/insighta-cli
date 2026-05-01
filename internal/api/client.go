// internal/api/client.go
package api

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "time"

    "github.com/ikennarichard/insighta-cli/internal/auth"
)

var (
    ErrNotLoggedIn    = errors.New("not logged in — run: insighta-cli login")
    ErrSessionExpired = errors.New("session expired — run: insighta-cli login")
    ErrForbidden      = errors.New("you do not have permission to perform this action")
    ErrNotFound       = errors.New("resource not found")
    ErrServerError    = errors.New("server error — please try again later")
)

type Client struct {
    baseURL    string
    httpClient *http.Client
    creds      *auth.Credentials
}

func NewClient() (*Client, error) {
    creds, err := auth.LoadCredentials()
    if err != nil {
        return nil, ErrNotLoggedIn
    }

    // Optional: Check if token looks expired (if you store ExpiresAt later)
    return &Client{
        baseURL:    creds.BaseURL,
        httpClient: &http.Client{Timeout: 15 * time.Second},
        creds:      creds,
    }, nil
}

func (c *Client) do(method, path string, body any) (*http.Response, error) {
    resp, err := c.executeRequest(method, path, body)
    if err != nil {
        return nil, err
    }

    // Auto refresh on 401
    if resp.StatusCode == http.StatusUnauthorized {
        resp.Body.Close()

        if refreshErr := c.refresh(); refreshErr != nil {
            return nil, ErrSessionExpired
        }

        // Retry the request with new token
        resp, err = c.executeRequest(method, path, body)
        if err != nil {
            return nil, err
        }
    }

    if err := c.checkStatus(resp); err != nil {
        resp.Body.Close()
        return nil, err
    }

    return resp, nil
}

// executeRequest builds and sends the request (without status checking)
func (c *Client) executeRequest(method, path string, body any) (*http.Response, error) {
    var buf bytes.Buffer
    if body != nil {
        if err := json.NewEncoder(&buf).Encode(body); err != nil {
            return nil, fmt.Errorf("failed to encode request body: %w", err)
        }
    }

    req, err := http.NewRequest(method, c.baseURL+path, &buf)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Version", "1")
    req.Header.Set("Authorization", "Bearer "+c.creds.AccessToken)

    return c.httpClient.Do(req)
}

func (c *Client) refresh() error {
    payload := map[string]string{
        "refresh_token": c.creds.RefreshToken,
    }

    resp, err := c.httpClient.Post(
        c.baseURL+"/auth/cli/refresh",
        "application/json",
        bytes.NewReader(mustMarshal(payload)),
    )
    if err != nil {
        return fmt.Errorf("refresh request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("refresh failed with status %d", resp.StatusCode)
    }

    var result struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return fmt.Errorf("failed to decode refresh response: %w", err)
    }

    // Update in-memory credentials
    c.creds.AccessToken = result.AccessToken
    c.creds.RefreshToken = result.RefreshToken
    // Optionally set ExpiresAt if your backend returns it

    return auth.SaveCredentials(c.creds)
}

func (c *Client) checkStatus(resp *http.Response) error {
    switch resp.StatusCode {
    case http.StatusOK, http.StatusCreated, http.StatusNoContent:
        return nil
    case http.StatusUnauthorized:
        return ErrSessionExpired
    case http.StatusForbidden:
        return ErrForbidden
    case http.StatusNotFound:
        return ErrNotFound
    case http.StatusTooManyRequests:
        return fmt.Errorf("rate limit exceeded")
    case http.StatusBadRequest:
        var apiErr struct {
            Message string `json:"message,omitempty"`
        }
        json.NewDecoder(resp.Body).Decode(&apiErr) // safe even if body is small
        if apiErr.Message != "" {
            return fmt.Errorf("bad request: %s", apiErr.Message)
        }
        return fmt.Errorf("bad request")
    case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
        return ErrServerError
    default:
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }
}

// Helper functions
func mustMarshal(v any) []byte {
    b, _ := json.Marshal(v)
    return b
}

func (c *Client) Get(path string) (*http.Response, error) {
    return c.do(http.MethodGet, path, nil)
}

func (c *Client) Post(path string, body any) (*http.Response, error) {
    return c.do(http.MethodPost, path, body)
}

// Convenience method to get response body as JSON
func (c *Client) GetJSON(path string, v any) error {
    resp, err := c.Get(path)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) PostJSON(path string, body, v any) error {
    resp, err := c.Post(path, body)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return json.NewDecoder(resp.Body).Decode(v)
}