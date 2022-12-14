// Package weather provides functionality for interacting with the weather.gov
// API service.
package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiBaseURL = "https://api.weather.gov"
	userAgent  = "eric.wohltman@gmail.com"
)

// HTTPClient is an interface abstraction for HTTP clients.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// APIClient is a client for making weather.gov API requests.
type APIClient struct {
	httpClient  HTTPClient
	forecastURL string
}

// NewAPIClient returns a new *APIClient.
func NewAPIClient(ctx context.Context, httpClient HTTPClient, lat, long string) (*APIClient, error) {
	pointsURL := fmt.Sprintf("%s/points/%s,%s", apiBaseURL, lat, long)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pointsURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pointsResponse := &PointsResponse{}

	err = json.Unmarshal(body, pointsResponse)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		httpClient:  httpClient,
		forecastURL: pointsResponse.Properties.Forecast,
	}, nil
}

// QueryForecast returns the response from a weather.gov API request.
func (apiClient *APIClient) QueryForecast(ctx context.Context) (*GridPointsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiClient.forecastURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)

	resp, err := apiClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	gridPointsResponse := &GridPointsResponse{}

	err = json.Unmarshal(body, gridPointsResponse)
	if err != nil {
		return nil, err
	}

	return gridPointsResponse, nil
}
