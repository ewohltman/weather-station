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
	htmlTagFmtImg = `<span><img src="%s" alt="Weather icon" width="86" height="86"></span>`
	cellFmt       = "table%d_%d"
)

func run(ctx context.Context, apiClient *weather.APIClient, document js.Value) {
	ticker := time.NewTicker(refreshPeriod)
	defer ticker.Stop()

	for {
		/*for i := 0; i < tableRows; i++ {
			for j := 0; j < tableColumns; j++ {
				document.Call(getElementById, cell(i, j)).Set(innerHTML, nbsp)
			}
		}*/

		hourlyForecast, err := apiClient.QueryForecast(ctx)
		if err != nil {
			log.Printf("Error querying hourly forecast: %s", err)

			<-ticker.C

			continue
		}

		updateNow(document, hourlyForecast)

		err = updateToday(document, hourlyForecast)
		if err != nil {
			log.Printf("Error querying hourly forecast: %s", err)

			<-ticker.C

			continue
		}

		updateFiveDay()

		<-ticker.C
	}
}

func updateNow(document js.Value, hourlyForecast *weather.GridPointsResponse) {
	period := hourlyForecast.Properties.Periods[0]
	data := []string{
		fmt.Sprintf(htmlTagFmtImg, period.Icon),
		mustParseTime(time.Parse(time.RFC3339, period.StartTime)).Format(time.Kitchen),
		fmt.Sprintf("Temperature: %d F", period.Temperature),
		fmt.Sprintf("Wind: %s %s", period.WindSpeed, period.WindDirection),
		fmt.Sprintf("Forecast: %s", period.ShortForecast),
	}

	formatted := strings.Join(data, " <br>\n")

	document.Call(getElementById, "nowCard").Set(innerHTML, formatted)
}

func updateToday(document js.Value, hourlyForecast *weather.GridPointsResponse) error {
	const (
		card    = "todayCard"
		rows    = 5
		columns = 5
	)

	table, err := format.ExecuteTemplate(card, rows, columns)
	if err != nil {
		return fmt.Errorf("error creating today table: %w", err)
	}

	document.Call(getElementById, "todayCard").Set(innerHTML, table.String())

	/*for i := 0; i < rows; i++ {
		period := hourlyForecast.Properties.Periods[i*4]
		data := []string{
			fmt.Sprintf(htmlTagFmtImg, period.Icon),
			mustParseTime(time.Parse(time.RFC3339, period.StartTime)).Format(time.Kitchen),
			fmt.Sprintf("Temperature: %d F", period.Temperature),
			fmt.Sprintf("Wind: %s %s", period.WindSpeed, period.WindDirection),
			fmt.Sprintf("Forecast: %s", period.ShortForecast),
		}

		for j := 0; j < columns; j++ {
			document.Call(getElementById, format.CellID(card, i, j)).Set(innerHTML, data[j])
		}
	}*/

	for i := 0; i < rows; i++ {
		period := hourlyForecast.Properties.Periods[i*4]
		data := []string{
			fmt.Sprintf(htmlTagFmtImg, period.Icon),
			mustParseTime(time.Parse(time.RFC3339, period.StartTime)).Format(time.Kitchen),
			fmt.Sprintf("Temperature: %d F", period.Temperature),
			fmt.Sprintf("Wind: %s %s", period.WindSpeed, period.WindDirection),
			fmt.Sprintf("Forecast: %s", period.ShortForecast),
		}

		for j := 0; j < columns; j++ {
			document.Call(getElementById, format.CellID(card, j, i)).Set(innerHTML, data[j])
		}
	}

	return nil
}

func updateFiveDay() {

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
