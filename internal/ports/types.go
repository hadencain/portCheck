package ports

// PortEntry is one listening port with its owning process.
type PortEntry struct {
	Protocol    string // "TCP" or "UDP"
	LocalAddr   string // "127.0.0.1" or "0.0.0.0" or "::"
	Port        uint32
	PID         int32
	ProcessName string
	State       string // "LISTEN", "ESTABLISHED", etc.
}

// ProcessDetail is the expanded view for `inspect`.
type ProcessDetail struct {
	PID         int32
	Name        string
	ExePath     string
	MemoryMB    float64
	StartTime   string // formatted, empty if unavailable
	Cmdline     string
	Ports       []uint32
}
