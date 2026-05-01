package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ikennarichard/insighta-cli/internal/api"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
    Use:   "whoami",
    Short: "Show the current logged-in user",
    RunE:  runWhoami,
}

func runWhoami(cmd *cobra.Command, args []string) error {
    client, err := api.NewClient()
    if err != nil {
        return err
    }

    fmt.Print("Fetching user info... ")

    resp, err := client.Get("/api/v1/users/me")
    if err != nil {
        fmt.Println("Failed")
        return err
    }
    defer resp.Body.Close()
// Read raw body first
    var result struct {
        Data struct {
            Username string `json:"username"`
            Email    string `json:"email"`
            Role     string `json:"role"`
            IsActive bool   `json:"is_active"`
        } `json:"data"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return fmt.Errorf("failed to decode whoami response: %w", err)
    }

    u := result.Data

    fmt.Println("=== Logged in as ===")
    fmt.Printf("Username : @%s\n", u.Username)
    fmt.Printf("Email    : %s\n", u.Email)
    fmt.Printf("Role     : %s\n", u.Role)
    fmt.Printf("Status   : %s\n", map[bool]string{true: "Active", false: "Inactive"}[u.IsActive])

    return nil
}