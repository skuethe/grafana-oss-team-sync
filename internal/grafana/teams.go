package grafana

import (
	"log/slog"
	"os"

	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type team struct {
	client *client.GrafanaHTTPAPI
	log    slog.Logger
	form   models.CreateTeamCommand
}

func (t *team) doesTeamExist() bool {
	result, err := t.client.Teams.SearchTeams(&teams.SearchTeamsParams{
		Name: &t.form.Name,
	})
	if err != nil {
		t.log.Error("could not search for Grafana team", "error", err)
		os.Exit(1)
	}
	t.log.Debug("team search results", "searchName", t.form.Name, "searchResult", result)
	return len(result.Payload.Teams) > 0
}

func (t *team) createTeam() {
	_, err := t.client.Teams.CreateTeam(&models.CreateTeamCommand{
		Name:  t.form.Name,
		Email: t.form.Email,
	})
	if err != nil {
		t.log.Error("could not create Grafana team", "error", err)
	} else {
		t.log.Info("created Grafana team")
	}
}

func (g *GrafanaInstance) ProcessTeams(teamList *[]models.CreateTeamCommand) {
	teamsLog := slog.With(slog.String("package", "grafana.teams"))
	teamsLog.Info("processing Grafana teams")

	countSkipped := 0
	countCreated := 0

	for _, instance := range *teamList {

		teamLog := slog.With(
			slog.Group("team",
				slog.String("name", instance.Name),
			),
		)

		t := team{
			client: g.api,
			log:    *teamLog,
			form:   instance,
		}

		if t.doesTeamExist() {
			countSkipped++
			teamLog.Debug("skipped Grafana team")
		} else {
			t.createTeam()
			countCreated++
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
