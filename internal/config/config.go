package config

import (
	"os"
	"strings"

	"log/slog"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/skuethe/grafana-oss-team-sync/internal/config/types"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
)

var Instance *types.Config

func Load() {
	configLog := slog.With(slog.String("package", "config"))
	configLog.Info("loading config")

	// Handle .env files
	godotenv.Load()

	// Handle config file definitions from flag and environment variables
	// TODO: implement solution to load config from OS ENV / --config flag
	// mainConfig := os.Getenv("GOTS_CONFIG")

	// Global koanf instance for input handling
	var k = koanf.New(".")

	// Load main YAML config
	if err := k.Load(file.Provider(flags.FlagInputConfig), yaml.Parser()); err != nil {
		configLog.Error("could not load config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	// Load optional authfile as config and environment variables
	if k.String(types.ConfigParamAuthFile) != "" {
		godotenv.Load(k.String(types.ConfigParamAuthFile))
	}

	// Load env vars and merge (override) config
	k.Load(env.Provider("GOTS_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "GOTS_")), "_", ".", -1)
	}), nil)

	// Load flags and merge (override) config
	// This will also load default values specified in the flags package
	k.Load(posflag.Provider(flags.FlagSet, ".", k), nil)

	// Create new Config instance from the different inputs
	k.UnmarshalWithConf("", &Instance, koanf.UnmarshalConf{
		Tag: "yaml",
	})

	// Validate source input
	if err := Instance.ValdidateSourcePlugin(); err != nil {
		configLog.Error("error validating source config")
		panic(err)
	}

	// Give warning about empty team input
	if len(Instance.Teams) == 0 {
		configLog.Warn("your teams input is empty")
	}

}
