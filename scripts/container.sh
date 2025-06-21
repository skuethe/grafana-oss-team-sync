#!/usr/bin/env bash

set -eo pipefail

# VARIABLES

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
ROOT_DIR="${SCRIPT_DIR}/../"
DEPLOY_DIR="${ROOT_DIR}/deploy/"

GOLANGCI_LINT="v2.1.6"

# HELPERS

function requireCommand() {
  command="$1"
  if ! command -v ${command} &> /dev/null; then
    echo "This script requires the command \"${command}\" to be available / installed. Terminating!"
    exit 1
  fi
}

function usage() {
  echo -e "Requirements:"
  echo -e "  - podman"
  echo -e "  - podman-compose"
  echo -e "\nUsage: ${0} <PARAM>"
  echo -e "\n  start"
  echo -e "\twill start the compose stack"
  echo -e "\n  stop"
  echo -e "\twill stop compose stack"
  echo -e "\n  logs"
  echo -e "\twill follow the logs of all containers of the compose stack"
  echo -e "\n  lint"
  echo -e "\twill run lint checks against the code"
  echo -e "\n"
}

# VALIDATION

## Requires
requireCommand "podman"
requireCommand "podman-compose"


case "${1}" in
  "start")
    cd ${DEPLOY_DIR} && podman compose up -d
    ;;
  "stop")
    cd ${DEPLOY_DIR} && podman compose down
    ;;
  "logs")
    cd ${DEPLOY_DIR} && podman compose logs -f
    ;;
  "lint")
    cd ${ROOT_DIR} && podman run --rm -v $(pwd):/app:ro -w /app docker.io/golangci/golangci-lint:${GOLANGCI_LINT} golangci-lint run
    ;;
  *)
    usage
    ;;
esac
