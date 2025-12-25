FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o airgradient-exporter ./cmd/airgradient-exporter

FROM alpine:3.19
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/airgradient-exporter .

USER 1000:1000

ENV PORT=9112
EXPOSE 9112

ENTRYPOINT ["./airgradient-exporter"]
