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



- tests
    - write _test files
        - cli flags
        - env flags
            -> move "config -> k.unmarschal" into own func and write test for "minimal" and "full" supported config

    - add mock data to test against grafana (teams, users, folders)
        -> allow mock to be used as source?


- add README info
    - improve entraid docs


- groups / teams
    - either make Azure group search results case sensitive OR
    - make Grafana team search case insensitive
