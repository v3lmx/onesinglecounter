#!/bin/bash
set -x
git config --global --add safe.directory /opt/osc
direnv allow
eval "$(direnv export bash)"
git pull
make deploy
