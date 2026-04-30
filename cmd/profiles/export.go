// cmd/profiles/export.go
package profiles

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ikennarichard/insighta-cli/internal/api"
	"github.com/ikennarichard/insighta-cli/internal/display"
	"github.com/spf13/cobra"
)

var (
	flagFormat        string
	flagExportGender  string
	flagExportCountry string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export profiles to a file",
	Example: `  insighta profiles export --format csv
  insighta profiles export --format csv --gender male --country NG`,
	RunE: runExport,
}

func init() {
	exportCmd.Flags().StringVar(&flagFormat, "format", "csv", "Export format (csv)")
	exportCmd.Flags().StringVar(&flagExportGender, "gender", "", "Filter by gender")
	exportCmd.Flags().StringVar(&flagExportCountry, "country", "", "Filter by country code")
}

func runExport(cmd *cobra.Command, args []string) error {
	if flagFormat != "csv" {
		return fmt.Errorf("unsupported format %q — only csv is supported", flagFormat)
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("format", flagFormat)
	if flagExportGender != "" { params.Set("gender", flagExportGender) }
	if flagExportCountry != "" { params.Set("country_id", flagExportCountry) }

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Exporting profiles..."
	s.Start()

	resp, err := client.Get("/api/v1/profiles/export?" + params.Encode())
	s.Stop()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("export failed with status %d", resp.StatusCode)
	}

	// Save to current working directory
	filename := fmt.Sprintf("profiles_%s.csv", time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	display.Success(fmt.Sprintf("Exported to %s", filename))
	return nil
}