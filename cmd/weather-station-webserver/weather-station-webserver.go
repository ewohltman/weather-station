// Package main is the entry point to the webserver program.
package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	tableRows    = 14
	tableColumns = 4
)

var (
	//go:embed web
	web embed.FS
)

// Table is used in an HTML template.
type Table struct {
	Rows []Row
}

// Row is a dimension of a Table.
type Row struct {
	Columns []Column
}

// Column is a dimension of a Table.
type Column struct {
	ID string
}

func runServer(bufIndex *bytes.Buffer) {
	index := bufIndex.Bytes()

	css, err := web.ReadFile("web/static/weather-station.css")
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

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	http.HandleFunc("/weather-station.css", func(w http.ResponseWriter, r *http.Request) {
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
	index := &bytes.Buffer{}
	table := Table{Rows: make([]Row, tableRows)}

	for i := 0; i < tableRows; i++ {
		table.Rows[i].Columns = make([]Column, tableColumns)

		for j := 0; j < tableColumns; j++ {
			table.Rows[i].Columns[j].ID = fmt.Sprintf("table%d_%d", i, j)
		}
	}

	tmpl, err := template.ParseFS(web, "web/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(index, table)
	if err != nil {
		log.Fatal(err)
	}

	runServer(index)
}
