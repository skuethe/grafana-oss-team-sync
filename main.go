package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources"
)

//
// Entra ID
//

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
	slog.Info("starting Grafana OSS Team Sync")

	// Handle .env files
	godotenv.Load()

	// Initialize config
	config.Load()

	// Further configure logger
	logLevelFromConfig := config.GetLogLevel()
	loggerLevel.Set(logLevelFromConfig)
	slog.SetDefault(logger)

	// Initialize Grafana
	grafana.New()

	// Call the configured source plugin
	grafanaTeamList, grafanaUserList := sources.CallPlugin()

	// Grafana: continue to process teams
	if len(*grafanaTeamList) > 0 {
		grafanaTeamList.ProcessTeams()
	} else {
		slog.Warn("no groups to process, skipping Grafana teams package")
	}

	// Grafana: continue to process users
	if !config.Feature.DisableUserSync {
		if len(*grafanaUserList) > 0 {
			grafanaUserList.ProcessUsers()
		} else {
			slog.Warn("no users to process, skipping Grafana users package")
		}
	} else {
		slog.Info("user sync disabled in config, skipping Grafana users package")
	}

	// Grafana: continue to process folders
	if !config.Feature.DisableFolders {
		grafana.Instance.ProcessFolders()
	} else {
		slog.Info("folder feature disabled in config, skipping Grafana folders package")
	}

	//
	//
	// Temporary data inputs
	//
	//

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

}
