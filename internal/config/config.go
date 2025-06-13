package config

import (
	"errors"
	"os"
	"strings"

	"log/slog"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/skuethe/grafana-oss-team-sync/internal/config/configtypes"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
	"github.com/spf13/pflag"
)

var Instance *configtypes.Config

func getConfigFilePath() (*string, error) {
	var config string

	// Load possible env vars first
	config = os.Getenv(configtypes.ConfigVariable)

	// Merge possible input from flags
	if flags.Config != "" {
		config = flags.Config
	}

	if config == "" {
		return nil, errors.New("no config file defined")
	}

	return &config, nil
}

func loadYAMLFile(k *koanf.Koanf) error {
	// Handle config file definitions from flag and environment variables
	configFile, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Load main YAML config
	if err := k.Load(file.Provider(*configFile), yaml.Parser()); err != nil {
		return err
	}
	return nil
}

func loadEnvironmentVariables(k *koanf.Koanf) {
	k.Load(env.ProviderWithValue("GOTS_", ".", func(k string, v string) (string, any) {
		key := strings.Replace(strings.ToLower(strings.TrimPrefix(k, "GOTS_")), "_", ".", -1)

		switch key {
		// "teams" input -> allow comma seperated list
		case configtypes.TeamsParameter:
			return key, strings.Split(v, ",")
		// Features - addLocalAdminToTeams: return full parameter, to map it to the config object
		case configtypes.FeaturesAddLocalAdminToTeamsOptimized:
			return configtypes.FeaturesAddLocalAdminToTeamsParameter, v
		// Features - disableFolders: return full parameter, to map it to the config object
		case configtypes.FeaturesDisableFoldersOptimized:
			return configtypes.FeaturesDisableFoldersParameter, v
		// Features - disableUserSync: return full parameter, to map it to the config object
		case configtypes.FeaturesDisableUsersOptimized:
			return configtypes.FeaturesDisableUsersParameter, v
		}

		return key, v
	}), nil)
}

func loadCLIParameter(k *koanf.Koanf) {
	k.Load(posflag.ProviderWithFlag(flags.Instance, ".", k, func(f *pflag.Flag) (string, any) {
		key := f.Name
		val := posflag.FlagVal(flags.Instance, f)

		slog.Warn("DEBUG", "key", key, "val", val)

		switch key {
		// Features - addLocalAdminToTeams: return full parameter, to map it to the config object
		case configtypes.FeaturesAddLocalAdminToTeamsOptimized:
			return configtypes.FeaturesAddLocalAdminToTeamsParameter, val
		// Features - disableFolders: return full parameter, to map it to the config object
		case configtypes.FeaturesDisableFoldersOptimized:
			return configtypes.FeaturesDisableFoldersParameter, val
		// Features - disableUserSync: return full parameter, to map it to the config object
		case configtypes.FeaturesDisableUsersOptimized:
			return configtypes.FeaturesDisableUsersParameter, val
		// Special handling for "teams" input -> allow comma seperated list
		case configtypes.TeamsParameter:
			return key, strings.Split(val.(string), ",")
		}

		return key, val
	}), nil)
}

func loadOptionalAuthFile(k *koanf.Koanf) {
	authfile := k.String(configtypes.AuthFileParameter)
	if authfile != "" {
		godotenv.Load(authfile)
		loadEnvironmentVariables(k)
	}
}

func Load() {
	configLog := slog.With(slog.String("package", "config"))
	configLog.Info("loading config")

	// Handle .env files
	godotenv.Load()

	// Global koanf instance for input handling
	var k = koanf.New(".")

	// Load main YAML config
	if err := loadYAMLFile(k); err != nil {
		configLog.Error("could not load config file")
		panic(err)

	}

	// Load env vars and merge (override) config
	loadEnvironmentVariables(k)

	// Load flags and merge (override) config
	// This will also load default values specified in the flags package
	loadCLIParameter(k)

	// Load optional authfile as config and environment variables
	loadOptionalAuthFile(k)

	// Create new Config instance from the different inputs
	k.UnmarshalWithConf("", &Instance, koanf.UnmarshalConf{
		Tag: "yaml",
	})

	slog.Warn("DEBUG", "features", Instance.Features)

	// Validate Grafana authtype input
	if err := Instance.ValdidateGrafanaAuthType(); err != nil {
		configLog.Error("error validating Grafana authtype input")
		panic(err)
	}

	// Validate Grafana connection scheme input
	if err := Instance.ValdidateGrafanaScheme(); err != nil {
		configLog.Error("error validating Grafana connection scheme input")
		panic(err)
	}

	// Validate source input
	if err := Instance.ValdidateSourcePlugin(); err != nil {
		configLog.Error("error validating source input")
		panic(err)
	}

	// Give warning about empty team input
	if len(Instance.Teams) == 0 {
		configLog.Warn("your teams input is empty")
	}

}
