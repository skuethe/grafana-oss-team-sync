package grafana

import (
	"log/slog"

	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type Team models.CreateTeamCommand
type Teams []Team

func (t *Team) doesTeamExist() (bool, error) {
	result, err := Instance.api.Teams.SearchTeams(&teams.SearchTeamsParams{
		Name: &t.Name,
	})
	if err != nil {
		return false, err
	}
	return len(result.Payload.Teams) > 0, nil
}

func (t *Team) createTeam() error {
	_, err := Instance.api.Teams.CreateTeam(&models.CreateTeamCommand{
		Name:  t.Name,
		Email: t.Email,
	})
	if err != nil {
		return err
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
				slog.String("name", team.Name),
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
				} else {
					teamLog.Info("created Grafana team")
					countCreated++
				}
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
