package weather_test

import (
	"encoding/json"
	"testing"

	"github.com/ewohltman/weather-station/internal/weather"
)

func TestPointsResponse_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		body          []byte
		expectedError bool
	}

	testCases := []testCase{
		{
			name:          "good response",
			body:          newPointsResponse(),
			expectedError: false,
		},
		{
			name:          "bad response",
			body:          newBadResponse(),
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := json.Unmarshal(tc.body, &weather.PointsResponse{})
			if (err != nil) != tc.expectedError {
				t.Error("unexpected result")
			}
		})
	}
}
