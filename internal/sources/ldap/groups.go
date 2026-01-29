// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package ldap

import (
	"log/slog"
	"strings"

	"github.com/go-ldap/ldap/v3"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/helpers"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/sourcetypes"
)

type groups struct {
	client       *sourcetypes.LDAPClient
	grafanaTeams *grafana.Teams
}

func (g *groups) processGroupResult(result *ldap.SearchResult) {
	for _, group := range result.Entries {
		groupDisplayName := group.GetAttributeValue("cn")
		groupId := group.GetAttributeValue("cn")
		// groupMail := group.GetMail()

		// var mail string
		// if groupMail != nil {
		// 	mail = strings.ToLower(*groupMail)
		// }

		groupLog := slog.With(
			slog.Group("group",
				slog.String("displayname", groupDisplayName),
				slog.String("id", groupId),
			),
		)
		groupLog.Info("found LDAP group")

		// Process users
		grafanaUserList := g.ProcessUsers(group)

		*g.grafanaTeams = append(*g.grafanaTeams, grafana.Team{
			Parameter: &grafana.TeamParameter{
				Name: &groupDisplayName,
				// Email: mail,
			},
			Users: grafanaUserList,
		})
		config.Instance.Teams = helpers.RemoveFromSlice(config.Instance.Teams, groupDisplayName, false)
	}
}

func ProcessGroups(instance *sourcetypes.SourcePlugin) *grafana.Teams {
	groupsLog := slog.With(slog.String("package", "ldap.groups"))
	groupsLog.Info("processing LDAP groups")

	g := groups{
		client:       instance.LDAP,
		grafanaTeams: &grafana.Teams{},
	}

	s := search{
		connection: instance.LDAP.Connection,
		baseDN:     instance.LDAP.BaseDN,
		filter:     instance.LDAP.GroupFilter,
		attributes: []string{"cn", "member"},
	}

	sr, err := s.perform()
	if err != nil {
		groupsLog.Error("could not get group result from LDAP")
		panic(err)
	}

	countFound := len(sr.Entries)

	// Handle group result
	g.processGroupResult(sr)

	if len(config.Instance.Teams) > 0 {
		groupsLog.Warn("could not find the following groups in LDAP", "skipped", strings.Join(config.Instance.Teams, ","))
	}

	groupsLog.Info("finished processing LDAP groups",
		slog.Group("groups",
			slog.Int("found", countFound),
			slog.Int("skipped", len(config.Instance.Teams)),
		),
	)

	return g.grafanaTeams
}
