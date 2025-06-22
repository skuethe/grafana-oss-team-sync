//go:build !integration && !e2e

// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import (
	"errors"
	"testing"
)

func TestValidateGrafanaPermission(t *testing.T) {

	type addTest struct {
		name     string
		input    GrafanaPermission
		expected error
	}

	var tests = []addTest{
		{"permission viewer via const", GrafanaPermissionViewer, nil},
		{"permission viewer via int", GrafanaPermission(1), nil},
		{"permission editor via const", GrafanaPermissionEditor, nil},
		{"permission editor via int", GrafanaPermission(2), nil},
		{"permission admin via const", GrafanaPermissionAdmin, nil},
		{"permission admin via int", GrafanaPermission(4), nil},
		{"invalid permission", GrafanaPermission(0), ErrInvalidPermission},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := ValidateGrafanaPermission(test.input); !errors.Is(output, test.expected) {
				t.Errorf("got error: %v, wanted error: %v", output, test.expected)
			}
		})
	}
}

func TestValdidateGrafanaAuthType(t *testing.T) {

	type addTest struct {
		name     string
		input    Config
		expected error
	}

	var tests = []addTest{
		{"auth type token via const", Config{Grafana: Grafana{AuthType: GrafanaAuthTypeToken}}, nil},
		{"auth type token via string", Config{Grafana: Grafana{AuthType: "token"}}, nil},
		{"auth type basicauth via const", Config{Grafana: Grafana{AuthType: GrafanaAuthTypeBasicAuth}}, nil},
		{"auth type basicauth via string", Config{Grafana: Grafana{AuthType: "basicauth"}}, nil},
		{"invalid auth type", Config{Grafana: Grafana{AuthType: "invalid"}}, ErrInvalidAuthType},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.ValdidateGrafanaAuthType(); !errors.Is(output, test.expected) {
				t.Errorf("got error: %v, wanted error: %v", output, test.expected)
			}
		})
	}
}

func TestValdidateGrafanaScheme(t *testing.T) {

	type addTest struct {
		name     string
		input    Config
		expected error
	}

	var tests = []addTest{
		{"scheme HTTP via const", Config{Grafana: Grafana{Connection: GrafanaConnection{Scheme: GrafanaConnectionSchemeAllowedHTTP}}}, nil},
		{"scheme HTTP via string", Config{Grafana: Grafana{Connection: GrafanaConnection{Scheme: "http"}}}, nil},
		{"scheme HTTP via string, uppercase", Config{Grafana: Grafana{Connection: GrafanaConnection{Scheme: "HTTP"}}}, nil},
		{"scheme HTTPS via const", Config{Grafana: Grafana{Connection: GrafanaConnection{Scheme: GrafanaConnectionSchemeAllowedHTTPS}}}, nil},
		{"scheme HTTPS via string", Config{Grafana: Grafana{Connection: GrafanaConnection{Scheme: "https"}}}, nil},
		{"scheme HTTPS via string, uppercase", Config{Grafana: Grafana{Connection: GrafanaConnection{Scheme: "HTTPS"}}}, nil},
		{"invalid scheme", Config{Grafana: Grafana{Connection: GrafanaConnection{Scheme: "invalid"}}}, ErrInvalidConnectionScheme},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.ValdidateGrafanaScheme(); !errors.Is(output, test.expected) {
				t.Errorf("got error: \"%v\", wanted error: \"%v\"", output, test.expected)
			}
		})
	}
}
