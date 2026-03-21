FROM golang:1.25@sha256:dfae680962532eeea67ab297f1166c2c4e686edb9a8f05f9d02d96fc9191833e

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o bin/airgradient-exporter main.go

ENTRYPOINT ["/app/bin/airgradient-exporter"]
