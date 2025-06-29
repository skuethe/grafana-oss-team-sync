#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

set -eo pipefail

# VARIABLES

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
ROOT_DIR="${SCRIPT_DIR}/../"
DEPLOY_DIR="${ROOT_DIR}/deploy/"

# renovate: github-releases=golangci/golangci-lint
GOLANGCI_LINT_VERSION="2.2.0"
# renovate: github-releases=fsfe/reuse-tool
REUSE_VERSION="5.0.2"

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
  echo -e "\nUsage: ${0} <ARGS>"

  echo -e "\n  go-lint"
  echo -e "\twill run go lint checks against the code"

  echo -e "\n  licenses"
  echo -e "\twill run reuse lint checks against the code"

  echo -e "\n  renovate"
  echo -e "\twill run renovate config check"

  echo -e "\n  e2e-start"
  echo -e "\twill start the end to end testing compose stack: Grafana + devproxy. Devproxy is serving mock data for EntraID"
  echo -e "\n  e2e-stop"
  echo -e "\twill stop the end to end testing compose stack"
  echo -e "\n  e2e-logs"
  echo -e "\twill follow the logs of all containers of the end to end testing compose stack"

  echo -e "\n  integration-start [version-tag]"
  echo -e "\twill start the integration testing compose stack: Grafana. You have to pass the Grafana version you want to test (supported: '11.1.0', '12.0.0' or 'latest')"
  echo -e "\n  integration-stop"
  echo -e "\twill stop the integration testing compose stack"
  echo -e "\n  integration-logs"
  echo -e "\twill follow the logs of all containers of the integration testing compose stack"

  echo -e "\n"
}

# VALIDATION

## Requires
requireCommand "podman"
requireCommand "podman-compose"


case "${1}" in
  "go-lint")
    cd ${ROOT_DIR} && podman run --rm -v $(pwd):/app:ro -w /app docker.io/golangci/golangci-lint:v${GOLANGCI_LINT_VERSION} golangci-lint run
    ;;
  "licenses")
    cd ${ROOT_DIR} && podman run --rm -v $(pwd):/data:ro docker.io/fsfe/reuse:${REUSE_VERSION} lint
    ;;
  "renovate")
    cd ${ROOT_DIR} && podman run -e RENOVATE_CONFIG_FILE=/data/renovate.json --rm -v $(pwd):/data:ro ghcr.io/renovatebot/renovate renovate-config-validator --strict
    ;;
  "e2e-start")
    podman compose -f ${DEPLOY_DIR}/e2e-tests_docker-compose.yaml up -d
    ;;
  "e2e-stop")
    podman compose -f ${DEPLOY_DIR}/e2e-tests_docker-compose.yaml down
    ;;
  "e2e-logs")
    podman compose -f ${DEPLOY_DIR}/e2e-tests_docker-compose.yaml logs -f
    ;;
  "integration-start")
    echo "Using Grafana version: ${2}"
    podman compose -f ${DEPLOY_DIR}/integration-tests_docker-compose.yaml -f ${DEPLOY_DIR}/integration-tests_override_grafana-${2}.yaml up -d
    ;;
  "integration-stop")
    podman compose -f ${DEPLOY_DIR}/integration-tests_docker-compose.yaml down
    ;;
  "integration-logs")
    podman compose -f ${DEPLOY_DIR}/integration-tests_docker-compose.yaml logs -f
    ;;
  *)
    usage
    ;;
esac
