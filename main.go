package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/skuethe/grafana-oss-team-sync/internal/azure"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
)

//
// Entra ID
//

// client:
// https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go

// search groups
// https://learn.microsoft.com/en-gb/graph/api/group-list?view=graph-rest-1.0&tabs=go#optional-query-parameters

// list group members:
// https://learn.microsoft.com/en-gb/graph/api/group-list-members?view=graph-rest-1.0&tabs=go

// get user:
// https://learn.microsoft.com/en-gb/graph/api/user-get?view=graph-rest-1.0&tabs=http

// app access needed?
// https://learn.microsoft.com/en-gb/graph/permissions-reference?view=graph-rest-1.0#groupmemberreadall
// https://learn.microsoft.com/en-gb/graph/permissions-reference?view=graph-rest-1.0#userreadbasicall
// current app has:
//  - User.Read
//  - Directory.Read.All

func main() {
	// Initialize logger
	loggerLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))
	slog.SetDefault(logger)
	slog.Info("Running Grafana Team Sync")

	// Handle .env files
	godotenv.Load()

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

	// var teamList []models.CreateTeamCommand
	// teams := config.K.MapKeys("teams")
	// for _, teamName := range teams {
	// 	teamList = append(teamList, models.CreateTeamCommand{
	// 		Name:  teamName,
	// 		Email: strings.ToLower(config.K.String("teams." + teamName + ".email")),
	// 	})
	// }

	// userList := []models.AdminCreateUserForm{
	// 	{
	// 		Email:    strings.ToLower("test@test.com"),
	// 		Login:    strings.ToLower("testUser"),
	// 		Name:     "testUser",
	// 		Password: "testPassword",
	// 	},
	// 	{
	// 		Email:    strings.ToLower("second@test.com"),
	// 		Login:    strings.ToLower("secondUser"),
	// 		Name:     "secondUserName",
	// 		Password: "secondPassword",
	// 	},
	// }

	//
	//
	//
	//

	// Run Azure related packages
	azure.Start()

	// Run Grafana related packages
	// grafana.Start(&userList, &teamList)
}
