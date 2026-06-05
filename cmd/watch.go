package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"portwatch/internal/ports"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Continuously monitor ports — highlights changes",
	RunE:  runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func runWatch(_ *cobra.Command, _ []string) error {
	fmt.Println()
	color.New(color.FgCyan, color.Bold).Print("  PortWatch — watch mode (Ctrl+C to exit)")
	fmt.Println()

	prev := make(map[string]bool)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		entries, err := ports.ListeningPorts()
		if err != nil {
			return err
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Port < entries[j].Port
		})

		curr := make(map[string]bool)
		for _, e := range entries {
			key := fmt.Sprintf("%s:%d:%d", e.Protocol, e.Port, e.PID)
			curr[key] = true
		}

		// Detect opened
		for key := range curr {
			if !prev[key] && len(prev) > 0 {
				color.New(color.FgGreen, color.Bold).Printf("  + OPENED  %s\n", key)
			}
		}
		// Detect closed
		for key := range prev {
			if !curr[key] {
				color.New(color.FgRed, color.Bold).Printf("  - CLOSED  %s\n", key)
			}
		}

		prev = curr
		<-ticker.C
	}
}
