FROM golang:1.26@sha256:978e0a4e3d953662577d784471db27742f543ab09eead06839a838188c73314b

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o bin/airgradient-exporter main.go

ENTRYPOINT ["/app/bin/airgradient-exporter"]
