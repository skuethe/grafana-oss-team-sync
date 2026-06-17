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

// AppendUnique appends addThis to appendHere only if it is not already present.
// The comparison is case-insensitive to avoid syncing the same group twice.
func AppendUnique(appendHere []string, addThis string) []string {
	for _, v := range appendHere {
		if strings.EqualFold(v, addThis) {
			return appendHere
		}
	}
	return append(appendHere, addThis)
}

// RemoveEmptyStrings returns a copy of the slice without any empty string entries.
// Empty entries can be introduced by the default value of comma-separated flags
// and would otherwise produce invalid (empty) source filters.
func RemoveEmptyStrings(slice []string) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}
