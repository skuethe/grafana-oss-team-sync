package grafana

import (
	"log/slog"
	"strconv"

	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
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
	return len(result.Payload.Teams) == 1, nil
}

func (t *Team) getTeamID() (*int64, error) {
	result, err := t.searchTeam()
	if err != nil {
		return nil, err
	}
	return &result.Payload.Teams[0].ID, nil
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

func (t *Team) addUsersToTeam() (*[]string, error) {
	teamID, tErr := t.getTeamID()
	if tErr != nil {
		return nil, tErr
	}

	adminMemberList := &[]string{}
	teamMemberList := &[]string{}

	if config.Instance.Features.AddLocalAdminToTeams {
		*adminMemberList = append(*adminMemberList, "admin@localhost")
	}

	for _, user := range *t.Users {
		*teamMemberList = append(*teamMemberList, user.Email)
	}

	_, mErr := Instance.api.Teams.SetTeamMemberships(strconv.FormatInt(*teamID, 10), &models.SetTeamMembershipsCommand{
		Admins:  *adminMemberList,
		Members: *teamMemberList,
	})
	if mErr != nil {
		return nil, mErr
	}

	return teamMemberList, nil
}

func (t *Teams) ProcessTeams() {
	teamsLog := slog.With(slog.String("package", "grafana.teams"))

	if len(*t) == 0 {
		teamsLog.Info("no teams to process, skipping")
	} else {
		teamsLog.Info("processing teams")

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
				teamLog.Error("could not search for team",
					slog.Any("error", err),
				)
			} else {
				if exists {
					countSkipped++
					teamLog.Debug("skipping already existing team")
				} else {
					err := team.createTeam()
					if err != nil {
						teamLog.Error("could not create team",
							slog.Any("error", err),
						)
						continue
					} else {
						teamLog.Info("created team")
						countCreated++
					}
				}
				// Add users to team, if userSync feature is enabled
				if !config.Instance.Features.DisableUserSync {
					teamLog.Info("processing team members")

					if userList, err := team.addUsersToTeam(); err != nil {
						teamLog.Error("could not add users to team",
							slog.Any("error", err),
						)
					} else {
						teamLog.Debug("added users to team",
							slog.Any("list", *userList),
						)
					}

					teamLog.Info("finished processing team members")
				}
			}
		}

		teamsLog.Info("finished processing teams",
			slog.Group("teams",
				slog.Int("created", countCreated),
				slog.Int("existing", countSkipped),
			),
		)
	}
}
