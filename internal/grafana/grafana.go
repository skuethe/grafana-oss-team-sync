package grafana

import (
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/config/types"
)

type GrafanaInstance struct {
	api *client.GrafanaHTTPAPI
}

var Instance *GrafanaInstance

func New() {
	grafanaLog := slog.With(slog.String("package", "grafana"))
	grafanaLog.Info("initializing Grafana")

	cfg := client.TransportConfig{
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

	if config.Instance.Grafana.AuthType == types.GrafanaAuthTypeToken {
		// APIKey is an API key or service account token
		cfg.APIKey = os.Getenv(types.GrafanaAuthVariableToken)

		if !config.Instance.Features.DisableUserSync {
			grafanaLog.Warn("token auth does not support creating new users. Switch to basic auth or disable the user sync feature")
		}
	} else {
		// BasicAuth is basic auth credentials.
		cfg.BasicAuth = url.UserPassword(os.Getenv(types.GrafanaAuthVariableUsername), os.Getenv(types.GrafanaAuthVariablePassword))
	}

	client := client.NewHTTPClientWithConfig(strfmt.Default, &cfg)

	grafanaLog.Info("connecting to Grafana instance",
		slog.Int("retry", config.Instance.Grafana.Connection.Retry),
	)
	health, err := client.Health.GetHealth()
	if err != nil {
		grafanaLog.Error("Grafana instance is not healthy",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	grafanaLog.Info("validated instance health",
		slog.String("version", health.Payload.Version),
	)

	Instance = &GrafanaInstance{
		api: client,
	}
}
