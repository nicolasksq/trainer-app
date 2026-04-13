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
- Garmin Connect — Custom Go client using mobile SSO OAuth

## Conventions

- Go standard project layout
- Error handling: return errors, don't panic
- MCP tools return JSON-formatted content
- Auth tokens stored in `~/.trainer-app/`
- **Pace must always be displayed as min:sec per km (e.g. 4:30/km)**. Never use km/h or m/s when presenting data to the user. Garmin API uses m/s internally but all user-facing output must be converted to min/km.
