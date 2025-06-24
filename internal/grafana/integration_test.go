//go:build integration

// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package grafana

import (
	"errors"
	"os"
	"testing"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
)

var (
	ErrCouldNotSetConfigVariable = errors.New("could not set required environment variable for GOTS_CONFIG")
	ErrCouldNotSetHostVariable   = errors.New("could not set required environment variable for GOTS_HOST")
)

func TestIntegrationGrafanaFolders(t *testing.T) {

	// 1. set config via env var
	os.Clearenv()
	if err := os.Setenv("GOTS_CONFIG", "../../test/data/integration-tests_config.yaml"); err != nil {
		t.Fatal(ErrCouldNotSetConfigVariable, err)
	}
	if err := os.Setenv("GOTS_AUTHFILE", "../../test/data/integration-tests_authfile.env"); err != nil {
		t.Fatal(ErrCouldNotSetConfigVariable, err)
	}

	// Load flags to not fail
	flags.Load()

	// 2. parse config
	config.Load()

	// 3. init Grafana
	New()

	// Run ProcessFolders
	Instance.ProcessFolders()

}

// func TestIntegrationGrafanaTeams(t *testing.T) {
// }

// func TestIntegrationGrafanaUsers(t *testing.T) {
// }
