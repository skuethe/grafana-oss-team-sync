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

- main
    - print "features" on startup, if they are NOT set to default


- tests
    - write _test files
        - cli flags
        - env flags
            -> move "config -> k.unmarschal" into own func and write test for "minimal" and "full" supported config
    - integration / end2end tests
        - run against multiple Grafana versions in CI
    - add mock data to test against grafana (teams, users, folders)
        -> allow mock to be used as source?


- add README info
    - how to setup Grafana
        - "Azure AD" auth
            - minimum permissions required?
        - disable "allow sign up"
            -> If not enabled, only existing Grafana users can log in using OAuth.
        - ...
    - how to setup EntraID app with permissions


- entraid
    - modify graph sdk via kiota and verify if that gives us a smaller package size
        -> https://learn.microsoft.com/en-gb/graph/sdks/customize-client?tabs=go
        -> https://stackoverflow.com/questions/78355878/how-to-disable-backingstore-and-dirty-tracking-in-graph-beta-sdk-for-java


- groups / teams
    - either make Azure group search results case sensitive OR
    - make Grafana team search case insensitive


- logging
    - use "-v" / "--verbose" flag for debugging?
    - make detailed results from entraid all "DEBUG"
    - make "skipped" results from Grafana "DEBUG"
    - make "created" results from Grafana "INFO"
    - Use OpenTelemetry standards / go libs?


- look into:
    - Go "data streams" aka channels?
    - Go contexts and their benefits


- future features
    - allow folder permissions to add individual users
    - allow folder permissions to add roles
    - allow to assign admin permissions to team members
