//go:build !integration && !e2e

// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import (
	"reflect"
	"testing"
)

func TestSanitizeTeamPrefixes(t *testing.T) {

	type addTest struct {
		name     string
		input    TeamPrefixes
		expected TeamPrefixes
	}

	var tests = []addTest{
		{"one empty entry is removed", TeamPrefixes{"one-", "", "two-"}, TeamPrefixes{"one-", "two-"}},
		{"multiple empty entries are removed", TeamPrefixes{"", "one-", "", "two-", ""}, TeamPrefixes{"one-", "two-"}},
		{"only empty entries are removed", TeamPrefixes{"", ""}, TeamPrefixes{}},
		{"slice without empty entries is unchanged", TeamPrefixes{"one-", "two-"}, TeamPrefixes{"one-", "two-"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := Config{TeamPrefixes: test.input}
			c.SanitizeTeamPrefixes()
			if !reflect.DeepEqual(c.TeamPrefixes, test.expected) {
				t.Errorf("got %q, wanted %q", c.TeamPrefixes, test.expected)
			}
		})
	}
}
