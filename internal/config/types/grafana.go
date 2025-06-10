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
	GrafanaAuthVariableToken    string = "GOTS_TOKEN"
	GrafanaAuthVariableUsername string = "GOTS_USERNAME"
	GrafanaAuthVariablePassword string = "GOTS_PASSWORD"

	GrafanaConnectionSchemeHTTP  string = "http"
	GrafanaConnectionSchemeHTTPS string = "https"

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
	case GrafanaConnectionSchemeHTTP:
		return nil
	case GrafanaConnectionSchemeHTTPS:
		return nil
	}

	return errors.New("invalid source: " + c.Source.String())
}
