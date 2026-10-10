FROM golang:1.26@sha256:d7722066f0b60ceccb6c0643cbed1f5f9e15506ac237146c95333504d7805d89

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o bin/airgradient-exporter main.go

ENTRYPOINT ["/app/bin/airgradient-exporter"]
