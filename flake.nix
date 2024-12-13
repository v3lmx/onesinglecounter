{
  description = "osc";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    utils,
  }:
    utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
      };
    in {
      packages = {
        default = pkgs.buildGoModule {
          pname = "osc";
          version = "0.0.1";
          src = ./.;
          vendorHash = "sha256-2NsSRaFiFu7ZKl/OS07s0RK+094sIRyuuYXZzOQFsIs=";
          proxyVendor = true;

          meta = {
            description = "osc server";
            homepage = "https://github.com/v3lmx/onesinglecounter";
            license = pkgs.lib.licenses.gpl3Plus;
            maintainers = with pkgs.lib.maintainers; [v3lmx];
          };
        };
      };

      devShells = {
        default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_23
            golangci-lint

            mprocs
            wgo
            websocat
            falkon
          ];
        };
      };
    });
}
