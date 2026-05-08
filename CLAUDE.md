# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o fmeflow
go build -ldflags="-X github.com/safesoftware/fmeflow-cli/cmd.appVersion=$(git rev-parse --short HEAD)" -o fmeflow

# Run all tests
go test ./...

# Run a single test
go test -run TestBackupV4 ./cmd

# Run directly during development
go run main.go

# Regenerate command docs
go run docgen/docgen.go
```

## Architecture

This is a Cobra-based CLI for interacting with the FME Flow REST API. All commands live in `/cmd/`, one file per command (e.g., `backup.go`, `login.go`).

### Command pattern

Every command follows the same structure:

```go
func newBackupCmd() *cobra.Command {
    f := backupFlags{}
    cmd := &cobra.Command{
        Use:   "backup",
        Short: "...",
        RunE:  backupRun(&f),
    }
    cmd.Flags().StringVar(&f.outputPath, "output-path", "", "...")
    return cmd
}

func backupRun(f *backupFlags) func(*cobra.Command, []string) error {
    return func(cmd *cobra.Command, args []string) error { ... }
}
```

Commands are registered in `cmd/root.go`.

### API versioning

FME Flow has v3 and v4 REST APIs. Each command detects the correct version via a build number threshold stored in the Viper config key `build`:

```go
const backupV4BuildThreshold = 25208

if viper.GetInt("build") >= backupV4BuildThreshold {
    return backupRunV4(f, cmd, args)
}
return backupRunV3(f, cmd, args)
```

Split test files (`backup_v3_test.go`, `backup_v4_test.go`) cover each version independently.

### Config & HTTP

- Config is stored at `$XDG_CONFIG_HOME/.fmeflow-cli.yaml` (fields: `url`, `token`, `build`)
- Viper manages config; `buildFmeFlowRequest()` in `cmd/functions.go` injects the auth token into every outbound request

### Output formatting

- Tables use the `go-pretty` library
- Commands support `--json` for raw JSON output and `--no-headers` / custom column flags for table output

### Testing

Tests are table-driven and live in `*_test.go` files alongside each command. The shared infrastructure is in `cmd/testFunctions.go`:

- `testCase` struct defines inputs (args, HTTP status, request body) and expected outputs (stdout regex, stderr, file contents)
- `runTests()` spins up an `httptest` server, configures Viper, executes the command, and validates all outputs

Assertions use `testify/require`.

## MCP Server (`mcp/`)

A TypeScript MCP server that exposes FME Flow operations as tools, wrapping the `fmeflow` CLI binary.

```bash
cd mcp
npm install
npm run build      # compile TypeScript -> dist/
npm run dev        # run with tsx (no build step needed)
```

**Binary resolution:** tries `fmeflow` on PATH first, falls back to `../fmeflow` (the repo root build output).

**Tools exposed:** `fmeflow_healthcheck`, `fmeflow_workspaces_list`, `fmeflow_run`, `fmeflow_jobs`. Each tool appends `--json` and returns stdout as the tool result.

**To wire into Claude Desktop**, add to `~/Library/Application Support/Claude/claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "fmeflow": {
      "command": "node",
      "args": ["/path/to/fmeflow-cli/mcp/dist/index.js"]
    }
  }
}
```
