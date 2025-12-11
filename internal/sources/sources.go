// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import (
	"context"
	"log/slog"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/config/configtypes"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/entraid"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/keycloak"
)

type SourcerPlugin interface {
	ProcessGroups(ctx context.Context) *grafana.Teams
}

func CallPlugin(ctx context.Context) *grafana.Teams {
	pluginLog := slog.With(slog.String("package", "sources"))

	var instance SourcerPlugin
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
		case configtypes.SourcePluginKeycloak:
			var err error
			instance, err = keycloak.New(ctx)
			if err != nil {
				pluginLog.Error("unable to create keycloak instance")
			}
		default:
			pluginLog.Error("invalid source plugin defined")
		}
		grafanaTeamList = instance.ProcessGroups(ctx)
		pluginLog.Info("finished processing source plugin")
	}

	return grafanaTeamList
}
