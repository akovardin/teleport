dir = $(shell pwd)
.PHONY: all public clean

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ./bin/teleport ./cmd/teleport

macos:
	go build -o ./bin/teleport ./cmd/teleport

clean:
	rm -rf ./bin
