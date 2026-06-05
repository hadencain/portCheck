package health

import (
	"testing"
	"portwatch/internal/ports"
)

func TestDetectAllInterfaceListeners(t *testing.T) {
	entries := []ports.PortEntry{
		{Protocol: "TCP", LocalAddr: "0.0.0.0", Port: 3000, PID: 1, ProcessName: "node.exe", State: "LISTEN"},
		{Protocol: "TCP", LocalAddr: "127.0.0.1", Port: 5173, PID: 2, ProcessName: "node.exe", State: "LISTEN"},
	}
	findings := Scan(entries)
	found := false
	for _, f := range findings {
		if f.Kind == KindAllInterfaces && f.Port == 3000 {
			found = true
		}
	}
	if !found {
		t.Error("expected AllInterfaces finding for 0.0.0.0:3000")
	}
}

func TestDetectCommonDevPort(t *testing.T) {
	entries := []ports.PortEntry{
		{Protocol: "TCP", LocalAddr: "127.0.0.1", Port: 3000, PID: 1, ProcessName: "node.exe", State: "LISTEN"},
	}
	findings := Scan(entries)
	found := false
	for _, f := range findings {
		if f.Kind == KindCommonDevPort && f.Port == 3000 {
			found = true
		}
	}
	if !found {
		t.Error("expected CommonDevPort finding for port 3000")
	}
}

func TestNoFindingsForCleanState(t *testing.T) {
	entries := []ports.PortEntry{
		{Protocol: "TCP", LocalAddr: "127.0.0.1", Port: 63000, PID: 4, ProcessName: "svchost.exe", State: "LISTEN"},
	}
	findings := Scan(entries)
	for _, f := range findings {
		t.Errorf("unexpected finding: %+v", f)
	}
}
