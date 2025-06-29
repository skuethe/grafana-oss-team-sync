// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import "log/slog"

type LogLevel int

const (
	LogLevelDefault   int    = 0
	LogLevelFlagHelp  string = "the log-level.\nAllowed: 0 (INFO), 1 (WARN), 2 (ERROR), 99 (DEBUG)\nDefault: 0 (INFO)"
	LogLevelFlagShort string = "l"
	LogLevelParameter string = "loglevel"
	LogLevelVariable  string = "GOTS_LOGLEVEL"
)

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
