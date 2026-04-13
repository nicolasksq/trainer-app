package strava

import (
	"fmt"
	"os"

	gostrava "github.com/strava/go.strava"
)

type Client struct {
	stravaClient   *gostrava.Client
	token          *Token
	clientID       string
	clientSecret   string
	currentAthlete *gostrava.CurrentAthleteService
	athletes       *gostrava.AthletesService
	activities     *gostrava.ActivitiesService
	streams        *gostrava.ActivityStreamsService
}

func NewClient() (*Client, error) {
	clientID := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("STRAVA_CLIENT_ID and STRAVA_CLIENT_SECRET environment variables are required")
	}

	token, err := LoadToken()
	if err != nil {
		return nil, fmt.Errorf("loading token (run 'strava-mcp auth' first): %w", err)
	}

	if token.Expired() {
		token, err = RefreshAccessToken(clientID, clientSecret, token)
		if err != nil {
			return nil, fmt.Errorf("refreshing expired token: %w", err)
		}
	}

	sc := gostrava.NewClient(token.AccessToken)

	return &Client{
		stravaClient:   sc,
		token:          token,
		clientID:       clientID,
		clientSecret:   clientSecret,
		currentAthlete: gostrava.NewCurrentAthleteService(sc),
		athletes:       gostrava.NewAthletesService(sc),
		activities:     gostrava.NewActivitiesService(sc),
		streams:        gostrava.NewActivityStreamsService(sc),
	}, nil
}

func (c *Client) ensureValidToken() error {
	if !c.token.Expired() {
		return nil
	}

	token, err := RefreshAccessToken(c.clientID, c.clientSecret, c.token)
	if err != nil {
		return err
	}

	c.token = token
	sc := gostrava.NewClient(token.AccessToken)
	c.stravaClient = sc
	c.currentAthlete = gostrava.NewCurrentAthleteService(sc)
	c.athletes = gostrava.NewAthletesService(sc)
	c.activities = gostrava.NewActivitiesService(sc)
	c.streams = gostrava.NewActivityStreamsService(sc)

	return nil
}

func (c *Client) GetCurrentAthlete() (*gostrava.AthleteDetailed, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}
	return c.currentAthlete.Get().Do()
}

func (c *Client) GetAthleteStats(athleteID int64) (*gostrava.AthleteStats, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}
	return c.athletes.Stats(athleteID).Do()
}
