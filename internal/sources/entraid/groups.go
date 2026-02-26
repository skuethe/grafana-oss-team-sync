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
	"github.com/skuethe/grafana-oss-team-sync/internal/helpers"
)

type groups struct {
	client       *graph.GraphServiceClient
	headers      *abstractions.RequestHeaders
	grafanaTeams *grafana.Teams
}

func (g *groups) processGroupResult(result *models.GroupCollectionResponseable) {
	for _, group := range (*result).GetValue() {
		groupDisplayName := group.GetDisplayName()
		groupId := *group.GetId()
		groupMail := group.GetMail()

		var mail string
		if groupMail != nil {
			mail = strings.ToLower(*groupMail)
		}

		groupLog := slog.With(
			slog.Group("group",
				slog.String("displayname", *groupDisplayName),
				slog.String("id", groupId),
			),
		)
		groupLog.Info("found EntraID group")

		// Process users
		grafanaUserList := g.ProcessUsers(&groupId)

		*g.grafanaTeams = append(*g.grafanaTeams, grafana.Team{
			Parameter: &grafana.TeamParameter{
				Name:  groupDisplayName,
				Email: mail,
			},
			Users: grafanaUserList,
		})
		config.Instance.Teams = helpers.RemoveFromSlice(config.Instance.Teams, *groupDisplayName, false)
	}
}

func (g *groups) handleGroupPagination(nextLink *string) (*models.GroupCollectionResponseable, error) {
	configuration := &graphgroups.GroupsRequestBuilderGetRequestConfiguration{
		Headers: g.headers,
	}
	result, err := g.client.Groups().WithUrl(*nextLink).Get(context.Background(), configuration)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (g *groups) getInitialGroupRequest() (*models.GroupCollectionResponseable, error) {
	requestCount := true
	requestFilter := "displayName in ('" + strings.Join(config.Instance.Teams, "', '") + "')"
	requestParams := &graphgroups.GroupsRequestBuilderGetQueryParameters{
		Filter: &requestFilter,
		Select: []string{"id", "displayName", "mail"},
		Count:  &requestCount,
	}
	configuration := &graphgroups.GroupsRequestBuilderGetRequestConfiguration{
		Headers:         g.headers,
		QueryParameters: requestParams,
	}
	result, err := g.client.Groups().Get(context.Background(), configuration)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
