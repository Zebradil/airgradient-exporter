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

## Installation and Running

Choose the method that best fits your environment:

- **[Running Locally](docs/running-locally.md)** - Run directly with Go
- **[Docker](docs/docker.md)** - Run using Docker containers
- **[NixOS](docs/nixos.md)** - Install and configure with NixOS module
- **[Debian/Ubuntu](docs/installation-debian-ubuntu.md)** - Install using `.deb` packages
- **[RPM](docs/installation-rpm.md)** - Install using `.rpm` packages (RHEL, CentOS, Fedora, etc.)
- **[APK](docs/installation-apk.md)** - Install using `.apk` packages (Alpine Linux)

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
