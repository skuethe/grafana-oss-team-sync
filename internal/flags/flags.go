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
	FlagGrafanaConnectionHostShort   string = "h"
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
	Instance *flag.FlagSet

	Config string
	Source string

	FeatureDisableFolders       bool
	FeatureDisableUserss        bool
	FeatureAddLocalAdminToTeams bool

	GrafanaConnectionScheme   string
	GrafanaConnectionHost     string
	GrafanaConnectionBasePath string
	GrafanaConnectionRetry    int

	Loglevel int
	Version  bool
)

func Load() {
	Instance = flag.NewFlagSet("grafana-oss-team-sync", flag.ExitOnError)
	Instance.SortFlags = false

	// Add "--config", "-c" flag
	Instance.StringVarP(&Config, FlagConfigFull, FlagConfigShort, FlagConfigDefault, FlagConfigHelp)

	// Add "--source", "-s" flag
	Instance.StringVarP(&Source, FlagSourceFull, FlagSourceShort, FlagSourceDefault, FlagSourceHelp)

	// Add the feature specific flags
	Instance.BoolVar(&FeatureDisableFolders, FlagDisableFoldersFull, FlagDisableFoldersDefault, FlagDisableFoldersHelp)
	Instance.BoolVar(&FeatureDisableUserss, FlagDisableUsersFull, FlagDisableUsersDefault, FlagDisableUsersHelp)
	Instance.BoolVar(&FeatureAddLocalAdminToTeams, FlagAddLocalAdminToTeamsFull, FlagAddLocalAdminToTeamsDefault, FlagAddLocalAdminToTeamsHelp)

	// Add the Grafana.connection specific flags
	Instance.StringVar(&GrafanaConnectionScheme, FlagGrafanaConnectionSchemeFull, FlagGrafanaConnectionSchemeDefault, FlagGrafanaConnectionSchemeHelp)
	Instance.StringVarP(&GrafanaConnectionHost, FlagGrafanaConnectionHostFull, FlagGrafanaConnectionHostShort, FlagGrafanaConnectionHostDefault, FlagGrafanaConnectionHostHelp)
	Instance.StringVar(&GrafanaConnectionBasePath, FlagGrafanaConnectionBasePathFull, FlagGrafanaConnectionBasePathDefault, FlagGrafanaConnectionBasePathHelp)
	Instance.IntVarP(&GrafanaConnectionRetry, FlagGrafanaConnectionRetryFull, FlagGrafanaConnectionRetryShort, FlagGrafanaConnectionRetryDefault, FlagGrafanaConnectionRetryHelp)

	// Add "--loglevel", "-l" flag
	Instance.IntVarP(&Loglevel, FlagLoglevelFull, FlagLoglevelShort, FlagLoglevelDefault, FlagLoglevelHelp)

	// Add "--version", "-v" flag
	Instance.BoolVarP(&Version, FlagVersionFull, FlagVersionShort, FlagVersionDefault, FlagVersionHelp)
}
