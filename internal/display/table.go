package display

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"portwatch/internal/ports"
)

// RenderPortTable prints the port list to stdout.
func RenderPortTable(entries []ports.PortEntry) {
	if len(entries) == 0 {
		ColorMuted.Println("No listening ports found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PROTO", "ADDR", "PORT", "PID", "PROCESS", "STATE"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)
	table.SetColumnSeparator("  ")
	table.SetCenterSeparator("")
	table.SetRowSeparator("")

	for _, e := range entries {
		procName := e.ProcessName
		if isDevProcessName(procName) {
			procName = ColorDev.Sprint(procName)
		}
		addr := e.LocalAddr
		if addr == "0.0.0.0" || addr == "::" {
			addr = ColorWarn.Sprint(addr)
		}
		table.Append([]string{
			e.Protocol,
			addr,
			fmt.Sprintf("%d", e.Port),
			fmt.Sprintf("%d", e.PID),
			procName,
			e.State,
		})
	}
	table.Render()
}

// RenderDetailBlock prints the inspect view for a single process.
func RenderDetailBlock(w io.Writer, d *ports.ProcessDetail, hint string) {
	if w == nil {
		w = os.Stdout
	}
	if d == nil {
		return
	}
	fmt.Fprintln(w)
	ColorHeader.Fprintf(w, "  Process: %s (PID %d)\n", d.Name, d.PID)
	fmt.Fprintf(w, "  Path:     %s\n", d.ExePath)
	fmt.Fprintf(w, "  Memory:   %.1f MB\n", d.MemoryMB)
	if d.StartTime != "" {
		fmt.Fprintf(w, "  Started:  %s\n", d.StartTime)
	}
	if d.Cmdline != "" {
		cmdPreview := d.Cmdline
		if len(cmdPreview) > 80 {
			cmdPreview = cmdPreview[:77] + "..."
		}
		fmt.Fprintf(w, "  Command:  %s\n", cmdPreview)
	}
	if hint != "" {
		ColorDev.Fprintf(w, "  Likely:   %s\n", hint)
	}
	if len(d.Ports) > 0 {
		portStrs := make([]string, len(d.Ports))
		for i, p := range d.Ports {
			portStrs[i] = fmt.Sprintf("http://localhost:%d", p)
		}
		ColorURL.Fprintf(w, "  URLs:     %s\n", strings.Join(portStrs, "  "))
	}
	fmt.Fprintln(w)
}

// isDevProcessName returns true if the process name matches a known dev runtime.
func isDevProcessName(name string) bool {
	name = strings.ToLower(name)
	devNames := []string{"node.exe", "python.exe", "python3.exe", "bun.exe", "deno.exe", "java.exe", "docker.exe", "dockerd.exe", "go.exe"}
	for _, d := range devNames {
		if name == d {
			return true
		}
	}
	return false
}
