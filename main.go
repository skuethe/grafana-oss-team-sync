package main

import (
	"log"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-team-sync/internal/grafana"
)

// ref:
// https://github.com/grafana/grafana-foundation-sdk/blob/main/examples/go/grafana-openapi-client-go/main.go

// create user:
// https://pkg.go.dev/github.com/grafana/grafana-openapi-client-go@v0.0.0-20250428202209-be3a35ff1dac/client/admin_users#Client.AdminCreateUser

// create team:
// https://pkg.go.dev/github.com/grafana/grafana-openapi-client-go@v0.0.0-20250428202209-be3a35ff1dac/client/teams#Client.CreateTeam

// add team member:
// https://pkg.go.dev/github.com/grafana/grafana-openapi-client-go@v0.0.0-20250428202209-be3a35ff1dac/client/teams#Client.AddTeamMember

func main() {
	log.SetFlags(log.Ldate + log.Ltime + log.Lmsgprefix)
	log.Println("Running Grafana Team Sync")

	api := grafana.InitClient()

	teamList := []models.CreateTeamCommand{
		{
			Name: strings.ToLower("testTeam"),
		},
		{
			Name:  strings.ToLower("testTeam2"),
			Email: "test2@test.com",
		},
	}

	userList := []models.AdminCreateUserForm{
		{
			Email:    strings.ToLower("test@test.com"),
			Login:    strings.ToLower("testUser"),
			Name:     "testUser",
			Password: "testPassword",
		},
		{
			Email:    strings.ToLower("second@test.com"),
			Login:    strings.ToLower("secondUser"),
			Name:     "secondUserName",
			Password: "secondPassword",
		},
	}

	folderList := []models.CreateFolderCommand{
		{
			UID:         strings.ToLower("testFolder"),
			Title:       "testFolder",
			Description: "A folder to test things",
		},
		{
			UID:         "testFolder",
			Title:       "secondTest",
			Description: "A second folder to test things",
		},
		{
			UID:         "testfolder",
			Title:       "thirdTest",
			Description: "A third folder to test things",
		},
	}

	grafana.ProcessUsers(api, userList)
	grafana.ProcessTeams(api, teamList)
	grafana.ProcessFolders(api, folderList)

}
