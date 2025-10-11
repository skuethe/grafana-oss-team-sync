# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
COPY --from=docker.io/golang:alpine@sha256:06cdd34bd531b810650e47762c01e025eb9b1c7eadd191553b91c9f2d549fae8 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
