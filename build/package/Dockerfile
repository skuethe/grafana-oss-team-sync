# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
COPY --from=docker.io/golang:alpine@sha256:244baa35bcfaf9a5b3444519df6d42554a1f92dc33820bd98f0662df270d8a6a /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
