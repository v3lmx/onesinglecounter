#!/bin/bash
git config --global --add safe.directory /opt/osc
direnv allow
eval "$(direnv export bash)"
echo $PATH
git pull
make deploy
