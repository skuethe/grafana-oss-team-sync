package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// app access needed?
// https://learn.microsoft.com/en-gb/graph/permissions-reference?view=graph-rest-1.0#groupmemberreadall
// https://learn.microsoft.com/en-gb/graph/permissions-reference?view=graph-rest-1.0#userreadbasicall
// current app has:
//  - User.Read
//  - Directory.Read.All

func main() {
	// Handle Flags
	flags.LoadFlags()
	flags.FlagSet.Parse(os.Args[1:])

	// Handle version printing
	if flags.FlagInputVersion {
		fmt.Printf("Version: %s\nCommit: %s\nBuild date: %s\n", version, commit, date)
		os.Exit(0)
	}

	// Initialize logger
	loggerLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))
	slog.SetDefault(logger)
	slog.Info("starting Grafana OSS Team Sync",
		slog.Group("release",
			slog.String("version", version),
			slog.String("commit", commit),
			slog.String("date", date),
		),
	)

	// Initialize config
	config.Load()

	// Further configure logger
	logLevelFromConfig := config.Instance.GetLogLevel()
	loggerLevel.Set(logLevelFromConfig)
	slog.SetDefault(logger)

	// Initialize Grafana
	grafana.New()

	// Call the configured source plugin
	grafanaTeamList := sources.CallPlugin()

	// Grafana: continue to process teams
	if len(*grafanaTeamList) > 0 {
		if !config.Instance.Features.DisableUserSync {
			grafanaTeamList.ProcessUsers()
		} else {
			slog.Info("user sync disabled in config, skipping Grafana users package")
		}
		grafanaTeamList.ProcessTeams()
	} else {
		slog.Warn("no groups to process, skipping Grafana teams package")
	}

	// Grafana: continue to process folders
	if !config.Instance.Features.DisableFolders {
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

	slog.Info("finished Grafana OSS Team Sync")
}
