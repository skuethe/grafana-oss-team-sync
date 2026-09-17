// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import "strings"

// Team describes a single group/team to sync and the Grafana organization it should be synced into.
//
// In YAML, a team can either be given as a plain string (using the default organization)
// or as an object with "name" and "orgId" keys:
//
//	teams:
//	  - myTeamWithDefaultOrg
//	  - name: myTeamInOrg2
//	    orgId: 2
type Team struct {
	Name  string `yaml:"name"`
	OrgID int64  `yaml:"orgId"`
}

type Teams []Team

const (
	TeamsDefault   string = ""
	TeamsFlagHelp  string = "the comma-separated list of teams you want to sync"
	TeamsParameter string = "teams"

	// TeamsDefaultOrgID is used for display/validation purposes to describe "no explicit org configured".
	TeamsDefaultOrgID int64 = 0
)

// Names returns the plain list of team display names, e.g. to query a source plugin.
func (t Teams) Names() []string {
	names := make([]string, len(t))
	for i, team := range t {
		names[i] = team.Name
	}
	return names
}

// Find returns the team matching the given name (case-insensitive), if any.
func (t Teams) Find(name string) (Team, bool) {
	for _, team := range t {
		if strings.EqualFold(team.Name, name) {
			return team, true
		}
	}
	return Team{}, false
}

// Remove returns a copy of the teams list with the team matching the given name (case-insensitive) removed.
func (t Teams) Remove(name string) Teams {
	result := make(Teams, 0, len(t))
	for _, team := range t {
		if !strings.EqualFold(team.Name, name) {
			result = append(result, team)
		}
	}
	return result
}
