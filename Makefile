.PHONY: build clean install test run-relay

build:
	go build -o zipline

install:
	go install

clean:
	rm -f zipline zipline-* *.exe

test:
	go test -v ./...

run-relay:
	go run . relay

build-all:
	GOOS=linux GOARCH=amd64 go build -o dist/zipline-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o dist/zipline-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -o dist/zipline-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o dist/zipline-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -o dist/zipline-windows-amd64.exe
