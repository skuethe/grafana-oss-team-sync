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
)

var (
	Instance *flag.FlagSet
	Config   string
	Loglevel int
	Version  bool
)

func Load() {
	Instance = flag.NewFlagSet("grafana-oss-team-sync", flag.ExitOnError)

	// Add "--config", "-c" flag
	Instance.StringVarP(&Config, FlagConfigFull, FlagConfigShort, FlagConfigDefault, FlagConfigHelp)

	// Add "--loglevel", "-l" flag
	Instance.IntVarP(&Loglevel, FlagLoglevelFull, FlagLoglevelShort, FlagLoglevelDefault, FlagLoglevelHelp)

	// Add "--version", "-v" flag
	Instance.BoolVarP(&Version, FlagVersionFull, FlagVersionShort, FlagVersionDefault, FlagVersionHelp)
}
