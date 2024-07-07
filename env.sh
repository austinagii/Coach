#!/usr/bin/env sh

USAGE_MESSAGE=$(cat <<END
Usage: dev [options]

Starts up the devcontainers for the coach application 

Options:
  -r, --rebuid      Specifies whether to rebuild the devcontainer images. (Defaults to false)
END
)

OPTS=$(getopt -o rh -l rebuild,help --name "$0" -- "$@")
if [ $? -ne 0 ]; then 
  exit 1
fi
eval set -- "$OPTS"

SHOULD_REBUILD=0

while true; do
  case "$1" in 
    -r|--rebuild)
      SHOULD_REBUILD=1
      shift
      ;;
    -h|--help)
      echo "$USAGE_MESSAGE"
      exit 0
      ;;
    --)
      shift
      break
      ;;
    *)
      "Error: Invalid option '$1'"
      "See 'dev --help' for a list of valid options"
      exit 1
      ;;
  esac
done

if [ $SHOULD_REBUILD -eq 1 ]; then 
  docker compose -f docker-compose.yaml -f dev.docker-compose.yaml build --no-cache
  if [ $? -ne 0 ]; then 
    echo "Error: Failed to build dev environment. Exiting..." 
    exit 1
  fi
fi

docker compose -f docker-compose.yaml -f dev.docker-compose.yaml up -d
if [ $? -ne 0 ]; then 
  echo "Error: Failed to start dev environment. Exiting..."
  exit 1
fi

docker container attach coach

echo "Dev session complete, shutting down"
docker compose down --remove-orphans 
echo "Shutdown complete"
