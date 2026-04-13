package garmin

import (
	"encoding/json"
	"fmt"
)

// ListActivities returns recent activities as raw JSON.
func (c *Client) ListActivities(count int) (json.RawMessage, error) {
	path := fmt.Sprintf("/activitylist-service/activities/search/activities?start=0&limit=%d", count)
	return c.doRequest("GET", path, nil)
}

// GetActivity returns a single activity by ID as raw JSON.
func (c *Client) GetActivity(id int64) (json.RawMessage, error) {
	path := fmt.Sprintf("/activity-service/activity/%d", id)
	return c.doRequest("GET", path, nil)
}
