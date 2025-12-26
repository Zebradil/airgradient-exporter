# Running with Docker

To run the AirGradient Exporter using Docker:

```bash
docker build -t airgradient-exporter .
docker run -e AIRGRADIENT_MONITORS="192.168.1.50" -p 9112:9112 airgradient-exporter
```

## Using Pre-built Images

Pre-built Docker images are available from the [releases page](https://github.com/Zebradil/airgradient-exporter/releases).

You can run the exporter directly using a pre-built image:

```bash
docker run -e AIRGRADIENT_MONITORS="192.168.1.50" -p 9112:9112 ghcr.io/zebradil/airgradient-exporter:latest
```

## Configuration

Configure the exporter using environment variables:

| Variable               | Description                                                                     | Default  |
| ---------------------- | ------------------------------------------------------------------------------- | -------- |
| `AIRGRADIENT_MONITORS` | Comma-separated list of AirGradient monitor IPs, hostnames, or host:port.       | Required |
| `PORT`                 | Port to listen on.                                                              | `9112`   |
| `LOG_FORMAT`           | Log format: "text" or "json". Text format is colored when output is a terminal. | `text`   |

## Multiple Monitors

To monitor multiple AirGradient devices:

```bash
docker run -e AIRGRADIENT_MONITORS="192.168.1.50,192.168.1.51,192.168.1.52" -p 9112:9112 ghcr.io/zebradil/airgradient-exporter:latest
```
