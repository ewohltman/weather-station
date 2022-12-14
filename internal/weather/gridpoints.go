package weather

import (
	"encoding/json"
	"time"
)

// GridPointsResponse is the parsed result from the weather.gov API gridpoints
// endpoint.
type GridPointsResponse struct {
	ID         string                `json:"id,omitempty"`
	Type       string                `json:"type,omitempty"`
	Properties *GridPointsProperties `json:"properties,omitempty"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface for
// *GridPointsResponse.
func (response *GridPointsResponse) UnmarshalJSON(data []byte) error {
	tmp := &struct {
		ID             string                `json:"id,omitempty"`
		Type           string                `json:"type,omitempty"`
		Properties     *GridPointsProperties `json:"properties,omitempty"`
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

// GridPointsProperties contains the data from the weather.gov API gridpoints
// response.
type GridPointsProperties struct {
	Geometry          string               `json:"geometry,omitempty"`
	Units             string               `json:"units,omitempty"`
	ForecastGenerator string               `json:"forecastGenerator,omitempty"`
	GeneratedAt       time.Time            `json:"generatedAt,omitempty"`
	UpdateTime        time.Time            `json:"updateTime,omitempty"`
	Elevation         *GridPointsElevation `json:"elevation,omitempty"`
	Periods           []*GridPointsPeriods `json:"periods,omitempty"`
}

// GridPointsElevation contains elevation data from the weather.gov API
// gridpoints endpoint.
type GridPointsElevation struct {
	Value          float32 `json:"value,omitempty"`
	MaxValue       float32 `json:"maxValue,omitempty"`
	MinValue       float32 `json:"minValue,omitempty"`
	UnitCode       string  `json:"unitCode,omitempty"`
	QualityControl string  `json:"qualityControl,omitempty"`
}

// GridPointsPeriods contains time period data from the weather.gov API
// gridpoints endpoint.
type GridPointsPeriods struct {
	Number           int    `json:"number,omitempty"`
	Name             string `json:"name,omitempty"`
	StartTime        string `json:"startTime,omitempty"`
	EndTime          string `json:"endTime,omitempty"`
	IsDaytime        bool   `json:"isDaytime,omitempty"`
	Temperature      int    `json:"temperature,omitempty"`
	TemperatureUnit  string `json:"temperatureUnit,omitempty"`
	TemperatureTrend string `json:"temperatureTrend,omitempty"`
	WindSpeed        string `json:"windSpeed,omitempty"`
	WindDirection    string `json:"windDirection,omitempty"`
	Icon             string `json:"icon,omitempty"`
	ShortForecast    string `json:"shortForecast,omitempty"`
	DetailedForecast string `json:"detailedForecast,omitempty"`
}
