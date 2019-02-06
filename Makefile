dir = $(shell pwd)
.PHONY: all public clean

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ./bin/poster ./cmd/poster

macos:
	go build -o ./bin/poster ./cmd/poster

clean:
	rm -rf ./bin
