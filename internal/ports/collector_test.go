package ports

import "testing"

func TestProtoName(t *testing.T) {
	tests := []struct {
		input uint32
		want  string
	}{
		{1, "TCP"},
		{2, "UDP"},
		{99, "?"},
	}
	for _, tt := range tests {
		got := protoName(tt.input)
		if got != tt.want {
			t.Errorf("protoName(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestProcessNameZeroPID(t *testing.T) {
	// PID 0 is always invalid — should return "-"
	got := processName(0)
	if got != "-" {
		t.Errorf("processName(0) = %q, want \"-\"", got)
	}
}

func TestListeningPortsReturnsSlice(t *testing.T) {
	// Integration test: just verifies the call succeeds and returns a slice.
	entries, err := ListeningPorts()
	if err != nil {
		t.Fatalf("ListeningPorts() error: %v", err)
	}
	if entries == nil {
		t.Error("expected non-nil slice")
	}
}
