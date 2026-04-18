BINARY = matelabdev-watcher-agent
VERSION = 1.0.0

.PHONY: build build-all test clean

build:
	go build -ldflags="-s -w" -o $(BINARY) .

build-all:
	GOOS=linux GOARCH=amd64  go build -ldflags="-s -w" -o dist/$(BINARY)-linux-amd64  .
	GOOS=linux GOARCH=arm64  go build -ldflags="-s -w" -o dist/$(BINARY)-linux-arm64  .
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o dist/$(BINARY)-linux-armv7 .

test:
	go test ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/