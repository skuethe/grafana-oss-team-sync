# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
COPY --from=docker.io/golang:alpine@sha256:98e6cffc31ccc44c7c15d83df1d69891efee8115a5bb7ede2bf30a38af3e3c92 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
