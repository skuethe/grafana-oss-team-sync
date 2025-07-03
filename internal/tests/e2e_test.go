//go:build e2e

// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package tests

import (
	"fmt"
	"os"
	"testing"
)

func TestEndToEnd(t *testing.T) {

	type addTest struct {
		name          string
		inputsource   string
		inputauthfile string
		expected      error
	}

	var tests = []addTest{
		{"entraid - mock data via dev-proxy", "entraid", "entraid-mock", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Set required env var
			os.Clearenv()
			if err := os.Setenv("GOTS_CONFIG", "../../test/data/e2e-tests_config.yaml"); err != nil {
				t.Fatal(fmt.Errorf("%w (%v): %w", ErrCouldNotSetConfigVariable, "GOTS_CONFIG", err))
			}
			if err := os.Setenv("GOTS_AUTHFILE", "../../test/data/e2e-tests_"+test.fileinput+"_authfile.env"); err != nil {
				t.Fatal(fmt.Errorf("%w (%v): %w", ErrCouldNotSetConfigVariable, "GOTS_AUTHFILE", err))
			}
			if err := os.Setenv("HTTP_PROXY", "localhost:8000"); err != nil {
				t.Fatal(fmt.Errorf("%w (%v): %w", ErrCouldNotSetConfigVariable, "HTTP_PROXY", err))
			}
			if err := os.Setenv("HTTPS_PROXY", "localhost:8000"); err != nil {
				t.Fatal(fmt.Errorf("%w (%v): %w", ErrCouldNotSetConfigVariable, "HTTPS_PROXY", err))
			}

			// Run e2e routine
			if err := EndToEnd(); err != nil {
				t.Errorf("got error: %v, wanted error: %v", err, test.expected)
			}

			// TODO?: Validate group, user and folder where created
		})
	}

}
