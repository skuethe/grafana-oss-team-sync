// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import (
	"errors"
	"testing"
)

func TestString(t *testing.T) {

	type addTest struct {
		name     string
		input    Source
		expected string
	}

	var tests = []addTest{
		{"existing source", SourcePluginEntraID, "entraid"},
		{"non-existing source", Source("valid"), "valid"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.String(); output != test.expected {
				t.Errorf("got %v, wanted %v", output, test.expected)
			}
		})
	}
}

func TestGetSource(t *testing.T) {

	type addTest struct {
		name     string
		input    Config
		expected Source
	}

	var tests = []addTest{
		{"existing source via const", Config{Source: SourcePluginEntraID}, SourcePluginEntraID},
		{"existing source via string", Config{Source: SourcePluginEntraID}, Source("entraid")},
		{"non-existing source", Config{Source: Source("invalid")}, Source("invalid")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.GetSource(); output != test.expected {
				t.Errorf("got %v, wanted %v", output, test.expected)
			}
		})
	}
}

func TestValidateSourcePlugin(t *testing.T) {

	type addTest struct {
		name     string
		input    Config
		expected error
	}

	var tests = []addTest{
		{"existing source", Config{Source: SourcePluginEntraID}, nil},
		{"non-existing source", Config{Source: Source("invalid")}, ErrInvalidSourcePlugin},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.ValdidateSourcePlugin(); !errors.Is(output, test.expected) {
				t.Errorf("got %v, wanted %v", output, test.expected)
			}
		})
	}
}
