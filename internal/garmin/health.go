package garmin

import (
	"encoding/json"
	"fmt"
)

// GetTrainingReadiness returns training readiness data for the given date.
func (c *Client) GetTrainingReadiness(date string) (json.RawMessage, error) {
	path := fmt.Sprintf("/metrics-service/metrics/trainingreadiness/%s", date)
	return c.doRequest("GET", path, nil)
}

// GetTrainingStatus returns aggregated training status for the given date.
func (c *Client) GetTrainingStatus(date string) (json.RawMessage, error) {
	path := fmt.Sprintf("/metrics-service/metrics/trainingstatus/aggregated/%s", date)
	return c.doRequest("GET", path, nil)
}

// GetDailySummary returns the daily summary for the given date.
func (c *Client) GetDailySummary(date string) (json.RawMessage, error) {
	path := fmt.Sprintf("/usersummary-service/usersummary/daily/?calendarDate=%s", date)
	return c.doRequest("GET", path, nil)
}

// GetHeartRate returns daily heart rate data for the given date.
func (c *Client) GetHeartRate(date string) (json.RawMessage, error) {
	path := fmt.Sprintf("/wellness-service/wellness/dailyHeartRate/?date=%s", date)
	return c.doRequest("GET", path, nil)
}

// GetSleep returns daily sleep data for the given date.
func (c *Client) GetSleep(date string) (json.RawMessage, error) {
	path := fmt.Sprintf("/wellness-service/wellness/dailySleepData/?date=%s", date)
	return c.doRequest("GET", path, nil)
}

// GetBodyComposition returns body composition data for a date range.
func (c *Client) GetBodyComposition(startDate, endDate string) (json.RawMessage, error) {
	path := fmt.Sprintf("/weight-service/weight/dateRange?startDate=%s&endDate=%s", startDate, endDate)
	return c.doRequest("GET", path, nil)
}

// GetHRV returns HRV data for the given date.
func (c *Client) GetHRV(date string) (json.RawMessage, error) {
	path := fmt.Sprintf("/hrv-service/hrv/%s", date)
	return c.doRequest("GET", path, nil)
}

// GetVO2Max returns VO2 Max data for a date range.
func (c *Client) GetVO2Max(startDate, endDate string) (json.RawMessage, error) {
	path := fmt.Sprintf("/metrics-service/metrics/maxmet/daily/%s/%s", startDate, endDate)
	return c.doRequest("GET", path, nil)
}
