//go:build !integration && !e2e

// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package entraid

import (
	"testing"
)

func TestBuildGroupFilter(t *testing.T) {

	type addTest struct {
		name     string
		input    []string
		expected string
	}

	var tests = []addTest{
		{"one team", []string{"teamA"}, `displayName in ('teamA')`},
		{"two teams", []string{"teamA", "teamB"}, `displayName in ('teamA', 'teamB')`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := buildGroupFilter(test.input); output != test.expected {
				t.Errorf("got %q, wanted %q", output, test.expected)
			}
		})
	}
}

func TestBuildGroupSearch(t *testing.T) {

	type addTest struct {
		name     string
		input    []string
		expected string
	}

	var tests = []addTest{
		{"one team", []string{"teamA"}, `"displayName:teamA"`},
		{"two teams", []string{"teamA", "teamB"}, `"displayName:teamA" OR "displayName:teamB"`},
		{"team with spaces", []string{"team A"}, `"displayName:team A"`},
		{"double quote is escaped", []string{`team "A"`}, `"displayName:team \"A\""`},
		{"backslash is escaped", []string{`team\A`}, `"displayName:team\\A"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := buildGroupSearch(test.input); output != test.expected {
				t.Errorf("got %q, wanted %q", output, test.expected)
			}
		})
	}
}
