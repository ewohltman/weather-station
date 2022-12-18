//go:build js && wasm

// Package main is the entry point to the Wasm program.
package main

import (
	"context"
	"fmt"
	"github.com/ewohltman/weather-station/internal/weather"
	"log"
	"net/http"
	"strconv"
	"syscall/js"
	"time"
)

const (
	lat  = "40.738620"
	long = "-74.064810"
)

const refreshPeriod = time.Minute

const (
	tableRows    = 14
	tableColumns = 4
)

const (
	getElementById = "getElementById"
	innerHTML      = "innerHTML"
	nbsp           = "&nbsp;"
)

const htmlTagFmtImg = `<img src="%s" alt="%s" width="86" height="86">`

func populate(ctx context.Context, apiClient *weather.APIClient, document js.Value) {
	ticker := time.NewTicker(refreshPeriod)
	defer ticker.Stop()

	for {
		document.Call(getElementById, "updated").Set(innerHTML, time.Now().String())

		for i := 0; i < tableRows; i++ {
			for j := 0; j < tableColumns; j++ {
				tableElement := fmt.Sprintf("table%d_%d", i, j)

				document.Call(getElementById, tableElement).Set(innerHTML, nbsp)
			}
		}

		forecast, err := apiClient.QueryForecast(ctx)
		if err != nil {
			log.Fatalf("Error querying forecast: %s", err)
		}

		for i := 0; i < tableRows; i++ {
			period := forecast.Properties.Periods[i]
			rowData := []string{
				period.Name,
				fmt.Sprintf(htmlTagFmtImg,
					period.Icon,
					"Weather icon",
				),
				strconv.Itoa(period.Temperature),
				period.ShortForecast,
			}

			fmt.Println(fmt.Sprintf(htmlTagFmtImg, period.Icon, "Weather icon"))

			for j := 0; j < tableColumns; j++ {
				tableElement := fmt.Sprintf("table%d_%d", i, j)

				document.Call(getElementById, tableElement).Set(innerHTML, rowData[j])
			}
		}

		<-ticker.C
	}
}

func main() {
	ctx := context.Background()
	httpClient := &http.Client{}
	document := js.Global().Get("document")

	apiClient, err := weather.NewAPIClient(ctx, httpClient, lat, long)
	if err != nil {
		log.Fatalf("Error creating new API client: %s", err)
	}

	populate(ctx, apiClient, document)

	<-ctx.Done()
}
