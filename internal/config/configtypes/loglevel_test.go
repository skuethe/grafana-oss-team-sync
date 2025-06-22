// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import (
	"log/slog"
	"testing"
)

func TestGetLogLevel(t *testing.T) {

	type addTest struct {
		name     string
		input    Config
		expected slog.Level
	}

	var tests = []addTest{
		{"loglevel set to info", Config{LogLevel: 0}, slog.LevelInfo},
		{"loglevel set to warning", Config{LogLevel: 1}, slog.LevelWarn},
		{"loglevel set to error", Config{LogLevel: 2}, slog.LevelError},
		{"loglevel set to debug", Config{LogLevel: 99}, slog.LevelDebug},
		{"loglevel wrong input", Config{LogLevel: 1337}, slog.LevelInfo},
		{"loglevel undefined", Config{}, slog.LevelInfo},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.GetLogLevel(); output != test.expected {
				t.Errorf("got %v, wanted %v", output, test.expected)
			}
		})
	}
}
