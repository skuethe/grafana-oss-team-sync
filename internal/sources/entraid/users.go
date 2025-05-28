package entraid

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/plugin"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphgroups "github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type users struct {
	client        *graph.GraphServiceClient
	requestSelect []string
	fromGroup     string
}

func (u *users) getUserData() (models.UserCollectionResponseable, error) {

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	// requestTop := int32(5)
	requestCount := true
	requestParams := &graphgroups.ItemTransitiveMembersGraphUserRequestBuilderGetQueryParameters{
		Select: u.requestSelect,
		Count:  &requestCount,
		// Top: &requestTop,
	}
	configuraton := &graphgroups.ItemTransitiveMembersGraphUserRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParams,
	}
	result, err := u.client.Groups().ByGroupId(u.fromGroup).TransitiveMembers().GraphUser().Get(context.Background(), configuraton)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ProcessUsers(instance *plugin.SourceInstance, fromGroupList []string) *grafana.Users {
	usersLog := slog.With(slog.String("package", "entraid.users"))
	usersLog.Info("processing EntraID users")

	var grafanaUserList *grafana.Users = &grafana.Users{}

	globalDuplicateUsers := 0

	for _, groupID := range fromGroupList {
		groupLog := slog.With(
			slog.Group("group",
				slog.String("id", groupID),
			),
		)
		groupLog.Info("searching for users from group")

		usersAdded := 0
		usersDuplicate := 0

		u := users{
			client:        instance.EntraID,
			requestSelect: []string{"userPrincipalName", "displayName", "mail"},
			fromGroup:     groupID,
		}

		userList, err := u.getUserData()
		if err != nil {
			groupLog.Error("could not get user results from EntraID", "error", err)
			os.Exit(1)
		}

		countFound := *userList.GetOdataCount()
		groupLog.Info("processing users from group",
			slog.Group("users",
				slog.Int64("found", countFound),
			),
		)

		for _, user := range userList.GetValue() {
			userDisplayName := *user.GetDisplayName()
			userPrincipalName := *user.GetUserPrincipalName()
			userMail := user.GetMail()

			var mail string
			if userMail != nil {
				mail = strings.ToLower(*userMail)
			}

			userLog := groupLog.With(
				slog.Group("user",
					slog.String("principalname", userPrincipalName),
					slog.String("displayname", userDisplayName),
					slog.String("mail", mail),
				),
			)
			userLog.Debug("found EntraID user")

			grafanaUser := grafana.User{
				Login: userPrincipalName,
				Name:  userDisplayName,
				Email: mail,
			}
			userExists := false
			if slices.Contains(*grafanaUserList, grafanaUser) {
				userExists = true
				userLog.Debug("skipping duplicate user")
				usersDuplicate++
			}
			if !userExists {
				*grafanaUserList = append(*grafanaUserList, grafanaUser)
				usersAdded++
			}
		}

		globalDuplicateUsers += usersDuplicate

		groupLog.Info(
			"finished processing users from group",
			slog.Group("stats",
				slog.Int("unique", usersAdded),
				slog.Int("duplicate", usersDuplicate),
			),
		)

	}

	usersLog.Info(
		"finished processing EntraID users",
		slog.Group("stats",
			slog.Int("unique", len(*grafanaUserList)),
			slog.Int("duplicate", globalDuplicateUsers),
		),
	)

	return grafanaUserList
}
