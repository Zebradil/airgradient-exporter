FROM golang:1.27@sha256:ec47ea3725fea0cc88f6de94919aef19f40f12fc1ec0e2d1b3426dbf15ecca04

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o bin/airgradient-exporter main.go

ENTRYPOINT ["/app/bin/airgradient-exporter"]
