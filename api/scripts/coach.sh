#!/usr/bin/env bash

USAGE_MESSAGE=$(cat <<-END

Usage: coach [options] [<command> [command-options]]

Options:
    -h, --help         Show help message
    -v, --version      Show version information

Commands (required unless using -h or -v):
    chat               Chat with the coach agent 

For more information on a specific command, including available options, use:
    coach <command> -h

END
)

VERSION="Coach version 1.0.0"

show_help() {
  echo "$USAGE_MESSAGE"
}

show_version() {
  echo "$VERSION"
}

SCRIPT_DIR="$(dirname "$(realpath "$0")")"

# Parse the command line options.
while [[ $# -gt 0 ]]; do
  case "$1" in
    # Match short options first
    -[!-]*)
      # Match one letter at a time if multiple short options are combined
      # e.g. coach -hv
      for (( i=1; i<${#1}; i++ )); do
       # Execute whichever short option comes first
        case "${1:i:1}" in
          h)
            show_help
            exit 0
            ;;
          v)
            show_version
            exit 0
            ;;
          *)
            echo "Invalid option: -${1:i:1}" 1>&2
            show_help
            exit 1
            ;;
        esac
      done
      shift
      ;;
    # Match long options
    --help)
      show_help
      exit 0
      ;;
    --version)
      show_version
      exit 0
      ;;
    chat)
      source $SCRIPT_DIR/chat.sh "$@"
      break
      ;;
    *)
      echo "Invalid option: $1" 1>&2
      show_help
      exit 1
      ;;
  esac
done
