package processes

import "testing"

func TestIsDevProcess(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"node.exe", true},
		{"python.exe", true},
		{"bun.exe", true},
		{"deno.exe", true},
		{"java.exe", true},
		{"docker.exe", true},
		{"dockerd.exe", true},
		{"chrome.exe", false},
		{"explorer.exe", false},
		{"", false},
	}
	for _, tt := range tests {
		got := IsDevProcess(tt.name)
		if got != tt.want {
			t.Errorf("IsDevProcess(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFrameworkHint(t *testing.T) {
	tests := []struct {
		processName string
		cmdline     string
		port        uint32
		want        string
	}{
		{"node.exe", "node next dev", 3000, "Next.js"},
		{"node.exe", "node vite", 5173, "Vite"},
		{"node.exe", "node server.js", 3000, "Node.js (Express/custom)"},
		{"python.exe", "uvicorn main:app", 8000, "FastAPI"},
		{"python.exe", "manage.py runserver", 8000, "Django"},
		{"python.exe", "flask run", 5000, "Flask"},
		{"bun.exe", "bun run dev", 3000, "Bun"},
		{"deno.exe", "deno run server.ts", 8080, "Deno"},
		{"java.exe", "spring-boot", 8080, "Spring Boot"},
		{"java.exe", "-jar app.jar", 8080, "Java (Spring Boot / custom)"},
		{"chrome.exe", "", 0, ""},
	}
	for _, tt := range tests {
		got := FrameworkHint(tt.processName, tt.cmdline, tt.port)
		if got != tt.want {
			t.Errorf("FrameworkHint(%q, %q, %d) = %q, want %q",
				tt.processName, tt.cmdline, tt.port, got, tt.want)
		}
	}
}
