package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolas-andreoli/trainer-app/internal/strava"
	gostrava "github.com/strava/go.strava"
)

func RegisterStravaTools(s *server.MCPServer, client *strava.Client) {
	s.AddTool(getAthleteTool(), handleGetAthlete(client))
	s.AddTool(listActivitiesTool(), handleListActivities(client))
	s.AddTool(getActivityTool(), handleGetActivity(client))
	s.AddTool(getActivityStreamsTool(), handleGetActivityStreams(client))
	s.AddTool(getActivityLapsTool(), handleGetActivityLaps(client))
	s.AddTool(getActivityZonesTool(), handleGetActivityZones(client))
	s.AddTool(getAthleteStatsTool(), handleGetAthleteStats(client))
}

func toJSON(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to marshal response: %s"}`, err)
	}
	return string(data)
}

func errResult(msg string, err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(fmt.Sprintf("%s: %v", msg, err)), nil
}

// get_athlete

func getAthleteTool() mcp.Tool {
	return mcp.NewTool("get_athlete",
		mcp.WithDescription("Get the authenticated athlete's profile information including name, location, weight, FTP, bikes, and shoes"),
	)
}

func handleGetAthlete(client *strava.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		athlete, err := client.GetCurrentAthlete()
		if err != nil {
			return errResult("failed to get athlete", err)
		}
		return mcp.NewToolResultText(toJSON(athlete)), nil
	}
}

// list_activities

func listActivitiesTool() mcp.Tool {
	return mcp.NewTool("list_activities",
		mcp.WithDescription("List the authenticated athlete's activities with optional date filtering and pagination"),
		mcp.WithString("before",
			mcp.Description("Only return activities before this date (ISO 8601 format, e.g. 2024-01-15T00:00:00Z)"),
		),
		mcp.WithString("after",
			mcp.Description("Only return activities after this date (ISO 8601 format, e.g. 2024-01-01T00:00:00Z)"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (default: 1)"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of activities per page (default: 30, max: 200)"),
		),
	)
}

func handleListActivities(client *strava.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var before, after time.Time

		if s := request.GetString("before", ""); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return errResult("invalid 'before' date format", err)
			}
			before = t
		}
		if s := request.GetString("after", ""); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return errResult("invalid 'after' date format", err)
			}
			after = t
		}

		page := request.GetInt("page", 0)
		perPage := request.GetInt("per_page", 0)

		activities, err := client.ListActivities(before, after, page, perPage)
		if err != nil {
			return errResult("failed to list activities", err)
		}
		return mcp.NewToolResultText(toJSON(activities)), nil
	}
}

// get_activity

func getActivityTool() mcp.Tool {
	return mcp.NewTool("get_activity",
		mcp.WithDescription("Get detailed information about a specific activity including splits, segment efforts, and best efforts"),
		mcp.WithNumber("activity_id",
			mcp.Description("The unique identifier of the activity"),
			mcp.Required(),
		),
	)
}

func handleGetActivity(client *strava.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("activity_id")
		if err != nil {
			return errResult("missing parameter", err)
		}

		activity, err := client.GetActivity(int64(id))
		if err != nil {
			return errResult("failed to get activity", err)
		}
		return mcp.NewToolResultText(toJSON(activity)), nil
	}
}

// get_activity_streams

var validStreamTypes = map[string]gostrava.StreamType{
	"time":           gostrava.StreamTypes.Time,
	"latlng":         gostrava.StreamTypes.Location,
	"distance":       gostrava.StreamTypes.Distance,
	"altitude":       gostrava.StreamTypes.Elevation,
	"velocity_smooth": gostrava.StreamTypes.Speed,
	"heartrate":      gostrava.StreamTypes.HeartRate,
	"cadence":        gostrava.StreamTypes.Cadence,
	"watts":          gostrava.StreamTypes.Power,
	"temp":           gostrava.StreamTypes.Temperature,
	"moving":         gostrava.StreamTypes.Moving,
	"grade_smooth":   gostrava.StreamTypes.Grade,
}

func getActivityStreamsTool() mcp.Tool {
	return mcp.NewTool("get_activity_streams",
		mcp.WithDescription("Get time-series data streams for an activity (heart rate, power, cadence, speed, elevation, etc.)"),
		mcp.WithNumber("activity_id",
			mcp.Description("The unique identifier of the activity"),
			mcp.Required(),
		),
		mcp.WithArray("types",
			mcp.Description("Stream types to retrieve: time, latlng, distance, altitude, velocity_smooth, heartrate, cadence, watts, temp, moving, grade_smooth"),
			mcp.Required(),
			mcp.WithStringItems(),
		),
	)
}

func handleGetActivityStreams(client *strava.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("activity_id")
		if err != nil {
			return errResult("missing parameter", err)
		}

		typeNames := request.GetStringSlice("types", nil)
		if len(typeNames) == 0 {
			return mcp.NewToolResultError("at least one stream type is required"), nil
		}

		streamTypes := make([]gostrava.StreamType, 0, len(typeNames))
		for _, name := range typeNames {
			st, ok := validStreamTypes[name]
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("invalid stream type: %s", name)), nil
			}
			streamTypes = append(streamTypes, st)
		}

		streams, err := client.GetActivityStreams(int64(id), streamTypes)
		if err != nil {
			return errResult("failed to get activity streams", err)
		}
		return mcp.NewToolResultText(toJSON(streams)), nil
	}
}

// get_activity_laps

func getActivityLapsTool() mcp.Tool {
	return mcp.NewTool("get_activity_laps",
		mcp.WithDescription("Get lap data for a specific activity including time, distance, pace, and heart rate per lap"),
		mcp.WithNumber("activity_id",
			mcp.Description("The unique identifier of the activity"),
			mcp.Required(),
		),
	)
}

func handleGetActivityLaps(client *strava.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("activity_id")
		if err != nil {
			return errResult("missing parameter", err)
		}

		laps, err := client.GetActivityLaps(int64(id))
		if err != nil {
			return errResult("failed to get activity laps", err)
		}
		return mcp.NewToolResultText(toJSON(laps)), nil
	}
}

// get_activity_zones

func getActivityZonesTool() mcp.Tool {
	return mcp.NewTool("get_activity_zones",
		mcp.WithDescription("Get heart rate and power zone distribution data for a specific activity"),
		mcp.WithNumber("activity_id",
			mcp.Description("The unique identifier of the activity"),
			mcp.Required(),
		),
	)
}

func handleGetActivityZones(client *strava.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("activity_id")
		if err != nil {
			return errResult("missing parameter", err)
		}

		zones, err := client.GetActivityZones(int64(id))
		if err != nil {
			return errResult("failed to get activity zones", err)
		}
		return mcp.NewToolResultText(toJSON(zones)), nil
	}
}

// get_athlete_stats

func getAthleteStatsTool() mcp.Tool {
	return mcp.NewTool("get_athlete_stats",
		mcp.WithDescription("Get the authenticated athlete's aggregate statistics including recent, YTD, and all-time totals for rides and runs"),
	)
}

func handleGetAthleteStats(client *strava.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		athlete, err := client.GetCurrentAthlete()
		if err != nil {
			return errResult("failed to get athlete", err)
		}

		stats, err := client.GetAthleteStats(athlete.Id)
		if err != nil {
			return errResult("failed to get athlete stats", err)
		}
		return mcp.NewToolResultText(toJSON(stats)), nil
	}
}
