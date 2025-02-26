#! /usr/bin/env nix-shell
#! nix-shell -i sh -p git gnumake
git pull
make deploy
