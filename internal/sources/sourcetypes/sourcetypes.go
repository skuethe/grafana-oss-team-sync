// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package sourcetypes

import (
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
)

type SourcePlugin struct {
	EntraID *msgraph.GraphServiceClient
}
