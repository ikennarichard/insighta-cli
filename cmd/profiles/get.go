// cmd/profiles/get.go
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

var getCmd = &cobra.Command{
	Use:     "get <id>",
	Short:   "Get a profile by ID",
	Args:    cobra.ExactArgs(1), // enforces exactly one argument
	Example: "  insighta-cli profiles get abc-123-def",
	RunE:    runGet,
}

func runGet(cmd *cobra.Command, args []string) error {
	id := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Fetching profile..."
	s.Start()

	resp, err := client.Get("/api/v1/profiles/" + id)
	s.Stop()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("profile not found: %s", id)
	}

	var result struct {
		Data Profile `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	p := result.Data

	countryDisplay := display.OrDash(p.CountryName)
if p.CountryID != "" && p.CountryName != "" {
    countryDisplay = fmt.Sprintf("%s (%s)", p.CountryName, p.CountryID)
}


	// Display as key-value pairs
	rows := [][]string{
		{"ID", p.ID},
		{"Name", p.Name},
		{"Gender", display.OrDash(p.Gender)},
		{"Gender Probability", display.FormatProbability(p.GenderProbability)},
		{"Age", fmt.Sprintf("%d", p.Age)},
		{"Age Group", display.OrDash(p.AgeGroup)},
		{"Country", countryDisplay},
		{"Country Code", display.OrDash(p.CountryID)},
		{"Country Probability", display.FormatProbability(p.CountryProbability)},
		{"Created At", p.CreatedAt},
	}

	display.Table([]string{"FIELD", "VALUE"}, rows)
	return nil
}