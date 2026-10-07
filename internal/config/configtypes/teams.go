// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import (
	"slices"

	"github.com/skuethe/grafana-oss-team-sync/internal/helpers"
)

type Teams []string

const (
	TeamsDefault   string = ""
	TeamsFlagHelp  string = "the comma-separated list of teams you want to sync"
	TeamsParameter string = "teams"
)

// SanitizeTeams drops empty entries from the configured teams. Empty entries can
// be introduced by the default value of the comma-separated flag and would
// otherwise produce invalid (empty) source filters.
func (c *Config) SanitizeTeams() {
	for slices.Contains(c.Teams, "") {
		c.Teams = helpers.RemoveFromSlice(c.Teams, "", false)
	}
}
