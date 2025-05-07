package grafana

import (
	"log/slog"
	"os"

	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type TeamType struct {
	client *goapi.GrafanaHTTPAPI
	log    slog.Logger
	form   models.CreateTeamCommand
}

func (t *TeamType) doesTeamExist() bool {
	result, err := t.client.Teams.SearchTeams(&teams.SearchTeamsParams{
		Name: &t.form.Name,
	})
	if err != nil {
		t.log.Error("Could not search for Team", "error", err)
		os.Exit(1)
	}
	return len(result.Payload.Teams) > 0
}

func (t *TeamType) createTeam() {
	_, err := t.client.Teams.CreateTeam(&models.CreateTeamCommand{
		Name:  t.form.Name,
		Email: t.form.Email,
	})
	if err != nil {
		t.log.Error("Could not create Team", "error", err)
	} else {
		t.log.Info(
			"Created Team",
			slog.Group("team",
				slog.String("name", t.form.Name),
			),
		)
	}
}

func (g *GrafanaInstance) ProcessTeams(teamList *[]models.CreateTeamCommand) {
	teamsLog := slog.With(slog.String("package", "grafana.teams"))
	teamsLog.Info("Processing Teams")

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
				"Skipped Team",
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
		"Finished Teams",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
