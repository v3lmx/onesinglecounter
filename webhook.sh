#!/bin/bash
direnv allow
eval "$(direnv export bash)"
git pull
make deploy
