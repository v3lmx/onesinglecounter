{
  description = "osc";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    utils,
    ...
  }:
    utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
      };
    in {
      devShells = {
        default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_23
            golangci-lint

            uglify-js
            nodejs
            nodePackages.pnpm
            caddy

            mprocs
            wgo
            websocat

            ansible
            # kamal
            ruby
          ];

          # add kamal to path
          shellHook = ''
              export PATH="$PATH:$HOME/.local/share/gem/ruby/3.3.0/bin"
              alias ag="ansible-galaxy"
              alias ap="ansible-playbook"
            #   command -v kamal
            #   if [ $? -eq 0 ]; then
            #       gem install kamal
            #   fi
          '';
        };
      };

      nixosModules = {
        osc = import ./osc.nix;
      };
      nixosModule = self.nixosModules.osc;
    });
}
