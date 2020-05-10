FROM golang:1.14
MAINTAINER Bruno Farias "brunofarias1@gmail.com"

WORKDIR /app

COPY ./ /app

RUN go mod download

ENTRYPOINT go run main.go