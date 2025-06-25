// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

type Teams []string

const (
	TeamsDefault   string = ""
	TeamsFlagHelp  string = "the comma-seperated list of teams you want to sync"
	TeamsParameter string = "teams"
)
