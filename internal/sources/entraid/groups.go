package entraid

import (
	"context"
	"log/slog"
	"os"
	"strings"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphgroups "github.com/microsoftgraph/msgraph-sdk-go/groups"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/helpers"
	"github.com/skuethe/grafana-oss-team-sync/internal/plugin"
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

func (g *groups) getUsersFromGroup(groupID string) (models.UserCollectionResponseable, error) {

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	// requestTop := int32(5)
	requestCount := true
	requestParams := &graphgroups.ItemTransitiveMembersGraphUserRequestBuilderGetQueryParameters{
		Select: []string{"userPrincipalName", "displayName", "mail"},
		Count:  &requestCount,
		// Top: &requestTop,
	}
	configuraton := &graphgroups.ItemTransitiveMembersGraphUserRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParams,
	}
	result, err := g.client.Groups().ByGroupId(groupID).TransitiveMembers().GraphUser().Get(context.Background(), configuraton)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ProcessGroups(instance *plugin.SourceInstance) *grafana.Teams {
	groupsLog := slog.With(slog.String("package", "entraid.groups"))
	groupsLog.Info("processing entraid groups")

	var grafanaTeamList *grafana.Teams = &grafana.Teams{}
	var grafanaUserList *grafana.Users = &grafana.Users{}

	g := groups{
		client:        instance.EntraID,
		requestFilter: "displayName in ('" + strings.Join(config.Teams.List, "', '") + "')",
	}

	groupList, err := g.getGroupData()
	if err != nil {
		groupsLog.Error("could not get group results from EntraID", "error", err)
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
		groupLog.Info("found EntraID group")

		// Get all users from current group
		userList, err := g.getUsersFromGroup(groupId)
		if err != nil {
			groupLog.Error("could not get user results from EntraID", "error", err)
			os.Exit(1)
		}

		grafanaUserList = ProcessUsers(userList)

		*grafanaTeamList = append(*grafanaTeamList, grafana.Team{
			Parameter: &grafana.TeamParameter{
				Name:  groupDisplayName,
				Email: mail,
			},
			Users: grafanaUserList,
		})
		config.Teams.List = helpers.RemoveFromSlice(config.Teams.List, groupDisplayName, false)
	}

	if len(config.Teams.List) > 0 {
		groupsLog.Warn("could not find the following groups in EntraID", "skipped", strings.Join(config.Teams.List, ","))
	}

	groupsLog.Info(
		"finished processing EntraID groups",
		slog.Group("stats",
			slog.Int64("found", countFound),
			slog.Int("skipped", len(config.Teams.List)),
		),
	)

	return grafanaTeamList
}
