// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package config

import (
	"errors"
	"fmt"
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

var (
	Instance *configtypes.Config

	ErrNoConfigFileDefined              = errors.New("no config file defined")
	ErrCouldNotLoadCLIArgs              = errors.New("could not load CLI arguments")
	ErrCouldNotLoadYAMLFile             = errors.New("could not load YAML file")
	ErrCouldNotLoadAuthfile             = errors.New("could not process configured authfile")
	ErrCouldNotLoadEnvironmentVariables = errors.New("could not load environment variables")
	ErrCouldNotUnmarshalConfig          = errors.New("could not unmarshal config")
	ErrCouldNotSetRequiredVariable      = errors.New("could not set required environment variable")
)

func getConfigFilePath() (*string, error) {
	var config string

	// Load possible env vars first
	config = os.Getenv(configtypes.ConfigVariable)

	// Merge possible input from flags
	if flags.Config != "" {
		config = flags.Config
	}

	if config == "" {
		return nil, ErrNoConfigFileDefined
	}

	return &config, nil
}

func loadYAMLFile(k *koanf.Koanf) error {
	// Handle config file definitions from flag and environment variables
	configFile, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotLoadYAMLFile, err)
	}

	// Load main YAML config
	if err := k.Load(file.Provider(*configFile), yaml.Parser()); err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotLoadYAMLFile, err)
	}
	return nil
}

func loadEnvironmentVariables(k *koanf.Koanf) error {
	err := k.Load(env.ProviderWithValue("GOTS_", ".", func(k string, v string) (string, any) {
		key := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(k, "GOTS_")), "_", ".")

		switch key {
		// Grafana flags - return full parameter, to map it to the config object
		case configtypes.GrafanaAuthTypeOptimized:
			return configtypes.GrafanaAuthTypeParameter, v

		case configtypes.GrafanaConnectionSchemeOptimized:
			return configtypes.GrafanaConnectionSchemeParameter, v

		case configtypes.GrafanaConnectionHostOptimized:
			return configtypes.GrafanaConnectionHostParameter, v

		case configtypes.GrafanaConnectionBasePathOptimized:
			return configtypes.GrafanaConnectionBasePathParameter, v

		case configtypes.GrafanaConnectionRetryOptimized:
			return configtypes.GrafanaConnectionRetryParameter, v

		// Features - return full parameter, to map it to the config object
		case configtypes.FeaturesAddLocalAdminToTeamsOptimized:
			return configtypes.FeaturesAddLocalAdminToTeamsParameter, v

		case configtypes.FeaturesDisableFoldersOptimized:
			return configtypes.FeaturesDisableFoldersParameter, v

		case configtypes.FeaturesDisableUsersOptimized:
			return configtypes.FeaturesDisableUsersParameter, v

		// "teams" - respect comma seperated list
		case configtypes.TeamsParameter:
			return key, strings.Split(v, ",")
		}

		return key, v
	}), nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotLoadEnvironmentVariables, err)
	}
	return nil
}

func loadCLIParameter(k *koanf.Koanf, fs *pflag.FlagSet) error {
	err := k.Load(posflag.ProviderWithFlag(fs, ".", k, func(f *pflag.Flag) (string, any) {
		key := f.Name
		val := posflag.FlagVal(fs, f)

		switch key {
		// Grafana flags - return full parameter, to map it to the config object
		case configtypes.GrafanaAuthTypeOptimized:
			return configtypes.GrafanaAuthTypeParameter, val

		case configtypes.GrafanaConnectionSchemeOptimized:
			return configtypes.GrafanaConnectionSchemeParameter, val

		case configtypes.GrafanaConnectionHostOptimized:
			return configtypes.GrafanaConnectionHostParameter, val

		case configtypes.GrafanaConnectionBasePathOptimized:
			return configtypes.GrafanaConnectionBasePathParameter, val

		case configtypes.GrafanaConnectionRetryOptimized:
			return configtypes.GrafanaConnectionRetryParameter, val

		// Feature flags - return full parameter, to map it to the config object
		case configtypes.FeaturesAddLocalAdminToTeamsOptimized:
			return configtypes.FeaturesAddLocalAdminToTeamsParameter, val

		case configtypes.FeaturesDisableFoldersOptimized:
			return configtypes.FeaturesDisableFoldersParameter, val

		case configtypes.FeaturesDisableUsersOptimized:
			return configtypes.FeaturesDisableUsersParameter, val

		// "teams" - respect comma seperated list
		case configtypes.TeamsParameter:
			return key, strings.Split(val.(string), ",")
		}

		return key, val
	}), nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotLoadCLIArgs, err)
	}
	return nil
}

func loadOptionalAuthFile(k *koanf.Koanf) error {
	authfile := k.String(configtypes.AuthFileParameter)
	if authfile != "" {
		if err := godotenv.Load(authfile); err != nil {
			return fmt.Errorf("%w (%v): %w", ErrCouldNotLoadAuthfile, authfile, err)
		}
		if err := loadEnvironmentVariables(k); err != nil {
			return err
		}
	}
	return nil
}

func unmarshalIntoStruct(k *koanf.Koanf) error {
	if err := k.UnmarshalWithConf("", &Instance, koanf.UnmarshalConf{
		Tag: "yaml",
	}); err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotUnmarshalConfig, err)
	}
	return nil
}

func Load() error {
	configLog := slog.With(slog.String("package", "config"))
	configLog.Info("loading config")

	// Handle .env files
	_ = godotenv.Load()

	// Global koanf instance for input handling
	var k = koanf.New(".")

	// Load main YAML config
	if err := loadYAMLFile(k); err != nil {
		return err
	}

	// Load env vars and merge (override) config
	if err := loadEnvironmentVariables(k); err != nil {
		configLog.Warn(err.Error())
	}

	// Load flags and merge (override) config
	// This will also load default values specified in the flags package
	if err := loadCLIParameter(k, flags.Instance); err != nil {
		configLog.Warn(err.Error())
	}

	// Load optional authfile as config and environment variables
	if err := loadOptionalAuthFile(k); err != nil {
		configLog.Warn(err.Error())
	}

	// Create new Config instance from the different inputs
	if err := unmarshalIntoStruct(k); err != nil {
		return err
	}

	// Output feature configuration if non-defaults set
	configLog.Info("feature status",
		slog.Bool(configtypes.FeaturesDisableFoldersOptimized, Instance.Features.DisableFolders),
		slog.Bool(configtypes.FeaturesDisableUsersOptimized, Instance.Features.DisableUserSync),
		slog.Bool(configtypes.FeaturesAddLocalAdminToTeamsOptimized, Instance.Features.AddLocalAdminToTeams),
	)

	// Validate Grafana authtype input
	if err := Instance.ValdidateGrafanaAuthType(); err != nil {
		return err
	}

	// Validate Grafana connection scheme input
	if err := Instance.ValdidateGrafanaScheme(); err != nil {
		return err
	}

	// Validate source input
	if err := Instance.ValdidateSourcePlugin(); err != nil {
		return err
	}

	// Give warning about empty team input
	if len(Instance.Teams) == 0 {
		configLog.Warn("your teams input is empty")
	}

	return nil
}
