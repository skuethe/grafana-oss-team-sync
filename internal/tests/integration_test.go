//go:build integration

// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package tests

import (
	"os"
	"testing"

	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
)

func TestIntegrationGrafana(t *testing.T) {

	// 1. set config via env var
	os.Clearenv()
	if err := os.Setenv("GOTS_CONFIG", "../../test/data/integration-tests_config.yaml"); err != nil {
		t.Fatal(ErrCouldNotSetConfigVariable, err)
	}
	if err := os.Setenv("GOTS_AUTHFILE", "../../test/data/integration-tests_authfile.env"); err != nil {
		t.Fatal(ErrCouldNotSetConfigVariable, err)
	}

	// Create grafana.Teams for integration tests
	teamList := &grafana.Teams{
		grafana.Team{
			Parameter: &grafana.TeamParameter{
				Name: "group-integration-1",
			},
			Users: &grafana.Users{
				grafana.User{
					Email: "user1@example.com",
					Login: "user1@example.com",
					Name:  "User One",
				},
				grafana.User{
					Email: "user2@example.com",
					Login: "user2@example.com",
					Name:  "User Two",
				},
				grafana.User{
					Email: "",
					Login: "user3@example.com",
					Name:  "User Three",
				},
			},
		},
		grafana.Team{
			Parameter: &grafana.TeamParameter{
				Name: "group-integration-2",
			},
			Users: &grafana.Users{
				grafana.User{
					Email: "user1@example.com",
					Login: "user1@example.com",
					Name:  "User One",
				},
				grafana.User{
					Email: "user4@example.com",
					Login: "user4@example.com",
					Name:  "User Four",
				},
			},
		},
		grafana.Team{
			Parameter: &grafana.TeamParameter{
				Name: "group-integration-3",
			},
			Users: &grafana.Users{},
		},
	}

	// Run integration routine
	if err := IntegrationGrafana(teamList); err != nil {
		t.Fatal(err)
	}
}
