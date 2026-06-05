package ports

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// ListeningPorts returns all TCP/UDP connections in LISTEN state.
func ListeningPorts() ([]PortEntry, error) {
	conns, err := net.Connections("inet")
	if err != nil {
		return nil, fmt.Errorf("fetching connections: %w", err)
	}

	var entries []PortEntry
	for _, c := range conns {
		if c.Status != "LISTEN" && c.Status != "" {
			// UDP has no status; include those too
			if c.Type != 2 { // 2 = SOCK_DGRAM (UDP)
				continue
			}
		}
		name := processName(c.Pid)
		entries = append(entries, PortEntry{
			Protocol:    protoName(c.Type),
			LocalAddr:   c.Laddr.IP,
			Port:        c.Laddr.Port,
			PID:         c.Pid,
			ProcessName: name,
			State:       c.Status,
		})
	}
	return entries, nil
}

// GetPortDetail returns the ProcessDetail for the process listening on port.
// Returns an error if no process is found on that port.
func GetPortDetail(port uint32) (*ProcessDetail, error) {
	entries, err := ListeningPorts()
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.Port == port {
			return processDetail(e.PID, entries)
		}
	}
	return nil, fmt.Errorf("no process found listening on port %d", port)
}

func processName(pid int32) string {
	if pid == 0 {
		return "-"
	}
	p, err := process.NewProcess(pid)
	if err != nil {
		return "?"
	}
	name, err := p.Name()
	if err != nil {
		return "?"
	}
	return name
}

func processDetail(pid int32, all []PortEntry) (*ProcessDetail, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("opening process %d: %w", pid, err)
	}

	name, _ := p.Name()
	exe, _ := p.Exe()
	cmdline, _ := p.Cmdline()
	mem, _ := p.MemoryInfo()
	createTime, _ := p.CreateTime()

	var memMB float64
	if mem != nil {
		memMB = float64(mem.RSS) / 1024 / 1024
	}

	startStr := ""
	if createTime > 0 {
		t := time.Unix(createTime/1000, 0)
		startStr = t.Format("2006-01-02 15:04:05")
	}

	var ownedPorts []uint32
	for _, e := range all {
		if e.PID == pid {
			ownedPorts = append(ownedPorts, e.Port)
		}
	}

	return &ProcessDetail{
		PID:       pid,
		Name:      name,
		ExePath:   exe,
		MemoryMB:  memMB,
		StartTime: startStr,
		Cmdline:   cmdline,
		Ports:     ownedPorts,
	}, nil
}

func protoName(sockType uint32) string {
	switch sockType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "?"
	}
}
