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
    - fix "login" attribute -> migrate from ID to surname
    - add user to team after user was created

- config
    - allow to disable specific steps (f.e. via "skip.folders", "skip.teams" or "skip.users")

- logging
    - make detailed results from azure all "DEBUG"
    - make "skipped" results from Grafana "DEBUG"
    - make "created" results from Grafana "INFO"
