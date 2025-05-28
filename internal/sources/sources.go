package sources

import (
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/plugin"
	"github.com/skuethe/grafana-oss-team-sync/internal/sources/entraid"
)

func LoadPlugin() {
	var instance *plugin.SourceInstance

	switch config.GetSource() {
	case config.SourceEntraID:
		instance = entraid.New()
		entraid.ProcessGroups(instance)
	}

}
