package azure

import (
	"context"
	"log/slog"

	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphgroups "github.com/microsoftgraph/msgraph-sdk-go/groups"

	// "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
)

type group struct {
	client *graph.GraphServiceClient
	log    slog.Logger
	name   string
}

func (a *group) doesGroupExist() bool {
	requestFilter := "displayName eq " + a.name
	requestParams := &graphgroups.GroupsRequestBuilderGetQueryParameters{
		Filter: &requestFilter,
		Select: []string{"id", "displayName", "mail"},
	}
	configuraton := &graphgroups.GroupsRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParams,
	}
	groups, err := a.client.Groups().Get(context.Background(), configuraton)
	if groups != nil {
		a.log.Info("Found Azure Groups", "DEBUG", groups)
	}

	return err == nil
}

func (a *AzureInstance) processGroups() {
	groupsLog := slog.With(slog.String("package", "azure.groups"))
	groupsLog.Info("Initializing Azure Groups")

	countFound := 0
	countIgnored := 0

	// add logic

	teams := config.K.MapKeys("teams")
	for _, teamName := range teams {
		groupLog := slog.With(
			slog.Group("group",
				slog.String("name", teamName),
			),
		)
		groupLog.Info("Processing Azure Group")

		g := group{
			client: a.api,
			log:    *groupLog,
			name:   teamName,
		}

		if g.doesGroupExist() {
			countFound++
		} else {
			countIgnored++
			groupLog.Warn("Group not found - ignoring")
		}

		// teamList = append(teamList, models.CreateTeamCommand{
		// 	Name:  teamName,
		// 	Email: strings.ToLower(config.K.String("teams." + teamName + ".email")),
		// })
	}

	groupsLog.Info(
		"Finished Azure Groups",
		slog.Group("stats",
			slog.Int("found", countFound),
			slog.Int("ignored", countIgnored),
		),
	)

}
