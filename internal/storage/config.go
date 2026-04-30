// internal/storage/config.go
package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".insighta", "credentials.json")
}

func SaveConfig(cfg Config) error {
	path := GetConfigPath()
	os.MkdirAll(filepath.Dir(path), 0700)
	
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(path, data, 0600) // Restricted permissions
}