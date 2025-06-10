package grafana

import (
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
)

type GrafanaInstance struct {
	api *client.GrafanaHTTPAPI
}

var Instance *GrafanaInstance

func fromConfigOrDefault(configPath string, fallback string) string {
	configValue := config.K.String(configPath)
	if configValue == "" {
		return fallback
	}
	return configValue
}

func New() {
	grafanaLog := slog.With(slog.String("package", "grafana"))
	grafanaLog.Info("initializing Grafana")

	getAuth := fromConfigOrDefault(config.ConfigParamGrafana+".auth", "basicauth")
	getScheme := fromConfigOrDefault(config.ConfigParamGrafana+".connection.scheme", "http")
	getHost := fromConfigOrDefault(config.ConfigParamGrafana+".connection.host", "localhost:3000")
	getBasePath := fromConfigOrDefault(config.ConfigParamGrafana+".connection.basePath", "/api")

	getRetry := config.K.Int(config.ConfigParamGrafana + ".connection.retry")

	cfg := client.TransportConfig{
		// Host is the doman name or IP address of the host that serves the API.
		Host: getHost,
		// BasePath is the URL prefix for all API paths, relative to the host root.
		BasePath: getBasePath,
		// Schemes are the transfer protocols used by the API (http or https).
		Schemes: []string{getScheme},
		// NumRetries contains the optional number of attempted retries
		NumRetries: getRetry,
		// RetryTimeout sets an time to wait before retrying a request
		RetryTimeout: (2 * time.Second),
	}

	if getAuth == "token" {
		// APIKey is an API key or service account token
		cfg.APIKey = config.K.MustString(config.ConfigParamAuthToken)

		if !config.Feature.DisableUserSync {
			grafanaLog.Warn("token auth does not support creating new users. Switch to basic auth or disable the user sync feature")
		}
	} else {
		// BasicAuth is basic auth credentials.
		cfg.BasicAuth = url.UserPassword(config.K.MustString(config.ConfigParamAuthBasicUsername), config.K.MustString(config.ConfigParamAuthBasicPassword))
	}

	client := client.NewHTTPClientWithConfig(strfmt.Default, &cfg)

	grafanaLog.Info("connecting to Grafana instance",
		slog.Int("retry", getRetry),
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
