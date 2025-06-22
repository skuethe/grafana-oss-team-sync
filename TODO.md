<!--
SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
SPDX-License-Identifier: GPL-3.0-or-later
-->

# TODOs

- repo release
    - add golang specific github actions
        - ref: https://github.com/jkroepke/openvpn-auth-oauth2
        - codecov/codecov-action

    - add badges
        - `https://goreportcard.com/`
        - `https://github.com/badges/shields`
        - `https://reuse.software/dev/#api`
        - `https://www.bestpractices.dev/en` ?

    - add renovatebot
        - https://docs.renovatebot.com/golang/
        - https://docs.renovatebot.com/upgrade-best-practices/
        - https://docs.renovatebot.com/presets-group/#groupgoopenapi
        -

- CI
    - add service container for devproxy which runs integration tests on PRs


- CONTRIBUTING.md
    - add info about unit, integration and e2e tests


- tests
    - write _test files
        - res:
            - https://www.digitalocean.com/community/tutorials/how-to-write-unit-tests-in-go-using-go-test-and-the-testing-package#step-5-mdash-writing-coverage-tests-in-go
            - https://blog.jetbrains.com/go/2022/11/22/comprehensive-guide-to-testing-in-go/#errors-and-logs
        - cli flags
        - env flags
            -> move "config -> k.unmarschal" into own func and write test for "minimal" and "full" supported config

    - unit
        - "go test -v -coverprofile=c.out ./...
        - _test.go header: "//go:build !integration && !e2e"

    - integration test
        - don't test source (entraid)
        - send "grafana models" structs to grafana target
        - run complete "Grafana" package in one go (so: teams, users and folders)
        - run against multiple Grafana versions in CI
        - "go test -tags=integration  ./..."
        - _test.go header: "//go:build e2e"

    - end2end tests
        - run entraid against devproxy
        - "go test -tags=e2e  ./..."
        - _test.go header: "//go:build integration"

    - add mock data to test against grafana (teams, users, folders)
        -> allow mock to be used as source?


- add README info
    - reorg "configuration"
        -> improve entraid docs

    - how to setup Grafana
        - "Azure AD" auth
            - minimum permissions required?
        - disable "allow sign up"
            -> If not enabled, only existing Grafana users can log in using OAuth.
        - ...
    - how to setup EntraID app with permissions


- entraid
    - use "dev proxy" to test app and get more insights:
        - https://learn.microsoft.com/en-us/microsoft-cloud/dev/dev-proxy/how-to/test-my-app-with-random-errors
        - https://learn.microsoft.com/en-us/microsoft-cloud/dev/dev-proxy/how-to/simulate-errors-microsoft-graph-apis
        - https://learn.microsoft.com/en-us/microsoft-cloud/dev/dev-proxy/how-to/test-that-my-application-handles-throttling-properly
        - https://learn.microsoft.com/en-us/microsoft-cloud/dev/dev-proxy/how-to/simulate-rate-limit-api-responses
        - https://learn.microsoft.com/en-us/microsoft-cloud/dev/dev-proxy/how-to/detect-minimal-microsoft-graph-api-permissions

        - use in:
            - container: https://learn.microsoft.com/en-us/microsoft-cloud/dev/dev-proxy/how-to/use-dev-proxy-in-docker-container
            - ci/cd:     https://learn.microsoft.com/en-us/microsoft-cloud/dev/dev-proxy/how-to/use-dev-proxy-in-ci-cd-overview

        - mock data
            - https://learn.microsoft.com/en-us/microsoft-cloud/dev/dev-proxy/how-to/mock-responses#configure-dev-proxy-to-use-the-mock-responses
            - generator: https://learn.microsoft.com/en-us/microsoft-cloud/dev/dev-proxy/technical-reference/mockgeneratorplugin


    - modify graph sdk via kiota and verify if that gives us a smaller package size
        -> https://learn.microsoft.com/en-gb/graph/sdks/customize-client?tabs=go
        -> https://stackoverflow.com/questions/78355878/how-to-disable-backingstore-and-dirty-tracking-in-graph-beta-sdk-for-java


- groups / teams
    - either make Azure group search results case sensitive OR
    - make Grafana team search case insensitive


- logging
    - use "-v" / "--verbose" flag for debugging?
    - Use OpenTelemetry standards / go libs?


- look into:
    - Go "data streams" aka channels?
    - Go contexts and their benefits


- future features
    - allow folder permissions to add individual users
    - allow folder permissions to add roles
    - allow to assign admin permissions to team members
