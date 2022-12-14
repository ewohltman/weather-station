// Package weather provides functionality for interacting with the weather.gov
// API service.
package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiURL    = "https://api.weather.gov/points"
	userAgent = ""
)

// HTTPClient is an interface abstraction for HTTP clients.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Response is the parsed result from the weather.gov API.
type Response struct {
	ID         string      `json:"id,omitempty"`
	Type       string      `json:"type,omitempty"`
	Properties *Properties `json:"properties,omitempty"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface for *Response.
func (response *Response) UnmarshalJSON(data []byte) error {
	tmp := &struct {
		ID         string         `json:"id,omitempty"`
		Type       string         `json:"type,omitempty"`
		Properties *Properties    `json:"properties,omitempty"`
		Error      *ResponseError `json:",inline"`
	}{}

	err := json.Unmarshal(data, tmp)
	if err != nil {
		return err
	}

	if tmp.Error != nil {
		return tmp.Error
	}

	response.ID = tmp.ID
	response.Type = tmp.Type
	response.Properties = tmp.Properties

	return nil
}

// Properties contains the data from a weather.gov API response.
type Properties struct {
	Geometry            string `json:"geometry,omitempty"`
	ID                  string `json:"@id,omitempty"`
	Type                string `json:"@type,omitempty"`
	Cwa                 string `json:"cwa,omitempty"`
	ForecastOffice      string `json:"forecastOffice,omitempty"`
	GridID              string `json:"gridId,omitempty"`
	GridX               int    `json:"gridX,omitempty"`
	GridY               int    `json:"gridY,omitempty"`
	Forecast            string `json:"forecast,omitempty"`
	ForecastHourly      string `json:"forecastHourly,omitempty"`
	ForecastGridData    string `json:"forecastGridData,omitempty"`
	ObservationStations string `json:"observationStations,omitempty"`
	ForecastZone        string `json:"forecastZone,omitempty"`
	County              string `json:"county,omitempty"`
	FireWeatherZone     string `json:"fireWeatherZone,omitempty"`
	TimeZone            string `json:"timeZone,omitempty"`
	RadarStation        string `json:"radarStation,omitempty"`
}

// ResponseError is the parsed error result from the weather.gov API.
type ResponseError struct {
	Type          string `json:"type,omitempty"`
	Title         string `json:"title,omitempty"`
	Status        int    `json:"status,omitempty"`
	Detail        string `json:"detail,omitempty"`
	Instance      string `json:"instance,omitempty"`
	CorrelationID string `json:"correlationId,omitempty"`
}

// Error satisfies the error interface for *ResponseError.
func (err *ResponseError) Error() string {
	return fmt.Sprintf("[%d]: %s: %s", err.Status, err.Title, err.Detail)
}

// APIClient is a client for making weather.gov API requests.
type APIClient struct {
	httpClient HTTPClient
	zipCode    string
}

// NewAPIClient returns a new *APIClient.
func NewAPIClient(httpClient HTTPClient, zipCode string) *APIClient {
	return &APIClient{
		httpClient: httpClient,
		zipCode:    zipCode,
	}
}

// Query returns the response from a weather.gov API request.
func (apiClient *APIClient) Query() (*Response, error) {
	return nil, nil
}
