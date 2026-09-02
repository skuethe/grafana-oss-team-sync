# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
COPY --from=docker.io/golang:alpine@sha256:26402d86be3d72e6a9410afa0108f03529f51f0c1b5eb7f503d0bc44cc7857ac /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY $TARGETPLATFORM/grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
