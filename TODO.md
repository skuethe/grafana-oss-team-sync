# TODOs

- repo release
    - add golang specific github actions
        - ref: https://github.com/jkroepke/openvpn-auth-oauth2

    - add badges
        - `https://goreportcard.com/`
        - `https://github.com/badges/shields`
        - `https://www.bestpractices.dev/en` ?

    - add renovatebot
        - https://docs.renovatebot.com/golang/
        - https://docs.renovatebot.com/upgrade-best-practices/
        - https://docs.renovatebot.com/presets-group/#groupgoopenapi
        - 


- cli
    - optimize flags usage for "features" and "grafana"
    - add "token", "teams" and "folders
    - add "help"?


- add README info
    - add "build it yourself" docs
    - how to setup Grafana
        - "Azure AD" auth
            - minimum permissions required?
        - disable registration ?
        - ...
    - how to setup EntraID app with permissions
    - describe parameter hirarchy
        -> config.yaml
        -> env vars
        -> cli flags
        -> optional "auth file" if defined (WILL NOT OVERWRITE EXISTING)


- config
    - validate all possible ENV VAR input
    - implement solution to load config from OS ENV / --config flag


- entraid
    - modify graph sdk via kiota and verify if that gives us a smaller package size
        -> https://learn.microsoft.com/en-gb/graph/sdks/customize-client?tabs=go
        -> https://stackoverflow.com/questions/78355878/how-to-disable-backingstore-and-dirty-tracking-in-graph-beta-sdk-for-java


- groups / teams
    - either make Azure group search results case sensitive OR
    - make Grafana team search case insensitive
    - test entraid pagination behaviour -> request ~10 groups and see what happens


- logging
    - use "-v" / "--verbose" flag for debugging?
    - make detailed results from entraid all "DEBUG"
    - make "skipped" results from Grafana "DEBUG"
    - make "created" results from Grafana "INFO"
    - Use OpenTelemetry standards / go libs?


- tests
    - write _test files
    - integration tests
        - run against multiple Grafana versions in CI



- look into:
    - Go "data streams" aka channels?
    - Go contexts and their benefits
    - Go use "panic(err)" instead of "os.Exit()"?


- future features
    - allow folder permissions to add individual users
    - allow folder permissions to add roles
    - allow to assign admin permissions to team members
