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
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/sourcetypes"
)

type groups struct {
	client       *graph.GraphServiceClient
	headers      *abstractions.RequestHeaders
	grafanaTeams *grafana.Teams
}

func (g *groups) processGroupResult(result *models.GroupCollectionResponseable) {
	for _, group := range (*result).GetValue() {
		groupDisplayName := *group.GetDisplayName()
		groupId := *group.GetId()
		groupMail := group.GetMail()

		var mail string
		if groupMail != nil {
			mail = strings.ToLower(*groupMail)
		}

		groupLog := slog.With(
			slog.Group("group",
				slog.String("displayname", groupDisplayName),
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
		config.Instance.Teams = helpers.RemoveFromSlice(config.Instance.Teams, groupDisplayName, false)
	}
}

func (g *groups) handleGroupPagination(nextLink *string) (*models.GroupCollectionResponseable, error) {
	configuraton := &graphgroups.GroupsRequestBuilderGetRequestConfiguration{
		Headers: g.headers,
	}
	result, err := g.client.Groups().WithUrl(*nextLink).Get(context.Background(), configuraton)
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
	configuraton := &graphgroups.GroupsRequestBuilderGetRequestConfiguration{
		Headers:         g.headers,
		QueryParameters: requestParams,
	}
	result, err := g.client.Groups().Get(context.Background(), configuraton)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func ProcessGroups(instance *sourcetypes.SourcePlugin) *grafana.Teams {
	groupsLog := slog.With(slog.String("package", "entraid.groups"))
	groupsLog.Info("processing EntraID groups")

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	g := groups{
		client:       instance.EntraID,
		headers:      headers,
		grafanaTeams: &grafana.Teams{},
	}

	gr, err := g.getInitialGroupRequest()
	if err != nil {
		groupsLog.Error("could not get initial group result from EntraID")
		panic(err)
	}

	countFound := (*gr).GetOdataCount()

	for {
		// Handle group result
		g.processGroupResult(gr)

		// Handle possible pagination
		nextPageUrl := (*gr).GetOdataNextLink()
		if nextPageUrl != nil {
			groupsLog.Debug("processing paginated group result")
			gr, err = g.handleGroupPagination(nextPageUrl)
			if err != nil {
				groupsLog.Error("could not get paged group result from EntraID")
				panic(err)
			}
		} else {
			break
		}
	}

	if len(config.Instance.Teams) > 0 {
		groupsLog.Warn("could not find the following groups in EntraID", "skipped", strings.Join(config.Instance.Teams, ","))
	}

	groupsLog.Info("finished processing EntraID groups",
		slog.Group("groups",
			slog.Int64("found", *countFound),
			slog.Int("skipped", len(config.Instance.Teams)),
		),
	)

	return g.grafanaTeams
}
