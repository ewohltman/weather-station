package weather_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/ewohltman/weather-station/internal/weather"
)

type mockHTTPClient struct {
	PointsResponse     []byte
	GridPointsResponse []byte
}

func (httpClient *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	switch {
	case strings.Contains(req.URL.Path, "gridpoints"):
		return &http.Response{Body: io.NopCloser(bytes.NewReader(httpClient.GridPointsResponse))}, nil
	case strings.Contains(req.URL.Path, "points"):
		return &http.Response{Body: io.NopCloser(bytes.NewReader(httpClient.PointsResponse))}, nil
	default:
		return &http.Response{Body: io.NopCloser(bytes.NewReader(newBadResponse()))}, nil
	}
}

func TestNewAPIClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	httpClient := &mockHTTPClient{PointsResponse: newPointsResponse()}

	_, err := weather.NewAPIClient(ctx, httpClient, "", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAPIClient_QueryForecast(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		httpClient    weather.HTTPClient
		expectedError bool
	}

	testCases := []testCase{
		{
			name: "good response",
			httpClient: &mockHTTPClient{
				PointsResponse:     newPointsResponse(),
				GridPointsResponse: newGridPointsResponse(),
			},
			expectedError: false,
		},
		{
			name: "bad response",
			httpClient: &mockHTTPClient{
				PointsResponse:     newPointsResponse(),
				GridPointsResponse: newBadResponse(),
			},
			expectedError: true,
		},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			apiClient, err := weather.NewAPIClient(ctx, tc.httpClient, "", "")
			if err != nil {
				t.Fatal(err)
			}

			_, err = apiClient.QueryForecast(ctx)
			if (err != nil) != tc.expectedError {
				t.Error("unexpected result")
			}
		})
	}
}

func newPointsResponse() []byte {
	return newResponse("testdata/pointsResponse.json")
}

func newGridPointsResponse() []byte {
	return newResponse("testdata/gridPointsResponse.json")
}

func newBadResponse() []byte {
	return newResponse("testdata/badResponse.json")
}

func newResponse(filename string) []byte {
	contents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return contents
}
