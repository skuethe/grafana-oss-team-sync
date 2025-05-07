package grafana

import (
	"log/slog"
	"net/url"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-openapi-client-go/client"
)

func initClient() *client.GrafanaHTTPAPI {
	clientLog := slog.With(slog.String("package", "grafana.client"))
	clientLog.Info("Initializing Grafana Client")

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

	return client.NewHTTPClientWithConfig(strfmt.Default, &cfg)
}
