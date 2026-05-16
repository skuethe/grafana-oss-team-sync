// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package tests

import (
	"context"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources"
)

func EndToEnd() error {
	// Load flags to not fail
	if err := flags.Load(); err != nil {
		return err
	}

	// Parse config
	if err := config.Load(); err != nil {
		return err
	}

	// Init Grafana
	if err := grafana.New(); err != nil {
		return err
	}

	// Call the configured source plugin
	grafanaTeamList := sources.CallPlugin(context.TODO())

	// Grafana: continue to process users
	grafanaTeamList.ProcessUsers()

	// Grafana: continue to process teams
	grafanaTeamList.ProcessTeams()

	// Grafana: continue to process folders
	grafana.Instance.ProcessFolders()

	return nil
}
