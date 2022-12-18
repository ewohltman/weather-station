package weather

import (
	"encoding/json"
)

// PointsResponse is the parsed result from the weather.gov API points
// endpoint.
type PointsResponse struct {
	ID         string            `json:"id,omitempty"`
	Type       string            `json:"type,omitempty"`
	Properties *PointsProperties `json:"properties,omitempty"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface for *PointsResponse.
func (response *PointsResponse) UnmarshalJSON(data []byte) error {
	tmp := &struct {
		ID             string            `json:"id,omitempty"`
		Type           string            `json:"type,omitempty"`
		Properties     *PointsProperties `json:"properties,omitempty"`
		*ResponseError `json:",inline"`
	}{}

	err := json.Unmarshal(data, tmp)
	if err != nil {
		return err
	}

	if tmp.ResponseError != nil {
		return tmp.ResponseError
	}

	response.ID = tmp.ID
	response.Type = tmp.Type
	response.Properties = tmp.Properties

	return nil
}

// PointsProperties contains the data from the weather.gov API points response.
type PointsProperties struct {
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
