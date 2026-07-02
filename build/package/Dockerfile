# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
COPY --from=docker.io/golang:alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY $TARGETPLATFORM/grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
