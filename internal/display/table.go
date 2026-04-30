// internal/display/table.go
package display

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Table renders a nice formatted table
func Table(headers []string, rows [][]string) {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetStyle(table.StyleRounded)

    // Add header
    headerRow := make([]interface{}, len(headers))
    for i, h := range headers {
        headerRow[i] = h
    }
    t.AppendHeader(headerRow)

    // Add rows
    for _, row := range rows {
        dataRow := make([]interface{}, len(row))
        for i, cell := range row {
            dataRow[i] = cell
        }
        t.AppendRow(dataRow)
    }

    t.Render()
    fmt.Println() // extra newline
}

// Empty shows message when no data is found
func Empty(message string) {
    fmt.Printf("ℹ️  %s\n", message)
}

// Truncate shortens long strings
func Truncate(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-3] + "..."
}

// OrDash returns dash if string is empty
func OrDash(s string) string {
    if s == "" {
        return "-"
    }
    return s
}

// FormatProbability formats float as percentage
func FormatProbability(p float64) string {
    if p == 0 {
        return "-"
    }
    return fmt.Sprintf("%.1f%%", p*100)
}

func Success(message string) {
    fmt.Printf("✅ %s\n", message)
}

// Info prints an informational message
func Info(message string) {
    fmt.Printf("ℹ️  %s\n", message)
}

// Warning prints a warning message
func Warning(message string) {
    fmt.Printf("⚠️  %s\n", message)
}

// Error prints an error message (but note: errors are usually returned, not printed here)
func Error(message string) {
    fmt.Printf("❌ %s\n", message)
}
