#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

set -eo pipefail

# VARIABLES

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
ROOT_DIR="${SCRIPT_DIR}/../"

# GRAFANA_VERSIONS="11.1.0 12.0.0"
GRAFANA_VERSIONS="11.1.0"

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
  echo -e "  - go"
  echo -e "\nUsage: ${0} <PARAM>"
  echo -e "\n  unit"
  echo -e "\twill run unit tests"
  echo -e "\n  coverage"
  echo -e "\twill parse coverage into HTML format and open in browser. Need to run 'unit' before to generate the coverage file"
  echo -e "\n  integration"
  echo -e "\twill run integration tests. Requires integration services to be avilable"
  echo -e "\n  e2e"
  echo -e "\twill run e2e tests. Requires e2e services to be available"
  echo -e "\n"
}

# VALIDATION

## Requires
requireCommand "go"


case "${1}" in
  "unit")
    cd ${ROOT_DIR} && go test -race -v -covermode=atomic -coverprofile=coverage.out ./...
    ;;
  "coverage")
    cd ${ROOT_DIR} && go tool cover -html=coverage.out
    ;;
  "integration")
    cd ${ROOT_DIR}
    for i in ${GRAFANA_VERSIONS}; do
      scripts/container.sh integration-test-start ${i}
      go clean -testcache -tags=integration
      go test -v -tags=integration  ./...
      scripts/container.sh integration-test-stop ${i}
    done
    ;;
  "e2e")
    cd ${ROOT_DIR} && go test -tags=e2e  ./...
    ;;
  *)
    usage
    ;;
esac
