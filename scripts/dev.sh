#!/usr/bin/env sh

docker compose -f docker-compose.yaml -f docker-compose.devcontainer.yaml up --force-recreate -d
