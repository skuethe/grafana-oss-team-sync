// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import (
	"log/slog"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/config/configtypes"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/entraid"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/ldap"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/sourcetypes"
)

func CallPlugin() *grafana.Teams {
	pluginLog := slog.With(slog.String("package", "sources"))

	var instance *sourcetypes.SourcePlugin
	grafanaTeamList := &grafana.Teams{}

	if len(config.Instance.Teams) == 0 {
		pluginLog.Info("your team input is empty, skipping")
	} else {
		pluginLog.Info("processing source plugin")

		// Execute specific source plugin, which need to return a *grafana.Teams instance
		switch config.Instance.GetSource() {
		case configtypes.SourcePluginEntraID:
			// EntraID: create new msgraph client
			instance = entraid.New()
			// EntraID: search for all specified groups and users
			grafanaTeamList = entraid.ProcessGroups(instance)
		case configtypes.SourcePluginLDAP:
			instance = ldap.New()
			grafanaTeamList = ldap.ProcessGroups(instance)
		}

		pluginLog.Info("finished processing source plugin")
	}

	return grafanaTeamList
}
