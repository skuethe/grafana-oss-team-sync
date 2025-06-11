package flags

import (
	"github.com/skuethe/grafana-oss-team-sync/internal/config/types"
	flag "github.com/spf13/pflag"
)

const (
	flagBasicAuthUsernameFull    string = "username"
	flagBasicAuthUsernameShort   string = "u"
	flagBasicAuthUsernameDefault string = ""
	flagBasicAuthUsernameHelp    string = "the basic auth user you want to use to authenticate against Grafana"

	flagBasicAuthPasswordFull    string = "password"
	flagBasicAuthPasswordShort   string = "p"
	flagBasicAuthPasswordDefault string = ""
	flagBasicAuthPasswordHelp    string = "the basic auth password you want to use to authenticate against Grafana"
)

var (
	Instance *flag.FlagSet

	Config string

	basicAuthUsername string
	basicAuthPassword string

	Version bool
)

func Load() {
	// TODO: add "authfile", "token", "teams", "folders"
	Instance = flag.NewFlagSet("grafana-oss-team-sync", flag.ExitOnError)
	Instance.SortFlags = false

	// Add "config" flag
	Instance.StringVarP(&Config, types.ConfigParameter, types.ConfigFlagShort, types.ConfigDefault, types.ConfigFlagHelp)

	// Add "authfile" flag
	Instance.StringP(types.AuthFileParameter, types.AuthFileFlagShort, types.AuthFileDefault, types.AuthFileFlagHelp)

	// Add basic auth flags
	Instance.StringVarP(&basicAuthUsername, flagBasicAuthUsernameFull, flagBasicAuthUsernameShort, flagBasicAuthUsernameDefault, flagBasicAuthUsernameHelp)
	Instance.StringVarP(&basicAuthPassword, flagBasicAuthPasswordFull, flagBasicAuthPasswordShort, flagBasicAuthPasswordDefault, flagBasicAuthPasswordHelp)

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
}
