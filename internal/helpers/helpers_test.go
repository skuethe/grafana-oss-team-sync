//go:build !integration && !e2e

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

func TestAppendUnique(t *testing.T) {

	type addTest struct {
		name     string
		input    []string
		add      string
		expected []string
	}

	var tests = []addTest{
		{"new value gets appended", []string{"one", "two"}, "three", []string{"one", "two", "three"}},
		{"existing value is not duplicated", []string{"one", "two"}, "two", []string{"one", "two"}},
		{"existing value is case-insensitive", []string{"one", "two"}, "Two", []string{"one", "two"}},
		{"append to empty slice", []string{}, "one", []string{"one"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := AppendUnique(test.input, test.add); !reflect.DeepEqual(output, test.expected) {
				t.Errorf("got %q, wanted %q", output, test.expected)
			}
		})
	}
}
