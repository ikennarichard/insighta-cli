// internal/auth/store.go
package auth

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type Credentials struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    Username     string    `json:"username"`
    Email        string    `json:"email,omitempty"`
    Role         string    `json:"role,omitempty"`
    BaseURL      string    `json:"base_url"`
    SavedAt      time.Time `json:"saved_at"`
    ExpiresAt    time.Time `json:"expires_at,omitempty"`
}

func credentialsPath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(home, ".insighta", "credentials.json"), nil
}

func SaveCredentials(creds *Credentials) error {
    path, err := credentialsPath()
    if err != nil {
        return err
    }

    if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
        return err
    }

    creds.SavedAt = time.Now()
    data, err := json.MarshalIndent(creds, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0600)
}

func LoadCredentials() (*Credentials, error) {
    path, err := credentialsPath()
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("not logged in — run: insighta login")
        }
        return nil, fmt.Errorf("failed to read credentials: %w", err)
    }

    var creds Credentials
    if err := json.Unmarshal(data, &creds); err != nil {
        return nil, fmt.Errorf("corrupted credentials file: %w", err)
    }
    return &creds, nil
}

func ClearCredentials() error {
    path, err := credentialsPath()
    if err != nil {
        return err
    }
    if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
        return err
    }
    return nil
}