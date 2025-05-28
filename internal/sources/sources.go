package sources

import (
	"log/slog"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/plugin"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/entraid"
)

func CallPlugin() {
	pluginLog := slog.With(slog.String("package", "entraid.groups"))
	pluginLog.Info("processing entraid groups")

	var instance *plugin.SourceInstance
	var grafanaTeamList *grafana.Teams

	// Execute specific source plugin, which need to return a *grafana.Teams instance
	switch config.GetSource() {
	case config.SourceEntraID:
		instance = entraid.New()
		grafanaTeamList = entraid.ProcessGroups(instance)
	}

	// Continue to process available data against our Grafana instance
	if len(*grafanaTeamList) > 0 {
		// ProcessUsers(instance, groupIDList)
		grafanaTeamList.ProcessTeams()
	} else {
		pluginLog.Warn("no groups to process, skipping Grafana teams package")
	}

}
