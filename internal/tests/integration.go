// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package tests

import (
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
)

func IntegrationGrafana(teamList *grafana.Teams) error {

	// Load flags to not fail
	if err := flags.Load(); err != nil {
		return err
	}

	// 2. parse config
	if err := config.Load(); err != nil {
		return err
	}

	// 3. init Grafana
	if err := grafana.New(); err != nil {
		return err
	}

	// Run ProcessUsers
	teamList.ProcessUsers()

	// Run ProcessTeams
	teamList.ProcessTeams()

	// Run ProcessFolders
	grafana.Instance.ProcessFolders()

	return nil
}
