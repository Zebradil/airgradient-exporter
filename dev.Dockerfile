FROM golang:1.26@sha256:6e59e17a20f13c102e49d82dffab66d41cffe6c266b225d82794328aea859ab5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o bin/airgradient-exporter main.go

ENTRYPOINT ["/app/bin/airgradient-exporter"]
