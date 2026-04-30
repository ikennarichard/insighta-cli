// cmd/profiles/search.go
package profiles

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ikennarichard/insighta-cli/internal/api"
	"github.com/ikennarichard/insighta-cli/internal/display"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	Short:   "Search profiles using natural language",
	Args:    cobra.ExactArgs(1),
	Example: `  insighta profiles search "young males from nigeria"`,
	RunE:    runSearch,
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("q", query)

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Searching..."
	s.Start()

	resp, err := client.Get("/api/v1/profiles/search?" + params.Encode())
	s.Stop()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result PaginatedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Data) == 0 {
		display.Empty("No profiles matched your query")
		return nil
	}

	headers := []string{"NAME", "GENDER", "PROB", "AGE", "COUNTRY"}
	rows := make([][]string, len(result.Data))
	for i, p := range result.Data {
		rows[i] = []string{
			display.Truncate(p.Name, 20),
			display.OrDash(p.Gender),
			display.FormatProbability(p.GenderProbability),
			fmt.Sprintf("%d", p.Age),
			display.OrDash(p.CountryName),
		}
	}
	display.Table(headers, rows)

	fmt.Printf("\n%d result(s) for: %q\n", result.Total, query)
	return nil
}