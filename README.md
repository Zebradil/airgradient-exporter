# AirGradient Exporter

A Prometheus exporter for AirGradient monitors.

## Features

- Scrapes AirGradient monitors for current measurements.
- Scrapes AirGradient monitors for configuration.
- Supports multiple monitors.
- Exposes metrics in Prometheus format.

## Usage

### Configuration

The exporter is configured via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `AIRGRADIENT_MONITORS` | Comma-separated list of AirGradient monitor IPs, hostnames, or host:port. | Required |
| `PORT` | Port to listen on. | `9112` |
| `LOG_FORMAT` | Log format: "text" or "json". Text format is colored when output is a terminal. | `text` |

### Running Locally

```bash
export AIRGRADIENT_MONITORS="192.168.1.50,192.168.1.51"
go run ./cmd
```

You can also specify port numbers if your monitors are running on non-default ports:

```bash
export AIRGRADIENT_MONITORS="192.168.1.50:8080,192.168.1.51:8080"
go run ./cmd
```

To use JSON logging format:

```bash
export AIRGRADIENT_MONITORS="192.168.1.50"
export LOG_FORMAT="json"
go run ./cmd
```

Visit `http://localhost:9112/metrics` to see the metrics.

### Docker

```bash
docker build -t airgradient-exporter .
docker run -e AIRGRADIENT_MONITORS="192.168.1.50" -p 9112:9112 airgradient-exporter
```

## Metrics

metrics are prefixed with `airgradient_`.

- `airgradient_pm01`: PM1.0 (ug/m3)
- `airgradient_pm02`: PM2.5 (ug/m3)
- `airgradient_pm10`: PM10 (ug/m3)
- `airgradient_rco2`: CO2 (ppm)
- `airgradient_atmp`: Ambient Temperature (Celsius)
- `airgradient_rhum`: Relative Humidity (%)
- ... and many more.

Configuration metrics:

- `airgradient_config_info`: Metric with labels containing configuration details (country, model, etc.) with value 1.
- `airgradient_config_display_brightness`
- `airgradient_config_led_bar_brightness`
