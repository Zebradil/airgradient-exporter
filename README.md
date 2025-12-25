# AirGradient Exporter

A Prometheus exporter for AirGradient monitors.

## Features

- Scrapes AirGradient monitors for current measurements.
- Scrapes AirGradient monitors for configuration.
- Supports multiple monitors.
- Exposes metrics in Prometheus format.

## Usage

### Command Line Flags

- `--help`: Show help message
- `--version`: Show version information (version, commit, and build date)

### Configuration

The exporter is configured via environment variables:

| Variable               | Description                                                                     | Default  |
| ---------------------- | ------------------------------------------------------------------------------- | -------- |
| `AIRGRADIENT_MONITORS` | Comma-separated list of AirGradient monitor IPs, hostnames, or host:port.       | Required |
| `PORT`                 | Port to listen on.                                                              | `9112`   |
| `LOG_FORMAT`           | Log format: "text" or "json". Text format is colored when output is a terminal. | `text`   |

### Running Locally

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

### Docker

```bash
docker build -t airgradient-exporter .
docker run -e AIRGRADIENT_MONITORS="192.168.1.50" -p 9112:9112 airgradient-exporter
```

### NixOS Module

The package includes a NixOS module for easy integration into NixOS configurations.

#### Using with Flakes

Add the flake to your `flake.nix`:

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    airgradient-exporter.url = "github:zebradil/airgradient-exporter";
  };

  outputs = { self, nixpkgs, airgradient-exporter, ... }: {
    nixosConfigurations.your-host = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        airgradient-exporter.nixosModules.default
        {
          services.airgradient-exporter = {
            enable = true;
            monitors = [ "192.168.1.50" "192.168.1.51" ];
            port = 9112;
            logFormat = "json";
            # Optional: override package if needed
            # package = airgradient-exporter.packages.x86_64-linux.default;
          };
        }
      ];
    };
  };
}
```

#### Using without Flakes

Import the module in your `configuration.nix`:

```nix
{ config, pkgs, ... }:

{
  imports = [
    /path/to/airgradient-exporter/nix/module.nix
  ];

  services.airgradient-exporter = {
    enable = true;
    monitors = [ "192.168.1.50" "192.168.1.51" ];
    port = 9112;
    logFormat = "json";
  };
}
```

The service will automatically start on boot when enabled. The `multi-user.target` ensures the service starts when the system reaches the normal multi-user state (standard boot).

### Debian/Ubuntu Package Installation

#### Installation

1. Download the `.deb` package from the [releases page](https://github.com/zebradil/airgradient-exporter/releases)
2. Install the package:

```bash
sudo dpkg -i airgradient-exporter_*.deb
```

The package installs:
- Binary: `/usr/bin/airgradient-exporter`
- Systemd unit: `/usr/lib/systemd/system/airgradient-exporter.service`

#### Configuration

The systemd unit file does not include environment variables by default. You need to configure them using a systemd override file.

**Option 1: Using `systemctl edit` (Recommended)**

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

**Option 2: Manual override file**

Create the override directory and file:

```bash
sudo mkdir -p /etc/systemd/system/airgradient-exporter.service.d
sudo nano /etc/systemd/system/airgradient-exporter.service.d/override.conf
```

Add the same configuration as above, then reload systemd and restart the service:

```bash
sudo systemctl daemon-reload
sudo systemctl restart airgradient-exporter
```

**Enable and start the service:**

```bash
sudo systemctl enable airgradient-exporter
sudo systemctl start airgradient-exporter
```

**Check service status:**

```bash
sudo systemctl status airgradient-exporter
```

**View logs:**

```bash
sudo journalctl -u airgradient-exporter -f
```

## Metrics

All metrics are prefixed with `airgradient_`.

### Status Metrics

- `airgradient_up`: Was the last scrape of the AirGradient monitor successful (1 = success, 0 = failure). Labels: `host`

### Particulate Matter (PM) Metrics

- `airgradient_pm01`: PM1.0 (ug/m3). Labels: `host`, `serialno`, `firmware`
- `airgradient_pm02`: PM2.5 (ug/m3). Labels: `host`, `serialno`, `firmware`
- `airgradient_pm10`: PM10 (ug/m3). Labels: `host`, `serialno`, `firmware`
- `airgradient_pm01_standard`: PM1.0 Standard (ug/m3). Labels: `host`, `serialno`, `firmware`
- `airgradient_pm02_standard`: PM2.5 Standard (ug/m3). Labels: `host`, `serialno`, `firmware`
- `airgradient_pm10_standard`: PM10 Standard (ug/m3). Labels: `host`, `serialno`, `firmware`
- `airgradient_pm02_compensated`: PM2.5 Compensated (ug/m3). Labels: `host`, `serialno`, `firmware`
- `airgradient_pm003_count`: PM0.3 Count. Labels: `host`, `serialno`, `firmware`
- `airgradient_pm005_count`: PM0.5 Count. Labels: `host`, `serialno`, `firmware`
- `airgradient_pm01_count`: PM1.0 Count. Labels: `host`, `serialno`, `firmware`
- `airgradient_pm02_count`: PM2.5 Count. Labels: `host`, `serialno`, `firmware`
- `airgradient_pm50_count`: PM5.0 Count. Labels: `host`, `serialno`, `firmware`
- `airgradient_pm10_count`: PM10 Count. Labels: `host`, `serialno`, `firmware`

### Temperature and Humidity Metrics

- `airgradient_atmp`: Ambient Temperature (Celsius). Labels: `host`, `serialno`, `firmware`
- `airgradient_atmp_compensated`: Ambient Temperature Compensated (Celsius). Labels: `host`, `serialno`, `firmware`
- `airgradient_rhum`: Relative Humidity (%). Labels: `host`, `serialno`, `firmware`
- `airgradient_rhum_compensated`: Relative Humidity Compensated (%). Labels: `host`, `serialno`, `firmware`

### Air Quality Metrics

- `airgradient_rco2`: CO2 (ppm). Labels: `host`, `serialno`, `firmware`
- `airgradient_tvoc_index`: TVOC Index. Labels: `host`, `serialno`, `firmware`
- `airgradient_tvoc_raw`: TVOC Raw. Labels: `host`, `serialno`, `firmware`
- `airgradient_nox_index`: NOx Index. Labels: `host`, `serialno`, `firmware`
- `airgradient_nox_raw`: NOx Raw. Labels: `host`, `serialno`, `firmware`

### System Metrics

- `airgradient_boot_uptime_raw`: Boot uptime raw value. Labels: `host`, `serialno`, `firmware`
- `airgradient_boot_count`: Boot Count. Labels: `host`, `serialno`, `firmware`
- `airgradient_wifi_signal_dbm`: WiFi Signal (dBm). Labels: `host`, `serialno`, `firmware`

### Configuration Metrics

- `airgradient_config_info`: Configuration Info (value is always 1). Labels: `host`, `model`, `country`, `pm_standard`, `led_bar_mode`, `temperature_unit`
- `airgradient_config_display_brightness`: Display Brightness. Labels: `host`, `model`
- `airgradient_config_led_bar_brightness`: LED Bar Brightness. Labels: `host`, `model`
