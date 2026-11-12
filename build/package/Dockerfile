# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
COPY --from=docker.io/golang:alpine@sha256:eb61c3873c97416fff0116553fff65d67a3549f72e11d9758849a689cdd8b523 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY $TARGETPLATFORM/grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
