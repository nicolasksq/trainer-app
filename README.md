![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Claude Code](https://img.shields.io/badge/Claude_Code-Agent-blueviolet)

# Trainer App

**Your AI-powered endurance coach -- real data, real plans, straight to your watch.**

Trainer App is a [Claude Code](https://docs.anthropic.com/en/docs/claude-code) agent that acts as a personal endurance coach. It reads your actual training data from Strava and Garmin Connect, analyzes your fitness, recovery, and performance trends, then creates periodized training plans and pushes structured workouts directly to your Garmin watch.

## Features

- **Strava integration** -- activities, time-series streams, zone distributions, athlete stats, and profile data
- **Garmin integration** -- activities, training readiness, heart rate, HRV, sleep quality, and body composition
- **Workout creation and scheduling** -- build structured workouts (warmup, intervals, recovery, cooldown with pace/HR targets) and sync them to your Garmin device
- **Periodized coaching** -- base, build, peak, and taper phases with load monitoring and acute-to-chronic workload ratios
- **Multi-sport support** -- running, cycling, swimming, and strength training
- **Zone-based prescriptions** -- heart rate zones, pace zones, power zones, and RPE
- **Memory system** -- persists your athlete profile, goals, training plan, and coaching preferences across sessions

## Architecture

```
                          +------------------+
                          |    @trainer      |
                          |  Claude Code     |
                          |     Agent        |
                          +--------+---------+
                                   |
                      +------------+------------+
                      |                         |
               +------+------+          +-------+------+
               | strava-mcp  |          | garmin-mcp   |
               |  (stdio)    |          |  (stdio)     |
               +------+------+          +-------+------+
                      |                         |
               +------+------+          +-------+------+
               |  Strava API |          | Garmin       |
               |    v3       |          | Connect API  |
               +-------------+          +--------------+
```

## MCP Tools

### Strava Tools (7)

| Tool | Description |
|------|-------------|
| `get_athlete` | Get athlete profile (name, weight, FTP, bikes, shoes) |
| `list_activities` | List activities with date filtering and pagination |
| `get_activity` | Get detailed activity data (splits, segments, best efforts) |
| `get_activity_streams` | Get time-series data (HR, power, cadence, speed, elevation, temp) |
| `get_activity_laps` | Get lap data (time, distance, pace, HR per lap) |
| `get_activity_zones` | Get HR and power zone distribution for an activity |
| `get_athlete_stats` | Get aggregate stats (recent, YTD, all-time totals) |

### Garmin Tools (10)

| Tool | Description |
|------|-------------|
| `get_garmin_activities` | List recent activities with type, duration, distance, and metrics |
| `get_garmin_activity` | Get detailed information about a specific activity |
| `get_garmin_training_status` | Get training readiness and training status for a date |
| `get_garmin_body_composition` | Get weight, body fat %, BMI over a date range |
| `get_garmin_heart_rate` | Get heart rate and HRV data for a date |
| `get_garmin_sleep` | Get sleep duration, stages, and score |
| `create_garmin_workout` | Create a structured workout with pace/HR targets |
| `schedule_garmin_workout` | Schedule a workout to a specific date on your calendar |
| `list_garmin_workouts` | List existing workouts from Garmin Connect |
| `delete_garmin_workout` | Delete a workout from Garmin Connect |

## Quick Start

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Claude Code](https://docs.anthropic.com/en/docs/claude-code)
- A Strava account with activities
- A Garmin Connect account (optional but recommended)

### Setup

```bash
git clone https://github.com/nicolasksq/trainer-app.git
cd trainer-app
make setup
```

The interactive setup wizard will guide you through everything: Strava OAuth, Garmin credentials, and verification.

For manual setup or troubleshooting, see the detailed guides:
- [Setup Guide](docs/SETUP.md) -- step-by-step installation
- [Usage Guide](docs/USAGE.md) -- prompts, examples, and tips
- [Troubleshooting](docs/TROUBLESHOOTING.md) -- common issues and fixes

### 5. Use the agent

Open Claude Code in the project directory and invoke the trainer agent:

```
@trainer Analyze my last 2 weeks of training
```

The agent automatically starts both MCP servers and connects to your Strava and Garmin accounts.

## Usage Examples

```
@trainer Analyze my last 2 weeks of training and tell me if I'm overtraining

@trainer Create a 12-week 10K plan targeting sub-40 minutes

@trainer Schedule this week's workouts to my Garmin watch

@trainer How was my recovery this week? Check my HRV, sleep, and resting HR

@trainer Review my long run from Sunday -- was I in the right zones?

@trainer I have a half marathon in 8 weeks, adjust my plan
```

## Project Structure

```
trainer-app/
├── .claude/
│   └── agents/
│       └── trainer.md          # Coach agent definition and prompt
├── cmd/
│   ├── setup/
│   │   └── main.go             # Interactive setup wizard
│   ├── strava-mcp/
│   │   └── main.go             # Strava MCP server entry point + OAuth flow
│   └── garmin-mcp/
│       └── main.go             # Garmin MCP server entry point
├── docs/
│   ├── SETUP.md                # Detailed setup guide
│   ├── USAGE.md                # Usage guide with examples
│   └── TROUBLESHOOTING.md      # Common issues and fixes
├── internal/
│   ├── envutil/
│   │   └── dotenv.go           # .env file loader
│   ├── garmin/
│   │   ├── auth.go             # Garmin Connect authentication
│   │   ├── client.go           # HTTP client for Garmin Connect API
│   │   ├── activities.go       # Activity fetching
│   │   ├── health.go           # Health metrics (HR, HRV, sleep, body comp)
│   │   └── workout.go          # Workout creation, scheduling, and types
│   ├── strava/
│   │   ├── auth.go             # Strava OAuth2 token management
│   │   ├── client.go           # Strava API client wrapper
│   │   └── activities.go       # Activity and stream fetching
│   └── tools/
│       ├── strava_tools.go     # Strava MCP tool definitions
│       └── garmin_tools.go     # Garmin MCP tool definitions
├── .env.example                # Environment variable template
├── CLAUDE.md                   # Project conventions for Claude Code
├── Makefile                    # Build, setup, and test commands
├── go.mod
└── go.sum
```

## Tech Stack

- **Language** -- [Go](https://go.dev/)
- **MCP SDK** -- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)
- **Strava API** -- [strava/go.strava](https://github.com/strava/go.strava) (API v3)
- **Garmin Connect API** -- Custom Go client using Garmin's mobile SSO OAuth flow
- **Agent runtime** -- [Claude Code](https://docs.anthropic.com/en/docs/claude-code) with custom agent definitions

## License

MIT
