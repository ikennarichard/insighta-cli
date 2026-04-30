// cmd/profiles/create.go
package profiles

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ikennarichard/insighta-cli/internal/api"
	"github.com/ikennarichard/insighta-cli/internal/display"
	"github.com/spf13/cobra"
)

var flagName string

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new profile",
	Example: `  insighta profiles create --name "Harriet Tubman"`,
	RunE:    runCreate,
}

func init() {
	createCmd.Flags().StringVar(&flagName, "name", "", "Name to classify (required)")
	createCmd.MarkFlagRequired("name")
}

func runCreate(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Creating profile..."
	s.Start()

	resp, err := client.Post("/api/v1/profiles", map[string]string{
		"name": flagName,
	})
	s.Stop()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Status  string  `json:"status"`
		Message string  `json:"message"`
		Data    Profile `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Message != "" {
		// Profile already existed
		fmt.Printf("ℹ %s\n\n", result.Message)
	} else {
		display.Success(fmt.Sprintf("Profile created for %q", flagName))
		fmt.Println()
	}

	p := result.Data
	display.Table([]string{"FIELD", "VALUE"}, [][]string{
		{"Name", p.Name},
		{"Gender", display.OrDash(p.Gender)},
		{"Probability", display.FormatProbability(p.GenderProbability)},
		{"Age", fmt.Sprintf("%d", p.Age)},
		{"Country", display.OrDash(p.CountryName)},
	})

	return nil
}