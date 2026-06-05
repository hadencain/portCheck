package processes

import "strings"

var devProcessNames = map[string]bool{
	"node.exe":    true,
	"python.exe":  true,
	"python3.exe": true,
	"bun.exe":     true,
	"deno.exe":    true,
	"java.exe":    true,
	"docker.exe":  true,
	"dockerd.exe": true,
	"go.exe":      true,
}

// IsDevProcess returns true if the process name is a known dev runtime.
func IsDevProcess(name string) bool {
	return devProcessNames[strings.ToLower(name)]
}

// FrameworkHint returns a human-readable guess at the framework in use.
// Returns empty string if no match found.
func FrameworkHint(processName, cmdline string, port uint32) string {
	name := strings.ToLower(processName)
	cmd := strings.ToLower(cmdline)

	switch name {
	case "node.exe":
		switch {
		case strings.Contains(cmd, "next"):
			return "Next.js"
		case strings.Contains(cmd, "vite"):
			return "Vite"
		case strings.Contains(cmd, "nuxt"):
			return "Nuxt.js"
		case strings.Contains(cmd, "gatsby"):
			return "Gatsby"
		case strings.Contains(cmd, "remix"):
			return "Remix"
		case strings.Contains(cmd, "angular"):
			return "Angular"
		case strings.Contains(cmd, "react-scripts"):
			return "Create React App"
		default:
			return "Node.js (Express/custom)"
		}

	case "python.exe", "python3.exe":
		switch {
		case strings.Contains(cmd, "uvicorn") || strings.Contains(cmd, "fastapi"):
			return "FastAPI"
		case strings.Contains(cmd, "manage.py"):
			return "Django"
		case strings.Contains(cmd, "flask"):
			return "Flask"
		case strings.Contains(cmd, "streamlit"):
			return "Streamlit"
		case strings.Contains(cmd, "gradio"):
			return "Gradio"
		default:
			return "Python (custom)"
		}

	case "bun.exe":
		return "Bun"

	case "deno.exe":
		return "Deno"

	case "java.exe":
		switch {
		case strings.Contains(cmd, "spring"):
			return "Spring Boot"
		default:
			return "Java (Spring Boot / custom)"
		}

	case "docker.exe", "dockerd.exe":
		return "Docker"
	}

	return ""
}
