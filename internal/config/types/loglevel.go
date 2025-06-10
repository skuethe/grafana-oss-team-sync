package types

import "log/slog"

type LogLevel int

func (c *Config) GetLogLevel() slog.Level {
	switch c.LogLevel {
	case 0:
		return slog.LevelInfo
	case 1:
		return slog.LevelWarn
	case 2:
		return slog.LevelError
	case 99:
		return slog.LevelDebug
	default:
		slog.Warn("undefined log level set, falling back to INFO", "loglevel", c.LogLevel)
		return slog.LevelInfo
	}
}
