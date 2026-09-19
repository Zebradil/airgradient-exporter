FROM golang:1.26@sha256:ef46d02b02ed68aa94b6a49ba5d3ec62003b8d5ce4b7497fe83e056f2e8fd378

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o bin/airgradient-exporter main.go

ENTRYPOINT ["/app/bin/airgradient-exporter"]
