{
  pkgs,
}:
pkgs.mkShell {
  packages = (
    with pkgs;
    [
      go
      go-task
      gofumpt
      goimports-reviser
      golangci-lint
      goreleaser
      gosec
      lefthook

      # for updating hash of the Go module in the flake
      nix-update
      gnused

      # for generating GoReleaser configuration
      ytt
    ]
  );
}
