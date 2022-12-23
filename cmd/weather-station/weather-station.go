//go:build js && wasm

// Package main is the entry point to the Wasm program.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"syscall/js"
	"time"

	"github.com/ewohltman/weather-station/internal/format"
	"github.com/ewohltman/weather-station/internal/weather"
)

const (
	lat  = "40.7386"
	long = "-74.0648"
)

const refreshPeriod = time.Minute

const (
	getElementById = "getElementById"
	innerHTML      = "innerHTML"
)

const (
	htmlTagFmtImg = `<img src="%s" alt="Weather icon" width="86" height="86">`
)

func run(ctx context.Context, apiClient *weather.APIClient, document js.Value) {
	ticker := time.NewTicker(refreshPeriod)
	defer ticker.Stop()

	for {
		document.Call(getElementById, "timestamp").Set(innerHTML, time.Now().Format(time.Kitchen))

		hourlyForecast, err := apiClient.QueryHourlyForecast(ctx)
		if err != nil {
			log.Printf("Error querying hourly forecast: %s", err)

			<-ticker.C

			continue
		}

		dailyForecast, err := apiClient.QueryDailyForecast(ctx)
		if err != nil {
			log.Printf("Error querying hourly forecast: %s", err)

			<-ticker.C

			continue
		}

		updateNow(document, hourlyForecast)

		err = updateToday(document, hourlyForecast)
		if err != nil {
			log.Printf("Error updating today card: %s", err)
		}

		err = updateFiveDay(document, dailyForecast)
		if err != nil {
			log.Printf("Error updating five day card: %s", err)
		}

		<-ticker.C
	}
}

func updateNow(document js.Value, hourlyForecast *weather.GridPointsResponse) {
	period := hourlyForecast.Properties.Periods[0]
	data := []string{
		fmt.Sprintf("%d F", period.Temperature),
		fmt.Sprintf("Wind %s %s", period.WindSpeed, period.WindDirection),
		fmt.Sprintf("%s", period.ShortForecast),
	}

	formatted := fmt.Sprintf("%s\n%s",
		fmt.Sprintf(htmlTagFmtImg, period.Icon),
		strings.Join(data, " <br>\n"),
	)

	document.Call(getElementById, "nowCard").Set(innerHTML, formatted)
}

func updateToday(document js.Value, hourlyForecast *weather.GridPointsResponse) error {
	const (
		card    = "todayCard"
		rows    = 3
		columns = 5
	)

	table, err := format.ExecuteTemplate(card, rows, columns)
	if err != nil {
		return fmt.Errorf("error creating today table: %w", err)
	}

	document.Call(getElementById, card).Set(innerHTML, table.String())

	for c := 0; c < columns; c++ {
		period := hourlyForecast.Properties.Periods[c*4]
		data := []string{
			fmt.Sprintf(htmlTagFmtImg, period.Icon),
			mustParseTime(time.Parse(time.RFC3339, period.StartTime)).Format(time.Kitchen),
			fmt.Sprintf("%d F", period.Temperature),
			fmt.Sprintf("%s", period.ShortForecast),
		}

		for r := 0; r < rows; r++ {
			document.Call(getElementById, format.CellID(card, r, c)).Set(innerHTML, data[r])
		}
	}

	return nil
}

func updateFiveDay(document js.Value, dailyForecast *weather.GridPointsResponse) error {
	const (
		card    = "fiveDayCard"
		rows    = 3
		columns = 5
	)

	table, err := format.ExecuteTemplate(card, rows, columns)
	if err != nil {
		return fmt.Errorf("error creating today table: %w", err)
	}

	document.Call(getElementById, card).Set(innerHTML, table.String())

	for c := 0; c < columns; c++ {
		period := dailyForecast.Properties.Periods[c*2]
		data := []string{
			fmt.Sprintf(htmlTagFmtImg, period.Icon),
			period.Name,
			fmt.Sprintf("%d F", period.Temperature),
		}

		for r := 0; r < rows; r++ {
			document.Call(getElementById, format.CellID(card, r, c)).Set(innerHTML, data[r])
		}
	}

	return nil
}

func mustParseTime(t time.Time, err error) time.Time {
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func main() {
	ctx := context.Background()
	httpClient := &http.Client{}
	document := js.Global().Get("document")

	apiClient, err := weather.NewAPIClient(ctx, httpClient, lat, long)
	if err != nil {
		log.Fatalf("Error creating new API client: %s", err)
	}

	run(ctx, apiClient, document)
}
