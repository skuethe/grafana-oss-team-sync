package flags

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagVersionFull    string = "version"
	FlagVersionShort   string = "v"
	FlagVersionDefault bool   = false
	FlagVersionHelp    string = "print the version, commit hash and build date. Then exit"

	FlagConfigFull    string = "config"
	FlagConfigShort   string = "c"
	FlagConfigDefault string = "configs/config.yaml"
	FlagConfigHelp    string = "the path to your config.yaml file"

	FlagLoglevelFull    string = "loglevel"
	FlagLoglevelShort   string = "l"
	FlagLoglevelDefault int    = 0
	FlagLoglevelHelp    string = "configure the log level.\nAllowed: 0 (INFO), 1 (WARN), 2 (ERROR), 99 (DEBUG)\nDefault: 0 (INFO)"

	FlagSourceFull    string = "source"
	FlagSourceShort   string = "s"
	FlagSourceDefault string = "entraid"
	FlagSourceHelp    string = "specify the source `plugin` you want to use\nAllowed: entraid"

	FlagBasicAuthUsernameFull    string = "username"
	FlagBasicAuthUsernameShort   string = "u"
	FlagBasicAuthUsernameDefault string = ""
	FlagBasicAuthUsernameHelp    string = "specify the basic auth `user` you want to use to authenticate against Grafana"

	FlagBasicAuthPasswordFull    string = "password"
	FlagBasicAuthPasswordShort   string = "p"
	FlagBasicAuthPasswordDefault string = ""
	FlagBasicAuthPasswordHelp    string = "specify the basic auth `password` you want to use to authenticate against Grafana"

	FlagDisableFoldersFull    string = "features.disableFolders"
	FlagDisableFoldersDefault bool   = false
	FlagDisableFoldersHelp    string = "disable the folder sync feature"

	FlagDisableUsersFull    string = "features.disableUsers"
	FlagDisableUsersDefault bool   = false
	FlagDisableUsersHelp    string = "disable the user sync feature"

	FlagAddLocalAdminToTeamsFull    string = "features.addLocalAdminToTeams"
	FlagAddLocalAdminToTeamsDefault bool   = true
	FlagAddLocalAdminToTeamsHelp    string = "add the local Grafana admin user to each team you create"

	FlagGrafanaConnectionSchemeFull    string = "grafana.connection.scheme"
	FlagGrafanaConnectionSchemeDefault string = "http"
	FlagGrafanaConnectionSchemeHelp    string = "the scheme of your Grafana instance\nAllowed: http or https"

	FlagGrafanaConnectionHostFull    string = "grafana.connection.host"
	FlagGrafanaConnectionHostShort   string = "H"
	FlagGrafanaConnectionHostDefault string = "localhost:3000"
	FlagGrafanaConnectionHostHelp    string = "the host of your Grafana instance"

	FlagGrafanaConnectionBasePathFull    string = "grafana.connection.basePath"
	FlagGrafanaConnectionBasePathDefault string = "/api"
	FlagGrafanaConnectionBasePathHelp    string = "the base path of your Grafana instance"

	FlagGrafanaConnectionRetryFull    string = "grafana.connection.retry"
	FlagGrafanaConnectionRetryShort   string = "r"
	FlagGrafanaConnectionRetryDefault int    = 0
	FlagGrafanaConnectionRetryHelp    string = "configure the amount of retries to connect to your Grafana instance\nRetries are paused by 2 seconds\nDefault: 0"
)

var (
	FlagSet *flag.FlagSet

	FlagInputConfig string
	FlagInputSource string

	FlagInputBasicAuthUsername string
	FlagInputBasicAuthPassword string

	FlagInputFeatureDisableFolders       bool
	FlagInputFeatureDisableUserss        bool
	FlagInputFeatureAddLocalAdminToTeams bool

	FlagInputGrafanaConnectionScheme   string
	FlagInputGrafanaConnectionHost     string
	FlagInputGrafanaConnectionBasePath string
	FlagInputGrafanaConnectionRetry    int

	FlagInputLoglevel int
	FlagInputVersion  bool
)

func LoadFlags() {
	FlagSet = flag.NewFlagSet("grafana-oss-team-sync", flag.ExitOnError)
	FlagSet.SortFlags = false

	// Add "--config", "-c" flag
	FlagSet.StringVarP(&FlagInputConfig, FlagConfigFull, FlagConfigShort, FlagConfigDefault, FlagConfigHelp)

	// Add "--source", "-s" flag
	FlagSet.StringVarP(&FlagInputSource, FlagSourceFull, FlagSourceShort, FlagSourceDefault, FlagSourceHelp)

	// Add basic auth flags
	FlagSet.StringVarP(&FlagInputBasicAuthUsername, FlagBasicAuthUsernameFull, FlagBasicAuthUsernameShort, FlagBasicAuthUsernameDefault, FlagBasicAuthUsernameHelp)
	FlagSet.StringVarP(&FlagInputBasicAuthPassword, FlagBasicAuthPasswordFull, FlagBasicAuthPasswordShort, FlagBasicAuthPasswordDefault, FlagBasicAuthPasswordHelp)

	// Add the feature specific flags
	FlagSet.BoolVar(&FlagInputFeatureDisableFolders, FlagDisableFoldersFull, FlagDisableFoldersDefault, FlagDisableFoldersHelp)
	FlagSet.BoolVar(&FlagInputFeatureDisableUserss, FlagDisableUsersFull, FlagDisableUsersDefault, FlagDisableUsersHelp)
	FlagSet.BoolVar(&FlagInputFeatureAddLocalAdminToTeams, FlagAddLocalAdminToTeamsFull, FlagAddLocalAdminToTeamsDefault, FlagAddLocalAdminToTeamsHelp)

	// Add the Grafana.connection specific flags
	FlagSet.StringVar(&FlagInputGrafanaConnectionScheme, FlagGrafanaConnectionSchemeFull, FlagGrafanaConnectionSchemeDefault, FlagGrafanaConnectionSchemeHelp)
	FlagSet.StringVarP(&FlagInputGrafanaConnectionHost, FlagGrafanaConnectionHostFull, FlagGrafanaConnectionHostShort, FlagGrafanaConnectionHostDefault, FlagGrafanaConnectionHostHelp)
	FlagSet.StringVar(&FlagInputGrafanaConnectionBasePath, FlagGrafanaConnectionBasePathFull, FlagGrafanaConnectionBasePathDefault, FlagGrafanaConnectionBasePathHelp)
	FlagSet.IntVarP(&FlagInputGrafanaConnectionRetry, FlagGrafanaConnectionRetryFull, FlagGrafanaConnectionRetryShort, FlagGrafanaConnectionRetryDefault, FlagGrafanaConnectionRetryHelp)

	// Add "--loglevel", "-l" flag
	FlagSet.IntVarP(&FlagInputLoglevel, FlagLoglevelFull, FlagLoglevelShort, FlagLoglevelDefault, FlagLoglevelHelp)

	// Add "--version", "-v" flag
	FlagSet.BoolVarP(&FlagInputVersion, FlagVersionFull, FlagVersionShort, FlagVersionDefault, FlagVersionHelp)
}
