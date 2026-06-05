package main_test

import (
	"encoding/json"
	"portwatch/internal/health"
	"portwatch/internal/ports"
	"testing"
)

func TestJSONOutputNoFindings(t *testing.T) {
	cleanEntries := []ports.PortEntry{
		{Protocol: "TCP", LocalAddr: "127.0.0.1", Port: 63000, PID: 4, ProcessName: "svchost.exe", State: "LISTEN"},
	}
	findings := health.Scan(cleanEntries)

	type jsonOut struct {
		Ports    []ports.PortEntry `json:"ports"`
		Findings []health.Finding  `json:"findings"`
	}
	out := jsonOut{Ports: cleanEntries, Findings: findings}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("expected non-empty JSON output")
	}

	if findings == nil {
		t.Fatal("findings should be empty slice, not nil")
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}
