# Documentation Index

This directory contains detailed documentation for the AirGradient Exporter.

## Installation and Running

- **[Running Locally](running-locally.md)** - Instructions for running the exporter directly with Go
- **[Docker](docker.md)** - Guide for running the exporter in Docker containers
- **[NixOS](nixos.md)** - NixOS module configuration and usage
- **[Debian/Ubuntu](installation-debian-ubuntu.md)** - Installation guide for Debian and Ubuntu systems using `.deb` packages

## Quick Start

1. Choose your preferred installation method from the links above
2. Configure the `AIRGRADIENT_MONITORS` environment variable with your monitor IPs
3. Access metrics at `http://localhost:9112/metrics`

## Configuration

All methods support the same environment variables:

| Variable               | Description                                                                     | Default  |
| ---------------------- | ------------------------------------------------------------------------------- | -------- |
| `AIRGRADIENT_MONITORS` | Comma-separated list of AirGradient monitor IPs, hostnames, or host:port.       | Required |
| `PORT`                 | Port to listen on.                                                              | `9112`   |
| `LOG_FORMAT`           | Log format: "text" or "json". Text format is colored when output is a terminal. | `text`   |

## Additional Resources

- [Main README](../README.md) - Overview and metrics documentation
- [Releases](https://github.com/Zebradil/airgradient-exporter/releases) - Download pre-built binaries and packages
