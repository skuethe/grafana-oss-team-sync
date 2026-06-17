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

func (g *groups) getInitialPrefixRequest() (*models.GroupCollectionResponseable, error) {

	// Build an OData filter that matches any group whose displayName starts with
	// one of the configured prefixes, e.g. "startswith(displayName, 'a') or startswith(displayName, 'b')"
	prefixFilters := make([]string, len(config.Instance.TeamPrefixes))
	for i, prefix := range config.Instance.TeamPrefixes {
		prefixFilters[i] = "startswith(displayName, '" + prefix + "')"
	}

	requestCount := true
	requestFilter := strings.Join(prefixFilters, " or ")
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

// resolveGroupPrefixes lists every EntraID group matching one of the configured
// prefixes and merges their display names into config.Instance.Teams. The merge
// is unique (case-insensitive) so groups already listed explicitly in "teams" -
// or matched by multiple prefixes - are not synced twice.
func (g *groups) resolveGroupPrefixes(groupsLog *slog.Logger) {
	groupsLog.Info("resolving EntraID groups by prefix", "prefixes", strings.Join(config.Instance.TeamPrefixes, ","))

	gr, err := g.getInitialPrefixRequest()
	if err != nil {
		groupsLog.Error("could not get initial group prefix result from EntraID")
		panic(err)
	}

	for {
		for _, group := range (*gr).GetValue() {
			groupDisplayName := group.GetDisplayName()
			if groupDisplayName != nil {
				config.Instance.Teams = helpers.AppendUnique(config.Instance.Teams, *groupDisplayName)
			}
		}

		nextPageUrl := (*gr).GetOdataNextLink()
		if nextPageUrl != nil {
			groupsLog.Debug("processing paginated group prefix result")
			gr, err = g.handleGroupPagination(nextPageUrl)
			if err != nil {
				groupsLog.Error("could not get paged group prefix result from EntraID")
				panic(err)
			}
		} else {
			break
		}
	}
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

	// Drop empty entries (e.g. introduced by default flag values) so we never
	// build an empty source filter
	config.Instance.Teams = helpers.RemoveEmptyStrings(config.Instance.Teams)
	config.Instance.TeamPrefixes = helpers.RemoveEmptyStrings(config.Instance.TeamPrefixes)

	// Resolve any configured prefixes into concrete group names and merge them
	// (uniquely) into the list of teams to sync
	if len(config.Instance.TeamPrefixes) > 0 {
		g.resolveGroupPrefixes(groupsLog)
	}

	// Nothing to do if neither an explicit team nor a prefix match remains
	if len(config.Instance.Teams) == 0 {
		groupsLog.Warn("no teams left to process after resolving prefixes")
		return g.grafanaTeams
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
