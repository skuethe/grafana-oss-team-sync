package flags

import (
	"fmt"
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

	// Add basic auth flags
	Instance.StringVarP(&BasicAuthUsername, configtypes.GrafanaBasicAuthUsernameParameter, configtypes.GrafanaBasicAuthUsernameFlagShort, configtypes.GrafanaBasicAuthUsernameDefault, configtypes.GrafanaBasicAuthUsernameFlagHelp)
	Instance.StringVarP(&BasicAuthPassword, configtypes.GrafanaBasicAuthPasswordParameter, configtypes.GrafanaBasicAuthPasswordFlagShort, configtypes.GrafanaBasicAuthPasswordDefault, configtypes.GrafanaBasicAuthPasswordFlagHelp)

	// Add token auth flags
	Instance.StringVarP(&Token, configtypes.GrafanaTokenAuthParameter, configtypes.GrafanaTokenAuthFlagShort, configtypes.GrafanaTokenAuthDefault, configtypes.GrafanaTokenAuthFlagHelp)

	// Add "authtype" flag
	Instance.String(configtypes.GrafanaAuthTypeParameter, configtypes.GrafanaAuthTypeDefault, configtypes.GrafanaAuthTypeFlagHelp)

	// Add the Grafana.connection specific flags
	Instance.String(configtypes.GrafanaConnectionSchemeParameter, configtypes.GrafanaConnectionSchemeDefault, configtypes.GrafanaConnectionSchemeFlagHelp)
	Instance.StringP(configtypes.GrafanaConnectionHostParameter, configtypes.GrafanaConnectionHostFlagShort, configtypes.GrafanaConnectionHostDefault, configtypes.GrafanaConnectionHostFlagHelp)
	Instance.String(configtypes.GrafanaConnectionBasePathParameter, configtypes.GrafanaConnectionBasePathDefault, configtypes.GrafanaConnectionBasePathFlagHelp)
	Instance.IntP(configtypes.GrafanaConnectionRetryParameter, configtypes.GrafanaConnectionRetryFlagShort, configtypes.GrafanaConnectionRetryDefault, configtypes.GrafanaConnectionRetryFlagHelp)

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
	Instance.Parse(os.Args[1:])

	// Handle help printing
	if Help {
		fmt.Printf("Usage of %s:\r\n\r\n", Instance.Name())
		fmt.Println(Instance.FlagUsages())
		os.Exit(0)
	}
}
