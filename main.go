package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
)

// ref:
// https://github.com/grafana/grafana-foundation-sdk/blob/main/examples/go/grafana-openapi-client-go/main.go

// add team member:
// https://pkg.go.dev/github.com/grafana/grafana-openapi-client-go@v0.0.0-20250428202209-be3a35ff1dac/client/teams#Client.AddTeamMember

func main() {
	// Initialize logger
	loggerLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))
	slog.SetDefault(logger)
	slog.Info("Running Grafana Team Sync")

	// Initialize config
	config.Start()

	// Further configure logger
	logLevelFromConfig := config.GetLogLevel()
	loggerLevel.Set(logLevelFromConfig)
	slog.SetDefault(logger)

	//
	//
	// Temporary data inputs
	//
	//

	var teamList []models.CreateTeamCommand
	teams := config.K.MapKeys("teams")
	for _, teamName := range teams {
		teamList = append(teamList, models.CreateTeamCommand{
			Name:  teamName,
			Email: strings.ToLower(config.K.String("teams." + teamName + ".email")),
		})
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

	var folderList []models.CreateFolderCommand
	folders := config.K.MapKeys("folders")
	for _, folderName := range folders {
		folderList = append(folderList, models.CreateFolderCommand{
			UID:         strings.ToLower(folderName),
			Title:       config.K.MustString("folders." + folderName + ".title"),
			Description: config.K.String("folders." + folderName + ".description"),
		})
	}

	//
	//
	//
	//

	// Run Grafana related packages
	grafana.Start(&userList, &teamList, &folderList)
}
