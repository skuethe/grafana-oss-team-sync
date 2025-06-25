// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

import (
	"errors"
	"fmt"
	"strings"
)

type Grafana struct {
	AuthType   string            `yaml:"authtype"`
	Connection GrafanaConnection `yaml:"connection"`
}

type GrafanaConnection struct {
	Scheme   string `yaml:"scheme"`
	Host     string `yaml:"host"`
	BasePath string `yaml:"basepath"`
	Retry    int    `yaml:"retry"`
}

type GrafanaPermission int

const (
	GrafanaAuthTypeDefault   string = "basicauth"
	GrafanaAuthTypeFlagHelp  string = "the auth type you want to use to authenticate against Grafana\nAllowed: basicauth or token"
	GrafanaAuthTypeOptimized string = "authtype"
	GrafanaAuthTypeParameter string = "grafana.authtype"

	GrafanaBasicAuthUsernameDefault   string = ""
	GrafanaBasicAuthUsernameFlagHelp  string = "the basic auth user you want to use to authenticate against Grafana\nOnly used if authtype is set to basicauth"
	GrafanaBasicAuthUsernameFlagShort string = "u"
	GrafanaBasicAuthUsernameParameter string = "username"
	GrafanaBasicAuthUsernameVariable  string = "GOTS_USERNAME"

	GrafanaBasicAuthPasswordDefault   string = ""
	GrafanaBasicAuthPasswordFlagHelp  string = "the basic auth password you want to use to authenticate against Grafana\nOnly used if authtype is set to basicauth"
	GrafanaBasicAuthPasswordFlagShort string = "p"
	GrafanaBasicAuthPasswordParameter string = "password"
	GrafanaBasicAuthPasswordVariable  string = "GOTS_PASSWORD"

	GrafanaTokenAuthDefault   string = ""
	GrafanaTokenAuthFlagHelp  string = "the service account token you want to use to authenticate against Grafana\nOnly used if authtype is set to token"
	GrafanaTokenAuthFlagShort string = "t"
	GrafanaTokenAuthParameter string = "token"
	GrafanaTokenAuthVariable  string = "GOTS_TOKEN"

	GrafanaConnectionSchemeDefault   string = "http"
	GrafanaConnectionSchemeFlagHelp  string = "the scheme of your Grafana instance\nAllowed: http or https"
	GrafanaConnectionSchemeOptimized string = "scheme"
	GrafanaConnectionSchemeParameter string = "grafana.connection.scheme"

	GrafanaConnectionSchemeAllowedHTTP  string = "http"
	GrafanaConnectionSchemeAllowedHTTPS string = "https"

	GrafanaConnectionHostDefault   string = "localhost:3000"
	GrafanaConnectionHostFlagHelp  string = "the host of your Grafana instance"
	GrafanaConnectionHostFlagShort string = "H"
	GrafanaConnectionHostOptimized string = "host"
	GrafanaConnectionHostParameter string = "grafana.connection.host"

	GrafanaConnectionBasePathDefault   string = "/api"
	GrafanaConnectionBasePathFlagHelp  string = "the base path of your Grafana instance"
	GrafanaConnectionBasePathOptimized string = "basepath"
	GrafanaConnectionBasePathParameter string = "grafana.connection.basepath"

	GrafanaConnectionRetryDefault   int    = 0
	GrafanaConnectionRetryFlagHelp  string = "the amount of retries to connect to your Grafana instance\nRetry timeout is set to 2 seconds\nDefault: 0"
	GrafanaConnectionRetryFlagShort string = "r"
	GrafanaConnectionRetryOptimized string = "retry"
	GrafanaConnectionRetryParameter string = "grafana.connection.retry"

	GrafanaAuthTypeToken     string = "token"
	GrafanaAuthTypeBasicAuth string = "basicauth"
)

const (
	GrafanaPermissionViewer GrafanaPermission = 1 << iota
	GrafanaPermissionEditor
	GrafanaPermissionAdmin
)

var (
	ErrInvalidPermission       = errors.New("invalid permission defined")
	ErrInvalidAuthType         = errors.New("invalid authtype defined")
	ErrInvalidConnectionScheme = errors.New("invalid connection scheme defined")
)

func ValidateGrafanaPermission(in GrafanaPermission) error {
	switch in {
	case GrafanaPermissionViewer:
		return nil
	case GrafanaPermissionEditor:
		return nil
	case GrafanaPermissionAdmin:
		return nil
	}

	return fmt.Errorf("%w: %d", ErrInvalidPermission, in)
}

func (c *Config) ValdidateGrafanaAuthType() error {
	switch c.Grafana.AuthType {
	case GrafanaAuthTypeToken:
		return nil
	case GrafanaAuthTypeBasicAuth:
		return nil
	}

	return fmt.Errorf("%w: %q", ErrInvalidAuthType, c.Grafana.AuthType)
}

func (c *Config) ValdidateGrafanaScheme() error {
	switch strings.ToLower(c.Grafana.Connection.Scheme) {
	case GrafanaConnectionSchemeAllowedHTTP:
		return nil
	case GrafanaConnectionSchemeAllowedHTTPS:
		return nil
	}

	return fmt.Errorf("%w: %q", ErrInvalidConnectionScheme, c.Grafana.Connection.Scheme)
}
