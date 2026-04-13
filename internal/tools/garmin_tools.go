package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolas-andreoli/trainer-app/internal/garmin"
)

// RegisterGarminTools registers all Garmin-related MCP tools on the server.
func RegisterGarminTools(s *server.MCPServer, client *garmin.Client) {
	s.AddTool(getGarminActivitiesTool(), getGarminActivitiesHandler(client))
	s.AddTool(getGarminActivityTool(), getGarminActivityHandler(client))
	s.AddTool(getGarminTrainingStatusTool(), getGarminTrainingStatusHandler(client))
	s.AddTool(getGarminBodyCompositionTool(), getGarminBodyCompositionHandler(client))
	s.AddTool(getGarminHeartRateTool(), getGarminHeartRateHandler(client))
	s.AddTool(getGarminSleepTool(), getGarminSleepHandler(client))
	s.AddTool(createGarminWorkoutTool(), createGarminWorkoutHandler(client))
	s.AddTool(scheduleGarminWorkoutTool(), scheduleGarminWorkoutHandler(client))
	s.AddTool(listGarminWorkoutsTool(), listGarminWorkoutsHandler(client))
	s.AddTool(deleteGarminWorkoutTool(), deleteGarminWorkoutHandler(client))
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func getDateParam(request mcp.CallToolRequest, key string) (string, *mcp.CallToolResult) {
	dateStr := request.GetString(key, today())
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		return "", mcp.NewToolResultError(fmt.Sprintf("invalid date format, use YYYY-MM-DD: %v", err))
	}
	return dateStr, nil
}

// --- get_garmin_activities ---

func getGarminActivitiesTool() mcp.Tool {
	return mcp.NewTool("get_garmin_activities",
		mcp.WithDescription("Get recent activities from Garmin Connect. Returns a list of activities with type, duration, distance, and other metrics."),
		mcp.WithNumber("count",
			mcp.Description("Number of recent activities to fetch (default 10, max 100)"),
		),
	)
}

func getGarminActivitiesHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		count := request.GetInt("count", 10)
		if count < 1 {
			count = 1
		}
		if count > 100 {
			count = 100
		}

		data, err := client.ListActivities(count)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// --- get_garmin_activity ---

func getGarminActivityTool() mcp.Tool {
	return mcp.NewTool("get_garmin_activity",
		mcp.WithDescription("Get detailed information about a specific Garmin activity by its ID."),
		mcp.WithNumber("activity_id",
			mcp.Required(),
			mcp.Description("The Garmin activity ID"),
		),
	)
}

func getGarminActivityHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("activity_id")
		if err != nil {
			return mcp.NewToolResultError("activity_id is required and must be a number"), nil
		}

		data, err := client.GetActivity(int64(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// --- get_garmin_training_status ---

func getGarminTrainingStatusTool() mcp.Tool {
	return mcp.NewTool("get_garmin_training_status",
		mcp.WithDescription("Get training readiness and training status from Garmin Connect for a given date."),
		mcp.WithString("date",
			mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
		),
	)
}

func getGarminTrainingStatusHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		date, errResult := getDateParam(request, "date")
		if errResult != nil {
			return errResult, nil
		}

		readiness, err := client.GetTrainingReadiness(date)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get training readiness: %v", err)), nil
		}

		status, err := client.GetTrainingStatus(date)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get training status: %v", err)), nil
		}

		result := map[string]json.RawMessage{
			"trainingReadiness": readiness,
			"trainingStatus":    status,
		}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

// --- get_garmin_body_composition ---

func getGarminBodyCompositionTool() mcp.Tool {
	return mcp.NewTool("get_garmin_body_composition",
		mcp.WithDescription("Get body composition data from Garmin Connect (weight, body fat %, BMI, etc)."),
		mcp.WithString("start_date",
			mcp.Description("Start date in YYYY-MM-DD format (defaults to 30 days ago)"),
		),
		mcp.WithString("end_date",
			mcp.Description("End date in YYYY-MM-DD format (defaults to today)"),
		),
	)
}

func getGarminBodyCompositionHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		endDate := request.GetString("end_date", today())
		startDate := request.GetString("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))

		data, err := client.GetBodyComposition(startDate, endDate)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to fetch body composition: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// --- get_garmin_heart_rate ---

func getGarminHeartRateTool() mcp.Tool {
	return mcp.NewTool("get_garmin_heart_rate",
		mcp.WithDescription("Get heart rate and HRV data from Garmin Connect for a given date."),
		mcp.WithString("date",
			mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
		),
	)
}

func getGarminHeartRateHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		date, errResult := getDateParam(request, "date")
		if errResult != nil {
			return errResult, nil
		}

		hr, err := client.GetHeartRate(date)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to fetch heart rate: %v", err)), nil
		}

		hrv, err := client.GetHRV(date)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to fetch HRV: %v", err)), nil
		}

		result := map[string]json.RawMessage{
			"heartRate": hr,
			"hrv":       hrv,
		}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

// --- get_garmin_sleep ---

func getGarminSleepTool() mcp.Tool {
	return mcp.NewTool("get_garmin_sleep",
		mcp.WithDescription("Get sleep data from Garmin Connect including sleep duration, sleep stages, and sleep score."),
		mcp.WithString("date",
			mcp.Description("Date in YYYY-MM-DD format (defaults to today). Returns sleep data for the night ending on this date."),
		),
	)
}

func getGarminSleepHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		date, errResult := getDateParam(request, "date")
		if errResult != nil {
			return errResult, nil
		}

		data, err := client.GetSleep(date)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to fetch sleep data: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// --- create_garmin_workout ---

func createGarminWorkoutTool() mcp.Tool {
	return mcp.NewTool("create_garmin_workout",
		mcp.WithDescription(`Create a structured workout on Garmin Connect. The steps parameter accepts a simplified JSON format:
[
  {"type": "warmup", "duration": "10:00"},
  {"type": "repeat", "iterations": 4, "steps": [
    {"type": "interval", "duration": "3:00", "pace_min": "4:15", "pace_max": "4:25"},
    {"type": "recovery", "duration": "2:00"}
  ]},
  {"type": "cooldown", "duration": "5:00"}
]
Duration format is MM:SS. Pace format is M:SS per km.`),
		mcp.WithString("name", mcp.Required(), mcp.Description("Workout name")),
		mcp.WithString("description", mcp.Description("Workout description")),
		mcp.WithString("sport", mcp.Description("Sport type: running, cycling, or swimming (default: running)")),
		mcp.WithString("steps", mcp.Required(), mcp.Description("JSON array of step definitions")),
	)
}

// simplifiedStep is the user-facing step format.
type simplifiedStep struct {
	Type       string           `json:"type"`
	Duration   string           `json:"duration,omitempty"`
	PaceMin    string           `json:"pace_min,omitempty"`
	PaceMax    string           `json:"pace_max,omitempty"`
	HRMin      *float64         `json:"hr_min,omitempty"`
	HRMax      *float64         `json:"hr_max,omitempty"`
	Iterations int              `json:"iterations,omitempty"`
	Steps      []simplifiedStep `json:"steps,omitempty"`
}

// parseDuration parses "MM:SS" to seconds.
func parseDuration(s string) (float64, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid duration format %q, expected MM:SS", s)
	}
	min, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes in %q: %v", s, err)
	}
	sec, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds in %q: %v", s, err)
	}
	return float64(min*60 + sec), nil
}

// parsePace parses "M:SS" pace per km to seconds per km.
func parsePace(s string) (float64, error) {
	return parseDuration(s) // same format
}

func convertStep(s simplifiedStep) (garmin.WorkoutStep, error) {
	switch s.Type {
	case "warmup":
		dur, err := parseDuration(s.Duration)
		if err != nil {
			return garmin.WorkoutStep{}, err
		}
		step := garmin.NewWarmupStep(dur)
		applyTargets(&step, s)
		return step, nil

	case "cooldown":
		dur, err := parseDuration(s.Duration)
		if err != nil {
			return garmin.WorkoutStep{}, err
		}
		step := garmin.NewCooldownStep(dur)
		applyTargets(&step, s)
		return step, nil

	case "interval":
		dur, err := parseDuration(s.Duration)
		if err != nil {
			return garmin.WorkoutStep{}, err
		}
		step := garmin.NewIntervalStep(dur)
		applyTargets(&step, s)
		return step, nil

	case "recovery":
		dur, err := parseDuration(s.Duration)
		if err != nil {
			return garmin.WorkoutStep{}, err
		}
		step := garmin.NewRecoveryStep(dur)
		applyTargets(&step, s)
		return step, nil

	case "repeat":
		if s.Iterations < 1 {
			return garmin.WorkoutStep{}, fmt.Errorf("repeat must have iterations >= 1")
		}
		var subSteps []garmin.WorkoutStep
		for _, sub := range s.Steps {
			ws, err := convertStep(sub)
			if err != nil {
				return garmin.WorkoutStep{}, err
			}
			subSteps = append(subSteps, ws)
		}
		return garmin.NewRepeatGroup(s.Iterations, subSteps), nil

	default:
		return garmin.WorkoutStep{}, fmt.Errorf("unknown step type: %q", s.Type)
	}
}

func applyTargets(step *garmin.WorkoutStep, s simplifiedStep) {
	if s.PaceMin != "" && s.PaceMax != "" {
		minPace, err1 := parsePace(s.PaceMin)
		maxPace, err2 := parsePace(s.PaceMax)
		if err1 == nil && err2 == nil {
			garmin.WithPaceTarget(step, minPace, maxPace)
		}
	}
	if s.HRMin != nil && s.HRMax != nil {
		garmin.WithHRTarget(step, *s.HRMin, *s.HRMax)
	}
}

func sportTypeFromString(s string) garmin.SportType {
	switch strings.ToLower(s) {
	case "cycling":
		return garmin.SportCycling
	case "swimming":
		return garmin.SportSwimming
	default:
		return garmin.SportRunning
	}
}

func createGarminWorkoutHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError("name is required"), nil
		}
		description := request.GetString("description", "")
		sport := request.GetString("sport", "running")
		stepsJSON, err := request.RequireString("steps")
		if err != nil {
			return mcp.NewToolResultError("steps is required"), nil
		}

		var simplified []simplifiedStep
		if err := json.Unmarshal([]byte(stepsJSON), &simplified); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid steps JSON: %v", err)), nil
		}

		sportType := sportTypeFromString(sport)
		var workoutSteps []garmin.WorkoutStep
		var totalDuration float64

		for i, s := range simplified {
			ws, err := convertStep(s)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("step %d: %v", i, err)), nil
			}
			ws.StepOrder = i + 1
			workoutSteps = append(workoutSteps, ws)
			totalDuration += ws.EndConditionValue
		}

		workout := &garmin.Workout{
			WorkoutName: name,
			Description: description,
			SportType:   sportType,
			WorkoutSegments: []garmin.WorkoutSegment{
				{
					SegmentOrder: 1,
					SportType:    sportType,
					WorkoutSteps: workoutSteps,
				},
			},
			EstimatedDuration: totalDuration,
		}

		data, err := client.CreateWorkout(workout)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create workout: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// --- schedule_garmin_workout ---

func scheduleGarminWorkoutTool() mcp.Tool {
	return mcp.NewTool("schedule_garmin_workout",
		mcp.WithDescription("Schedule an existing Garmin workout to a specific date."),
		mcp.WithString("workout_id", mcp.Required(), mcp.Description("The workout ID to schedule")),
		mcp.WithString("date", mcp.Required(), mcp.Description("Date in YYYY-MM-DD format")),
	)
}

func scheduleGarminWorkoutHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workoutID, err := request.RequireString("workout_id")
		if err != nil {
			return mcp.NewToolResultError("workout_id is required"), nil
		}
		date, err := request.RequireString("date")
		if err != nil {
			return mcp.NewToolResultError("date is required"), nil
		}
		if _, parseErr := time.Parse("2006-01-02", date); parseErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid date format, use YYYY-MM-DD: %v", parseErr)), nil
		}

		data, err := client.ScheduleWorkout(workoutID, date)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to schedule workout: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// --- list_garmin_workouts ---

func listGarminWorkoutsTool() mcp.Tool {
	return mcp.NewTool("list_garmin_workouts",
		mcp.WithDescription("List existing workouts from Garmin Connect."),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of workouts to return (default 20)"),
		),
	)
}

func listGarminWorkoutsHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limit := request.GetInt("limit", 20)
		if limit < 1 {
			limit = 1
		}

		data, err := client.ListWorkouts(0, limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list workouts: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// --- delete_garmin_workout ---

func deleteGarminWorkoutTool() mcp.Tool {
	return mcp.NewTool("delete_garmin_workout",
		mcp.WithDescription("Delete a workout from Garmin Connect."),
		mcp.WithString("workout_id", mcp.Required(), mcp.Description("The workout ID to delete")),
	)
}

func deleteGarminWorkoutHandler(client *garmin.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workoutID, err := request.RequireString("workout_id")
		if err != nil {
			return mcp.NewToolResultError("workout_id is required"), nil
		}

		if err := client.DeleteWorkout(workoutID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to delete workout: %v", err)), nil
		}

		return mcp.NewToolResultText(`{"status":"deleted"}`), nil
	}
}
