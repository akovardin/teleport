FROM golang:latest as builder

LABEL maintainer="Artem Kovadin <artem.kovardin@gmail.com>"

WORKDIR /go

COPY . /go/src/github.com/horechek/teleport

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o ./teleport github.com/horechek/teleport/cmd/teleport

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/teleport /bin/teleport

EXPOSE 8080
