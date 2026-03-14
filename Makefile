VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/nimsforestwebviewer ./cmd/nimsforestwebviewer

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/nimsforestwebviewer-linux-amd64 ./cmd/nimsforestwebviewer

docker:
	docker build --build-arg VERSION=$(VERSION) -t registry.nimsforest.com/nimsforestwebviewer:latest .

.PHONY: build build-linux docker
