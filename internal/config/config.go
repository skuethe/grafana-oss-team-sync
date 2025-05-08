package config

import (
	"os"
	"strings"

	"log/slog"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// type Config struct {
// 	k *koanf.Koanf
// }

// Global koanf instance for input handling
var K = koanf.New(".")

func GetLogLevel() slog.Level {
	var level slog.Level
	configLevel := K.Int("loglevel")
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
		slog.Warn("Undefined log level set, falling back to INFO", "loglevel", configLevel)
		level = slog.LevelInfo
	}
	return level
}

func Start() {
	configLog := slog.With(slog.String("package", "config"))
	configLog.Info("Initializing Config")

	// Load YAML config
	if err := K.Load(file.Provider("configs/config.yaml"), yaml.Parser()); err != nil {
		configLog.Error("Could not load Config", "error", err)
		os.Exit(1)
	}

	// Load env vars and merge (override) config
	K.Load(env.Provider("GOTS_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "GOTS_")), "_", ".", -1)
	}), nil)
}
