# Trainer App

Personal trainer Claude Code agent that reads training data from Strava and Garmin to generate personalized training plans.

## Project Structure

- `cmd/strava-mcp/` — Strava MCP server (stdio transport)
- `cmd/garmin-mcp/` — Garmin MCP server (stdio transport)
- `internal/strava/` — Strava API client wrapper
- `internal/garmin/` — Garmin Connect client wrapper
- `internal/tools/` — Shared MCP tool definitions
- `.claude/agents/trainer.md` — Claude Code agent definition

## Build & Run

```bash
go build ./...                          # Build all
go run ./cmd/strava-mcp                 # Run Strava MCP server
go run ./cmd/garmin-mcp                 # Run Garmin MCP server
```

## Key Dependencies

- `github.com/mark3labs/mcp-go` — MCP server SDK
- `github.com/strava/go.strava` — Official Strava API client
- `github.com/abrander/garmin-connect` — Garmin Connect client

## Conventions

- Go standard project layout
- Error handling: return errors, don't panic
- MCP tools return JSON-formatted content
- Auth tokens stored in `~/.trainer-app/`
