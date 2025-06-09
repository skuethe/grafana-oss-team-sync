package grafana

import (
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/health"
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

	getScheme := fromConfigOrDefault("grafana.connection.scheme", "http")
	getHost := fromConfigOrDefault("grafana.connection.host", "localhost:3000")
	getBasePath := fromConfigOrDefault("grafana.connection.basePath", "/api")

	cfg := client.TransportConfig{
		// Host is the doman name or IP address of the host that serves the API.
		Host: getHost,
		// BasePath is the URL prefix for all API paths, relative to the host root.
		BasePath: getBasePath,
		// Schemes are the transfer protocols used by the API (http or https).
		Schemes: []string{getScheme},
		// APIKey is an optional API key or service account token.
		APIKey: os.Getenv("API_ACCESS_TOKEN"),
		// BasicAuth is optional basic auth credentials.
		BasicAuth: url.UserPassword("admin", "admin"),
	}

	client := client.NewHTTPClientWithConfig(strfmt.Default, &cfg)

	var health *health.GetHealthOK
	var healthErr error
	retryLoop := 0
	retryMax := config.K.Int(config.ConfigPathGrafana + ".connection.retry")

	grafanaLog.Info("connecting to Grafana instance",
		slog.Int("retry", retryMax),
	)
	for {
		grafanaLog.Debug("trying to establish connection",
			slog.Group("retry",
				slog.Int("loop", retryLoop),
				slog.Int("max", retryMax),
			),
		)
		health, healthErr = client.Health.GetHealth()
		if healthErr == nil || retryLoop == retryMax {
			break
		}
		retryLoop++
		time.Sleep(2 * time.Second)
	}
	if healthErr != nil {
		grafanaLog.Error("Grafana instance is not healthy",
			slog.Any("error", healthErr),
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
