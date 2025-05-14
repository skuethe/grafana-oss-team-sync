package azure

import (
	"log/slog"
	"os"

	graph "github.com/microsoftgraph/msgraph-sdk-go"
)

type azureInstance struct {
	api *graph.GraphServiceClient
}

func Start() {
	azureLog := slog.With(slog.String("package", "azure"))
	azureLog.Info("Initializing Azure")

	instance, err := new()
	if err != nil {
		azureLog.Error("Could not create MS Graph client", "error", err)
		os.Exit(1)
	}

	// azureLog.Info("Validating Azure Health")
	// _, err := instance.api.Health.GetHealth()
	// if err != nil {
	// 	grafanaLog.Error("Grafana instance is not healthy", "error", err)
	// 	os.Exit(1)
	// }

	// instance.processUsers(userList)
	teams := instance.processGroups()
	azureLog.Info("DEBUG", "return", teams)
}
