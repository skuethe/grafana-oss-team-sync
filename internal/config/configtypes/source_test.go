//go:build !integration && !e2e

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
		{"existing source entraid", SourcePluginEntraID, "entraid"},
		{"existing source ldap", SourcePluginLDAP, "ldap"},
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
		{"existing source via const entraid", Config{Source: SourcePluginEntraID}, SourcePluginEntraID},
		{"existing source via const ldap", Config{Source: SourcePluginLDAP}, SourcePluginLDAP},
		{"existing source via string entraid", Config{Source: SourcePluginEntraID}, Source("entraid")},
		{"existing source via string ldap", Config{Source: SourcePluginLDAP}, Source("ldap")},
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
		{"existing source entraid", Config{Source: SourcePluginEntraID}, nil},
		{"existing source ldap", Config{Source: SourcePluginLDAP}, nil},
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
