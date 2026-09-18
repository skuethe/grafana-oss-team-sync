# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
COPY --from=docker.io/golang:alpine@sha256:4cb7ac979db5fcc41cae44b2227ba5ab8a51e8807f40d9ba4dee20a0ad960b5b /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY $TARGETPLATFORM/grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
