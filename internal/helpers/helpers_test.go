// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package helpers

import (
	"reflect"
	"testing"
)

func TestRemoveFromSlice(t *testing.T) {

	type addTest struct {
		name          string
		input         []string
		remove        string
		casesensitive bool
		expected      []string
	}

	var tests = []addTest{
		{"two should be removed", []string{"one", "two"}, "two", false, []string{"one"}},
		{"three does not exist", []string{"one", "two"}, "three", false, []string{"one", "two"}},
		{"case sensitive Two should not be removed", []string{"one", "two"}, "Two", true, []string{"one", "two"}},
		{"empty slice", []string{}, "two", false, []string{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := RemoveFromSlice(test.input, test.remove, test.casesensitive); !reflect.DeepEqual(output, test.expected) {
				t.Errorf("got %q, wanted %q", output, test.expected)
			}
		})
	}
}
