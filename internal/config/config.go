package config

import (
	"os"
	"strings"

	"log/slog"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
)

// Global koanf instance for input handling
var K = koanf.New(".")

const (
	ConfigParamAuthFile string = "authFile"
	ConfigParamFeatures string = "features"
	ConfigParamFolders  string = "folders"
	ConfigParamGrafana  string = "grafana"
	ConfigParamLogLevel string = "loglevel"
	ConfigParamSource   string = "source"

	ConfigParamAuthBasicUsername string = "username"
	ConfigParamAuthBasicPassword string = "password"
	ConfigParamAuthToken         string = "token"
)

func GetLogLevel() slog.Level {
	var level slog.Level
	configLevel := K.Int(ConfigParamLogLevel)
	switch configLevel {
	case 0:
		level = slog.LevelInfo
	case 1:
		level = slog.LevelWarn
	case 2:
		level = slog.LevelError
	case 99:
		level = slog.LevelDebug
	default:
		slog.Warn("undefined log level set, falling back to INFO", "loglevel", configLevel)
		level = slog.LevelInfo
	}
	return level
}

func isAuthFileSet() bool {
	authFile := K.String(ConfigParamAuthFile)
	return authFile != ""
}

func Load() {
	configLog := slog.With(slog.String("package", "config"))
	configLog.Info("loading config")

	// Handle .env files
	godotenv.Load()

	// Load YAML config
	if err := K.Load(file.Provider(flags.Config), yaml.Parser()); err != nil {
		configLog.Error("could not load config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	// Load optional authFile YAML
	if isAuthFileSet() {
		if err := K.Load(file.Provider(K.String(ConfigParamAuthFile)), dotenv.Parser()); err != nil {
			configLog.Error("could not load "+ConfigParamAuthFile,
				slog.Any("error", err),
			)
			os.Exit(1)
		}
		godotenv.Load(K.String(ConfigParamAuthFile))
	}

	// Load env vars and merge (override) config
	K.Load(env.Provider("GOTS_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "GOTS_")), "_", ".", -1)
	}), nil)

	// Load flags and merge (override) config
	K.Load(posflag.Provider(flags.Instance, ".", K), nil)

	// Validate source input
	sourceErr := valdidateSource()
	if sourceErr != nil {
		configLog.Error("error parsing source config",
			slog.Any("error", sourceErr),
		)
		os.Exit(1)
	}

	// Validate feature flags
	featureErr := valdidateFeatureSchema()
	if featureErr != nil {
		configLog.Error("error parsing feature config",
			slog.Any("error", featureErr),
		)
		os.Exit(1)
	}

	// Validate teams input
	teamsErr := validateTeamsSchema()
	if teamsErr != nil {
		configLog.Error("error parsing teams config",
			slog.Any("error", teamsErr),
		)
		os.Exit(1)
	}
	if len(Teams.List) == 0 {
		configLog.Warn("your teams input is empty")
	}

}
