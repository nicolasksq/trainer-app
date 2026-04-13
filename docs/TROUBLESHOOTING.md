# Troubleshooting

Common issues and how to fix them.

## Common Issues

| Problem | Cause | Solution |
|---------|-------|----------|
| `No Strava token found` | OAuth flow has not been run yet | Run `make auth-strava` or `make setup` |
| `Strava Authorization Error` | Access token expired or revoked | Run `make auth-strava` to re-authorize |
| `Garmin returned 403 Forbidden` | CAPTCHA or bot detection triggered | Log in at [connect.garmin.com](https://connect.garmin.com) in your browser first, then retry |
| `Missing GARMIN_EMAIL` | `.env` file not configured | Run `make setup` or add credentials to `.env` manually |
| `Missing STRAVA_CLIENT_ID` | `.env` file not configured | Run `make setup` or add credentials to `.env` manually |
| MCP server won't start | Missing `.env` or invalid credentials | Run `make test-strava` and `make test-garmin` to diagnose |
| Port 8080 in use during Strava auth | Another process is using the port | Stop the other process (`lsof -i :8080`) or wait and retry |
| `command not found: claude` | Claude Code is not installed | Install from [docs.anthropic.com](https://docs.anthropic.com/en/docs/claude-code) |
| Agent doesn't remember my goals | Memory was not saved | Ask `@trainer` to save your profile explicitly |
| Garmin workouts don't appear on watch | Device has not synced | Open the Garmin Connect app on your phone and sync your device |
| `go: command not found` | Go is not installed | Install from [go.dev/dl](https://go.dev/dl/) |
| Rate limit errors from Strava | Too many API calls in a short period | Wait a few minutes and retry |
| `Failed - run 'make setup' first` | Test command failed | Run `make setup` to configure credentials and auth |
| `invalid 'before' date format` | Wrong date format in Strava queries | Use ISO 8601 format: `2024-01-15T00:00:00Z` |
| Garmin login fails repeatedly | Password changed or 2FA enabled | Update your password in `.env` and check 2FA settings |

## Diagnostic Commands

Run these to narrow down the issue:

```bash
# Check Go is installed and at the right version
go version

# Test Strava MCP server and list tools
make test-strava

# Test Garmin MCP server and list tools
make test-garmin

# Re-run Strava OAuth if tokens are expired
make auth-strava

# Re-run full interactive setup
make setup
```

## Resetting Your Setup

If things are in a bad state, you can start fresh:

```bash
# Remove saved tokens
rm -rf ~/.trainer-app/

# Remove local environment file
rm .env

# Run setup again
make setup
```

This will not affect your Strava or Garmin accounts -- it only removes local tokens and configuration.

## Still Stuck?

Open an issue on GitHub: [github.com/nicolasksq/trainer-app/issues](https://github.com/nicolasksq/trainer-app/issues)

Include the output of `make test-strava` and `make test-garmin` in your issue to help with debugging.
