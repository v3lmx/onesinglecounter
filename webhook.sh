#! /usr/bin/env nix-shell
#! nix-shell -i sh -p git gnumake go_1_23 uglifyjs podman
git config --global --add safe.directory /opt/osc
git pull
make deploy
#/bin/bash
