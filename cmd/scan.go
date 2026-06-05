package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"portwatch/internal/display"
	"portwatch/internal/health"
	"portwatch/internal/ports"
)

var scanJSON bool

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Health scan — detect orphaned servers and suspicious listeners",
	RunE:  runScan,
}

func init() {
	scanCmd.Flags().BoolVar(&scanJSON, "json", false, "Output results as JSON")
	rootCmd.AddCommand(scanCmd)
}

func runScan(_ *cobra.Command, _ []string) error {
	entries, err := ports.ListeningPorts()
	if err != nil {
		return fmt.Errorf("listing ports: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Port < entries[j].Port
	})

	findings := health.Scan(entries)

	if scanJSON {
		type jsonOut struct {
			Ports    []ports.PortEntry `json:"ports"`
			Findings []health.Finding  `json:"findings"`
		}
		out := jsonOut{Ports: entries, Findings: findings}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	display.ColorHeader.Printf("\n  PortWatch Scan — %d port(s) found\n\n", len(entries))
	display.RenderPortTable(entries)

	if len(findings) == 0 {
		fmt.Println()
		display.ColorSuccess.Print("  No issues detected.")
		fmt.Println()
		return nil
	}

	display.ColorWarn.Printf("\n  Findings (%d)\n\n", len(findings))
	for _, f := range findings {
		icon := "⚠"
		c := display.ColorWarn
		if f.Kind == health.KindAllInterfaces {
			icon = "⛔"
			c = display.ColorDanger
		}
		c.Printf("  %s  [%s] port %d  %s\n", icon, f.Process, f.Port, f.Message)
	}
	fmt.Println()
	return nil
}
