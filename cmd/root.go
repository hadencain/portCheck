package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	fmt.Println("PortWatch — listing ports (not yet implemented)")
	return nil
}
