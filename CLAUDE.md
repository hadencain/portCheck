# PortWatch

Windows CLI tool for localhost dev server management.

## Build

```bash
go build -o portwatch.exe .
```

## Test

```bash
go test ./...
```

## Stack
- Go 1.22+
- cobra (CLI), gopsutil (process/network), tablewriter (tables), fatih/color (colors)

## Key Files
- `cmd/` — one file per subcommand
- `internal/ports/` — all port + process data collection
- `internal/processes/` — dev process detection + framework hints
- `internal/health/` — scan findings
- `internal/display/` — all terminal rendering
