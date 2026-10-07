// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package helpers

import (
	"slices"
	"strings"
)

func RemoveFromSlice(searchHere []string, removeThis string, caseSensitive bool) []string {
	for k, v := range searchHere {
		var validated bool
		if caseSensitive {
			validated = v == removeThis
		} else {
			validated = strings.EqualFold(v, removeThis)
		}
		if validated {
			return slices.Delete(searchHere, k, k+1)
		}
	}
	return searchHere
}
