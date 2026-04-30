// cmd/login.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ikennarichard/insighta-cli/internal/auth"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "Authenticate with GitHub",
    Long:  "Starts OAuth 2.0 flow with PKCE and stores tokens locally.",
    RunE:  runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
    state, err := auth.GenerateState()
    if err != nil {
        return fmt.Errorf("failed to generate state: %w", err)
    }
    codeVerifier, err := auth.GenerateCodeVerifier()
    if err != nil {
        return fmt.Errorf("failed to generate code verifier: %w", err)
    }
    codeChallenge := auth.GenerateCodeChallenge(codeVerifier)

    listener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        return fmt.Errorf("failed to start local server: %w", err)
    }
    port := listener.Addr().(*net.TCPAddr).Port
    redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

    type callbackResult struct {
        accessToken  string
        refreshToken string
        username     string
        email        string
        role         string
    }
    resultCh := make(chan callbackResult, 1)

    mux := http.NewServeMux()
    server := &http.Server{Handler: mux}

    mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
        // Tokens come back as query params — no second exchange needed
        fmt.Fprintf(w, `<html><body>
            <h2>Authentication successful!</h2>
            <p>You can close this window and return to the terminal.</p>
        </body></html>`)

        resultCh <- callbackResult{
            accessToken:  r.URL.Query().Get("access_token"),
            refreshToken: r.URL.Query().Get("refresh_token"),
            username:     r.URL.Query().Get("username"),
            email:        r.URL.Query().Get("email"),
            role:         r.URL.Query().Get("role"),
        }
    })

    go server.Serve(listener)
    defer server.Shutdown(context.Background())

    // Send code_verifier so backend can do the exchange in the callback
    params := url.Values{}
    params.Set("state", state)
    params.Set("code_challenge", codeChallenge)
    params.Set("code_verifier", codeVerifier)
    params.Set("redirect_uri", redirectURI)

    loginURL := fmt.Sprintf("%s/auth/github/login?%s", BaseURL, params.Encode())

    fmt.Println("Opening GitHub login in your browser...")
    openBrowser(loginURL)

    s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
    s.Suffix = " Waiting for authentication..."
    s.Start()

    var result callbackResult
    select {
    case result = <-resultCh:
        s.Stop()
    case <-time.After(5 * time.Minute):
        s.Stop()
        return fmt.Errorf("login timed out — please try again")
    }

    if result.accessToken == "" {
        return fmt.Errorf("no token received — authentication may have failed")
    }

    if err := auth.SaveCredentials(&auth.Credentials{
        AccessToken:  result.accessToken,
        RefreshToken: result.refreshToken,
        Username:     result.username,
        Email:        result.email,
        Role:         result.role,
        BaseURL:      BaseURL,
    }); err != nil {
        return fmt.Errorf("failed to save credentials: %w", err)
    }

    fmt.Printf("\n✓ Logged in as @%s\n", result.username)
    return nil
}

func mustMarshal(v any) []byte {
    b, _ := json.Marshal(v)
    return b
}

func openBrowser(url string) error {
    var cmdName string
    var args []string

    switch runtime.GOOS {
    case "darwin":
        cmdName = "open"
        args = []string{url}
    case "windows":
        cmdName = "rundll32"
        args = []string{"url.dll,FileProtocolHandler", url}
    default:
        cmdName = "xdg-open"
        args = []string{url}
    }

    return exec.Command(cmdName, args...).Start()
}