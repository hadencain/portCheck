package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"portwatch/internal/display"
	"portwatch/internal/ports"
	"portwatch/internal/processes"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <port>",
	Short: "Show details for the process on a port",
	Args:  cobra.ExactArgs(1),
	RunE:  runInspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

func runInspect(_ *cobra.Command, args []string) error {
	portNum, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid port %q: must be a number", args[0])
	}

	detail, err := ports.GetPortDetail(uint32(portNum))
	if err != nil {
		display.ColorDanger.Fprintf(os.Stderr, "\n  %v\n\n", err)
		return nil
	}

	hint := processes.FrameworkHint(detail.Name, detail.Cmdline, uint32(portNum))
	display.RenderDetailBlock(os.Stdout, detail, hint)
	return nil
}
