package grafana

import (
	"log/slog"
	"net/url"
	"os"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-openapi-client-go/client"
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
		Host: "localhost:3000",
		// BasePath is the URL prefix for all API paths, relative to the host root.
		BasePath: "/api",
		// Schemes are the transfer protocols used by the API (http or https).
		Schemes: []string{"http"},
		// BasicAuth is optional basic auth credentials.
		BasicAuth: url.UserPassword("admin", "admin"),
	}

	client := client.NewHTTPClientWithConfig(strfmt.Default, &cfg)

	health, err := client.Health.GetHealth()
	if err != nil {
		grafanaLog.Error("Grafana instance is not healthy", "error", err)
		os.Exit(1)
	}
	grafanaLog.Info("validated instance health", "version", health.Payload.Version)

	Instance = &GrafanaInstance{
		api: client,
	}
}
