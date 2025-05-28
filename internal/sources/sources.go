package sources

import (
	"log/slog"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	"github.com/skuethe/grafana-oss-team-sync/internal/plugin"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/entraid"
)

func CallPlugin() (*grafana.Teams, *grafana.Users) {
	pluginLog := slog.With(slog.String("package", "sources"))
	pluginLog.Info("processing source plugin")

	var instance *plugin.SourceInstance
	var grafanaTeamList *grafana.Teams = &grafana.Teams{}
	var grafanaUserList *grafana.Users = &grafana.Users{}

	// Execute specific source plugin, which need to return a *grafana.Teams instance
	switch config.GetSource() {
	case config.SourceEntraID:
		var entraidGroupIDList []string

		// EntraID: create new msgraph client
		instance = entraid.New()

		// EntraID: search for all specified groups
		grafanaTeamList, entraidGroupIDList = entraid.ProcessGroups(instance)

		// EntraID: fetch all users which are a member of any found group
		if !config.Feature.DisableUserSync {
			grafanaUserList = entraid.ProcessUsers(instance, entraidGroupIDList)
		}
	}

	pluginLog.Info("finished processing source plugin")

	return grafanaTeamList, grafanaUserList
}
