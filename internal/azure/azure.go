package azure

import (
	"log/slog"

	graph "github.com/microsoftgraph/msgraph-sdk-go"
)

type AzureInstance struct {
	api *graph.GraphServiceClient
}

func Start() {
	azureLog := slog.With(slog.String("package", "azure"))
	azureLog.Info("Initializing Azure")

	instance := AzureInstance{
		api: initClient(),
	}

	// azureLog.Info("Validating Azure Health")
	// _, err := instance.api.Health.GetHealth()
	// if err != nil {
	// 	grafanaLog.Error("Grafana instance is not healthy", "error", err)
	// 	os.Exit(1)
	// }

	// instance.processUsers(userList)
	// instance.processTeams(teamList)
	instance.processGroups()
}
