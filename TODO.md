# TODOs

- repo release
    - add golang specific github actions
    - add badges
        - `https://goreportcard.com/`
        - `https://github.com/badges/shields`
        - `https://www.bestpractices.dev/en` ?



- future features
    - allow folder permissions to add individual users
    - allow folder permissions to add roles
    - allow to assign admin permissions to team members


- cli
    - optimize flags usage for "features" and "grafana"


- add README info
    - add flags documentation
    - add hint that passwords are generated
    - add "build it yourself" docs
    - how to setup Grafana
        - "Azure AD" auth
            - minimum permissions required?
        - disable registration ?
        - ...
    - how to setup EntraID app with permissions


- config
    - add "authfile" feature to load authentication variables from an optional file instead of env variables or the config.yaml
        -> need to load these as ENV variables (for source auth) AND as config input (for Grafana auth)
    - validate all possible ENV VAR input


- grafana
    - validate possible authentication possibilities (besides basic auth) and see if we can / want to support them
        -> https://grafana.com/docs/grafana/latest/developers/http_api/authentication/
        - basic auth
        - service account token


- entraid
    - modify graph sdk via kiota and verify if that gives us a smaller package size
        -> https://learn.microsoft.com/en-gb/graph/sdks/customize-client?tabs=go
        -> https://stackoverflow.com/questions/78355878/how-to-disable-backingstore-and-dirty-tracking-in-graph-beta-sdk-for-java


- groups / teams
    - either make Azure group search results case sensitive OR
    - make Grafana team search case insensitive
    - test entraid pagination behaviour -> request ~10 groups and see what happens


- logging
    - make detailed results from entraid all "DEBUG"
    - make "skipped" results from Grafana "DEBUG"
    - make "created" results from Grafana "INFO"
    - Use OpenTelemetry standards / go libs?


- tests
    - write _test files
    - integration tests
        - run against multiple Grafana versions in CI
