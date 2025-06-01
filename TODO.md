# TODOs

- repo release
    - add golang specific github actions
    - add badges
        - `https://goreportcard.com/`

- future features
    - allow team permissions to add individual users
    - allow team permissions to add roles

- add README info
    - how to setup Grafana
        - "Azure AD" auth
            - minimum permissions required?
        - disable registration ?
        - ...
    - how to setup EntraID app with permissions

- grafana
    - validate possible authentication possibilities (besides basic auth) and see if we can / want to support them

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

- CI
    - add github actions

- workflow:
    1. plugin: fetch groups and save it an the grafana.Teams object
    2. plugin: fetch users per group and also save it an the grafana.Teams object
    3. 