dir = $(shell pwd)
.PHONY: all public clean statik

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o ./bin/teleport-linux ./cmd/teleport

macos:
	go build -o ./bin/teleport-macos ./cmd/teleport

statik:
	statik -src=./web/dist/web

clean:
	rm -rf ./bin
	rm -rf ./statik
