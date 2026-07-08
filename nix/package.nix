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
  baseVersion = "1.2.8";
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
  vendorHash = "sha256-TmYUfVaaXEwLTMxcKM0fsOufo40VM7hpk7jyvoc0zuA=";
  version = version;

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

  # Include systemd unit file in the package
  postInstall = ''
    mkdir -p $out/lib/systemd/system
    cp ${../systemd/airgradient-exporter.service} $out/lib/systemd/system/airgradient-exporter.service
  '';
}
