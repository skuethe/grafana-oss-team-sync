package entraid

import (
	"log/slog"
	"strings"

	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func ProcessUsers(userCollection models.UserCollectionResponseable) *grafana.Users {
	usersLog := slog.With(slog.String("package", "entraid.users"))

	var grafanaUserList *grafana.Users = &grafana.Users{}

	countFound := userCollection.GetOdataCount()
	usersLog.Info("processing EntraID users for group",
		slog.Group("users",
			slog.Int64("found", *countFound),
		),
	)

	for _, user := range userCollection.GetValue() {
		userDisplayName := *user.GetDisplayName()
		userPrincipalName := *user.GetUserPrincipalName()
		userMail := user.GetMail()

		var mail string
		if userMail != nil {
			mail = strings.ToLower(*userMail)
		}

		userLog := usersLog.With(
			slog.Group("user",
				slog.String("principalname", userPrincipalName),
				slog.String("displayname", userDisplayName),
				slog.String("mail", mail),
			),
		)
		userLog.Debug("found EntraID user")

		*grafanaUserList = append(*grafanaUserList, grafana.User{
			Login: userPrincipalName,
			Name:  userDisplayName,
			Email: mail,
		})
	}

	usersLog.Info(
		"finished processing EntraID users for group",
		slog.Group("stats",
			slog.Int("unique", len(*grafanaUserList)),
		),
	)

	return grafanaUserList
}
