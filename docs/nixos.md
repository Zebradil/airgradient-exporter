# NixOS Module

The package includes a NixOS module for easy integration into NixOS configurations.

## Using with Flakes

Add the flake to your `flake.nix`:

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    airgradient-exporter.url = "github:Zebradil/airgradient-exporter";
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

## Using without Flakes

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

## Service Management

The service will automatically start on boot when enabled. The `multi-user.target` ensures the service starts when the system reaches the normal multi-user state (standard boot).

Check service status:

```bash
systemctl status airgradient-exporter
```

View logs:

```bash
journalctl -u airgradient-exporter -f
```

## Configuration Options

| Option       | Type           | Description                                                    | Default  |
| ------------ | -------------- | -------------------------------------------------------------- | -------- |
| `enable`     | boolean        | Enable the AirGradient Exporter service                         | `false`  |
| `monitors`   | list of string | List of AirGradient monitor IPs, hostnames, or host:port       | Required |
| `port`       | integer        | Port to listen on                                              | `9112`   |
| `logFormat`  | string         | Log format: "text" or "json"                                   | `"text"` |
| `package`    | package        | The airgradient-exporter package to use                        | default  |
