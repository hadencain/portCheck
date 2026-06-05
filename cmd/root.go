package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"portwatch/internal/display"
	"portwatch/internal/ports"
)

var rootCmd = &cobra.Command{
	Use:   "portwatch",
	Short: "See what's listening on your localhost",
	Long: `PortWatch — localhost service manager for developers.

Lists active ports, identifies dev servers, and lets you kill
processes you forgot to stop.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runList(_ *cobra.Command, _ []string) error {
	entries, err := ports.ListeningPorts()
	if err != nil {
		return fmt.Errorf("listing ports: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Port < entries[j].Port
	})

	color.New(color.FgCyan, color.Bold).Printf("\n  PortWatch — %d listening port(s)\n\n", len(entries))
	display.RenderPortTable(entries)
	return nil
}
