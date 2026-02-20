.PHONY: build test lint clean install

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/repokit ./cmd/repokit

test:
	go test -race -coverprofile=coverage.txt ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/ coverage.txt

install:
	go install $(LDFLAGS) ./cmd/repokit
