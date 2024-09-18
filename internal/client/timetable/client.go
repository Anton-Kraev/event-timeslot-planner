package timetable

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule/timetable"
)

const (
	api = "%s/api"
)

type Client struct {
	host       string
	httpClient *http.Client
}

func NewClient(host string, httpClient *http.Client) Client {
	return Client{host: host, httpClient: httpClient}
}

func (c Client) doHTTP(
	ctx context.Context, method string, url string, queryParams map[string]string,
) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	for k, v := range queryParams {
		req.URL.Query().Add(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Error: %w", err)
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respB, nil
}

func (c Client) getID(ctx context.Context, url string, queryParams map[string]string) (uint64, error) {
	respB, err := c.doHTTP(ctx, http.MethodGet, url, queryParams)
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseUint(string(respB), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response body: %w", err)
	}

	return id, nil
}

func (c Client) getEvents(ctx context.Context, url string) ([]timetable.Event, error) {
	respB, err := c.doHTTP(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var events []timetable.Event

	if err = json.Unmarshal(respB, &events); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return events, nil
}
