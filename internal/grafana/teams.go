package grafana

import (
	"log"

	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
)

func doesTeamExist(c *goapi.GrafanaHTTPAPI, name *string) bool {
	result, err := c.Teams.SearchTeams(&teams.SearchTeamsParams{
		Name: name,
	})
	if err != nil {
		log.Fatal(err)
	}
	return len(result.Payload.Teams) > 0
}

func createTeam(c *goapi.GrafanaHTTPAPI, team models.CreateTeamCommand) {
	_, err := c.Teams.CreateTeam(&models.CreateTeamCommand{
		Name:  team.Name,
		Email: team.Email,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("  Created: \"%v\"", team.Name)
	}
}

func ProcessTeams(c *goapi.GrafanaHTTPAPI, teamList []models.CreateTeamCommand) {
	log.SetPrefix("[Grafana.Teams] ")
	log.Println("Processing ...")

	countSkipped := 0
	countCreated := 0

	for _, team := range teamList {
		if doesTeamExist(c, &team.Name) {
			countSkipped++
		} else {
			createTeam(c, team)
			countCreated++
		}
	}
	log.Printf("Created: %v; Skipped: %v", countCreated, countSkipped)
}
