// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

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

func main() {
	// Handle Flags
	flags.Load()

	// Handle version printing
	if flags.Version {
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
		),
	)
	// Print build details for debug output
	slog.Info("build details",
		slog.Group("build",
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

	// Grafana: continue to process users
	grafanaTeamList.ProcessUsers()

	// Grafana: continue to process teams
	grafanaTeamList.ProcessTeams()

	// Grafana: continue to process folders
	grafana.Instance.ProcessFolders()

	slog.Info("finished Grafana OSS Team Sync")
}
