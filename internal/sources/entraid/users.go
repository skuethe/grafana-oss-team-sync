// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package entraid

import (
	"context"
	"log/slog"
	"strings"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphgroups "github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
)

type users struct {
	groupid      *string
	client       *graph.GraphServiceClient
	headers      *abstractions.RequestHeaders
	grafanaUsers *grafana.Users
}

func (u *users) processUserResult(result *models.UserCollectionResponseable) {
	for _, user := range (*result).GetValue() {
		userDisplayName := *user.GetDisplayName()
		userPrincipalName := *user.GetUserPrincipalName()
		userMail := user.GetMail()

		var mail string
		if userMail != nil {
			mail = strings.ToLower(*userMail)
		}

		userLog := slog.With(
			slog.Group("user",
				slog.String("principalname", userPrincipalName),
				slog.String("displayname", userDisplayName),
				slog.String("mail", mail),
			),
		)
		if userMail == nil {
			userLog.Warn("user is missing the required email - skipping")
			continue
		}
		userLog.Debug("found EntraID user")

		*u.grafanaUsers = append(*u.grafanaUsers, grafana.User{
			Login: userPrincipalName,
			Name:  userDisplayName,
			Email: mail,
		})
	}
}

func (u *users) handleUserPagination(nextLink *string) (*models.UserCollectionResponseable, error) {
	configuraton := &graphgroups.ItemTransitiveMembersGraphUserRequestBuilderGetRequestConfiguration{
		Headers: u.headers,
	}
	result, err := u.client.Groups().ByGroupId(*u.groupid).TransitiveMembers().GraphUser().WithUrl(*nextLink).Get(context.Background(), configuraton)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *users) getInitialUsersFromGroup() (*models.UserCollectionResponseable, error) {

	requestCount := true
	requestParams := &graphgroups.ItemTransitiveMembersGraphUserRequestBuilderGetQueryParameters{
		Select: []string{"userPrincipalName", "displayName", "mail"},
		Count:  &requestCount,
	}
	configuraton := &graphgroups.ItemTransitiveMembersGraphUserRequestBuilderGetRequestConfiguration{
		Headers:         u.headers,
		QueryParameters: requestParams,
	}
	result, err := u.client.Groups().ByGroupId(*u.groupid).TransitiveMembers().GraphUser().Get(context.Background(), configuraton)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (g *groups) ProcessUsers(groupID *string) *grafana.Users {
	usersLog := slog.With(
		slog.String("package", "entraid.users"),
		slog.String("group", *groupID),
	)

	grafanaUsers := &grafana.Users{}

	if config.Instance.Features.DisableUserSync {
		usersLog.Debug("usersync feature disabled, skipping")
		return grafanaUsers
	} else {
		usersLog.Info("processing EntraID users for group")

		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		u := users{
			groupid:      groupID,
			client:       g.client,
			headers:      headers,
			grafanaUsers: grafanaUsers,
		}

		ur, err := u.getInitialUsersFromGroup()
		if err != nil {
			usersLog.Error("could not get initial user result from EntraID")
			panic(err)
		}

		countFound := (*ur).GetOdataCount()

		for {
			// Handle user result
			u.processUserResult(ur)

			// Handle possible pagination
			nextPageUrl := (*ur).GetOdataNextLink()
			if nextPageUrl != nil {
				usersLog.Debug("processing paginated user result")
				ur, err = u.handleUserPagination(nextPageUrl)
				if err != nil {
					usersLog.Error("could not get paged user result from EntraID")
					panic(err)
				}
			} else {
				break
			}
		}

		usersLog.Info("finished processing EntraID users for group",
			slog.Group("users",
				slog.Int64("found", *countFound),
			),
		)
		return u.grafanaUsers
	}
}
