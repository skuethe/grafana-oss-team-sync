package sources

import (
	"log/slog"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/plugin"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/entraid"
)

func CallPlugin() *grafana.Teams {
	pluginLog := slog.With(slog.String("package", "sources"))
	pluginLog.Info("processing source plugin")

	var instance *plugin.SourceInstance
	var grafanaTeamList *grafana.Teams = &grafana.Teams{}

	// Execute specific source plugin, which need to return a *grafana.Teams instance
	switch config.GetSource() {
	case config.SourceEntraID:
		// EntraID: create new msgraph client
		instance = entraid.New()
		// EntraID: search for all specified groups and users
		grafanaTeamList = entraid.ProcessGroups(instance)
	}

	pluginLog.Info("finished processing source plugin")

	return grafanaTeamList
}
