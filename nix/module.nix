{
  config,
  lib,
  pkgs,
  ...
}:
with lib;
let
  cfg = config.services.airgradient-exporter;
in
{
  options.services.airgradient-exporter = {
    enable = mkEnableOption "AirGradient Prometheus exporter";

    monitors = mkOption {
      type = types.listOf types.str;
      description = ''
        List of AirGradient monitor IPs, hostnames, or host:port.
        Example: [ "192.168.1.50" "192.168.1.51" ] or [ "192.168.1.50:8080" ]
      '';
      example = [ "192.168.1.50" ];
    };

    port = mkOption {
      type = types.port;
      default = 9112;
      description = "Port to listen on";
    };

    logFormat = mkOption {
      type = types.enum [
        "text"
        "json"
      ];
      default = "text";
      description = "Log format: 'text' or 'json'";
    };

    package = mkOption {
      type = types.package;
      default =
        pkgs.airgradient-exporter or (import ../nix/package.nix {
          inherit pkgs;
          self = {
            shortRev = "unknown";
            dirtyShortRev = "unknown";
          };
        });
      defaultText = "pkgs.airgradient-exporter or package from flake";
      description = ''
        The airgradient-exporter package to use.
        When using this module from a flake, you can set it to:
        package = inputs.airgradient-exporter.packages.${pkgs.system}.default;
      '';
    };
  };

  config = mkIf cfg.enable {
    systemd.services.airgradient-exporter = {
      description = "Prometheus exporter for AirGradient air quality sensors";
      documentation = [
        "https://github.com/Zebradil/airgradient-exporter/blob/master/README.md"
      ];
      # multi-user.target is the standard systemd target for services that should
      # start automatically when the system reaches the normal multi-user state (boot).
      # This ensures the service starts on boot when enabled.
      wantedBy = [ "multi-user.target" ];
      wants = [ "network-online.target" ];
      after = [
        "network-online.target"
        "nss-lookup.target"
      ];

      serviceConfig = {
        Type = "simple";
        ExecStart = "${cfg.package}/bin/airgradient-exporter";
        Restart = "on-failure";
        RestartSec = "5s";
        DynamicUser = true;
        ProtectSystem = "strict";
        ProtectHome = true;
        PrivateTmp = true;
        NoNewPrivileges = true;
        AmbientCapabilities = "";
        CapabilityBoundingSet = "";
      };

      environment = {
        AIRGRADIENT_MONITORS = concatStringsSep "," cfg.monitors;
        PORT = toString cfg.port;
        LOG_FORMAT = cfg.logFormat;
      };
    };
  };
}
