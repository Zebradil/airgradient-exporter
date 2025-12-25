{
  pkgs,
  self,
}:
let
  fs = pkgs.lib.fileset;
  sourceFiles = fs.unions [
    ../go.mod
    ../go.sum
    ../pkg
    ../main.go
  ];
  baseVersion = "0.1.0";
  commit = self.shortRev or self.dirtyShortRev or "unknown";
  version = "${baseVersion}-${commit}";
  name = "airgradient-exporter";
in
pkgs.buildGoModule {
  pname = name;
  src = fs.toSource {
    root = ./..;
    fileset = sourceFiles;
  };
  vendorHash = "sha256-FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF";
  version = version;

  CGO_ENABLED = 0;
  ldflags = [
    "-s"
    "-w"
    "-X=main.version=${baseVersion}"
    "-X=main.commit=${commit}"
    "-X=main.date=1970-01-01"
  ];

  meta = {
    changelog = "https://github.com/Zebradil/${name}/blob/${baseVersion}/CHANGELOG.md";
    description = "Prometheus exporter for Airgradient DIY air quality monitors";
    homepage = "https://github.com/Zebradil/${name}";
    license = pkgs.lib.licenses.mit;
  };
}
