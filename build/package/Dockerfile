# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
COPY --from=docker.io/golang:alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY $TARGETPLATFORM/grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
