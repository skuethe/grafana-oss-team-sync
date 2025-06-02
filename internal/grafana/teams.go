package grafana

import (
	"log/slog"

	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type TeamParameter models.CreateTeamCommand
type Team struct {
	Parameter *TeamParameter
	Users     *Users
}
type Teams []Team

func (t *Team) searchTeam() (*teams.SearchTeamsOK, error) {
	result, err := Instance.api.Teams.SearchTeams(&teams.SearchTeamsParams{
		Name: &t.Parameter.Name,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *Team) doesTeamExist() (bool, error) {
	result, err := t.searchTeam()
	if err != nil {
		return false, err
	}
	return len(result.Payload.Teams) > 0, nil
}

func (t *Team) getTeamUID() (*string, error) {
	result, err := t.searchTeam()
	if err != nil {
		return nil, err
	}
	return &result.Payload.Teams[0].UID, nil
}

func (t *Team) createTeam() error {
	_, err := Instance.api.Teams.CreateTeam(&models.CreateTeamCommand{
		Name:  t.Parameter.Name,
		Email: t.Parameter.Email,
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *Team) addUsersToTeam() error {
	teamID, getIDErr := t.getTeamUID()
	if getIDErr != nil {
		return getIDErr
	}
	// TODO: call GetTeamMembers on the team id and in the loop only run "AddTeamMember" if user is not yet part of it
	// ref: https://pkg.go.dev/github.com/grafana/grafana-openapi-client-go@v0.0.0-20250516123951-83fcd32d7bbe/client/teams#Client.GetTeamMembers
	for _, user := range *t.Users {
		userID, getUserIDErr := user.getUserID()
		if getUserIDErr != nil {
			slog.Error("could not get UID for user", "error", getUserIDErr)
			continue
		}
		// TODO: eval if using "SetTeamMemberships" is better here...
		// ref: https://pkg.go.dev/github.com/grafana/grafana-openapi-client-go@v0.0.0-20250428202209-be3a35ff1dac/client/teams#Client.SetTeamMemberships
		_, userAddErr := Instance.api.Teams.AddTeamMember(*teamID, &models.AddTeamMemberCommand{
			UserID: *userID,
		})
		if userAddErr != nil {
			slog.Error("could not add user to team", "error", userAddErr)
			continue
		}
		slog.Debug("added Grafana user to Grafana team")
	}
	return nil
}

func (t *Teams) ProcessTeams() {
	teamsLog := slog.With(slog.String("package", "grafana.teams"))
	teamsLog.Info("processing Grafana teams")

	countSkipped := 0
	countCreated := 0

	for _, team := range *t {

		teamLog := slog.With(
			slog.Group("team",
				slog.String("name", team.Parameter.Name),
			),
		)

		exists, err := team.doesTeamExist()
		if err != nil {
			teamLog.Error("could not search for Grafana team", "error", err)
		} else {
			if exists {
				countSkipped++
				teamLog.Debug("skipped Grafana team")
			} else {
				err := team.createTeam()
				if err != nil {
					teamLog.Error("could not create Grafana team", "error", err)
					continue
				} else {
					teamLog.Info("created Grafana team")
					countCreated++
				}
			}
			err := team.addUsersToTeam()
			if err != nil {
				teamLog.Error("could not add Grafana users to Grafana team", "error", err)
			}
		}
	}
	teamsLog.Info(
		"finished processing Grafana teams",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
