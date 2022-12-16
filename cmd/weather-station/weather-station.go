//go:build js && wasm
// +build js,wasm

// Package main is the entry point to the program.
package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ewohltman/weather-station/internal/weather"
	"log"
	"net/http"
	"syscall/js"
)

const (
	lat  = "40.738620"
	long = "-74.064810"
)

func registerCallbacks() {
	js.Global().
		Get("document").
		Call("getElementById", "addButton").
		Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
			go func() {
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

				buffer := &bytes.Buffer{}

				for _, period := range forecast.Properties.Periods {
					_, _ = buffer.WriteString(fmt.Sprintf("%s: Temp: %d, %s\n", period.Name, period.Temperature, period.ShortForecast))
				}

				js.Global().Get("document").Call("getElementById", "output").Set("innerHTML", buffer.String())
				fmt.Println("button clicked")

			}()

			return nil
		}))
}

func main() {
	c := make(chan struct{}, 0)

	registerCallbacks()

	<-c
}
