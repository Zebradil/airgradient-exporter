FROM golang:1.26@sha256:0f063af2d465d8dcae54cce04278ada488b96f77b42449c8d071e47d016cc65a

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o bin/airgradient-exporter main.go

ENTRYPOINT ["/app/bin/airgradient-exporter"]
