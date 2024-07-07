#!/usr/bin/env bash

docker container exec -it --workdir /coach/cmd/cli coach-api-devcontainer go run .
