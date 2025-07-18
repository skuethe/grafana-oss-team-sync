// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package grafana

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/config/configtypes"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
)

type GrafanaInstance struct {
	api *client.GrafanaHTTPAPI
}

var (
	Instance *GrafanaInstance

	ErrCouldNotEnableAuthentication = errors.New("could not enable authentication")
	ErrInstanceNotHealthy           = errors.New("instance is not healthy")
	ErrCouldNotFetchOrgDetails      = errors.New("could not fetch org details from specified auth")
	ErrAuthTokenMissing             = errors.New("token auth specified, but token is missing")
	ErrAuthBasicUsernameMissing     = errors.New("basic auth specified, but username is missing")
	ErrAuthBasicPasswordMissing     = errors.New("basic auth specified, but password is missing")
	ErrAuthUnsupported              = errors.New("unsupported authentication type defined")
)

// We are explicitly handling auth data here, because we do not want to add it to our global config.Instance
func setAuthData(c *client.TransportConfig) error {
	switch config.Instance.Grafana.AuthType {
	// Handle token auth
	case configtypes.GrafanaAuthTypeToken:
		token := ""

		// First fetch from OS env
		token = os.Getenv(configtypes.GrafanaTokenAuthVariable)

		// Override if flag is set
		if flags.Token != "" {
			token = flags.Token
		}

		// Return error if token not defined
		if token == "" {
			return ErrAuthTokenMissing
		}

		// APIKey is an API key (deprecated) or service account token
		c.APIKey = token

		if !config.Instance.Features.DisableUserSync {
			slog.Warn("token auth does not support creating new users. Switch to basic auth or disable the user sync feature")
		}
	// Handle basic auth
	case configtypes.GrafanaAuthTypeBasicAuth:
		username := ""
		password := ""

		// First fetch from OS env
		username = os.Getenv(configtypes.GrafanaBasicAuthUsernameVariable)
		password = os.Getenv(configtypes.GrafanaBasicAuthPasswordVariable)

		// Override if flag is set
		if flags.BasicAuthUsername != "" {
			username = flags.BasicAuthUsername
		}
		if flags.BasicAuthPassword != "" {
			password = flags.BasicAuthPassword
		}

		// Return error if token not defined
		if username == "" {
			return ErrAuthBasicUsernameMissing
		}
		if password == "" {
			return ErrAuthBasicPasswordMissing
		}

		// BasicAuth is basic auth credentials.
		c.BasicAuth = url.UserPassword(username, password)

	// Something went wrong, this should not happen...
	default:
		return fmt.Errorf("%w: %q", ErrAuthUnsupported, config.Instance.Grafana.AuthType)
	}

	return nil
}

func New() error {
	grafanaLog := slog.With(slog.String("package", "grafana"))
	grafanaLog.Info("initializing Grafana")

	cfg := &client.TransportConfig{
		// Host is the doman name or IP address of the host that serves the API.
		Host: config.Instance.Grafana.Connection.Host,
		// BasePath is the URL prefix for all API paths, relative to the host root.
		BasePath: config.Instance.Grafana.Connection.BasePath,
		// Schemes are the transfer protocols used by the API (http or https).
		Schemes: []string{config.Instance.Grafana.Connection.Scheme},
		// NumRetries contains the optional number of attempted retries
		NumRetries: config.Instance.Grafana.Connection.Retry,
		// RetryTimeout sets an time to wait before retrying a request
		RetryTimeout: (2 * time.Second),
	}

	// Add authentication data based on config input
	if err := setAuthData(cfg); err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotEnableAuthentication, err)
	}

	client := client.NewHTTPClientWithConfig(strfmt.Default, cfg)

	grafanaLog.Info("connecting to Grafana instance",
		slog.Int("retry", config.Instance.Grafana.Connection.Retry),
	)

	// Validate Grafana health
	if health, err := client.Health.GetHealth(); err != nil {
		return fmt.Errorf("%w: %w", ErrInstanceNotHealthy, err)
	} else {
		grafanaLog.Info("validated instance health",
			slog.String("version", health.Payload.Version),
		)
	}

	// Fetching current org here for additional information AND to fail fast on auth errors
	if currentOrg, err := client.Org.GetCurrentOrg(); err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotFetchOrgDetails, err)
	} else {
		grafanaLog.Info("successfully authenticated against Grafana",
			slog.Group("org",
				slog.Int64("id", currentOrg.Payload.ID),
				slog.String("name", currentOrg.Payload.Name),
			),
		)
	}

	Instance = &GrafanaInstance{
		api: client,
	}
	return nil
}
