// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package grafana

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/config/configtypes"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
)

type GrafanaInstance struct {
	api *client.GrafanaHTTPAPI
	// defaultOrgID is the org the configured credentials authenticate into by default.
	defaultOrgID int64
	// currentOrgID tracks the org the sync user is currently switched into, to avoid redundant API calls.
	currentOrgID int64
	// syncUserLogin is the login of the authenticated (basic auth) sync user, used to grant it
	// membership in other orgs so it can be switched into them. Empty for token auth.
	syncUserLogin string
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
	ErrCouldNotSwitchOrg            = errors.New("could not switch active organization")
	ErrCouldNotEnsureOrgMembership  = errors.New("could not ensure sync user org membership")
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
	currentOrg, err := client.Org.GetCurrentOrg()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotFetchOrgDetails, err)
	}
	grafanaLog.Info("successfully authenticated against Grafana",
		slog.Group("org",
			slog.Int64("id", currentOrg.Payload.ID),
			slog.String("name", currentOrg.Payload.Name),
		),
	)

	Instance = &GrafanaInstance{
		api:          client,
		defaultOrgID: currentOrg.Payload.ID,
		currentOrgID: currentOrg.Payload.ID,
	}

	// Only basic auth can switch orgs (API keys/service account tokens are pinned to one org).
	if config.Instance.Grafana.AuthType == configtypes.GrafanaAuthTypeBasicAuth {
		if signedInUser, err := client.SignedInUser.GetSignedInUser(); err != nil {
			grafanaLog.Warn("could not fetch signed-in user, syncing teams/folders into other orgs will fail",
				slog.Any("error", err),
			)
		} else {
			Instance.syncUserLogin = signedInUser.Payload.Login
		}
	}

	return nil
}

// EnsureOrgContext switches the Grafana API session to operate against orgID.
//
// Grafana's Team and Folder APIs always act on the caller's "current organization" and
// provide no per-request orgId parameter. Grafana OSS does not honor an X-Grafana-Org-Id
// header on Basic Auth requests, so the only way to target another org is to explicitly
// switch the authenticated user's active org (which Grafana persists server-side) before
// issuing further requests. orgID of 0 means "the default org the credentials log into".
func (g *GrafanaInstance) EnsureOrgContext(orgID int64) error {
	if orgID == 0 {
		orgID = g.defaultOrgID
	}
	if orgID == g.currentOrgID {
		return nil
	}

	if err := g.ensureSyncUserOrgMembership(orgID); err != nil {
		return fmt.Errorf("%w (org %d): %w", ErrCouldNotEnsureOrgMembership, orgID, err)
	}

	if _, err := g.api.SignedInUser.UserSetUsingOrg(orgID); err != nil {
		return fmt.Errorf("%w (org %d): %w", ErrCouldNotSwitchOrg, orgID, err)
	}
	g.currentOrgID = orgID
	return nil
}

// ensureSyncUserOrgMembership makes sure the authenticated sync user is a member (with Admin
// role) of orgID, which is required both to switch into that org and to manage its teams/folders.
func (g *GrafanaInstance) ensureSyncUserOrgMembership(orgID int64) error {
	if g.syncUserLogin == "" {
		return nil
	}

	orgUsers, err := g.api.Orgs.GetOrgUsers(orgID)
	if err != nil {
		return err
	}
	for _, orgUser := range orgUsers.Payload {
		if strings.EqualFold(orgUser.Login, g.syncUserLogin) {
			return nil
		}
	}

	_, err = g.api.Orgs.AddOrgUser(orgID, &models.AddOrgUserCommand{
		LoginOrEmail: g.syncUserLogin,
		Role:         "Admin",
	})
	return err
}
