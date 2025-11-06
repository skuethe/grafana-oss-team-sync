# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
COPY --from=docker.io/golang:alpine@sha256:8b6b77a5e6a9dda591e864e1a2856d436d94219befa5f54d7ce76d2a77cc7a06 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
