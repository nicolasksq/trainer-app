# Setup Guide

Complete guide to get Trainer App running on your machine.

## Prerequisites

Before you start, make sure you have:

- **Go 1.25+** -- Install from [go.dev/dl](https://go.dev/dl/)
- **Claude Code** -- Install from [docs.anthropic.com](https://docs.anthropic.com/en/docs/claude-code) (requires an Anthropic account)
- **Strava account** -- With recorded activities. Sign up at [strava.com](https://www.strava.com/)
- **Garmin Connect account** -- Optional but recommended. Requires a compatible Garmin device for health metrics and workout scheduling

## Quick Setup (Recommended)

The fastest way to get started. The interactive setup CLI handles everything for you.

```bash
git clone https://github.com/nicolasksq/trainer-app.git
cd trainer-app
make setup
```

`make setup` runs an interactive CLI that will:

1. Check that Go is installed and at the correct version
2. Walk you through creating a Strava API application (if needed)
3. Create your `.env` file with the required credentials
4. Run the Strava OAuth flow to authorize access to your data
5. Configure your Garmin Connect credentials
6. Verify that both integrations work

After setup completes, skip ahead to [Verifying Your Setup](#verifying-your-setup).

## Manual Setup

If you prefer to configure things yourself, follow these steps.

### 1. Clone the Repository

```bash
git clone https://github.com/nicolasksq/trainer-app.git
cd trainer-app
```

### 2. Install Go Dependencies

```bash
go mod download
```

### 3. Create a Strava API Application

You need a Strava API app to authorize access to your training data.

1. Go to [strava.com/settings/api](https://www.strava.com/settings/api)
2. Fill in the application form:
   - **Application Name**: Any name (e.g., "Trainer App")
   - **Category**: Choose "Training"
   - **Club**: Leave blank
   - **Website**: Any URL (e.g., `http://localhost`)
   - **Authorization Callback Domain**: `localhost`
3. Click **Create**
4. Note your **Client ID** and **Client Secret** -- you will need these in the next step

### 4. Create Your Environment File

Copy the example file and fill in your credentials:

```bash
cp .env.example .env
```

Edit `.env` with your values:

```
# From step 3
STRAVA_CLIENT_ID=12345
STRAVA_CLIENT_SECRET=abc123def456...

# Your Garmin Connect login credentials
GARMIN_EMAIL=your-email@example.com
GARMIN_PASSWORD=your-password
```

### 5. Authorize Strava Access

Run the OAuth flow to grant the app access to your Strava data:

```bash
go run ./cmd/strava-mcp auth
```

This will:

1. Open your default browser to Strava's authorization page
2. Ask you to approve access to your activities
3. Save the OAuth tokens locally to `~/.trainer-app/`

You only need to do this once. Tokens are refreshed automatically.

### 6. Verify Strava Works

```bash
make test-strava
```

If successful, you will see a list of available Strava tools:

```
  get_athlete                    Get the authenticated athlete's profile information
  list_activities                List the authenticated athlete's activities with
  get_activity                   Get detailed information about a specific activity
  get_activity_streams           Get time-series data streams for an activity
  get_activity_laps              Get lap data for a specific activity
  get_activity_zones             Get heart rate and power zone distribution data
  get_athlete_stats              Get the authenticated athlete's aggregate statistics
```

### 7. Verify Garmin Works

```bash
make test-garmin
```

If successful, you will see a list of available Garmin tools:

```
  get_garmin_activities          Get recent activities from Garmin Connect
  get_garmin_activity            Get detailed information about a specific Garmin
  get_garmin_training_status     Get training readiness and training status from
  get_garmin_body_composition    Get body composition data from Garmin Connect
  get_garmin_heart_rate          Get heart rate and HRV data from Garmin Connect
  get_garmin_sleep               Get sleep data from Garmin Connect including sleep
  create_garmin_workout          Create a structured workout on Garmin Connect
  schedule_garmin_workout        Schedule an existing Garmin workout to a specific
  list_garmin_workouts           List existing workouts from Garmin Connect
  delete_garmin_workout          Delete a workout from Garmin Connect
```

### 8. Start Using the Agent

Open Claude Code in the project directory and invoke the trainer:

```
@trainer Analyze my last 2 weeks of training
```

The agent automatically starts both MCP servers and connects to your Strava and Garmin accounts. No additional configuration is needed.

## Verifying Your Setup

Run the test commands to confirm everything is working:

```bash
# Test Strava connection
make test-strava

# Test Garmin connection
make test-garmin
```

**What success looks like**: Both commands print a table of available tools (7 for Strava, 10 for Garmin) without any errors.

**What failure looks like**: You will see `Failed - run 'make setup' first`. Common causes:

- Missing `.env` file or empty credentials -- run `make setup` or check your `.env`
- Strava tokens not yet authorized -- run `make auth-strava`
- Garmin credentials incorrect -- verify your email and password in `.env`

## Token Storage

Auth tokens are stored in `~/.trainer-app/`. This directory is created automatically during the Strava OAuth flow. Do not commit this directory to version control -- it contains your access tokens.

## Updating

To update to the latest version:

```bash
git pull
go mod download
```

Your `.env` file and auth tokens in `~/.trainer-app/` are preserved across updates.
