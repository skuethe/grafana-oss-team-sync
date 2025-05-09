package grafana

import (
	"log/slog"
	"os"

	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type TeamType struct {
	client *client.GrafanaHTTPAPI
	log    slog.Logger
	form   models.CreateTeamCommand
}

func (t *TeamType) doesTeamExist() bool {
	result, err := t.client.Teams.SearchTeams(&teams.SearchTeamsParams{
		Name: &t.form.Name,
	})
	if err != nil {
		t.log.Error("Could not search for Grafana Team", "error", err)
		os.Exit(1)
	}
	t.log.Debug("Team search results", "search", "Name: "+t.form.Name, "result", result)
	return len(result.Payload.Teams) > 0
}

func (t *TeamType) createTeam() {
	_, err := t.client.Teams.CreateTeam(&models.CreateTeamCommand{
		Name:  t.form.Name,
		Email: t.form.Email,
	})
	if err != nil {
		t.log.Error("Could not create Grafana Team", "error", err)
	} else {
		t.log.Info(
			"Created Grafana Team",
			slog.Group("team",
				slog.String("name", t.form.Name),
			),
		)
	}
}

func (g *GrafanaInstance) processTeams(teamList *[]models.CreateTeamCommand) {
	teamsLog := slog.With(slog.String("package", "grafana.teams"))
	teamsLog.Info("Processing Grafana Teams")

	countSkipped := 0
	countCreated := 0

	for _, team := range *teamList {
		t := TeamType{
			client: g.api,
			log:    *teamsLog,
			form:   team,
		}
		if t.doesTeamExist() {
			countSkipped++
			teamsLog.Debug(
				"Skipped Grafana Team",
				slog.Group("team",
					slog.String("name", team.Name),
				),
			)
		} else {
			t.createTeam()
			countCreated++
		}
	}
	teamsLog.Info(
		"Finished Grafana Teams",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
