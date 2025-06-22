// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package flags

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/skuethe/grafana-oss-team-sync/internal/config/configtypes"
	flag "github.com/spf13/pflag"
)

var (
	Instance *flag.FlagSet

	Config  string
	Version bool
	Help    bool

	Token             string
	BasicAuthUsername string
	BasicAuthPassword string
)

func Load() {
	Instance = flag.NewFlagSet("grafana-oss-team-sync", flag.ExitOnError)
	Instance.SortFlags = false

	// Add "config" flag
	Instance.StringVarP(&Config, configtypes.ConfigParameter, configtypes.ConfigFlagShort, configtypes.ConfigDefault, configtypes.ConfigFlagHelp)

	// Add "authfile" flag
	Instance.StringP(configtypes.AuthFileParameter, configtypes.AuthFileFlagShort, configtypes.AuthFileDefault, configtypes.AuthFileFlagHelp)

	// Add "authtype" flag
	Instance.String(configtypes.GrafanaAuthTypeOptimized, configtypes.GrafanaAuthTypeDefault, configtypes.GrafanaAuthTypeFlagHelp)

	// Add basic auth flags
	Instance.StringVarP(&BasicAuthUsername, configtypes.GrafanaBasicAuthUsernameParameter, configtypes.GrafanaBasicAuthUsernameFlagShort, configtypes.GrafanaBasicAuthUsernameDefault, configtypes.GrafanaBasicAuthUsernameFlagHelp)
	Instance.StringVarP(&BasicAuthPassword, configtypes.GrafanaBasicAuthPasswordParameter, configtypes.GrafanaBasicAuthPasswordFlagShort, configtypes.GrafanaBasicAuthPasswordDefault, configtypes.GrafanaBasicAuthPasswordFlagHelp)

	// Add token auth flags
	Instance.StringVarP(&Token, configtypes.GrafanaTokenAuthParameter, configtypes.GrafanaTokenAuthFlagShort, configtypes.GrafanaTokenAuthDefault, configtypes.GrafanaTokenAuthFlagHelp)

	// Add the Grafana.connection specific flags
	Instance.String(configtypes.GrafanaConnectionSchemeOptimized, configtypes.GrafanaConnectionSchemeDefault, configtypes.GrafanaConnectionSchemeFlagHelp)
	Instance.StringP(configtypes.GrafanaConnectionHostOptimized, configtypes.GrafanaConnectionHostFlagShort, configtypes.GrafanaConnectionHostDefault, configtypes.GrafanaConnectionHostFlagHelp)
	Instance.String(configtypes.GrafanaConnectionBasePathOptimized, configtypes.GrafanaConnectionBasePathDefault, configtypes.GrafanaConnectionBasePathFlagHelp)
	Instance.IntP(configtypes.GrafanaConnectionRetryOptimized, configtypes.GrafanaConnectionRetryFlagShort, configtypes.GrafanaConnectionRetryDefault, configtypes.GrafanaConnectionRetryFlagHelp)

	// Add the feature specific flags
	Instance.Bool(configtypes.FeaturesAddLocalAdminToTeamsOptimized, configtypes.FeaturesAddLocalAdminToTeamsDefault, configtypes.FeaturesAddLocalAdminToTeamsFlagHelp)
	Instance.Bool(configtypes.FeaturesDisableFoldersOptimized, configtypes.FeaturesDisableFoldersDefault, configtypes.FeaturesDisableFoldersFlagHelp)
	Instance.Bool(configtypes.FeaturesDisableUsersOptimized, configtypes.FeaturesDisableUsersDefault, configtypes.FeaturesDisableUsersFlagHelp)

	// Add "source" flag
	Instance.StringP(configtypes.SourceParameter, configtypes.SourceFlagShort, configtypes.SourceDefault, configtypes.SourceFlagHelp)

	// Add "teams" flag
	Instance.String(configtypes.TeamsParameter, configtypes.TeamsDefault, configtypes.TeamsFlagHelp)

	// Add "loglevel" flag
	Instance.IntP(configtypes.LogLevelParameter, configtypes.LogLevelFlagShort, configtypes.LogLevelDefault, configtypes.LogLevelFlagHelp)

	// Add "version" flag
	Instance.BoolVarP(&Version, configtypes.VersionParameter, configtypes.VersionFlagShort, configtypes.VersionDefault, configtypes.VersionFlagHelp)

	// Add "help" flag
	Instance.BoolVarP(&Help, configtypes.HelpParameter, configtypes.HelpFlagShort, configtypes.HelpDefault, configtypes.HelpFlagHelp)

	// Parse cli input
	if err := Instance.Parse(os.Args[1:]); err != nil {
		slog.Warn("could not parse CLI arguments")
	}

	// Handle help printing
	if Help {
		fmt.Printf("Usage of %s:\r\n\r\n", Instance.Name())
		fmt.Println(Instance.FlagUsages())
		os.Exit(0)
	}
}
