# SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
COPY --from=docker.io/golang:alpine@sha256:a86c313035ea07727c53a9037366a63c2216f3c5690c613179f37ee33ea71301 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY grafana-oss-team-sync /usr/bin/grafana-oss-team-sync
ENTRYPOINT ["/usr/bin/grafana-oss-team-sync"]
