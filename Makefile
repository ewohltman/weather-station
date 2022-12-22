.PHONY: tidy fmt vet lint test test-report build

project=weather-station
projectDirectories=$(shell go list -f "{{.Dir}}" ./...)

tidy:
	go mod tidy

fmt: tidy
	@echo 'gofmt -s -w'
	@gofmt -s -w ${projectDirectories}
	@echo 'goimports -w'
	@goimports -local github.com/ewohltman/${project} -w ${projectDirectories}

vet: fmt
	go vet ./...

lint: vet
	golangci-lint run ./...

test:
	go test -v -race -bench=. -coverprofile=coverage.out -covermode atomic ./...
	@echo "all tests passed"

test-report:
	@make test | grep -E -A 1 '^coverage|^Benchmark' | grep -E -v '^PASS'
	@echo "all tests passed"

build:
	CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -o cmd/${project}-webserver/web/app/${project}.wasm cmd/${project}/${project}.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o build/package/${project}-webserver/${project}-webserver.arm cmd/${project}-webserver/${project}-webserver.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/package/${project}-webserver/${project}-webserver.amd64 cmd/${project}-webserver/${project}-webserver.go

run: build
	-/usr/bin/google-chrome-stable --chrome --kiosk http://localhost:8080 > /dev/null 2>&1
	cd build/package/weather-station-webserver && ./weather-station-webserver.amd64
