// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

type Teams []string

type TeamPrefixes []string

const (
	TeamsDefault   string = ""
	TeamsFlagHelp  string = "the comma-separated list of teams you want to sync"
	TeamsParameter string = "teams"

	TeamPrefixesDefault   string = ""
	TeamPrefixesFlagHelp  string = "the comma-separated list of team name prefixes you want to sync"
	TeamPrefixesParameter string = "teamPrefixes"
	TeamPrefixesOptimized string = "teamprefixes"
)
