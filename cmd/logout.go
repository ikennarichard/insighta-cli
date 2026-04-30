// cmd/logout.go
package cmd

import (
    "fmt"

    "github.com/ikennarichard/insighta-cli/internal/auth"
    "github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
    Use:   "logout",
    Short: "Log out and clear saved credentials",
    RunE:  runLogout,
}

func runLogout(cmd *cobra.Command, args []string) error {
    // Clear local credentials first (main purpose of logout for CLI)
    if err := auth.ClearCredentials(); err != nil {
        return fmt.Errorf("failed to clear credentials: %w", err)
    }

    fmt.Println("✓ Logged out successfully")
    return nil
}