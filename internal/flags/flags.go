package flags

import (
	"fmt"
	"os"

	"github.com/skuethe/grafana-oss-team-sync/internal/config/types"
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
	// TODO: add "teams", "folders"
	Instance = flag.NewFlagSet("grafana-oss-team-sync", flag.ExitOnError)
	Instance.SortFlags = false

	// Add "config" flag
	Instance.StringVarP(&Config, types.ConfigParameter, types.ConfigFlagShort, types.ConfigDefault, types.ConfigFlagHelp)

	// Add "authfile" flag
	Instance.StringP(types.AuthFileParameter, types.AuthFileFlagShort, types.AuthFileDefault, types.AuthFileFlagHelp)

	// Add basic auth flags
	Instance.StringVarP(&BasicAuthUsername, types.GrafanaBasicAuthUsernameParameter, types.GrafanaBasicAuthUsernameFlagShort, types.GrafanaBasicAuthUsernameDefault, types.GrafanaBasicAuthUsernameFlagHelp)
	Instance.StringVarP(&BasicAuthPassword, types.GrafanaBasicAuthPasswordParameter, types.GrafanaBasicAuthPasswordFlagShort, types.GrafanaBasicAuthPasswordDefault, types.GrafanaBasicAuthPasswordFlagHelp)

	// Add token auth flags
	Instance.StringVarP(&Token, types.GrafanaTokenAuthParameter, types.GrafanaTokenAuthFlagShort, types.GrafanaTokenAuthDefault, types.GrafanaTokenAuthFlagHelp)

	// Add "authtype" flag
	Instance.String(types.GrafanaAuthTypeParameter, types.GrafanaAuthTypeDefault, types.GrafanaAuthTypeFlagHelp)

	// Add the Grafana.connection specific flags
	Instance.String(types.GrafanaConnectionSchemeParameter, types.GrafanaConnectionSchemeDefault, types.GrafanaConnectionSchemeFlagHelp)
	Instance.StringP(types.GrafanaConnectionHostParameter, types.GrafanaConnectionHostFlagShort, types.GrafanaConnectionHostDefault, types.GrafanaConnectionHostFlagHelp)
	Instance.String(types.GrafanaConnectionBasePathParameter, types.GrafanaConnectionBasePathDefault, types.GrafanaConnectionBasePathFlagHelp)
	Instance.IntP(types.GrafanaConnectionRetryParameter, types.GrafanaConnectionRetryFlagShort, types.GrafanaConnectionRetryDefault, types.GrafanaConnectionRetryFlagHelp)

	// Add the feature specific flags
	Instance.Bool(types.FeaturesDisableFoldersParameter, types.FeaturesDisableFoldersDefault, types.FeaturesDisableFoldersFlagHelp)
	Instance.Bool(types.FeaturesDisableUsersParameter, types.FeaturesDisableUsersDefault, types.FeaturesDisableUsersFlagHelp)
	Instance.Bool(types.FeaturesAddLocalAdminToTeamsParameter, types.FeaturesAddLocalAdminToTeamsDefault, types.FeaturesAddLocalAdminToTeamsFlagHelp)

	// Add "source" flag
	Instance.StringP(types.SourceParameter, types.SourceFlagShort, types.SourceDefault, types.SourceFlagHelp)

	// Add "loglevel" flag
	Instance.IntP(types.LogLevelParameter, types.LogLevelFlagShort, types.LogLevelDefault, types.LogLevelFlagHelp)

	// Add "version" flag
	Instance.BoolVarP(&Version, types.VersionParameter, types.VersionFlagShort, types.VersionDefault, types.VersionFlagHelp)

	// Add "help" flag
	Instance.BoolVarP(&Help, types.HelpParameter, types.HelpFlagShort, types.HelpDefault, types.HelpFlagHelp)

	// Parse cli input
	Instance.Parse(os.Args[1:])

	// Handle help printing
	if Help {
		fmt.Printf("Usage of %s:\r\n\r\n", Instance.Name())
		fmt.Println(Instance.FlagUsages())
		os.Exit(0)
	}
}
