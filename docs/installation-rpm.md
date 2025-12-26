# RPM Package Installation

## Installation

1. Download the `.rpm` package from the [releases page](https://github.com/Zebradil/airgradient-exporter/releases)
2. Install the package:

```bash
sudo rpm -i airgradient-exporter_*.rpm
```

Or on systems with `yum` or `dnf`:

```bash
sudo yum install airgradient-exporter_*.rpm
# or
sudo dnf install airgradient-exporter_*.rpm
```

The package installs:
- Binary: `/usr/bin/airgradient-exporter`
- Systemd unit: `/usr/lib/systemd/system/airgradient-exporter.service`
- Example override file: `/usr/share/doc/airgradient-exporter/override.conf.example`

## Configuration

The systemd unit file does not include environment variables by default. You need to configure them using a systemd override file.

### Option 1: Using `systemctl edit` (Recommended)

```bash
sudo systemctl edit airgradient-exporter
```

This opens an editor. Add the following configuration:

```ini
[Service]
Environment="AIRGRADIENT_MONITORS=192.168.1.50,192.168.1.51"
Environment="PORT=9112"
Environment="LOG_FORMAT=json"
```

### Option 2: Using the example file

Copy the example override file and edit it:

```bash
sudo mkdir -p /etc/systemd/system/airgradient-exporter.service.d
sudo cp /usr/share/doc/airgradient-exporter/override.conf.example /etc/systemd/system/airgradient-exporter.service.d/override.conf
sudo nano /etc/systemd/system/airgradient-exporter.service.d/override.conf
```

Edit the file with your actual configuration, then reload systemd and restart the service:

```bash
sudo systemctl daemon-reload
sudo systemctl restart airgradient-exporter
```

### Option 3: Manual override file

Create the override directory and file:

```bash
sudo mkdir -p /etc/systemd/system/airgradient-exporter.service.d
sudo nano /etc/systemd/system/airgradient-exporter.service.d/override.conf
```

Add the configuration:

```ini
[Service]
Environment="AIRGRADIENT_MONITORS=192.168.1.50,192.168.1.51"
Environment="PORT=9112"
Environment="LOG_FORMAT=json"
```

Then reload systemd and restart the service:

```bash
sudo systemctl daemon-reload
sudo systemctl restart airgradient-exporter
```

## Service Management

Enable and start the service:

```bash
sudo systemctl enable airgradient-exporter
sudo systemctl start airgradient-exporter
```

Check service status:

```bash
sudo systemctl status airgradient-exporter
```

View logs:

```bash
sudo journalctl -u airgradient-exporter -f
```

## Configuration Options

The exporter is configured via environment variables:

| Variable               | Description                                                                     | Default  |
| ---------------------- | ------------------------------------------------------------------------------- | -------- |
| `AIRGRADIENT_MONITORS` | Comma-separated list of AirGradient monitor IPs, hostnames, or host:port.       | Required |
| `PORT`                 | Port to listen on.                                                              | `9112`   |
| `LOG_FORMAT`           | Log format: "text" or "json". Text format is colored when output is a terminal. | `text`   |

