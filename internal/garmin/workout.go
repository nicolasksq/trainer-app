package garmin

import (
	"encoding/json"
	"fmt"
)

// Workout represents a Garmin workout structure.
type Workout struct {
	WorkoutName       string           `json:"workoutName"`
	Description       string           `json:"description,omitempty"`
	SportType         SportType        `json:"sportType"`
	WorkoutSegments   []WorkoutSegment `json:"workoutSegments"`
	EstimatedDuration float64          `json:"estimatedDurationInSecs,omitempty"`
}

// SportType identifies the sport for a workout.
type SportType struct {
	SportTypeId  int    `json:"sportTypeId"`
	SportTypeKey string `json:"sportTypeKey"`
	DisplayOrder int    `json:"displayOrder"`
}

// WorkoutSegment is a segment within a workout.
type WorkoutSegment struct {
	SegmentOrder int           `json:"segmentOrder"`
	SportType    SportType     `json:"sportType"`
	WorkoutSteps []WorkoutStep `json:"workoutSteps"`
}

// WorkoutStep represents a single step or repeat group within a workout.
type WorkoutStep struct {
	Type               string        `json:"type"`
	StepOrder          int           `json:"stepOrder"`
	StepType           StepType      `json:"stepType"`
	EndCondition       EndCondition  `json:"endCondition"`
	EndConditionValue  float64       `json:"endConditionValue,omitempty"`
	TargetType         *TargetType   `json:"targetType,omitempty"`
	TargetValueOne     *float64      `json:"targetValueOne,omitempty"`
	TargetValueTwo     *float64      `json:"targetValueTwo,omitempty"`
	NumberOfIterations int           `json:"numberOfIterations,omitempty"`
	SmartRepeat        bool          `json:"smartRepeat,omitempty"`
	WorkoutSteps       []WorkoutStep `json:"workoutSteps,omitempty"`
}

// StepType identifies the type of workout step.
type StepType struct {
	StepTypeId  int    `json:"stepTypeId"`
	StepTypeKey string `json:"stepTypeKey"`
}

// EndCondition specifies how a step ends.
type EndCondition struct {
	ConditionTypeId  int    `json:"conditionTypeId"`
	ConditionTypeKey string `json:"conditionTypeKey"`
}

// TargetType specifies the target metric for a step.
type TargetType struct {
	WorkoutTargetTypeId  int    `json:"workoutTargetTypeId"`
	WorkoutTargetTypeKey string `json:"workoutTargetTypeKey"`
}

// Step type constants.
var (
	StepTypeWarmup   = StepType{StepTypeId: 1, StepTypeKey: "warmup"}
	StepTypeCooldown = StepType{StepTypeId: 2, StepTypeKey: "cooldown"}
	StepTypeInterval = StepType{StepTypeId: 3, StepTypeKey: "interval"}
	StepTypeRecovery = StepType{StepTypeId: 4, StepTypeKey: "recovery"}
	StepTypeRest     = StepType{StepTypeId: 5, StepTypeKey: "rest"}
	StepTypeRepeat   = StepType{StepTypeId: 6, StepTypeKey: "repeat"}
)

// End condition constants.
var (
	EndConditionDistance   = EndCondition{ConditionTypeId: 1, ConditionTypeKey: "distance"}
	EndConditionTime      = EndCondition{ConditionTypeId: 2, ConditionTypeKey: "time"}
	EndConditionHeartRate = EndCondition{ConditionTypeId: 3, ConditionTypeKey: "heart.rate"}
	EndConditionIterations = EndCondition{ConditionTypeId: 7, ConditionTypeKey: "iterations"}
)

// Target type constants.
var (
	TargetTypeNoTarget  = &TargetType{WorkoutTargetTypeId: 1, WorkoutTargetTypeKey: "no.target"}
	TargetTypePowerZone = &TargetType{WorkoutTargetTypeId: 2, WorkoutTargetTypeKey: "power.zone"}
	TargetTypeCadence   = &TargetType{WorkoutTargetTypeId: 3, WorkoutTargetTypeKey: "cadence"}
	TargetTypeHRZone    = &TargetType{WorkoutTargetTypeId: 4, WorkoutTargetTypeKey: "heart.rate.zone"}
	TargetTypeSpeedZone = &TargetType{WorkoutTargetTypeId: 5, WorkoutTargetTypeKey: "speed.zone"}
)

// Sport type presets.
var (
	SportRunning  = SportType{SportTypeId: 1, SportTypeKey: "running", DisplayOrder: 1}
	SportCycling  = SportType{SportTypeId: 2, SportTypeKey: "cycling", DisplayOrder: 2}
	SportSwimming = SportType{SportTypeId: 5, SportTypeKey: "swimming", DisplayOrder: 5}
)

// NewWarmupStep creates a warmup step with a time-based end condition.
func NewWarmupStep(durationSecs float64) WorkoutStep {
	return WorkoutStep{
		Type:              "ExecutableStepDTO",
		StepType:          StepTypeWarmup,
		EndCondition:      EndConditionTime,
		EndConditionValue: durationSecs,
		TargetType:        TargetTypeNoTarget,
	}
}

// NewCooldownStep creates a cooldown step with a time-based end condition.
func NewCooldownStep(durationSecs float64) WorkoutStep {
	return WorkoutStep{
		Type:              "ExecutableStepDTO",
		StepType:          StepTypeCooldown,
		EndCondition:      EndConditionTime,
		EndConditionValue: durationSecs,
		TargetType:        TargetTypeNoTarget,
	}
}

// NewIntervalStep creates an interval step with a time-based end condition.
func NewIntervalStep(durationSecs float64) WorkoutStep {
	return WorkoutStep{
		Type:              "ExecutableStepDTO",
		StepType:          StepTypeInterval,
		EndCondition:      EndConditionTime,
		EndConditionValue: durationSecs,
		TargetType:        TargetTypeNoTarget,
	}
}

// NewRecoveryStep creates a recovery step with a time-based end condition.
func NewRecoveryStep(durationSecs float64) WorkoutStep {
	return WorkoutStep{
		Type:              "ExecutableStepDTO",
		StepType:          StepTypeRecovery,
		EndCondition:      EndConditionTime,
		EndConditionValue: durationSecs,
		TargetType:        TargetTypeNoTarget,
	}
}

// NewRepeatGroup creates a repeat group containing the given steps.
func NewRepeatGroup(iterations int, steps []WorkoutStep) WorkoutStep {
	for i := range steps {
		steps[i].StepOrder = i + 1
	}
	return WorkoutStep{
		Type:               "RepeatGroupDTO",
		StepType:           StepTypeRepeat,
		EndCondition:       EndConditionIterations,
		NumberOfIterations: iterations,
		WorkoutSteps:       steps,
	}
}

// WithPaceTarget sets a speed zone target on a step. Pace is in sec/km, converted to m/s.
// targetValueOne = faster pace (higher m/s), targetValueTwo = slower pace (lower m/s).
func WithPaceTarget(step *WorkoutStep, minPaceSecPerKm, maxPaceSecPerKm float64) {
	step.TargetType = TargetTypeSpeedZone
	faster := 1000.0 / minPaceSecPerKm  // min pace = faster
	slower := 1000.0 / maxPaceSecPerKm  // max pace = slower
	step.TargetValueOne = &faster
	step.TargetValueTwo = &slower
}

// WithHRTarget sets a heart rate zone target on a step.
func WithHRTarget(step *WorkoutStep, minHR, maxHR float64) {
	step.TargetType = TargetTypeHRZone
	step.TargetValueOne = &minHR
	step.TargetValueTwo = &maxHR
}

// CreateWorkout creates a workout on Garmin Connect.
func (c *Client) CreateWorkout(workout *Workout) (json.RawMessage, error) {
	return c.doRequest("POST", "/workout-service/workout", workout)
}

// ScheduleWorkout schedules a workout to a specific date.
func (c *Client) ScheduleWorkout(workoutId string, date string) (json.RawMessage, error) {
	path := fmt.Sprintf("/workout-service/schedule/%s", workoutId)
	body := map[string]string{"date": date}
	return c.doRequest("POST", path, body)
}

// ListWorkouts lists workouts with pagination.
func (c *Client) ListWorkouts(start, limit int) (json.RawMessage, error) {
	path := fmt.Sprintf("/workout-service/workouts?start=%d&limit=%d", start, limit)
	return c.doRequest("GET", path, nil)
}

// DeleteWorkout deletes a workout by ID.
func (c *Client) DeleteWorkout(workoutId string) error {
	path := fmt.Sprintf("/workout-service/workout/%s", workoutId)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
