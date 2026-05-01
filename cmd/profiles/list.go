// cmd/profiles/list.go
package profiles

import (
    "fmt"
    "net/url"

    "github.com/ikennarichard/insighta-cli/internal/api"
    "github.com/ikennarichard/insighta-cli/internal/display"
    "github.com/spf13/cobra"
)

type Profile struct {
    ID                 string  `json:"id"`
    Name               string  `json:"name"`
    Gender             string  `json:"gender"`
    GenderProbability  float64 `json:"gender_probability"`
    Age                int     `json:"age"`
    AgeGroup           string  `json:"age_group"`
    CountryID          string  `json:"country_id"`
    CountryName        string  `json:"country_name"`
    CountryProbability float64 `json:"country_probability"`
    CreatedAt          string  `json:"created_at"`
}

type PaginatedResponse struct {
    Data  []Profile `json:"data"`
    Total int       `json:"total"`
    Page  int       `json:"page"`
    Limit int       `json:"limit"`
}

var (
    flagGender   string
    flagCountry  string
    flagAgeGroup string
    flagMinAge   int
    flagMaxAge   int
    flagSortBy   string
    flagOrder    string
    flagPage     int
    flagLimit    int
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List profiles with optional filters",
    Example: `  insighta-cli profiles list
  insighta-cli profiles list --gender male
  insighta-cli profiles list --country NG --age-group adult
  insighta-cli profiles list --min-age 25 --max-age 40
  insighta-cli profiles list --sort-by age --order desc
  insighta-cli profiles list --page 2 --limit 20`,
    RunE: runList,
}

func init() {
    listCmd.Flags().StringVar(&flagGender, "gender", "", "Filter by gender (male/female)")
    listCmd.Flags().StringVar(&flagCountry, "country", "", "Filter by country code e.g. NG")
    listCmd.Flags().StringVar(&flagAgeGroup, "age-group", "", "Filter by age group")
    listCmd.Flags().IntVar(&flagMinAge, "min-age", 0, "Minimum age")
    listCmd.Flags().IntVar(&flagMaxAge, "max-age", 0, "Maximum age")
    listCmd.Flags().StringVar(&flagSortBy, "sort-by", "", "Sort by (age, name, created_at)")
    listCmd.Flags().StringVar(&flagOrder, "order", "asc", "Sort order: asc or desc")
    listCmd.Flags().IntVar(&flagPage, "page", 1, "Page number")
    listCmd.Flags().IntVar(&flagLimit, "limit", 20, "Results per page (max 50)")
}

func runList(cmd *cobra.Command, args []string) error {
    client, err := api.NewClient()
    if err != nil {
        return err
    }

    // Build query parameters
    params := url.Values{}
    params.Set("page", fmt.Sprintf("%d", flagPage))
    params.Set("limit", fmt.Sprintf("%d", flagLimit))
    if flagGender != "" {
        params.Set("gender", flagGender)
    }
    if flagCountry != "" {
        params.Set("country_id", flagCountry)
    }
    if flagAgeGroup != "" {
        params.Set("age_group", flagAgeGroup)
    }
    if flagMinAge > 0 {
        params.Set("min_age", fmt.Sprintf("%d", flagMinAge))
    }
    if flagMaxAge > 0 {
        params.Set("max_age", fmt.Sprintf("%d", flagMaxAge))
    }
    if flagSortBy != "" {
        params.Set("sort_by", flagSortBy)
    }
    if flagOrder != "" {
        params.Set("order", flagOrder)
    }

    display.ShowLoader("Fetching profiles...")

    var result PaginatedResponse
    err = client.GetJSON("/api/v1/profiles?"+params.Encode(), &result)
    display.StopLoader()

    if err != nil {
        return err
    }

    if len(result.Data) == 0 {
        display.Empty("No profiles found matching your criteria.")
        return nil
    }

    // Prepare table data
    headers := []string{"ID", "NAME", "GENDER", "PROB", "AGE", "GROUP", "COUNTRY"}
    rows := make([][]string, len(result.Data))

    for i, p := range result.Data {
        rows[i] = []string{
            display.Truncate(p.ID, 8),
            display.Truncate(p.Name, 25),
            display.OrDash(p.Gender),
            display.FormatProbability(p.GenderProbability),
            fmt.Sprintf("%d", p.Age),
            display.OrDash(p.AgeGroup),
            display.OrDash(p.CountryName),
        }
    }

    display.Table(headers, rows)

    fmt.Printf("\nShowing %d of %d profiles (Page %d)\n", len(result.Data), result.Total, result.Page)
    return nil
}