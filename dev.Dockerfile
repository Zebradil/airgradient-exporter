FROM golang:1.26@sha256:9d2f36f06329b2a141b9db99ffa32765cf695ee57b813ca29e245e8670bcbfff

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o bin/airgradient-exporter main.go

ENTRYPOINT ["/app/bin/airgradient-exporter"]
