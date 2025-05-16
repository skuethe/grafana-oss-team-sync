package azure

import (
	"context"
	"log/slog"
	"os"
	"strings"

	grafanaModels "github.com/grafana/grafana-openapi-client-go/models"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphgroups "github.com/microsoftgraph/msgraph-sdk-go/groups"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/helpers"
)

type groups struct {
	client        *graph.GraphServiceClient
	requestFilter string
}

func (g *groups) getGroupData() (models.GroupCollectionResponseable, error) {

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	// requestSearch := "\"displayName:" + a.name + "\""
	// requestTop := int32(5)
	requestCount := true
	requestFilter := g.requestFilter
	requestParams := &graphgroups.GroupsRequestBuilderGetQueryParameters{
		Filter: &requestFilter,
		Select: []string{"id", "displayName", "mail"},
		Count:  &requestCount,
		// Expand: []string{"members($select=id,displayName)"},
		// Top: &requestTop,
		// Search: &requestSearch,
	}
	configuraton := &graphgroups.GroupsRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParams,
	}
	result, err := g.client.Groups().Get(context.Background(), configuraton)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *AzureInstance) processGroups() {
	groupsLog := slog.With(slog.String("package", "azure.groups"))
	groupsLog.Info("processing Azure groups")

	teams := config.K.Strings("teams")

	g := groups{
		client:        a.api,
		requestFilter: "displayName in ('" + strings.Join(teams, "', '") + "')",
	}

	groupList, err := g.getGroupData()
	if err != nil {
		groupsLog.Error("could not get group results from Azure", "error", err)
		os.Exit(1)
	}

	countFound := *groupList.GetOdataCount()

	var grafanaTeamList []grafanaModels.CreateTeamCommand
	var groupIDList []string

	for _, group := range groupList.GetValue() {

		groupDisplayName := *group.GetDisplayName()
		groupId := *group.GetId()
		groupMail := group.GetMail()

		var mail string
		if groupMail != nil {
			mail = strings.ToLower(*groupMail)
		}

		groupLog := slog.With(
			slog.Group("group",
				slog.String("id", groupId),
				slog.String("displayname", groupDisplayName),
				slog.String("mail", mail),
			),
		)
		groupLog.Info("found Azure group")

		grafanaTeamList = append(grafanaTeamList, grafanaModels.CreateTeamCommand{
			Name:  groupDisplayName,
			Email: mail,
		})
		groupIDList = append(groupIDList, groupId)
		teams = helpers.RemoveFromSlice(teams, groupDisplayName, false)
	}

	if len(teams) > 0 {
		groupsLog.Warn("could not find the following groups in Azure", "skipped", strings.Join(teams, ","))
	}

	groupsLog.Info(
		"finished processing Azure groups",
		slog.Group("stats",
			slog.Int64("found", countFound),
			slog.Int("skipped", len(teams)),
		),
	)

	if len(groupIDList) > 0 {
		a.processUsers(groupIDList)
	}

	// grafana.Instance.ProcessTeams(&grafanaTeamList)
}
