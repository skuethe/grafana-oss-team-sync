# TODOs

- search "TODO" in code

- repo release
    - add golang specific github actions
    - add badges
        - `https://goreportcard.com/`


- future features
    - allow folder permissions to add individual users
    - allow folder permissions to add roles
    - allow to assign admin permissions to team members


- add README info
    - how to setup Grafana
        - "Azure AD" auth
            - minimum permissions required?
        - disable registration ?
        - ...
    - how to setup EntraID app with permissions


- grafana
    - validate possible authentication possibilities (besides basic auth) and see if we can / want to support them
    - add "labels" to users / groups / folders to mark them "source: grafana-oss-team-sync"?


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
