# portWatch

Windows CLI tool for monitoring and managing localhost dev server processes by port.

## Build

```bash
go build -o portwatch.exe .
```

Requires Go 1.22+.

## Commands

### List all listening ports (default)

```bash
portwatch
```

Shows a table of all TCP/UDP listening ports — protocol, address, port number, PID, process name, and state. Dev runtimes (node, python, bun, etc.) are highlighted in green. Addresses bound to all interfaces (0.0.0.0, ::) are highlighted in yellow.

---

### Inspect a port

```bash
portwatch inspect <port>
```

Shows detailed info for the process listening on the given port: executable path, memory usage, start time, command line, framework guess, and localhost URLs.

```
portwatch inspect 3000
portwatch inspect 8080
```

---

### Show dev servers

```bash
portwatch dev
```

Filters to only dev runtime processes (node, python, bun, deno, java, docker) and shows their ports and `http://localhost:<port>` URLs. Includes a framework hint (Next.js, Vite, FastAPI, Django, Flask, etc.) where detectable from the command line.

---

### Kill a process

```bash
portwatch kill <pid>
portwatch kill <pid> --force
```

Terminates a process by PID. Prompts for confirmation unless `--force` / `-f` is passed.

```bash
portwatch kill 12345          # prompts [y/N]
portwatch kill 12345 --force  # no prompt
```

---

### Health scan

```bash
portwatch scan
portwatch scan --json
```

Scans all listening ports for potential issues:

- **All-interface listeners** — processes bound to `0.0.0.0` or `::` (accessible from the network, not just localhost)
- **Common dev ports in use** — well-known ports (3000, 5173, 8000, 8080, etc.) occupied by dev runtimes
- **Multiple listeners on the same port** — possible port conflict

`--json` outputs a machine-readable JSON object with `ports` and `findings` arrays.

---

### Watch mode

```bash
portwatch watch
```

Polls every 2 seconds and prints changes as they happen:

```
+ OPENED  TCP:3000:18432
- CLOSED  TCP:5173:9120
```

Press `Ctrl+C` to exit.

## Framework Detection

`inspect` and `dev` infer the likely framework from the process name and command line:

| Runtime | Detected frameworks |
|---------|-------------------|
| node.exe | Next.js, Vite, Nuxt.js, Gatsby, Remix, Angular, Create React App, Node.js (Express/custom) |
| python.exe | FastAPI, Django, Flask, Streamlit, Gradio, Python (custom) |
| bun.exe | Bun |
| deno.exe | Deno |
| java.exe | Spring Boot, Java (Spring Boot / custom) |
| docker.exe | Docker |
