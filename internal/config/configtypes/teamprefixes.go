// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import (
	"slices"

	"github.com/skuethe/grafana-oss-team-sync/internal/helpers"
)

type TeamPrefixes []string

const (
	TeamPrefixesDefault   string = ""
	TeamPrefixesFlagHelp  string = "the comma-separated list of team name prefixes you want to sync"
	TeamPrefixesParameter string = "teamPrefixes"
	TeamPrefixesOptimized string = "teamprefixes"
)

// SanitizeTeamPrefixes drops empty entries from the configured team prefixes.
// Empty entries can be introduced by the default value of the comma-separated
// flag and would otherwise produce invalid (empty) source filters.
func (c *Config) SanitizeTeamPrefixes() {
	for slices.Contains(c.TeamPrefixes, "") {
		c.TeamPrefixes = helpers.RemoveFromSlice(c.TeamPrefixes, "", false)
	}
}
