{
  description = "buzzer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {inherit system;};
    in {
      devShells = {
        default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # backend
            go_1_23
            golangci-lint

            # frontend
            nodejs
            pnpm

            # run commands
            mprocs

            # live reload go files
            wgo

            # ws client
            websocat
          ];
        };
      };
    });
}
