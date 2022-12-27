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

const (
	refreshPeriod = time.Minute
	dateFormat    = "Mon 1/2"
	timeFormat    = "3:04PM"
)

const (
	getElementById = "getElementById"
	innerHTML      = "innerHTML"
)

const (
	idDocument    = "document"
	idDate        = "date"
	idTime        = "time"
	idFooter      = "footer"
	idNowCard     = "nowCard"
	idTodayCard   = "todayCard"
	idFiveDayCard = "fiveDayCard"
)

const htmlTagFmtImg = `<img src="%s" alt="Weather icon" width="86" height="86">`

func run(ctx context.Context, apiClient *weather.APIClient, document js.Value) {
	ticker := time.NewTicker(refreshPeriod)
	defer ticker.Stop()

	for {
		timeNow := time.Now()

		document.Call(getElementById, idFooter).Set(innerHTML, "")
		document.Call(getElementById, idDate).Set(innerHTML, timeNow.Format(dateFormat))
		document.Call(getElementById, idTime).Set(innerHTML, timeNow.Format(timeFormat))

		hourlyForecast, err := apiClient.QueryHourlyForecast(ctx)
		if err != nil {
			appendFooter(document, fmt.Sprintf("Error querying hourly forecast: %s", err))
		}

		dailyForecast, err := apiClient.QueryDailyForecast(ctx)
		if err != nil {
			appendFooter(document, fmt.Sprintf("Error querying daily forecast: %s", err))
		}

		if hourlyForecast == nil || dailyForecast == nil {
			<-ticker.C

			continue
		}

		updateNow(document, hourlyForecast)

		err = updateToday(document, hourlyForecast)
		if err != nil {
			appendFooter(document, fmt.Sprintf("Error updating today card: %s", err))
		}

		err = updateFiveDay(document, dailyForecast)
		if err != nil {
			appendFooter(document, fmt.Sprintf("Error updating five day card: %s", err))
		}

		<-ticker.C
	}
}

func appendFooter(document js.Value, msg string) {
	current := document.Call(getElementById, idFooter).Get(innerHTML).String()

	var updated string

	if current == "" {
		updated = fmt.Sprintf("%s", msg)
	} else {
		updated = fmt.Sprintf("%s <br>\n%s", current, msg)
	}

	document.Call(getElementById, idFooter).Set(innerHTML, updated)
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

	document.Call(getElementById, idNowCard).Set(innerHTML, formatted)
}

func updateToday(document js.Value, hourlyForecast *weather.GridPointsResponse) error {
	const (
		card    = idTodayCard
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
		card    = idFiveDayCard
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
	document := js.Global().Get(idDocument)

	apiClient, err := weather.NewAPIClient(ctx, httpClient, lat, long)
	if err != nil {
		log.Fatalf("Error creating new API client: %s", err)
	}

	run(ctx, apiClient, document)
}
