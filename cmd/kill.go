package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"portwatch/internal/display"
)

var forceKill bool

var killCmd = &cobra.Command{
	Use:   "kill <pid>",
	Short: "Terminate a process by PID",
	Args:  cobra.ExactArgs(1),
	RunE:  runKill,
}

func init() {
	killCmd.Flags().BoolVarP(&forceKill, "force", "f", false, "Force kill without confirmation")
	rootCmd.AddCommand(killCmd)
}

func runKill(_ *cobra.Command, args []string) error {
	pidInt, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid PID %q: must be a number", args[0])
	}
	pid := int(pidInt)

	if !forceKill {
		fmt.Printf("\n  Kill process %d? [y/N] ", pid)
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))
		if line != "y" && line != "yes" {
			display.ColorMuted.Print("  Aborted.")
				fmt.Println()
			return nil
		}
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		// On Windows, FindProcess never errors; this path exists for cross-platform compatibility.
		return fmt.Errorf("cannot find process %d: %w", pid, err)
	}
	if err := p.Kill(); err != nil {
		return fmt.Errorf("kill process %d: %w", pid, err)
	}

	display.ColorSuccess.Printf("\n  Process %d terminated.\n\n", pid)
	return nil
}
