// Package main is the entry point to the program.
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ewohltman/weather-station/internal/weather"
)

const (
	lat  = "40.738620"
	long = "-74.064810"
)

func main() {
	ctx := context.Background()
	httpClient := &http.Client{}

	apiClient, err := weather.NewAPIClient(ctx, httpClient, lat, long)
	if err != nil {
		log.Fatalf("Error creating new API client: %s", err)
	}

	forecast, err := apiClient.QueryForecast(ctx)
	if err != nil {
		log.Fatalf("Error querying forecast: %s", err)
	}

	for _, period := range forecast.Properties.Periods {
		log.Printf("%s: Temp: %d, %s\n", period.Name, period.Temperature, period.ShortForecast)
	}
}
