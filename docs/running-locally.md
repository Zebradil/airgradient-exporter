# Running Locally

To run the AirGradient Exporter locally using Go:

```bash
export AIRGRADIENT_MONITORS="192.168.1.50,192.168.1.51"
go run .
```

You can also specify port numbers if your monitors are running on non-default ports:

```bash
export AIRGRADIENT_MONITORS="192.168.1.50:8080,192.168.1.51:8080"
go run .
```

To use JSON logging format:

```bash
export AIRGRADIENT_MONITORS="192.168.1.50"
export LOG_FORMAT="json"
go run .
```

Visit `http://localhost:9112/metrics` to see the metrics.

## Configuration

The exporter is configured via environment variables:

| Variable               | Description                                                                     | Default  |
| ---------------------- | ------------------------------------------------------------------------------- | -------- |
| `AIRGRADIENT_MONITORS` | Comma-separated list of AirGradient monitor IPs, hostnames, or host:port.       | Required |
| `PORT`                 | Port to listen on.                                                              | `9112`   |
| `LOG_FORMAT`           | Log format: "text" or "json". Text format is colored when output is a terminal. | `text`   |
