package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
)

// ref:
// https://github.com/grafana/grafana-foundation-sdk/blob/main/examples/go/grafana-openapi-client-go/main.go

// add team member:
// https://pkg.go.dev/github.com/grafana/grafana-openapi-client-go@v0.0.0-20250428202209-be3a35ff1dac/client/teams#Client.AddTeamMember

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	// log.SetFlags(log.Ldate + log.Ltime + log.Lmsgprefix)
	slog.Info("Running Grafana Team Sync")

	// Temporary data inputs
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

	api := grafana.InitClient()

	grafana.ProcessUsers(api, userList)
	grafana.ProcessTeams(api, teamList)
	grafana.ProcessFolders(api, folderList)

}
