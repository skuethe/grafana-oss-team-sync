// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

type AuthFile string

const (
	AuthFileDefault   string = ""
	AuthFileFlagHelp  string = "an optional file to load authentication data from. File content needs to be in key=value syntax"
	AuthFileFlagShort string = "a"
	AuthFileParameter string = "authfile"
	AuthFileVariable  string = "GOTS_AUTHFILE"
)

func (c *Config) IsAuthFileSet() bool {
	return c.AuthFile != ""
}
