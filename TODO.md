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


- groups / teams
    - either make Azure group search results case sensitive OR
    - make Grafana team search case insensitive

- users
    0. Optional: enable / disable user sync in config
    1. for every source group add users to a list `[]models.AdminCreateUserForm`
    2. make sure each entry is unique
    3. call grafana.ProcessUsers
