package strava

import (
	"time"

	gostrava "github.com/strava/go.strava"
)

func (c *Client) ListActivities(before, after time.Time, page, perPage int) ([]*gostrava.ActivitySummary, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	call := c.currentAthlete.ListActivities()
	if !before.IsZero() {
		call = call.Before(int(before.Unix()))
	}
	if !after.IsZero() {
		call = call.After(int(after.Unix()))
	}
	if page > 0 {
		call = call.Page(page)
	}
	if perPage > 0 {
		call = call.PerPage(perPage)
	}

	return call.Do()
}

func (c *Client) GetActivity(id int64) (*gostrava.ActivityDetailed, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}
	return c.activities.Get(id).IncludeAllEfforts().Do()
}

func (c *Client) GetActivityStreams(id int64, types []gostrava.StreamType) (*gostrava.StreamSet, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}
	return c.streams.Get(id, types).Do()
}

func (c *Client) GetActivityLaps(id int64) ([]*gostrava.LapEffortSummary, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}
	return c.activities.ListLaps(id).Do()
}

func (c *Client) GetActivityZones(id int64) ([]*gostrava.ZonesSummary, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}
	return c.activities.ListZones(id).Do()
}
