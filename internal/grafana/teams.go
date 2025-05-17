package grafana

import (
	"log/slog"

	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type Teams struct {
	Teams *[]Team
}

type Team models.CreateTeamCommand

func doesTeamExist(team Team) (bool, error) {
	result, err := Instance.api.Teams.SearchTeams(&teams.SearchTeamsParams{
		Name: &team.Name,
	})
	if err != nil {
		return false, err
	}
	return len(result.Payload.Teams) > 0, nil
}

func createTeam(team *Team) error {
	_, err := Instance.api.Teams.CreateTeam(team)
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

	for _, instance := range *t.Teams {

		teamLog := slog.With(
			slog.Group("team",
				slog.String("name", instance.Name),
			),
		)

		exists, err := doesTeamExist(instance)
		if err != nil {
			teamLog.Error("could not search for Grafana team", "error", err)
		} else {
			if exists {
				countSkipped++
				teamLog.Debug("skipped Grafana team")
			} else {
				err := createTeam(instance)
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
