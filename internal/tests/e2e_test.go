//go:build e2e

// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package tests

import (
	"testing"
)

func TestEndToEnd(t *testing.T) {

	// Run e2e routine
	if err := EndToEnd(); err != nil {
		t.Errorf("got error: %v", err)
	}

	// TODO?: Validate group, user and folder where created
}
