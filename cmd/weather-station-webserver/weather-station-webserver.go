// Package main is the entry point to the webserver program.
package main

import (
	"embed"
	"log"
	"net/http"
)

var (
	//go:embed web
	web embed.FS
)

func runServer() {
	index, err := web.ReadFile("web/static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	css, err := web.ReadFile("web/static/theme.css")
	if err != nil {
		log.Fatal(err)
	}

	js, err := web.ReadFile("web/static/wasm_exec.js")
	if err != nil {
		log.Fatal(err)
	}

	wasm, err := web.ReadFile("web/app/weather-station.wasm")
	if err != nil {
		log.Fatal(err)
	}

	favicon, err := web.ReadFile("web/static/favicon.ico")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/theme.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")

		_, err = w.Write(css)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")

		_, err = w.Write(js)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/weather-station.wasm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")

		_, err = w.Write(wasm)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")

		_, err = w.Write(favicon)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		_, err = w.Write(index)
		if err != nil {
			log.Fatal(err)
		}
	})

	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	runServer()
}
