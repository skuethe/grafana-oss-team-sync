package grafana

import (
	"log/slog"
	"net/url"
	"os"

	"github.com/go-openapi/strfmt"
	goapi "github.com/grafana/grafana-openapi-client-go/client"
)

func InitClient() *goapi.GrafanaHTTPAPI {
	cfg := goapi.TransportConfig{
		// Host is the doman name or IP address of the host that serves the API.
		Host: "localhost:3000",
		// BasePath is the URL prefix for all API paths, relative to the host root.
		BasePath: "/api",
		// Schemes are the transfer protocols used by the API (http or https).
		Schemes: []string{"http"},
		// BasicAuth is optional basic auth credentials.
		BasicAuth: url.UserPassword("admin", "admin"),
	}

	api := goapi.NewHTTPClientWithConfig(strfmt.Default, &cfg)

	_, err := api.Health.GetHealth()
	if err != nil {
		slog.Error("Grafana instance not healthy", "error", err)
		os.Exit(1)
	}
	return api
}
