package display

import "github.com/fatih/color"

var (
	ColorHeader  = color.New(color.FgCyan, color.Bold)
	ColorDev     = color.New(color.FgGreen, color.Bold)
	ColorWarn    = color.New(color.FgYellow, color.Bold)
	ColorDanger  = color.New(color.FgRed, color.Bold)
	ColorMuted   = color.New(color.FgWhite)
	ColorURL     = color.New(color.FgCyan)
	ColorSuccess = color.New(color.FgGreen)
)

// Highlight wraps s in color if condition is true, else returns s plain.
func Highlight(c *color.Color, s string, condition bool) string {
	if condition {
		return c.Sprint(s)
	}
	return s
}
