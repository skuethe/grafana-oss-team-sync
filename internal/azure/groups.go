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

func (s *groups) getData() (models.GroupCollectionResponseable, error) {

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	// requestSearch := "\"displayName:" + a.name + "\""
	// requestTop := int32(5)
	requestCount := true
	requestFilter := s.requestFilter
	requestParams := &graphgroups.GroupsRequestBuilderGetQueryParameters{
		Filter: &requestFilter,
		Select: []string{"id", "displayName", "mail"},
		Count:  &requestCount,
		// Top: &requestTop,
		// Search: &requestSearch,
		// Expand: []string{"members($select=id,displayName)"},
	}
	configuraton := &graphgroups.GroupsRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParams,
	}
	result, err := s.client.Groups().Get(context.Background(), configuraton)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *azureInstance) processGroups() (grafanaTeamList []*grafanaModels.CreateTeamCommand) {
	groupsLog := slog.With(slog.String("package", "azure.groups"))
	groupsLog.Info("Initializing Azure Groups")

	teams := config.K.Strings("teams")

	g := groups{
		client:        a.api,
		requestFilter: "displayName in ('" + strings.Join(teams, "', '") + "')",
	}

	groupList, err := g.getData()
	if err != nil {
		groupsLog.Error("Could not get Group results from Azure", "error", err)
		os.Exit(1)
	}

	countFound := *groupList.GetOdataCount()

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
		groupLog.Info("Processing Azure Group")

		grafanaTeamList = append(grafanaTeamList, &grafanaModels.CreateTeamCommand{
			Name:  groupDisplayName,
			Email: mail,
		})
		teams = helpers.RemoveFromSlice(teams, groupDisplayName, false)
	}

	if len(teams) > 0 {
		groupsLog.Warn("Could not find the following groups in Azure", "ignored", strings.Join(teams, ","))
	}

	groupsLog.Info(
		"Finished Azure Groups",
		slog.Group("stats",
			slog.Int64("found", countFound),
			slog.Int("ignored", len(teams)),
		),
	)

	return
}
