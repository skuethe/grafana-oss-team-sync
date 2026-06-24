//go:build !integration && !e2e

// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import (
	"reflect"
	"testing"
)

func TestSanitizeTeams(t *testing.T) {

	type addTest struct {
		name     string
		input    Teams
		expected Teams
	}

	var tests = []addTest{
		{"one empty entry is removed", Teams{"one", "", "two"}, Teams{"one", "two"}},
		{"multiple empty entries are removed", Teams{"", "one", "", "two", ""}, Teams{"one", "two"}},
		{"only empty entries are removed", Teams{"", ""}, Teams{}},
		{"slice without empty entries is unchanged", Teams{"one", "two"}, Teams{"one", "two"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := Config{Teams: test.input}
			c.SanitizeTeams()
			if !reflect.DeepEqual(c.Teams, test.expected) {
				t.Errorf("got %q, wanted %q", c.Teams, test.expected)
			}
		})
	}
}
