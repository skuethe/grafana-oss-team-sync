// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import (
	"errors"
	"fmt"
)

type Source string

const (
	SourceDefault   string = "entraid"
	SourceFlagHelp  string = "the source `plugin` you want to use\nAllowed: entraid, ldap"
	SourceFlagShort string = "s"
	SourceParameter string = "source"
	SourceVariable  string = "GOTS_SOURCE"

	SourcePluginEntraID Source = "entraid"
	SourcePluginLDAP    Source = "ldap"
)

var (
	ErrInvalidSourcePlugin = errors.New("invalid source plugin defined")
)

func (s Source) String() string {
	return string(s)
}

func (c *Config) ValdidateSourcePlugin() error {
	switch c.Source {
	case SourcePluginEntraID:
		return nil
	case SourcePluginLDAP:
		return nil
	}

	return fmt.Errorf("%w: %q", ErrInvalidSourcePlugin, c.Source.String())
}

func (c *Config) GetSource() Source {
	return c.Source
}
