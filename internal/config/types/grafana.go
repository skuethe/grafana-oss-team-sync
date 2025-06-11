package types

import (
	"errors"
	"fmt"
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
	// GrafanaConnectionSchemeFlagShort string = ""
	// GrafanaConnectionSchemeVariable  string = "GOTS_"
	GrafanaConnectionSchemeDefault      string = "http"
	GrafanaConnectionSchemeFlagHelp     string = "the scheme of your Grafana instance\nAllowed: http or https"
	GrafanaConnectionSchemeParameter    string = "grafana.connection.scheme"
	GrafanaConnectionSchemeAllowedHTTP  string = "http"
	GrafanaConnectionSchemeAllowedHTTPS string = "https"

	// GrafanaConnectionHostVariable  string = "GOTS_"
	GrafanaConnectionHostDefault   string = "localhost:3000"
	GrafanaConnectionHostFlagHelp  string = "the host of your Grafana instance"
	GrafanaConnectionHostFlagShort string = "H"
	GrafanaConnectionHostParameter string = "grafana.connection.host"

	// GrafanaConnectionBasePathFlagShort string = ""
	// GrafanaConnectionBasePathVariable  string = "GOTS_"
	GrafanaConnectionBasePathDefault   string = "/api"
	GrafanaConnectionBasePathFlagHelp  string = "the base path of your Grafana instance"
	GrafanaConnectionBasePathParameter string = "grafana.connection.basepath"

	// GrafanaConnectionRetryVariable  string = "GOTS_"
	GrafanaConnectionRetryDefault   int    = 0
	GrafanaConnectionRetryFlagHelp  string = "the amount of retries to connect to your Grafana instance\nRetry timeout is set to 2 seconds\nDefault: 0"
	GrafanaConnectionRetryFlagShort string = "r"
	GrafanaConnectionRetryParameter string = "grafana.connection.retry"

	GrafanaAuthVariableToken    string = "GOTS_TOKEN"
	GrafanaAuthVariableUsername string = "GOTS_USERNAME"
	GrafanaAuthVariablePassword string = "GOTS_PASSWORD"

	GrafanaAuthTypeToken     string = "token"
	GrafanaAuthTypeBasicAuth string = "basicauth"
)

const (
	GrafanaPermissionViewer GrafanaPermission = 1 << iota
	GrafanaPermissionEditor
	GrafanaPermissionAdmin
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

	return errors.New("invalid permission defined: " + fmt.Sprint(in))
}

func (c *Config) ValdidateGrafanaScheme() error {
	switch c.Grafana.Connection.Scheme {
	case GrafanaConnectionSchemeAllowedHTTP:
		return nil
	case GrafanaConnectionSchemeAllowedHTTPS:
		return nil
	}

	return errors.New("invalid connection scheme defined: " + c.Grafana.Connection.Scheme)
}
