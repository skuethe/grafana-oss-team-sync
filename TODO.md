# TODOs

- repo release
    - add golang specific github actions
    - add badges
        - `https://goreportcard.com/`

- features
    - allow team permissions to add individual users
    - allow team permissions to add roles

- add README info
    - how to setup Grafana
        - "Azure AD" auth
        - disable registration ?
        - ...

- azure
    - rename azure to entraid
    - modify graph sdk via kiota and verify if that gives us a smaller package size
        -> https://learn.microsoft.com/en-gb/graph/sdks/customize-client?tabs=go
        -> https://stackoverflow.com/questions/78355878/how-to-disable-backingstore-and-dirty-tracking-in-graph-beta-sdk-for-java

- groups / teams
    - either make Azure group search results case sensitive OR
    - make Grafana team search case insensitive
    - test azure pagination behaviour -> request ~10 groups and see what happens

- users
    - add user to team after user was created (we need the user ID)
        -> https://pkg.go.dev/github.com/grafana/grafana-openapi-client-go@v0.0.0-20250428202209-be3a35ff1dac/client/teams#Client.AddTeamMember

- config
    - allow to disable specific steps (f.e. via "skip.folders", "skip.teams" or "skip.users")
    - add "retry" to connect to Grafana instance

- Grafana
    - version check via API?

- logging
    - make detailed results from azure all "DEBUG"
    - make "skipped" results from Grafana "DEBUG"
    - make "created" results from Grafana "INFO"

- tests
    - write _test files
    - integration tests
        - run against multiple Grafana versions in CI

- CI
    - add github actions
