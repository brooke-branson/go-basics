package mlbstats

import (
	"context"
	"encoding/json"
	"fmt"
	"go-basics/models"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const defaultBaseURL = "https://statsapi.mlb.com/api/v1"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) FetchPlayer(ctx context.Context, playerID string, season int) (models.Player, error) {
	endpoint, err := c.statsURL(playerID, season)
	if err != nil {
		return models.Player{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return models.Player{}, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.Player{}, fmt.Errorf("call mlb api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Player{}, fmt.Errorf("mlb api returned status %d", resp.StatusCode)
	}

	var payload statsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return models.Player{}, fmt.Errorf("decode mlb response: %w", err)
	}

	return mapStatsToPlayer(playerID, payload)
}

func (c *Client) statsURL(playerID string, season int) (string, error) {
	u, err := url.Parse(fmt.Sprintf("%s/people/%s/stats", c.baseURL, playerID))
	if err != nil {
		return "", fmt.Errorf("parse base url: %w", err)
	}

	q := u.Query()
	q.Set("stats", "season")
	q.Set("group", "hitting")
	q.Set("season", strconv.Itoa(season))
	u.RawQuery = q.Encode()

	return u.String(), nil
}
