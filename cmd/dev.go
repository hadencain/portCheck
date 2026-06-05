package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"portwatch/internal/display"
	"portwatch/internal/ports"
	"portwatch/internal/processes"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Show running development server processes",
	RunE:  runDev,
}

func init() {
	rootCmd.AddCommand(devCmd)
}

func runDev(_ *cobra.Command, _ []string) error {
	entries, err := ports.ListeningPorts()
	if err != nil {
		return fmt.Errorf("listing ports: %w", err)
	}

	// Index listening ports by PID for fast lookup
	listeningByPID := make(map[int32][]uint32)
	for _, e := range entries {
		listeningByPID[e.PID] = append(listeningByPID[e.PID], e.Port)
	}

	// Filter to only dev entries
	var devEntries []ports.PortEntry
	seen := make(map[int32]bool)
	for _, e := range entries {
		if processes.IsDevProcess(e.ProcessName) && !seen[e.PID] {
			seen[e.PID] = true
			devEntries = append(devEntries, e)
		}
	}

	sort.Slice(devEntries, func(i, j int) bool {
		return devEntries[i].Port < devEntries[j].Port
	})

	if len(devEntries) == 0 {
		fmt.Println()
		display.ColorMuted.Print("  No development servers detected.")
		fmt.Println()
		return nil
	}

	color.New(color.FgCyan, color.Bold).Printf("\n  Dev servers — %d found\n\n", len(devEntries))

	for _, e := range devEntries {
		detail, err := ports.GetPortDetail(e.Port)
		hint := ""
		if err == nil {
			hint = processes.FrameworkHint(detail.Name, detail.Cmdline, e.Port)
		}

		listening := listeningByPID[e.PID]
		portStrs := make([]string, len(listening))
		for i, p := range listening {
			portStrs[i] = fmt.Sprintf("%d", p)
		}

		display.ColorDev.Printf("  ● %s", e.ProcessName)
		if hint != "" {
			fmt.Printf(" (%s)", hint)
		}
		fmt.Printf("  PID %d  ports: %s\n", e.PID, strings.Join(portStrs, ", "))

		for _, p := range listening {
			display.ColorURL.Printf("    http://localhost:%d\n", p)
		}
		fmt.Println()
	}
	return nil
}
