// cmd/root.go
package cmd

import (
    "errors"
    "fmt"
    "os"

    "github.com/ikennarichard/insighta-cli/cmd/profiles"
    "github.com/ikennarichard/insighta-cli/internal/api"
    "github.com/spf13/cobra"
)

var (
    BaseURL string
    rootCmd = &cobra.Command{
        Use:   "insighta",
        Short: "Insighta CLI - Interact with Insighta+ Labs API",
        Long: `Insighta CLI allows you to manage profiles, authenticate with GitHub,
and interact with the Insighta+ backend from the terminal.`,
        SilenceUsage:  true,
        SilenceErrors: true,
        Version:       "0.1.0", // Update as needed
    }
)

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        switch {
        case errors.Is(err, api.ErrNotLoggedIn):
            fmt.Fprintln(os.Stderr, "✗ Not logged in. Please run: insighta login")
        case errors.Is(err, api.ErrSessionExpired):
            fmt.Fprintln(os.Stderr, "✗ Session expired. Please run: insighta login")
        case errors.Is(err, api.ErrForbidden):
            fmt.Fprintln(os.Stderr, "✗ Permission denied")
        case errors.Is(err, api.ErrNotFound):
            fmt.Fprintln(os.Stderr, "✗ Resource not found")
        case errors.Is(err, api.ErrServerError):
            fmt.Fprintln(os.Stderr, "✗ Server error — please try again later")
        default:
            fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
        }
        os.Exit(1)
    }
}

func init() {
    // Global flags
    rootCmd.PersistentFlags().StringVar(&BaseURL, "api-url", "http://genderize-plum.vercel.app", "Backend API base URL")

    // Add commands
    rootCmd.AddCommand(loginCmd)
    rootCmd.AddCommand(logoutCmd)
    rootCmd.AddCommand(whoamiCmd)
    rootCmd.AddCommand(profiles.ProfilesCmd)
}