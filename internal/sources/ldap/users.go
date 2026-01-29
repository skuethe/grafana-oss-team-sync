// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package ldap

import (
	"log/slog"
	"strings"

	"github.com/go-ldap/ldap/v3"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
)

type users struct {
	// groupid      *string
	// client       *ldap.Client
	grafanaUsers *grafana.Users
}

func (u *users) processUserResult(result *ldap.SearchResult) {
	for _, user := range result.Entries {
		userDisplayName := user.GetAttributeValue("displayName")
		userCN := user.GetAttributeValue("cn")
		userMail := strings.ToLower(user.GetAttributeValue("mail"))

		userLog := slog.With(
			slog.Group("user",
				slog.String("cn", userCN),
				slog.String("displayname", userDisplayName),
				slog.String("mail", userMail),
			),
		)
		// if userMail == nil {
		// 	userLog.Warn("user is missing the required email - skipping")
		// 	continue
		// }
		userLog.Debug("found LDAP user")

		*u.grafanaUsers = append(*u.grafanaUsers, grafana.User{
			Login: userMail,
			Name:  userDisplayName,
			Email: userMail,
		})
	}
}

func (g *groups) ProcessUsers(group *ldap.Entry) *grafana.Users {
	usersLog := slog.With(
		slog.String("package", "ldap.users"),
		slog.String("group", group.GetAttributeValue("cn")),
	)

	grafanaUsers := &grafana.Users{}

	if config.Instance.Features.DisableUserSync {
		usersLog.Debug("usersync feature disabled, skipping")
		return grafanaUsers
	} else {
		usersLog.Info("processing LDAP users for group")

		u := users{
			grafanaUsers: grafanaUsers,
		}

		userList := group.GetAttributeValues("member")
		countFound := len(userList)

		for _, member := range userList {
			usersLog.Debug("processing member", "member", member)
			s := search{
				connection: g.client.Connection,
				baseDN:     member,
				filter:     "(objectClass=*)",
				attributes: []string{"cn", "displayName", "mail"},
			}

			sr, err := s.perform()
			if err != nil {
				usersLog.Error("could not get user result from LDAP")
				panic(err)
			}

			u.processUserResult(sr)
		}

		usersLog.Info("finished processing LDAP users for group",
			slog.Group("users",
				slog.Int("found", countFound),
			),
		)
		return u.grafanaUsers
	}
}
