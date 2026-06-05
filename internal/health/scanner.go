package health

import (
	"portwatch/internal/ports"
	"portwatch/internal/processes"
)

type FindingKind string

const (
	KindAllInterfaces  FindingKind = "all_interfaces"
	KindCommonDevPort  FindingKind = "common_dev_port"
	KindMultipleOnPort FindingKind = "multiple_on_port"
)

type Finding struct {
	Kind    FindingKind
	Port    uint32
	PID     int32
	Process string
	Message string
}

var commonDevPorts = map[uint32]string{
	3000: "React/Next.js/Express",
	3001: "React (alt)",
	4200: "Angular",
	5000: "Flask/generic",
	5173: "Vite",
	5174: "Vite (alt)",
	6006: "Storybook",
	8000: "Django/FastAPI",
	8080: "Spring Boot/generic",
	8081: "generic alt",
	8888: "Jupyter",
	9000: "PHP-FPM/generic",
}

// Scan analyses a snapshot of listening ports and returns findings.
func Scan(entries []ports.PortEntry) []Finding {
	var findings []Finding
	portCount := make(map[uint32]int)

	for _, e := range entries {
		portCount[e.Port]++

		// Bound to all interfaces
		if e.LocalAddr == "0.0.0.0" || e.LocalAddr == "::" {
			findings = append(findings, Finding{
				Kind:    KindAllInterfaces,
				Port:    e.Port,
				PID:     e.PID,
				Process: e.ProcessName,
				Message: "Listening on all interfaces — accessible from network, not just localhost",
			})
		}

		// Common dev port occupied
		if label, ok := commonDevPorts[e.Port]; ok && processes.IsDevProcess(e.ProcessName) {
			findings = append(findings, Finding{
				Kind:    KindCommonDevPort,
				Port:    e.Port,
				PID:     e.PID,
				Process: e.ProcessName,
				Message: "Common dev port in use: " + label,
			})
		}
	}

	// Multiple processes on same port (shouldn't happen but flag it)
	for port, count := range portCount {
		if count > 1 {
			findings = append(findings, Finding{
				Kind:    KindMultipleOnPort,
				Port:    port,
				Message: "Multiple listeners on same port — possible conflict",
			})
		}
	}

	return findings
}
